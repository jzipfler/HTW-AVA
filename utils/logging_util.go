package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

var (
	loggingBuffer bytes.Buffer
	logger        *log.Logger
)

// Creates a new logger which uses a buffer where he collects the messages.
func InitializeLogger(logFile, preface string) {
	if logFile == "path/to/logfile.txt" {
		logger = log.New(os.Stdout, preface, log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = log.New(&loggingBuffer, preface, log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// This function is used as a wrapper for the logging functinality.
// If the logger is not initialized, it calls the log methods from the
// log package.
func PrintMessage(message interface{}) {
	if logger == nil {
		switch typeValue := message.(type) {
		case fmt.Stringer:
			log.Println(typeValue.String())
		default:
			log.Println(typeValue)
		}
	} else {
		switch typeValue := message.(type) {
		case fmt.Stringer:
			logger.Println(typeValue.String())
		default:
			logger.Println(typeValue)
		}
	}
}

// This function prints the content of the logging buffer to the STDOUT.
// TODO: Decide if it should be printed to console or wrote to a file or simmilar.
func PrintAndClearLoggerContent() {
	if loggingBuffer.Len() != 0 {
		fmt.Println(&loggingBuffer)
		loggingBuffer.Reset()
	}
}
