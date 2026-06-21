package main

import (
	"context"
	"log"
	"tennis-league/user-service/internal/consumer/playerpoint"
	"tennis-league/user-service/internal/repository/postgres"

	"os"
	"os/signal"
	"syscall"

	"tennis-league/common/consumer"
	"tennis-league/common/lib/database"

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

	playerEarcnedPointRepository := postgres.NewPlayerEarnedPointRepository(db)

	playerRepository := postgres.NewPlayerRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	updatePlayerPointsConsumer := playerpoint.NewEventConsumer(transactionManager, playerRepository, playerEarcnedPointRepository)
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
