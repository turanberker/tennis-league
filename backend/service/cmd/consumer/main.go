package main

import (
	"context"
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/turanberker/tennis-league-service/internal/delivery/message/consumer"

	"github.com/turanberker/tennis-league-service/internal/infrastructure/persistence/postgres"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
	"github.com/turanberker/tennis-league-service/internal/platform/messaging"
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
	setRepository := postgres.NewMatchSetRepository(db)
	scoreboardRepository := postgres.NewScoreBoardRepository(db)
	playerRepository := postgres.NewPlayerRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	leagueMatchApprovedConsumer := consumer.NewLeagueMatchApprovedEventConsumer(transactionManager, matchRepository, setRepository, scoreboardRepository)
	consumer.RegisterConsumer(rabbit, ctx, leagueMatchApprovedConsumer.Consumer)

	updatePlayerPointsConsumer := consumer.NewMatchScoreApprovedEventConsumer(transactionManager, matchRepository, playerRepository)
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
