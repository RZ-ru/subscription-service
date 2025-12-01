package logger

import (
	"log"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
}

func New() *Logger {
	return &Logger{
		info:  log.New(log.Writer(), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(log.Writer(), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.info.Printf(msg+" %v", args)
	} else {
		l.info.Println(msg)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.error.Printf(msg+" %v", args)
	} else {
		l.error.Println(msg)
	}
}
