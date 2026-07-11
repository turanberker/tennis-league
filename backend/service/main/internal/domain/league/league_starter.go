package league

import (
	"context"
	"fmt"
	"math/rand"

	customerror "tennis-league/common/lib/error"
	"tennis-league/service/internal/domain/match"
	"tennis-league/service/internal/domain/team"
)

type LeagueStarter interface {
	Start(ctx context.Context, leagueId string) error
}

func (u *Usecase) leagueStarterGetter(league League, teamRepo team.Repository,
	matchRepo match.Repository) (LeagueStarter, error) {
	switch league.ProcessType {
	case LeagueProcessType_FIXTURE:
		return &fixtureTypeLeagueStarter{teamRepo: teamRepo, matchRepo: matchRepo}, nil

	case LeagueProcessType_DEFI:
		return &challengingTypeLeagueStarter{}, nil
	default:
		err := fmt.Errorf("%s tipin implement edilmemiştir", league.ProcessType)
		return nil, customerror.NewInternalError(err)
	}
}

type challengingTypeLeagueStarter struct {
}

func (s *challengingTypeLeagueStarter) Start(ctx context.Context, leagueId string) error {
	return nil
}

type fixtureTypeLeagueStarter struct {
	teamRepo  team.Repository
	matchRepo match.Repository
}

func (u *fixtureTypeLeagueStarter) Start(txCtx context.Context, leagueId string) error {
	//TODO Burada ligin Single- double olmasına göre işlem yapılacak
	teams, err := u.teamRepo.GetByLeagueId(txCtx, leagueId)
	if err != nil {
	}
	var bulkInsert match.BulkInsertMatches
	var matches []match.SideIds
	var teamIds []string

	for i := 0; i < len(teams); i++ {

		teamIds = append(teamIds, teams[i].ID)
		for j := i + 1; j < len(teams); j++ { // j=i+1 → tekrar ve kendisiyle maç yok
			team1Id := teams[i].ID
			team2Id := teams[j].ID

			// 50% ihtimalle takımların yerini değiştir
			if rand.Intn(2) == 0 {
				team1Id, team2Id = team2Id, team1Id
			}

			match := match.SideIds{
				Side1: team1Id,
				Side2: team2Id,
			}

			matches = append(matches, match)
		}
	}
	//Maçların sırasını karıştır (Opsiyonel ama daha profesyonel bir fikstür sağlar)
	rand.Shuffle(len(matches), func(i, j int) {
		matches[i], matches[j] = matches[j], matches[i]
	})

	bulkInsert.Sides = matches
	bulkInsert.Type = match.MatchType{Id: &leagueId,
		Source: match.MatchSource_LEAGUE,
		Type:   match.MatchType_DOUBLE,
	}

	return u.matchRepo.SaveBulkMatches(txCtx, &bulkInsert)

}
