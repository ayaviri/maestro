package worker

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RejectAndExitOnError(err error, message string, d amqp.Delivery) {
	if err != nil {
		d.Reject(true) // Requeues delivery
		log.Fatalf(message+": %v", err.Error())
	}
}
