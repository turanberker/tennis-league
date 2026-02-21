package match

import (
	"context"
	"database/sql"
	"errors"
	"time"

	matchSet "github.com/turanberker/tennis-league-service/internal/domain/matchset"
)

type UseCase struct {
	db              *sql.DB
	repository      Repository
	scoreRepository matchSet.Repository
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

func NewUseCase(db *sql.DB, r Repository, scoreRepository matchSet.Repository) *UseCase {
	return &UseCase{db: db, repository: r, scoreRepository: scoreRepository}
}

func (u *UseCase) UpdateMatchDate(ctx context.Context, matchId string, matchDate *time.Time) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = u.repository.UpdateMatchDate(ctx, tx, UpdateMatchDate{Id: matchId, MatchDate: matchDate})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *UseCase) SaveMatchScore(ctx context.Context, score *SaveMatchScore) (*UpdateMatchScore, error) {

	macScore, err := calculateMatchScore(score)

	if err != nil {
		return nil, err
	}

	teamIds := u.repository.GetMatchTeamIds(ctx, score.MatchId)
	if teamIds == nil {
		return nil, errors.New("Maç bulunamadı")
	}

	if teamIds.Status == StatusApproved {
		return nil, errors.New("Maç skoru onaylandığı için güncelleyemezsiniz")
	}

	if macScore.Team1Score > macScore.Team2Score {
		macScore.WinnerTeamId = teamIds.Team1Id
	} else {
		macScore.WinnerTeamId = teamIds.Team2Id
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	u.scoreRepository.DeleteSetScores(ctx, tx, score.MatchId)
	u.scoreRepository.SaveSetScore(ctx, tx, &matchSet.UpdateSetScore{MatchId: score.MatchId,
		Set:        1,
		Team1Score: score.Set1.Team1Score,
		Team2Score: score.Set1.Team2Score})
	u.scoreRepository.SaveSetScore(ctx, tx, &matchSet.UpdateSetScore{MatchId: score.MatchId,
		Set:        2,
		Team1Score: score.Set2.Team1Score,
		Team2Score: score.Set2.Team2Score})
	if score.SuperTie != nil {
		u.scoreRepository.SaveSuperTieScore(ctx, tx, &matchSet.UpdateSuperTieScore{MatchId: score.MatchId,
			Team1Score: score.SuperTie.Team1Score,
			Team2Score: score.SuperTie.Team2Score})
	}

	u.repository.UpdateMatchScore(ctx, tx, macScore)
	tx.Commit()
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

type winner int
type totalScore struct {
	team1 int
	team2 int
}

const (
	team1 winner = iota
	team2
)
