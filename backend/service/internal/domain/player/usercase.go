package player

import (
	"context"
	"log"

	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type Usecase struct {
	tm              *database.TransactionManager
	repo            Repository
	matchRepository match.Repository
}

func (u *Usecase) GetImconimgMatches(ctx context.Context, dto PlayerIncomingMatchesRequest) ([]IncomingMatches, error) {

	matches, err := u.matchRepository.GetPlayerIncomingMatches(ctx, match.PlayerIncomingMatchesQueryParam{PlayerId: dto.PlayerId, Limit: dto.Limit})

	if err != nil {
		log.Printf("Error while fetching incoming matches for player %s: %v", dto.PlayerId, err)
		return nil, err
	}

	var incomingMatches []IncomingMatches
	for _, m := range matches {
		incomingMatches = append(incomingMatches,
			IncomingMatches{MatchId: m.MatchId,
				MatchDate:    m.MatchDate,
				MatchType:    m.MatchType,
				Source:       m.Source,
				LeagueId:     m.LeagueId,
				LeagueName:   m.LeagueName,
				OppenentId:   m.OppenentId,
				OppenentName: m.OppenentName,
			})
	}
	return incomingMatches, nil

}

func (u *Usecase) GetPlayerStatistics(context context.Context, request PlayerStatisticsRequest) (*PlayerStatistics, error) {

	return u.repo.GetPlayerStatistics(context, request)

}

func (u *Usecase) AssignToUser(ctx context.Context, playerId string, userId string) error {
	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.repo.AssignToUser(txCtx, playerId, userId)
	})

}

func NewUsecase(tm *database.TransactionManager, r Repository, matchRepository match.Repository) *Usecase {
	return &Usecase{tm: tm, repo: r, matchRepository: matchRepository}
}

func (u *Usecase) GetById(ctx context.Context, id int64) (*Player, error) {
	return u.repo.GetById(ctx, id)
}

func (u *Usecase) Save(ctx context.Context, persistPlayer *PersistPlayer) (*string, error) {
	var userId *string

	err := u.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		newUserId, err := u.repo.Save(txCtx, persistPlayer)
		if err == nil {
			userId = newUserId
			return nil
		} else {
			return nil
		}
	})

	return userId, err
}

func (u *Usecase) List(ctx context.Context, queryParams ListQueryParameters) ([]*Player, error) {
	return u.repo.List(ctx, queryParams)
}
