package internal

import (
	"log"
)

func HandleError(err error, message string) {
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
