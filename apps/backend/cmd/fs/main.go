package main

import (
	xhttp "maestro/internal/http"
	"net/http"
	"os"

	"github.com/ayaviri/goutils/timer"
)

func initialiseServer() {
	loggingHandler := xhttp.NewLoggingHandler(os.Stdout)
	// TODO: Throw this in an env var or something
	http.Handle("/", loggingHandler(http.FileServer(http.Dir("downloads"))))
	http.ListenAndServe(":8001", nil)
}

func main() {
	timer.WithTimer("running file server", initialiseServer)
}
