package consumer

import (
	"context"
	"log"

	"tennis-league/common/lib/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Queue       string
	RoutingName string

	Handler func(amqp091.Delivery) error
}

func RegisterConsumer(rabbit *messaging.RabbitMQ, ctx context.Context, consumer *Consumer) {

	err := rabbit.DeclareQueue(consumer.Queue, consumer.RoutingName)
	if err != nil {
		log.Fatal(err)
	}

	err = rabbit.Consume(ctx, consumer.Queue, consumer.Handler)
	if err != nil {
		log.Fatal(err)
	}
}
