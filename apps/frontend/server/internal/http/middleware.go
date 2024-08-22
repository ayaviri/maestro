package http

import (
	"io"
	"net/http"

	"github.com/gorilla/handlers"
)

func NewLoggingHandler(destination io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(destination, next)
	}
}
