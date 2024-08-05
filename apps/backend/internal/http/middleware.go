package http

import (
	"database/sql"
	xdb "maestro/internal/db"
	"net/http"
	"strings"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(w, request)
	})
}

type BearerTokenAuthMiddlewareFactory struct {
	DB *sql.DB
}

// TODO: Need to implement token expiration here
func (f BearerTokenAuthMiddlewareFactory) New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		var bearerToken string
		bearerToken, err = GetAuthBearerToken(request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		var isValidToken bool
		isValidToken, err = xdb.IsValidToken(f.DB, bearerToken)

		if err != nil {
			http.Error(w, "Could not validate token", http.StatusInternalServerError)
		}

		if !isValidToken {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}

		next.ServeHTTP(w, request)
	})
}
