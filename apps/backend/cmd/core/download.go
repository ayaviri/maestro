package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DownloadResourceHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// NOTE: There will be no check of whether the user making this request
	// is downloading their own cart items here
	var filePath string = strings.TrimPrefix(request.URL.Path, "/download/")
	// TODO: Make the file server host name and port an environment variable
	response, err := httpClient.Get("http://localhost:8001/" + filePath)

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Failed to request file server: %v\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	defer response.Body.Close()
	writer.WriteHeader(response.StatusCode)

	for key, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	_, err = io.Copy(writer, response.Body)

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf(
				"Failed to copy response body from file server: %v\n",
				err.Error(),
			),
			http.StatusInternalServerError,
		)
		return
	}
}
