package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"maestro/internal"
	xamqp "maestro/internal/amqp"
	xdb "maestro/internal/db"
	xworker "maestro/internal/worker"
	"os"
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
var CORE_SERVER_ADDRESS string
var DOWNLOAD_DIRECTORY string

func init() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		timer.WithTimer(
			"initialising connection to message broker server, declaring necessary queues",
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

	go func() {
		timer.WithTimer("reading from environment variables", func() {
			defer wg.Done()
			CORE_SERVER_ADDRESS = HealthCheckCoreServer()
			DOWNLOAD_DIRECTORY = os.Getenv("DOWNLOAD_DIRECTORY")

			if DOWNLOAD_DIRECTORY == "" {
				log.Fatalf("Read empty download directory")
			}
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
				var cartDownloadDirectory string = path.Join(
					DOWNLOAD_DIRECTORY, requestMessage.JobId,
				)
				var filePaths []string
				filePaths, err = DownloadCart(
					db,
					requestMessage.UserId,
					cartDownloadDirectory,
				)
				xworker.RejectOnError(err, "Failed to download cart contents", d)
				downloadUrls = make([]string, len(filePaths))

				for index, filePath := range filePaths {
					stripped := strings.TrimPrefix(filePath, DOWNLOAD_DIRECTORY+"/")
					downloadUrls[index] = fmt.Sprintf(
						"%s/download/%s",
						CORE_SERVER_ADDRESS,
						stripped,
					)
				}
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
