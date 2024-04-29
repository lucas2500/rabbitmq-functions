package messages

import (
	"context"
	"log"
	"time"

	conn "github.com/lucas2500/rabbitmq-functions/connections"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueProperties struct {
	ExchangeName string
	RoutingKey   string
	QueueName    string
}

type Message struct {
	QueueProperties
	Body []byte
}

func (q QueueProperties) DequeueMessage(MsgBody chan amqp.Delivery, WorkerName string) {

	forever := make(chan bool)

	conn := conn.RabbitConn

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Unable to open channel: %s\n", err)
	}

	defer ch.Close()

	err = ch.Qos(
		1,
		0,
		true,
	)

	if err != nil {
		log.Fatalf("Unable to define Qos: %s\n", err)
	}

	msgs, err := ch.Consume(
		q.QueueName,
		WorkerName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Unable to register consumer: %s\n", err)
	}

	go func() {
		for d := range msgs {
			MsgBody <- d
		}
	}()

	log.Printf("- %s: up and running...", WorkerName)

	<-forever
}

func (m Message) QueueMessage() error {

	conn := conn.RabbitConn

	ch, err := conn.Channel()

	if err != nil {
		log.Printf("Unable to open channel: %s\n", err)
	}

	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		m.ExchangeName,
		m.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        m.Body,
		},
	)

	if err != nil {
		log.Printf("Unable to publish message: %s\n", err)
	}

	return err
}
