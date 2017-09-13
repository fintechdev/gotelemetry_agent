package gotelemetry

import (
	"fmt"
	"strings"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = 100
	LogLevelLog   LogLevel = 200
	LogLevelError LogLevel = 300
)

type Error struct {
	StatusCode int         // HTTP status code of the error
	Message    string      // Error message
	Data       interface{} // Additional error data returned by the server
	LogLevel   LogLevel    // Whether this is a debug message
}

func NewError(statusCode int, message string) *Error {
	return &Error{statusCode, message, nil, LogLevelError}
}

func NewErrorWithData(statusCode int, message string, data interface{}) *Error {
	return &Error{statusCode, message, data, LogLevelError}
}

func NewErrorWithFormat(statusCode int, format string, data interface{}, args ...interface{}) *Error {
	return &Error{statusCode, fmt.Sprintf(format, args...), data, LogLevelError}
}

func NewDebugError(message string, args ...interface{}) *Error {
	result := NewErrorWithFormat(-1, message, nil, args...)
	result.SetLogLevel(LogLevelDebug)

	return result
}

func NewLogError(message string, args ...interface{}) *Error {
	result := NewErrorWithFormat(-1, message, nil, args...)
	result.SetLogLevel(LogLevelLog)

	return result
}

func (e *Error) Error() string {
	e.Message = strings.TrimSpace(e.Message)

	if data, ok := e.Data.([]byte); ok {
		e.Data = strings.TrimSpace(string(data))
	}

	if e.StatusCode > 0 {
		if e.Data != nil {
			return fmt.Sprintf("%d -> %s (%s)", e.StatusCode, e.Message, e.Data)
		}

		return fmt.Sprintf("%d -> %s", e.StatusCode, e.Message)
	}

	if e.Data == nil {
		return e.Message
	}

	return fmt.Sprintf("%s (%#v)", e.Message, e.Data)
}

func (e *Error) SetLogLevel(level LogLevel) {
	e.LogLevel = level
}

func (e *Error) GetLogLevel() LogLevel {
	return e.LogLevel
}
