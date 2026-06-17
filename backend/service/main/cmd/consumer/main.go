package main

import (
	"context"
	"log"
	"tennis-league/service/internal/delivery/message/consumer/match_score/leaguematch"
	postgres2 "tennis-league/user-service/internal/repository/postgres"

	"os"
	"os/signal"
	"syscall"

	"tennis-league/common/consumer"
	"tennis-league/common/lib/database"

	"tennis-league/service/internal/delivery/message/consumer/match_score/matchScoreApproved"

	"tennis-league/service/internal/infrastructure/persistence/postgres"

	"tennis-league/common/lib/messaging"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go handleShutdown(cancel)

	rabbit, err := messaging.NewRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Close()

	db, err := database.NewPostgres()
	transactionManager := database.NewTransactionManager(db)
	matchRepository := postgres.NewMatchRepository(db)
	playerEarcnedPointRepository := postgres.NewPlayerEarnedPointRepository(db)
	setRepository := postgres.NewMatchSetRepository(db)
	scoreboardRepository := postgres.NewScoreBoardRepository(db)
	playerRepository := postgres2.NewPlayerRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	leagueMatchApprovedConsumer := leaguematch.NewLeagueMatchApprovedEventConsumer(transactionManager, matchRepository, setRepository, scoreboardRepository)
	consumer.RegisterConsumer(rabbit, ctx, leagueMatchApprovedConsumer.Consumer)

	updatePlayerPointsConsumer := matchScoreApproved.NewEventConsumer(transactionManager, matchRepository, playerRepository, playerEarcnedPointRepository)
	consumer.RegisterConsumer(rabbit, ctx, updatePlayerPointsConsumer.Consumer)

	log.Println("📥 Consumer running...")

	<-ctx.Done()
}
func handleShutdown(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	cancel()
}
