package utils

import "log"

type Logger struct {
	// logger implementation details go here
}

func NewLogger() *Logger {
	return &Logger{
		// initialize logger settings
	}
}

func (l *Logger) Info(message string) {
	log.Println("[INFO]", message)
}

func (l *Logger) Error(message string) {
	log.Println("[ERROR]", message)
}
