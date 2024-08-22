package main

import (
	xhttp "maestro/internal/http"
	"net/http"
	"os"

	"github.com/ayaviri/goutils/timer"
)

func initialiseServer() {
	// TODO: Throw this in an env var or something
	loggingHandler := xhttp.NewLoggingHandler(os.Stdout)
	http.Handle("/", loggingHandler(http.FileServer(http.Dir("../static"))))
	http.ListenAndServe(":3000", nil)
}

func main() {
	timer.WithTimer("running file server", initialiseServer)
}
