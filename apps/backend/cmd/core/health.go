package main

import "net/http"

func HealthResourceHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writer.Write([]byte("maestro"))
}
