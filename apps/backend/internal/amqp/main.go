package amqp

import (
	"fmt"
	"maestro/internal"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var err error

func getBrokerServerUrl() string {
	brokerUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		os.Getenv("RABBIT_USER"),
		os.Getenv("RABBIT_PASSWORD"),
		os.Getenv("RABBIT_HOST"),
		os.Getenv("RABBIT_PORT"),
	)

	return brokerUrl
}

// Establishes a connection to the rabbitmq server and declares the
// queues for cart checkout request and completion
func SetupQueues(
	connectionPtr **amqp.Connection,
	checkoutRequestQueuePtr *amqp.Queue,
	checkoutCompletionQueuePtr *amqp.Queue,
) {
	// I didn't want to encapsulate these two function arguments, so I decided
	// to pass mutable pointers instead. Womp womp

	var url string = getBrokerServerUrl()
	*connectionPtr, err = amqp.Dial(url)
	internal.HandleError(
		err,
		fmt.Sprintf(
			"Could not connect to the Rabbit message broker, given URL: %s",
			url,
		),
	)
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
