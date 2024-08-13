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
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var messageQueueConnection *amqp.Connection
var checkoutMessageQueue amqp.Queue
var db *sql.DB
var err error

func init() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		internal.WithTimer(
			"initialising connection to rabbitmq server, declaring checkout queue",
			func() {
				defer wg.Done()
				xamqp.SetupCheckoutQueue(&messageQueueConnection, &checkoutMessageQueue)
			},
		)
	}()

	go func() {
		internal.WithTimer("initialising long-lived database connection", func() {
			defer wg.Done()
			xdb.EstablishConnection(&db)
		})
	}()

	wg.Wait()
}

func main() {

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
			var message xworker.DownloadCartMessageBody
			err = json.Unmarshal(d.Body, &message)

			internal.WithTimer("setting job status to received", func() {
				err = xdb.UpdateJobStatus(db, message.JobId, xdb.StatusReceived)
				RejectOnError(
					err,
					"Failed to set job status to received in database",
					d,
				)
			})

			var fileDownloadPaths []string

			internal.WithTimer("downloading cart contents", func() {
				// TODO: The download directory needs to be pulled into
				// some sort of environment file that the file server
				// can read from as well
				var cartDownloadDirectory string = path.Join("downloads", message.JobId)
				var fileNames []string
				fileNames, err = DownloadCart(
					db,
					message.UserId,
					cartDownloadDirectory,
				)
				fileDownloadPaths = make([]string, len(fileNames))

				for index, fileName := range fileNames {
					fileDownloadPaths[index] = path.Join(
						cartDownloadDirectory,
						fileName,
					)
				}

				RejectOnError(err, "Failed to download cart contents", d)
			})

			internal.WithTimer(
				"setting job status to finished + adding payload",
				func() {
					err = xdb.UpdateJobStatus(db, message.JobId, xdb.StatusFinished)
					RejectOnError(
						err,
						"Failed to set job status to finished in database",
						d,
					)
					// TODO: Update this to include the network location of the file server
					response := xworker.DownloadCartResponseBody{
						DownloadUrls: fileDownloadPaths,
					}
					var responsePayload []byte
					responsePayload, err = json.Marshal(response)
					RejectOnError(err, "Failed to marshal response to JSON", d)
					err = xdb.AddJobPayload(db, message.JobId, string(responsePayload))
					RejectOnError(err, "Failed to add response payload to database", d)
				},
			)

			d.Ack(false)
		}
	}()

	log.Printf("Waiting to receive messages...")
	<-forever
}
