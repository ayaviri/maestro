package utils

import (
	"fmt"
	"time"
)

func WithTimer(message string, task func()) {
    var startTime time.Time = time.Now()
    fmt.Println("Started " + message)
    task()
    var duration float64 = time.Since(startTime).Seconds()
    fmt.Printf("Finished %s in %.2f seconds\n", message, duration)
}
