package amqp

import (
	"maestro/internal"

	amqp "github.com/rabbitmq/amqp091-go"
)

var err error

// Establishes a connection to the rabbitmq server and declares the
// queues for cart checkout request and completion
func SetupQueues(
	connectionPtr **amqp.Connection,
	checkoutRequestQueuePtr *amqp.Queue,
	checkoutCompletionQueuePtr *amqp.Queue,
) {
	// I didn't want to encapsulate these two function arguments, so I decided
	// to pass mutable pointers instead. Womp womp

	*connectionPtr, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	internal.HandleError(err, "Could not connect to the Rabbit message broker")
	var channel *amqp.Channel
	channel, err = (*connectionPtr).Channel()
	internal.HandleError(
		err,
		"Could not construct channel of communication with Rabbit",
	)
	*checkoutRequestQueuePtr, err = channel.QueueDeclare(
		"checkout_request",
		true,
		false,
		false,
		true,
		nil,
	)
	*checkoutCompletionQueuePtr, err = channel.QueueDeclare(
		"checkout_completion",
		true,
		false,
		false,
		true,
		nil,
	)
	internal.HandleError(err, "Could not declare checkout queue")
}
