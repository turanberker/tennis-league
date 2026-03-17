package main

import (
	"context"
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
	transactionManager := database.NewTransactionManager(db)
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

			// 1. Adayları belirle (Kilitleme yok)
			ids, err := repository.GetPendingIDs(ctx, 10)
			if err != nil {
				log.Println("Adaylar çekilemedi:", err)
				continue
			}

			for _, id := range ids {
				// 2. Her aday için kendi küçük dünyasını (Transaction) kur
				_ = transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {

					// 3. Gerçekten kilitle (Eğer başka publisher kapmadıysa)
					event, err := repository.GetEventForUpdate(txCtx, id)
					if err != nil {
						// Eğer SKIP LOCKED yüzünden boş dönerse veya hata varsa sessizce geç
						return nil
					}

					// 4. Publish et ve statüyü PUBLISHED yap
					return processOutbox(txCtx, event, repository, rabbit)
				})
			}

		}
	}
}

func processOutbox(ctx context.Context, event *outbox.EventToPublish, repository outbox.Repository, rabbit *messaging.RabbitMQ) error {

	err := rabbit.Publish(ctx, event.EventType, event.Payload)
	if err != nil {
		log.Println("publish failed:", err)
		repository.IncreaseRetryCount(ctx, event.Id)
		return nil
	}
	err = repository.UpdateToPublished(ctx, event.Id)

	return nil
}

func handleShutdown(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	cancel()
}
