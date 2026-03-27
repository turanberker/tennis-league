package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/rabbitmq/amqp091-go"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type MatchScoreApprovedEventConsumer struct {
	*Consumer
	tm         *database.TransactionManager
	matchRepo  match.Repository
	playerRepo player.Repository
}

func NewMatchScoreApprovedEventConsumer(tm *database.TransactionManager,
	matchRepo match.Repository,
	playerRepo player.Repository) *MatchScoreApprovedEventConsumer {

	c := &MatchScoreApprovedEventConsumer{
		tm:         tm,
		matchRepo:  matchRepo,
		playerRepo: playerRepo,
	}

	c.Consumer = &Consumer{
		Queue:       "update_player_points_queue",
		RoutingName: "MatchApproved",
		Handler:     c.handle, // 👈 struct method
	}

	return c

}

func (c *MatchScoreApprovedEventConsumer) handle(msg amqp091.Delivery) error {
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
			participants, err := c.matchRepo.GetDoubleMatchParticipantsWithPoints(txCtx, event.MatchID)
			if err != nil {
				//Log yaz
				return err
			}

			change, err := calculateDoubleMatchElo(participants)
			if err != nil {
				//Log yaz
				return err
			}
			for _, p := range participants {
				if p.IsWinner {
					_, err = c.playerRepo.IncreaseDoublePoint(txCtx, p.PlayerID, *change)
				} else {
					_, err = c.playerRepo.DecreaseDoublePoint(txCtx, p.PlayerID, *change)
				}

				if err != nil {
					log.Printf("Puan güncelleme hatası (Player: %s): %v", p.PlayerID, err)
					return err
				}

			}

			// 4. History tablosuna kayıt at
			// history err := c.playerRepo.SavePointsHistory(txCtx, matchID, historyRecords)
		case match.MatchType_SINGLE:
			log.Panic("Not Implemented")
		}

		log.Println("Match Approved:", err)
		return nil
	})

}

func calculateDoubleMatchElo(participants []match.MatchParticipant) (*int, error) {
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
