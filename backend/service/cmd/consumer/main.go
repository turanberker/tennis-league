package main

import (
	"context"
	"log"

	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbitmq/amqp091-go"
	"github.com/turanberker/tennis-league-service/internal/delivery/message"
	matchapproved "github.com/turanberker/tennis-league-service/internal/delivery/message/match_approved"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
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
	if err != nil {
		log.Fatal(err)
	}

	matchApprovedConsumer := matchapproved.NewMatchApprovedEventConsumer(transactionManager, matchRepository, setRepository, scoreboardRepository)
	message.RegisterConsumer(rabbit, ctx, matchApprovedConsumer.Consumer)

	log.Println("📥 Consumer running...")

	<-ctx.Done()
}
func handleShutdown(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	cancel()
}
func handleMatchApproved(msg amqp091.Delivery) error {

	var event = &match.MatchApprovedEvent{}

	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		return err
	}

	log.Println("Match Approved:", event.MatchID)

	// burada DB update vs yapılabilir

	return nil
}
