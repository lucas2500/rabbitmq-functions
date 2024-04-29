package connections

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitConn *amqp.Connection
)

func ConnectToRabbitMQ() {

	var (
		dsn   string
		retry int = 0
		err   error
	)

	log.Printf("- Connecting to RabbitMQ...")

	dsn = os.Getenv("RABBITMQ_STRING_CONNECTION")

	for {

		retry++

		log.Printf("- %d connection attempt\n", retry)

		RabbitConn, err = amqp.Dial(dsn)

		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
	}
}
