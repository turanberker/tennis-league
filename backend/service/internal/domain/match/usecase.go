package match

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	matchSet "github.com/turanberker/tennis-league-service/internal/domain/matchset"
	"github.com/turanberker/tennis-league-service/internal/domain/outbox"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type UseCase struct {
	tm               *database.TransactionManager
	repository       Repository
	scoreRepository  matchSet.Repository
	outboxRepository outbox.Repository
}
type SaveScore struct {
	Team1Score int8
	Team2Score int8
}

type SaveMatchScore struct {
	MatchId  string
	Set1     SaveScore
	Set2     SaveScore
	SuperTie *SaveScore
}

func NewUseCase(tm *database.TransactionManager,
	r Repository,
	scoreRepository matchSet.Repository,
	outboxRepository outbox.Repository,
) *UseCase {
	return &UseCase{tm: tm,
		repository:      r,
		scoreRepository: scoreRepository, outboxRepository: outboxRepository}
}
func (u *UseCase) ApproveScore(ctx context.Context, matchId string) error {

	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		err := u.repository.ApproveScore(txCtx, matchId)

		if err != nil {
			return err
		}

		event := MatchApprovedEvent{
			MatchID: matchId,
		}
		payload, _ := json.Marshal(event)
		outboxEntity := &outbox.PersistEntity{
			AggregateType: "match",
			AggregateID:   matchId,
			EventType:     "MatchApproved",
			Payload:       payload,
		}
		err = u.outboxRepository.Save(txCtx, outboxEntity)
		if err != nil {
			return err
		}

		return nil
	})

}

func (u *UseCase) UpdateMatchDate(ctx context.Context, matchId string, matchDate *time.Time) error {

	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		err := u.repository.UpdateMatchDate(txCtx, UpdateMatchDate{Id: matchId, MatchDate: matchDate})
		if err != nil {
			return err
		}
		return nil
	})
}

func (u *UseCase) SaveMatchScore(ctx context.Context, score *SaveMatchScore) (*UpdateMatchScore, error) {

	macScore, err := calculateMatchScore(score)
	if err != nil {
		return nil, err
	}
	err = u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		if err != nil {
			return err
		}

		teamIds := u.repository.GetMatchTeamIds(txCtx, score.MatchId)
		if teamIds == nil {
			return errors.New("Maç bulunamadı")
		}

		if teamIds.Status == StatusApproved {
			return errors.New("Maç skoru onaylandığı için güncelleyemezsiniz")
		}

		if macScore.Team1Score > macScore.Team2Score {
			macScore.WinnerTeamId = teamIds.Team1Id
		} else {
			macScore.WinnerTeamId = teamIds.Team2Id
		}

		u.scoreRepository.DeleteSetScores(txCtx, score.MatchId)
		u.scoreRepository.SaveSetScore(txCtx, &matchSet.UpdateSetScore{MatchId: score.MatchId,
			Set:        1,
			Team1Score: score.Set1.Team1Score,
			Team2Score: score.Set1.Team2Score})
		u.scoreRepository.SaveSetScore(txCtx, &matchSet.UpdateSetScore{MatchId: score.MatchId,
			Set:        2,
			Team1Score: score.Set2.Team1Score,
			Team2Score: score.Set2.Team2Score})
		if score.SuperTie != nil {
			u.scoreRepository.SaveSuperTieScore(txCtx, &matchSet.UpdateSuperTieScore{MatchId: score.MatchId,
				Team1Score: score.SuperTie.Team1Score,
				Team2Score: score.SuperTie.Team2Score})
		}

		u.repository.UpdateMatchScore(txCtx, macScore)

		return nil
	})

	if err != nil {
		return nil, err
	}
	return macScore, nil
}

func (u *UseCase) GetSetScore(ctx context.Context, matchId string) []*matchSet.MatchSetScores {
	return u.scoreRepository.GetSetScoreList(ctx, matchId)
}

func calculateMatchScore(score *SaveMatchScore) (*UpdateMatchScore, error) {

	matchscore := &UpdateMatchScore{Id: score.MatchId, Team1Score: 0, Team2Score: 0}
	findSetWinner(score.Set1, matchscore)

	findSetWinner(score.Set2, matchscore)

	if matchscore.Team1Score == matchscore.Team2Score {
		if score.SuperTie == nil {
			return nil, errors.New("Süper Tie skoru eksik")
		} else {
			findSetWinner(*score.SuperTie, matchscore)
		}
	}

	return matchscore, nil
}

func findSetWinner(score SaveScore, matchScore *UpdateMatchScore) {
	if score.Team1Score > score.Team2Score {
		matchScore.Team1Score += 1
	} else {
		matchScore.Team2Score += 1
	}
}
