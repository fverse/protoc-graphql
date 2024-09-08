package plugin

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	// Set up logging to a file
	logFile, err := os.Create("logs/codegen.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating debug log file: %v\n", err)
		return nil
	}
	logger := log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &Logger{logger}
}

func (l *Logger) Log(format string, v ...interface{}) {
	if l.Logger != nil {
		l.Logger.Printf(format, v...)
	}
}
