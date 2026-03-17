package matchapproved

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/turanberker/tennis-league-service/internal/delivery/message"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
	matchSet "github.com/turanberker/tennis-league-service/internal/domain/matchset"
	"github.com/turanberker/tennis-league-service/internal/domain/scoreboard"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type MatchApprovedEventConsumer struct {
	*message.Consumer
	tm             *database.TransactionManager
	matchRepo      match.Repository
	setRepo        matchSet.Repository
	scoreboradRepo scoreboard.Repository
}

func NewMatchApprovedEventConsumer(tm *database.TransactionManager,
	matchRepo match.Repository,
	setRepo matchSet.Repository,
	scoreboradRepo scoreboard.Repository,
) *MatchApprovedEventConsumer {

	c := &MatchApprovedEventConsumer{
		tm:             tm,
		matchRepo:      matchRepo,
		setRepo:        setRepo,
		scoreboradRepo: scoreboradRepo,
	}

	c.Consumer = &message.Consumer{
		Queue:       "match_events_queue",
		RoutingName: "MatchApproved",
		Handler:     c.handle, // 👈 struct method
	}

	return c

}

func (c *MatchApprovedEventConsumer) handle(msg amqp091.Delivery) error {
	ctx := context.Background()

	var event = &match.MatchApprovedEvent{}

	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		return err
	}

	return c.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		matchTeams := c.matchRepo.GetMatchTeamIds(txCtx, event.MatchID)

		if matchTeams == nil {
			return fmt.Errorf("%s, Maç Bulunamadı", event.MatchID)
		}

		setScores := c.setRepo.GetSetScoreList(txCtx, event.MatchID)

		var team1Update = &scoreboard.IncreaseTeamScore{
			LeagueId:      matchTeams.LeagueId,
			TeamId:        matchTeams.Team1Id,
			Won:           false,
			WonSets:       0,
			LostSets:      0,
			WonGames:      0,
			LostGames:     0,
			IncreaseScore: 0,
		}

		var team2Update = &scoreboard.IncreaseTeamScore{
			LeagueId:      matchTeams.LeagueId,
			TeamId:        matchTeams.Team2Id,
			Won:           false,
			WonSets:       0,
			LostSets:      0,
			WonGames:      0,
			LostGames:     0,
			IncreaseScore: 0,
		}

		for _, set := range setScores {

			if set.SetNumber == 3 {
				if *set.Team1TiePoint > *set.Team2TiePoint {
					team1Update.WonSets += 1
					team2Update.LostSets += 1
				} else {
					team2Update.WonSets += 1
					team1Update.LostSets += 1
				}
			} else {
				if *set.Team1Game > *set.Team2Game {
					team1Update.WonSets += 1
					team2Update.LostSets += 1

				} else {
					team2Update.WonSets += 1
					team1Update.LostSets += 1
				}

				team1Update.WonGames += int16(*set.Team1Game)
				team1Update.LostGames += int16(*set.Team2Game)

				team2Update.WonGames += int16(*set.Team2Game)
				team2Update.LostGames += int16(*set.Team1Game)
			}
		}

		if team1Update.WonSets > team2Update.WonSets {
			team1Update.IncreaseScore = 20
			team1Update.Won = true
		} else {
			team2Update.IncreaseScore = 20
			team2Update.Won = true
		}

		c.scoreboradRepo.UpdateScore(txCtx, *team1Update)
		c.scoreboradRepo.UpdateScore(txCtx, *team2Update)
		log.Println("Match Approved:", event.MatchID)

		return nil
	})

}
