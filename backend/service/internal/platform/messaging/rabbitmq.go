package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/turanberker/tennis-league-service/internal/platform"
)

type RabbitMQ struct {
    dsn        string
    Connection *amqp091.Connection
    Channel    *amqp091.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {

	config := platform.LoadRabbitConfig()

	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.VHost,
	)

	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	log.Println("RabbitMQ connected")
	
	closeChan := make(chan *amqp091.Error)
	ch.NotifyClose(closeChan)

	go func() {
		err := <-closeChan
		log.Println("RabbitMQ channel closed:", err)
	}()

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
}

func (r *RabbitMQ) Publish(
    ctx context.Context,
    routingKey string,
    body []byte,
) error {

    if r.Channel == nil || r.Channel.IsClosed() {
        if err := r.reconnect(); err != nil {
            return err
        }
    }

    return r.Channel.PublishWithContext(
        ctx,
        "events_exchange",
        routingKey,
        false,
        false,
        amqp091.Publishing{
            ContentType:  "application/json",
            DeliveryMode: amqp091.Persistent,
            Body:         body,
        },
    )
}

func (r *RabbitMQ) DeclareQueue(name string, routingKey string) error {

	// 1️⃣ Exchange declare (EKLE)
	err := r.Channel.ExchangeDeclare(
		"events_exchange", // name
		"direct",          // type
		true,              // durable
		false,             // auto-delete
		false,             // internal
		false,             // no-wait
		nil,
	)
	if err != nil {
		return err
	}
	// 2️⃣ Queue declare
	_, err = r.Channel.QueueDeclare(
		name,
		true,  // durable
		false, // auto delete
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	// 3️⃣ Bind
	return r.Channel.QueueBind(
		name,
		routingKey,
		"events_exchange",
		false,
		nil,
	)
}
func (r *RabbitMQ) Consume(
	ctx context.Context,
	queueName string,
	handler func(amqp091.Delivery) error,
) error {

	err := r.Channel.Qos(10, 0, false)
	if err != nil {
		return err
	}

	msgs, err := r.Channel.Consume(
		queueName,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Consumer stopped")
				return

			case msg := <-msgs:

				err := handler(msg)
				if err != nil {
					log.Println("handler error:", err)
					msg.Nack(false, true)
					continue
				}

				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) reconnect() error {
    conn, err := amqp091.Dial(r.dsn)
    if err != nil {
        return err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return err
    }

    r.Connection = conn
    r.Channel = ch
    return nil
}
