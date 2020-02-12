package zlogging

import (
	"log"
)

type LogEvent struct {
	Level   string
	Message string
	Details interface{}
}

func WriteError(message string, err error) {
	log.Println(LogEvent{
		Level:   "error",
		Message: message,
		Details: err,
	})
}
