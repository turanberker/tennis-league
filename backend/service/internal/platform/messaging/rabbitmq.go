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
	// 1️⃣ Exchange declare (EKLE)
	err = ch.ExchangeDeclare(
		"events_exchange", // name
		"direct",          // type
		true,              // durable
		false,             // auto-delete
		false,             // internal
		false,             // no-wait
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("exchange declare error: %w", err)
	}

	log.Println("RabbitMQ connected and exchange declared")

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
	exchangeName := "events_exchange"
	// Exchange'in varlığını garanti altına al
	err := r.Channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	// Mesajı gönder
	err = r.Channel.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         body,
		},
	)

	if err != nil {
		return err
	}
	// ✅ Loglama: Hangi exchange, hangi routingKey ve mesaj boyutu
	log.Printf("[RabbitMQ] Mesaj başarıyla gönderildi | Exchange: %s | RoutingKey: %s | Payload Size: %d bytes",
		exchangeName,
		routingKey,
		len(body),
	)

	return nil
}

func (r *RabbitMQ) DeclareQueue(name string, routingKey string) error {

	// 2️⃣ Queue declare
	_, err := r.Channel.QueueDeclare(
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

			case msg, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ msgs channel closed")
					return
				}

				log.Printf("[RabbitMQ] Mesaj alındı | Kuyruk: %s | RoutingKey: %s | MessageID: %s",
					queueName,
					msg.RoutingKey,
					msg.MessageId,
				)
				err := handler(msg)
				if err != nil {
					log.Printf("[RabbitMQ] Handler hatası (Kuyruk: %s): %v", queueName, err)
					// requeue: true yaparak mesajı hatada tekrar kuyruğa gönderiyoruz
					msg.Nack(false, true)
					continue
				}

				// Başarılı işlem logu (Opsiyonel)
				log.Printf("[RabbitMQ] Mesaj başarıyla işlendi (Ack) | Kuyruk: %s", queueName)
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
