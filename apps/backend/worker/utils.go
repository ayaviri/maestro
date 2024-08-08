package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RejectOnError(err error, message string, d amqp.Delivery) {
	if err != nil {
		d.Reject(true) // Requeues delivery
		log.Fatalf(message+": %v", err.Error())
	}
}
