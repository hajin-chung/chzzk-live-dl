package main

import (
	"fmt"
	"time"
)

func FormatDate() string {
	currentTime := time.Now()
	return fmt.Sprintf(
		"%d%02d%02d-%02d%02d%02d",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second(),
	)
}

type ErrorData struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func ErrorResponse(err error) ErrorData {
	return ErrorData{
		Error: true,
		Message: err.Error(),
	}
}
