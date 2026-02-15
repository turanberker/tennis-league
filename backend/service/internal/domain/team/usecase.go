package team

import (
	"context"
	"database/sql"

	"github.com/turanberker/tennis-league-service/internal/domain/teamplayer"
)

type CreateTeamRequest struct {
	LeagueID  string
	Name      string
	PlayerIDs []string
}

type UseCase struct {
	db                   *sql.DB
	repository           Repository
	teamPlayerRepository teamplayer.Repository
}

func NewUseCase(db *sql.DB, repository Repository, teamPlayerRepository teamplayer.Repository) *UseCase {
	return &UseCase{db: db, repository: repository, teamPlayerRepository: teamPlayerRepository}
}

func (u *UseCase) GetById(ctx context.Context, id string) (*Team, error) {
	return u.repository.GetById(ctx, id)
}

func (u *UseCase) GetByLeagueId(ctx context.Context, leagueId string) ([]*Team, error) {
	return u.repository.GetByLeagueId(ctx, leagueId)
}

func (u *UseCase) Save(ctx context.Context, req *CreateTeamRequest) (*string, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	teamID, err := u.repository.Save(ctx, tx, &PersistTeam{
		LeagueID: req.LeagueID,
		Name:     req.Name,
	})

	if err != nil {
		return nil, err
	}

	for _, pid := range req.PlayerIDs {
		err = u.teamPlayerRepository.Save(ctx, tx, &teamplayer.PersistTeamPlayer{
			TeamID:   *teamID,
			PlayerID: pid,
		})
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return teamID, nil
}
