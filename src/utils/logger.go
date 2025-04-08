package utils

import (
	"log"
	"os"
)

type Logger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DebugLog *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		InfoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		DebugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(message string) {
	l.InfoLog.Println(message)
}

func (l *Logger) Error(message string) {
	l.ErrorLog.Println(message)
}

func (l *Logger) Debug(message string) {
	l.DebugLog.Println(message)
}