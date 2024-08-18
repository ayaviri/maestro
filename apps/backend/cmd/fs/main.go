package main

import (
	"maestro/internal"
	xhttp "maestro/internal/http"
	"net/http"
	"os"
)

func initialiseServer() {
	loggingHandler := xhttp.NewLoggingHandler(os.Stdout)
	// TODO: Throw this in an env var or something
	http.Handle("/", loggingHandler(http.FileServer(http.Dir("downloads"))))
	http.ListenAndServe(":8001", nil)
}

func main() {
	internal.WithTimer("running file server", initialiseServer)
}
