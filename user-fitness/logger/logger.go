package logger

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// we create a logger interface for flexible logging
type Logger interface {
	//handle and format message for regular errors
	Error(message string, args ...interface{})
	//http error handling with message format
	HttpError(w http.ResponseWriter, statusCode int, message string, args ...interface{})
	Info(message string, args ...interface{})
}

type StdLogger struct{}

// this function returns a StdLogger with the fields of the logger interface, or satisfies its contract
func NewLogger() Logger {
	return &StdLogger{}
}

// this will be called on the logger struct, which is a type of the logger interface
func (l *StdLogger) Error(format string, args ...interface{}) {
	//allows us to format messages with the errors, instead of just having the error by itself
	errMsg := fmt.Sprintf(format, args...)
	log.Printf("Error: %s", errMsg)
}

// HttpError logs an HTTP error with the specified status code and message
func (l *StdLogger) HttpError(w http.ResponseWriter, statusCode int, message string, args ...interface{}) {
	//handle nil pointers and prevent panic
	if w == nil {
		l.Error("Http Error (Status %d): %s (nil response writer)", statusCode, message)
		return
	}

	//log error
	l.Error("Http Error (Status %d): %s", statusCode, message)
	//write message to the client
	http.Error(w, fmt.Sprintf(message, args...), statusCode)
}

func (l *StdLogger) Info(format string, args ...interface{}) {
	infoMsg := fmt.Sprintf(format, args...)
	log.Printf("Info: %s", infoMsg)

}

//create
//loggerInfo
//loggerWarning
//loggerError
