package main

import (
	"fmt"
	"log"
	"maestro/internal"
	"net/http"
	"os"
)

// Reads the address of the core server from an environment variable,
// attempts to ping its /health endpoint. Exits with a non-zero status
// code if the request fails or receives a non 2xx status code. Returns
// the env var upon success
func HealthCheckCoreServer() string {
	var url string = os.Getenv("CORE_SERVER_ADDRESS")
	var (
		httpClient http.Client
		response   *http.Response
	)
	response, err = httpClient.Get(fmt.Sprintf("%s/health", url))
	healthCheckFailureMessage := fmt.Sprintf(
		"Could not perform health check on core server with URL: %s\n", url,
	)
	internal.HandleError(err, healthCheckFailureMessage)

	if response.StatusCode != http.StatusOK {
		log.Fatalf(healthCheckFailureMessage)
	}

	return url
}
