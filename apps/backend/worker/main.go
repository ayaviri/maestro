package main

import (
	"log"
	"maestro/internal"
	xamqp "maestro/internal/amqp"

	amqp "github.com/rabbitmq/amqp091-go"
)

var messageQueueConnection *amqp.Connection
var checkoutMessageQueue amqp.Queue
var err error

func main() {
	internal.WithTimer(
		"initialising connection to rabbitmq server, declaring checkout queue",
		func() {
			xamqp.SetupCheckoutQueue(&messageQueueConnection, &checkoutMessageQueue)
		},
	)

	defer messageQueueConnection.Close()
	var channel *amqp.Channel

	internal.WithTimer(
		"constructing channel of communciation with rabbitmq server",
		func() {
			channel, err = messageQueueConnection.Channel()
			internal.HandleError(
				err,
				"Could not construct channel of communication with Rabbit",
			)
		},
	)

	defer channel.Close()
	var messages <-chan amqp.Delivery

	internal.WithTimer(
		"opening channel for asynchronous message stream from queue",
		func() {
			messages, err = channel.Consume(
				checkoutMessageQueue.Name,
				"",
				true, // TODO: Change, don't want auto-ack
				false,
				false,
				false,
				nil,
			)
			internal.HandleError(err, "Could not read messages from the queue")
		},
	)

	var forever chan struct{}

	go func() {
		var d amqp.Delivery

		for d = range messages {
			log.Printf("Received a message: %s\n", d.Body)
		}
	}()

	log.Printf("Waiting to receive messages...")
	<-forever
}
