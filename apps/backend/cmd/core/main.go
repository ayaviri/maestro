package main

import (
	"context"
	"database/sql"
	"maestro/internal"
	xamqp "maestro/internal/amqp"
	xdb "maestro/internal/db"
	xhttp "maestro/internal/http"
	"net/http"
	"os"
	"sync"

	"github.com/ayaviri/goutils/timer"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/cors"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var httpClient http.Client
var youtubeService *youtube.Service
var db *sql.DB
var messageQueueConnection *amqp.Connection
var checkoutRequestQueue amqp.Queue
var checkoutCompletionQueue amqp.Queue
var err error

func init() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		timer.WithTimer("initialising youtube client", func() {
			defer wg.Done()
			ctx := context.Background()
			// TODO: Load from a dotenv file
			var apiKey string = os.Getenv("GCS_API_KEY")
			youtubeService, err = youtube.NewService(ctx, option.WithAPIKey(apiKey))
			internal.HandleError(err, "Could not initialise Youtube client")
		})
	}()

	go func() {
		timer.WithTimer("initialising long-lived database connection", func() {
			defer wg.Done()
			xdb.EstablishConnection(&db)
		})
	}()

	go func() {
		timer.WithTimer(
			"initialising connection with the rabbit message broker",
			func() {
				defer wg.Done()
				// TODO: Need to close the connection somewhere
				xamqp.SetupQueues(
					&messageQueueConnection,
					&checkoutRequestQueue,
					&checkoutCompletionQueue,
				)
			},
		)
	}()

	wg.Wait()
}

// ____   ___  _   _ _____ _____ ____
//|  _ \ / _ \| | | |_   _| ____|  _ \
//| |_) | | | | | | | | | |  _| | |_) |
//|  _ <| |_| | |_| | | | | |___|  _ <
//|_| \_\\___/ \___/  |_| |_____|_| \_\
//

func initialiseServer() {
	authMiddlewareFactory := xhttp.BearerTokenAuthMiddlewareFactory{DB: db}
	loggingHandler := xhttp.NewLoggingHandler(os.Stdout)
	c := cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle(
		"/health",
		c.Handler((loggingHandler(http.HandlerFunc(HealthResourceHandler)))),
	)
	http.Handle(
		"/videos",
		loggingHandler(
			c.Handler(
				authMiddlewareFactory.New(http.HandlerFunc(VideosResourceHandler)),
			),
		),
	)
	http.Handle(
		"/cart",
		loggingHandler(
			c.Handler(authMiddlewareFactory.New(http.HandlerFunc(CartResourceHandler))),
		),
	)
	http.Handle(
		"/register",
		loggingHandler(c.Handler(http.HandlerFunc(RegistrationResourceHandler))),
	)
	http.Handle(
		"/login",
		loggingHandler(c.Handler(http.HandlerFunc(LoginResourceHandler))),
	)
	http.Handle(
		"/checkout",
		loggingHandler(
			authMiddlewareFactory.New(http.HandlerFunc(CheckoutResourceHandler)),
		),
	)
	http.Handle(
		"/job/",
		loggingHandler(
			authMiddlewareFactory.New(http.HandlerFunc(JobResourceHandler)),
		),
	)
	http.Handle(
		"/download/",
		loggingHandler(
			authMiddlewareFactory.New(http.HandlerFunc(DownloadResourceHandler)),
		),
	)
	// TODO: Need to introduce TLS here
	http.ListenAndServe(":8000", nil)
}

func main() {
	timer.WithTimer("running HTTP server", initialiseServer)
}
