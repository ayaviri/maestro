package main

import (
	"context"
	"database/sql"
	"maestro/internal"
	xhttp "maestro/internal/http"
	"net/http"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var youtubeService *youtube.Service
var db *sql.DB
var err error

func init() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		internal.WithTimer("initialising youtube client", func() {
			defer wg.Done()
			ctx := context.Background()
			var apiKey string = os.Getenv("GCS_API_KEY")
			youtubeService, err = youtube.NewService(ctx, option.WithAPIKey(apiKey))
			internal.HandleError(err, "Could not initialise Youtube client")
		})
	}()

	go func() {
		internal.WithTimer("initialising DB object", func() {
			defer wg.Done()
			db, err = sql.Open("sqlite3", "./db/maestro.db")
			internal.HandleError(err, "Failed to connect to database")
			err = db.Ping()
			internal.HandleError(err, "Could not ping database")
		})
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

	http.Handle("/health", loggingHandler(http.HandlerFunc(HealthResourceHandler)))
	http.Handle(
		"/videos",
		loggingHandler(
			authMiddlewareFactory.New(http.HandlerFunc(VideosResourceHandler)),
		),
	)
	http.Handle(
		"/cart",
		loggingHandler(
			authMiddlewareFactory.New(http.HandlerFunc(CartResourceHandler)),
		),
	)
	http.Handle(
		"/register",
		loggingHandler(http.HandlerFunc(RegistrationResourceHandler)),
	)
	http.Handle("/login", loggingHandler(http.HandlerFunc(LoginResourceHandler)))
	// TODO: Need to introduce TLS here
	http.ListenAndServe(":8000", nil)
}

func main() {
	internal.WithTimer("running HTTP server", initialiseServer)
}
