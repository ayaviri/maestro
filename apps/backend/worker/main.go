package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"maestro/internal"
	xamqp "maestro/internal/amqp"
	xdb "maestro/internal/db"
	xworker "maestro/internal/worker"
	"os"
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
var FS_ADDRESS string
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
			CORE_SERVER_ADDRESS = os.Getenv("CORE_SERVER_ADDRESS")

			if CORE_SERVER_ADDRESS == "" {
				log.Fatalf("Read empty core web server address")
			}

			FS_ADDRESS = os.Getenv("FS_ADDRESS")

			if FS_ADDRESS == "" {
				log.Fatalf("Read empty file server address")
			}

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

			timer.WithTimer("unmarshalling message from JSON", func() {
				err = json.Unmarshal(d.Body, &requestMessage)
			})

			xworker.RejectAndExitOnError(
				err,
				"Failed to unmarshal cart checkout request message JSON",
				d,
			)

			var wg sync.WaitGroup
			wg.Add(2)

			var downloadUrls []string

			go func() {
				timer.WithTimer("downloading cart contents", func() {
					defer wg.Done()
					downloadUrls, err = DownloadCart(db, requestMessage)
					xworker.RejectAndExitOnError(
						err,
						"Failed to download cart contents",
						d,
					)
				})
			}()

			go func() {
				timer.WithTimer("writing job to DB", func() {
					defer wg.Done()
					// A worker owned, not a DB-owned utility, since all other
					// DB-owned utilities _create_ their own ID, but the ID
					// already exists here
					err = CreateNewJob(db, requestMessage.JobId, requestMessage.UserId)
					xworker.RejectAndExitOnError(err, "Failed to write job to DB", d)
				})
			}()

			wg.Wait()
			wg.Add(2)

			go func() {
				timer.WithTimer("publishing checkout completion message", func() {
					defer wg.Done()
					var clientMessage []byte
					clientMessage, _ = json.Marshal(
						xworker.CheckoutCompletionClientMessage{
							JobId: requestMessage.JobId,
							DownloadUrl: ConstructCoreServerDownloadUrlFromJob(
								requestMessage.JobId,
							),
						},
					)
					err = channel.Publish(
						"",
						checkoutCompletionQueue.Name,
						false,
						false,
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType:  "text/plain",
							Body:         clientMessage,
						},
					)
					xworker.RejectAndExitOnError(
						err,
						"Failed to publish checkout completion message to message queue",
						d,
					)
				})
			}()

			go func() {
				timer.WithTimer("writing job completion to DB", func() {
					defer wg.Done()
					var serverResponse []byte
					serverResponse, _ = json.Marshal(
						xworker.CheckoutCompletionResponse{
							DownloadUrls: downloadUrls,
						},
					)
					err = WriteJobCompletion(
						db,
						requestMessage.JobId,
						string(serverResponse),
					)
					xworker.RejectAndExitOnError(err, "Failed to update job in DB", d)
				})
			}()

			wg.Wait()
			d.Ack(false)
		}
	}()

	log.Printf("Waiting to receive messages...")
	<-forever
}
