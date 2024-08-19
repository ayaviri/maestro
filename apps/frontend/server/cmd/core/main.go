package main

import (
	"net/http"

	"github.com/ayaviri/goutils/timer"
)

func initialiseServer() {
	// TODO: Throw this in an env var or something
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.ListenAndServe(":3000", nil)
}

func main() {
	timer.WithTimer("running file server", initialiseServer)
}
