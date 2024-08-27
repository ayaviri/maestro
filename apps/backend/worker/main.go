package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"maestro/internal"
	xamqp "maestro/internal/amqp"
	xdb "maestro/internal/db"
	xworker "maestro/internal/worker"
	"path"
	"strings"
	"sync"

	"github.com/ayaviri/goutils/timer"
	amqp "github.com/rabbitmq/amqp091-go"
)

var messageQueueConnection *amqp.Connection
var checkoutRequestQueue amqp.Queue
var checkoutCompletionQueue amqp.Queue
var db *sql.DB
var err error

func init() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		timer.WithTimer(
			"initialising connection to rabbitmq server, declaring necessary queues",
			func() {
				defer wg.Done()
				xamqp.SetupQueues(
					&messageQueueConnection,
					&checkoutRequestQueue,
					&checkoutCompletionQueue,
				)
			},
		)
	}()

	go func() {
		timer.WithTimer("initialising long-lived database connection", func() {
			defer wg.Done()
			xdb.EstablishConnection(&db)
		})
	}()

	wg.Wait()
}

func main() {

	defer messageQueueConnection.Close()
	var channel *amqp.Channel

	timer.WithTimer(
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

	timer.WithTimer(
		"opening channel for asynchronous message stream from queue",
		func() {
			messages, err = channel.Consume(
				checkoutRequestQueue.Name,
				"",
				false,
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
			// TODO: We're going to have an issue here if we have other kinds of jobs
			// to pull out of band. We'll need some sort of controller
			var requestMessage xworker.CheckoutRequestMessage
			err = json.Unmarshal(d.Body, &requestMessage)
			xworker.RejectOnError(
				err,
				"Failed to unmarshal cart checkout request message JSON",
				d,
			)

			var downloadUrls []string

			timer.WithTimer("downloading cart contents", func() {
				// TODO: The download directory needs to be pulled into
				// some sort of environment file that the file server
				// can read from as well
				var cartDownloadDirectory string = path.Join(
					"downloads", requestMessage.JobId,
				)
				var filePaths []string
				filePaths, err = DownloadCart(
					db,
					requestMessage.UserId,
					cartDownloadDirectory,
				)
				downloadUrls = make([]string, len(filePaths))

				for index, filePath := range filePaths {
					// TODO: The address to the file server should also
					// be an environment variable
					stripped := strings.TrimPrefix(filePath, "downloads/")
					downloadUrls[index] = "http://localhost:8000/download/" + stripped
				}

				xworker.RejectOnError(err, "Failed to download cart contents", d)
			})

			// TODO: Update this to include the network location of the file server
			var completionMessage []byte

			timer.WithTimer("constructing checkout completion message", func() {
				completionMessage, err = json.Marshal(
					xworker.CheckoutCompletionMessage{
						JobId:        requestMessage.JobId,
						DownloadUrls: downloadUrls,
					},
				)
				xworker.RejectOnError(
					err,
					"Failed to marhsal checkout completion message to JSON",
					d,
				)
			})

			timer.WithTimer(
				"posting checkout completion message back to core web server",
				func() {
					err = channel.Publish(
						"",
						checkoutCompletionQueue.Name,
						false,
						false,
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType:  "text/plain",
							Body:         completionMessage,
						},
					)
					xworker.RejectOnError(
						err,
						"Failed to notify core web server of checkout completion",
						d,
					)
				},
			)

			d.Ack(false)
		}
	}()

	log.Printf("Waiting to receive messages...")
	<-forever
}
