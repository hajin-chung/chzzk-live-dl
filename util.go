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
