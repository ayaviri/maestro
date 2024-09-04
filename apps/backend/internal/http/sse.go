package http

import (
	"maestro/internal"
	"net/http"
	"time"
)

func SetSSEHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
}

func Heartbeat(
	d time.Duration,
	writer http.ResponseWriter,
	flusher http.Flusher,
	// terminationChannel chan int,
	terminationChannel *internal.SafeClosureChannel,
) {
	var heartbeat *time.Ticker = time.NewTicker(1 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-heartbeat.C:
			writer.Write([]byte("event: heartbeat\n\n"))
			flusher.Flush()
		case <-terminationChannel.C:
			return
		}
	}
}
