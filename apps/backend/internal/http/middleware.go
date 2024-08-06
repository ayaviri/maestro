package http

import (
	"database/sql"
	"fmt"
	"io"
	xdb "maestro/internal/db"
	"net/http"

	"github.com/gorilla/handlers"
)

type BearerTokenAuthMiddlewareFactory struct {
	DB *sql.DB
}

// TODO: Need to implement token expiration here
func (f BearerTokenAuthMiddlewareFactory) New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		var bearerToken string
		bearerToken, err = GetAuthBearerToken(request)

		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Could not obtain bearer token: %v\n", err.Error()),
				http.StatusUnauthorized,
			)
			return
		}

		var isValidToken bool
		isValidToken, err = xdb.IsValidToken(f.DB, bearerToken)

		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Could not validate token: %v\n", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		if !isValidToken {
			http.Error(
				w,
				fmt.Sprintf("Invalid token: %v\n", err.Error()),
				http.StatusUnauthorized,
			)
			return
		}

		next.ServeHTTP(w, request)
	})
}

func NewLoggingHandler(destination io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(destination, next)
	}
}
