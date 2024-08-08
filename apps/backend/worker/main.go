package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"maestro/internal"
	xamqp "maestro/internal/amqp"
	xdb "maestro/internal/db"
	xworker "maestro/internal/worker"
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
			// to pull out of band
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

			internal.WithTimer("downloading cart contents", func() {
				err = DownloadCart(db, message.UserId)
				RejectOnError(err, "Failed to download cart contents", d)
			})

			internal.WithTimer("setting job status to finished", func() {
				err = xdb.UpdateJobStatus(db, message.JobId, xdb.StatusFinished)
				RejectOnError(
					err,
					"Failed to set job status to finished in database",
					d,
				)
			})

			d.Ack(false)
		}
	}()

	log.Printf("Waiting to receive messages...")
	<-forever
}
