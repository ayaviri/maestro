package amqp

import (
	"maestro/internal"

	amqp "github.com/rabbitmq/amqp091-go"
)

var err error

// Establishes a connection to the rabbitmq server and declares the
// queue for cart checkouts
func SetupCheckoutQueue(connectionPtr **amqp.Connection, queuePtr *amqp.Queue) {
	// Look, I get this isn't testable. But I wasn't going to test it in the
	// first place. It's a block of initialisation. If anything fails, the server
	// or worker shouldn't start in the first place

	*connectionPtr, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	internal.HandleError(err, "Could not connect to the Rabbit message broker")
	var channel *amqp.Channel
	channel, err = (*connectionPtr).Channel()
	internal.HandleError(
		err,
		"Could not construct channel of communication with Rabbit",
	)
	*queuePtr, err = channel.QueueDeclare(
		"checkout",
		true,
		false,
		false,
		true,
		nil,
	)
	internal.HandleError(err, "Could not declare checkout queue")
}
