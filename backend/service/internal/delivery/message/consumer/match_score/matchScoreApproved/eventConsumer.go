package matchScoreApproved

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/rabbitmq/amqp091-go"
	"github.com/turanberker/tennis-league-service/internal/delivery/message/consumer"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type EventConsumer struct {
	*consumer.Consumer
	tm                    *database.TransactionManager
	matchRepo             match.Repository
	playerRepo            player.Repository
	earnedPointRepository Repository
}

func NewEventConsumer(tm *database.TransactionManager,
	matchRepo match.Repository,
	playerRepo player.Repository,
	earnedPointRepository Repository) *EventConsumer {

	c := &EventConsumer{
		tm:                    tm,
		matchRepo:             matchRepo,
		playerRepo:            playerRepo,
		earnedPointRepository: earnedPointRepository,
	}

	c.Consumer = &consumer.Consumer{
		Queue:       "update_player_points_queue",
		RoutingName: "MatchApproved",
		Handler:     c.handle, // 👈 struct method
	}

	return c

}

func (c *EventConsumer) handle(msg amqp091.Delivery) error {
	ctx := context.Background()

	var event = &match.MatchApprovedEvent{}

	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		return err
	}

	return c.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		matchType, err := c.matchRepo.GetMatchType(txCtx, event.MatchID)

		if err != nil {
			//Log yaz
			return err
		}

		switch *matchType {
		case match.MatchType_DOUBLE:
			c.updateUserDoublePoints(txCtx, event)
		case match.MatchType_SINGLE:
			log.Panic("Not Implemented")
		}

		log.Println("Match Approved:", event.MatchID)
		return nil
	})

}

func (c *EventConsumer) updateUserDoublePoints(txCtx context.Context, event *match.MatchApprovedEvent) error {
	participants, err := c.matchRepo.GetDoubleMatchParticipantsWithPoints(txCtx, event.MatchID)
	if err != nil {
		//Log yaz
		return err
	}

	change, err := c.calculateDoubleMatchElo(participants)
	if err != nil {
		//Log yaz
		return err
	}
	for _, p := range participants {
		if p.IsWinner {
			_, err = c.playerRepo.IncreaseDoublePoint(txCtx, p.PlayerID, *change)
			if err != nil {
				log.Printf("Kazanan Oyuncu Puanı Güncellenirken Hata Oluştu(Player: %s): %v", p.PlayerID, err)
				return err
			}
		} else {
			_, err = c.playerRepo.DecreaseDoublePoint(txCtx, p.PlayerID, *change)
			if err != nil {
				log.Printf("Kaybeden Oyuncu Puanı Güncellenirken Hata Oluştu(Player: %s): %v", p.PlayerID, err)
				return err
			}
		}

		pointChange := *change
		recordPoint := int32(pointChange)
		if !p.IsWinner {
			recordPoint = -recordPoint
		}

		err = c.earnedPointRepository.AddPlayerPoint(txCtx, &AddPlayerPoint{
			PlayerId:    p.PlayerID,
			EarnedPoint: recordPoint,
			MatchType:   match.MatchType_DOUBLE,
		})

		if err != nil {
			log.Printf("Puan geçmişi oluşturulurken hata oluştu (Player: %s): %v", p.PlayerID, err)
			return err
		}
	}
	return nil
}

func (c *EventConsumer) calculateDoubleMatchElo(participants []match.MatchParticipant) (*int, error) {
	// K-Factor: Bir maçta kazanılabilecek maksimum puan.
	// Rekabetin hızını artırmak istersen 40, daha stabil olsun dersen 24 yapabilirsin.
	const kFactor = 32.0

	var winnerSum, loserSum int
	var winnerCount, loserCount int

	// 1. Kazanan ve Kaybeden takımların toplam puanlarını ayır
	for _, p := range participants {
		if p.IsWinner {
			winnerSum += p.DoublePoint
			winnerCount++
		} else {
			loserSum += p.DoublePoint
			loserCount++
		}
	}

	// Güvenlik: Katılımcı listesi boşsa veya tek taraflıysa 0 dön
	if winnerCount != 2 || loserCount != 2 {
		return nil, fmt.Errorf("Kazanan ve kaybeden sayıları iki olmalı")
	}

	// Ortalamaları alıyoruz (Double maç olduğu için toplam/2)
	winnerAvg := float64(winnerSum) / 2.0
	loserAvg := float64(loserSum) / 2.0

	// KRİTİK NOKTA: (Kaybeden - Kazanan) / 400
	// Eğer Novak(3000) vs Sen(1000) ve Novak kazandıysa:
	// (1000 - 3000) / 400 = -5.0
	// 1 / (1 + 10^-5) = 0.999 (Beklenen durum)
	exponent := (loserAvg - winnerAvg) / 400.0
	expectedWin := 1.0 / (1.0 + math.Pow(10, exponent))

	// 4. Nihai Elo Değişimi
	// Beklenen durumda (1 - 0.999) * 32 = 0.03 puan (Çok az)
	// Sürpriz durumda (1 - 0.001) * 32 = 31.9 puan (Çok fazla)
	eloChange := kFactor * (1.0 - expectedWin)

	intEloChange := int(math.Round(eloChange))

	return &intEloChange, nil
}
