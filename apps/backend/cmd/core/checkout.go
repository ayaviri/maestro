package main

import (
	"encoding/json"
	"fmt"
	"maestro/internal"
	xdb "maestro/internal/db"
	xhttp "maestro/internal/http"
	xworker "maestro/internal/worker"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"net/http"
)

//  ____   ____ _   _ _____ __  __    _    ____
// / ___| / ___| | | | ____|  \/  |  / \  / ___|
// \___ \| |   | |_| |  _| | |\/| | / _ \ \___ \
//  ___) | |___|  _  | |___| |  | |/ ___ \ ___) |
// |____/ \____|_| |_|_____|_|  |_/_/   \_\____/
//

type CheckoutResponseBody struct {
	JobId string `json:"job_id"`
}

func CheckoutResourceHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user xdb.User

	internal.WithTimer("getting user from auth bearer token", func() {
		var bearerToken string
		bearerToken, _ = xhttp.GetAuthBearerToken(request)
		user, err = xdb.GetUserFromToken(db, bearerToken)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not get user from bearer token: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	var channel *amqp.Channel
	channel, err = messageQueueConnection.Channel()

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf(
				`Could not construct channel of communication with RabbitMQ server: %v\n`,
				err.Error(),
			),
			http.StatusInternalServerError,
		)
	}

	jobId := uuid.NewString()

	// var jobId string

	// internal.WithTimer("creating job ID and writing it to the database", func() {
	// 	jobId, err = xdb.CreateNewJob(db)
	// })

	// if err != nil {
	// 	http.Error(
	// 		writer,
	// 		fmt.Sprintf("Could not create job in database: %v\n", err.Error()),
	// 		http.StatusInternalServerError,
	// 	)
	// 	return
	// }

	internal.WithTimer("posting checkout request message for worker", func() {
		var requestMessage []byte
		requestMessage, err = json.Marshal(
			xworker.CheckoutRequestMessage{UserId: user.Id, JobId: jobId},
		)

		if err != nil {
			return
		}

		err = channel.Publish(
			"",
			checkoutRequestQueue.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         requestMessage,
			},
		)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf(
				"Could not publish message to checkout message queue: %v\n",
				err.Error(),
			),
			http.StatusInternalServerError,
		)
		return
	}

	var responseBody []byte
	responseBody, err = json.Marshal(CheckoutResponseBody{JobId: jobId})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not marshal response into JSON: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writer.Write(responseBody)
}
