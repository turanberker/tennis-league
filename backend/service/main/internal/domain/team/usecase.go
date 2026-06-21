package team

import (
	"context"

	"time"

	"tennis-league/common/lib/cache"
	"tennis-league/common/lib/database"
	"tennis-league/service/internal/domain/teamplayer"
)

type CreateTeamRequest struct {
	LeagueID  string
	Name      string
	PlayerIDs []string
}

type UseCase struct {
	tm                   *database.TransactionManager
	cacheManager         *cache.CacheManager
	repository           Repository
	teamPlayerRepository teamplayer.Repository
}

func NewUseCase(tm *database.TransactionManager, cacheManager *cache.CacheManager, repository Repository, teamPlayerRepository teamplayer.Repository) *UseCase {
	return &UseCase{tm: tm, cacheManager: cacheManager, repository: repository, teamPlayerRepository: teamPlayerRepository}
}

func (u *UseCase) GetById(ctx context.Context, id string) (*Team, error) {
	return u.repository.GetById(ctx, id)
}

func (u *UseCase) GetByLeagueId(ctx context.Context, leagueId string) ([]*LeagueTeam, error) {
	return u.repository.GetByLeagueId(ctx, leagueId)
}

func (u *UseCase) Save(ctx context.Context, req *CreateTeamRequest) (*string, error) {
	var teamId *string
	err := u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		teamID, err := u.repository.Save(txCtx, &PersistTeam{
			LeagueID: req.LeagueID,
			Name:     req.Name,
		})

		if err != nil {
			return err
		}

		for _, pid := range req.PlayerIDs {
			err = u.teamPlayerRepository.Save(txCtx, &teamplayer.PersistTeamPlayer{
				TeamID:   *teamID,
				PlayerID: pid,
			})
			if err != nil {
				return err
			}
		}
		teamId = teamID

		return nil
	})

	if err != nil {
		return nil, err
	}
	return teamId, nil
}

func (u *UseCase) GetTeamMembers(ctx context.Context, teamId string) ([]teamplayer.Player, error) {
	cacheKey := u.cacheManager.PrepareCacheKey("getTeamMembers", teamId)
	return cache.Cacheable(u.cacheManager, ctx, cacheKey, 1*time.Hour, func() ([]teamplayer.Player, error) {

		returnVal, err := u.teamPlayerRepository.GetByPlayersByTeamId(ctx, teamId)

		if err != nil {
			return nil, err
		}
		return returnVal, nil
	})
}
