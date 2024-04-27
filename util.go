package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
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
		Error:   true,
		Message: err.Error(),
	}
}

func CmdLogger(prefix string, cmd *exec.Cmd) error {
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go ReaderLogger(prefix, errPipe)
	go ReaderLogger(prefix, outPipe)
	return nil
}

func ReaderLogger(prefix string, reader io.Reader) error {
	for {
		buffer := make([]byte, 1024)
		_, err := reader.Read(buffer)
		if err != nil {
			return err
		}
		lines := strings.Split(string(buffer[:]), "\n")
		for _, line := range lines {
			fmt.Printf("[%s] %s\n", prefix, line)
		}
	}
}

func ReadLogger(reader io.ReadCloser) error {
	for {
	}
}
