package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

var err error

func GetAuthBearerToken(request *http.Request) (string, error) {
	var authHeader string = request.Header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("Authorization header required")
	}

	var bearerToken string = strings.TrimPrefix(authHeader, "Bearer ")

	if bearerToken == authHeader {
		return "", errors.New("Invalid token format")
	}

	return bearerToken, nil
}

// Reads the entirety of the given request's body and unmarshalls it into
// the given pointer to the JSON schema
func ReadUnmarshalRequestBody(request *http.Request, schema any) error {
	var requestBodyBytes []byte
	requestBodyBytes, err = io.ReadAll(request.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(requestBodyBytes, schema)

	if err != nil {
		return err
	}

	return nil
}
