package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/turanberker/tennis-league-service/internal/domain/outbox"
	"github.com/turanberker/tennis-league-service/internal/infrastructure/persistence/postgres"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
	"github.com/turanberker/tennis-league-service/internal/platform/messaging"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- Graceful shutdown ---
	go handleShutdown(cancel)
	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	repository := postgres.NewOutboxRepository(db)

	// --- Rabbit ---
	rabbit, err := messaging.NewRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Close()

	log.Println("🚀 Outbox publisher started")

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Shutting down outbox publisher...")
			return

		case <-ticker.C:
			if err := processOutbox(ctx, db, repository, rabbit); err != nil {
				log.Println("outbox error:", err)
			}
		}
	}
}

func processOutbox(ctx context.Context, db *sql.DB, repository outbox.Repository, rabbit *messaging.RabbitMQ) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	events, err := repository.GetEventsToPublish(ctx, tx)

	if err != nil {
		return err
	}

	for _, e := range events {

		err := rabbit.Publish(ctx, e.EventType, e.Payload)
		if err != nil {
			log.Println("publish failed:", err)
			repository.IncreaseRetryCount(ctx, tx, e.Id)
			continue
		}
		err = repository.UpdateToPublished(ctx, tx, e.Id)

		if err != nil {
			return err
		}

	}

	return tx.Commit()
}

func handleShutdown(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	cancel()
}
