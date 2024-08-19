package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	xhttp "maestro/internal/http"
	xworker "maestro/internal/worker"

	"github.com/ayaviri/goutils/timer"
	amqp "github.com/rabbitmq/amqp091-go"
)

//  ____   ____ _   _ _____ __  __    _    ____
// / ___| / ___| | | | ____|  \/  |  / \  / ___|
// \___ \| |   | |_| |  _| | |\/| | / _ \ \___ \
//  ___) | |___|  _  | |___| |  | |/ ___ \ ___) |
// |____/ \____|_| |_|_____|_|  |_/_/   \_\____/
//

func JobResourceHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestSegments []string = strings.Split(
		strings.TrimPrefix(request.URL.Path, "/"), "/",
	)

	if len(requestSegments) != 2 {
		http.Error(writer, "Expected path is /job/id", http.StatusBadRequest)
	}

	var jobId string = requestSegments[1]
	xhttp.SetSSEHeaders(writer)
	flusher, ok := writer.(http.Flusher)

	if !ok {
		http.Error(
			writer,
			"Response writer is not flushable",
			http.StatusInternalServerError,
		)
		return
	}

	var messageQueueChannel *amqp.Channel
	messageQueueChannel, err = messageQueueConnection.Channel()

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf(
				"Could not establish channel of communication with RabbitMQ server: %v\n",
				err.Error(),
			),
			http.StatusInternalServerError,
		)
		return
	}

	defer messageQueueChannel.Close()
	var messages <-chan amqp.Delivery

	timer.WithTimer(
		"opening channel for asynchronous message stream from queue",
		func() {
			messages, err = messageQueueChannel.Consume(
				checkoutCompletionQueue.Name,
				"",
				false,
				false,
				false,
				false,
				nil,
			)
		},
	)

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not read messages from the queue: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Used by the message stream goroutine to signal to the heartbeat goroutine
	// that it should terminate. Heartbeat goroutine will terminate when this
	// channel is closed
	heartbeatTerminationChannel := make(chan int)
	go xhttp.Heartbeat(1*time.Second, writer, flusher, heartbeatTerminationChannel)
	go awaitCheckoutCompletionMessage(
		messages,
		jobId,
		writer,
		flusher,
		heartbeatTerminationChannel,
		request.Context().Done(),
	)

	// NOTE: This makes the FRONTEND responsible for closing the connection
	// once it receives the required data
	<-request.Context().Done()
}

// Waits for checkout completion messages from the RabbitMQ server,
// closing the given termination channel if one of the following occur:
// 1) Correct message is received
// 2) Message stream is closed
func awaitCheckoutCompletionMessage(
	messages <-chan amqp.Delivery,
	jobId string,
	writer http.ResponseWriter,
	flusher http.Flusher,
	terminationChannel chan int,
	clientDisconnected <-chan struct{},
) {
	var delivery amqp.Delivery
	var ok bool

	for {
		select {
		case delivery, ok = <-messages:
			if !ok {
				http.Error(
					writer,
					"Channel of communication with RabbitMQ server severed prematurely",
					http.StatusInternalServerError,
				)
				close(terminationChannel)
				return
			}

			err = processCheckoutCompletionMessage(
				delivery,
				jobId,
				writer,
				flusher,
				terminationChannel,
			)

			if err != nil {
				delivery.Reject(true)
			}
		case <-clientDisconnected:
			return
		}
	}
}

// Processes the given message (delivery), ensure it matches the given job ID.
func processCheckoutCompletionMessage(
	delivery amqp.Delivery,
	jobId string,
	writer http.ResponseWriter,
	flusher http.Flusher,
	terminationChannel chan int,
) error {
	var completionMessage xworker.CheckoutCompletionMessage
	err = json.Unmarshal(delivery.Body, &completionMessage)

	if err != nil {
		return fmt.Errorf(
			"Failed to unmarhsal cart checkout completion message JSON: %v\n",
			err.Error(),
		)
	}

	if completionMessage.JobId != jobId {
		return errors.New("Picked up incorrect cart checkout completion message")
	}

	close(terminationChannel)
	io.WriteString(writer, fmt.Sprintf("event: urls\ndata: %s\n\n", delivery.Body))
	flusher.Flush()
	delivery.Ack(false)

	return nil
}
