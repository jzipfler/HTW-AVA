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
	logFileObject *os.File
)

// Creates a new logger which uses a buffer where he collects the messages.
func InitializeLogger(logFile, preface string) {
	if logFile == "path/to/logfile.txt" || logFile == "" {
		logger = log.New(os.Stdout, preface, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	} else {
		var err error
		if exists := CheckIfFileExists(logFile); exists {
			log.Println(fmt.Sprintf("The file %s already exists, append the logging output.", logFile))
			logFileObject, err = os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND, 0777)
		} else {
			logFileObject, err = os.Create(logFile)
			logFileObject.Chmod(0660)
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		logFileObject.WriteString("Begin logging...\n")
		logger = log.New(&loggingBuffer, preface, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
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
		printAndClearLoggerContent()
	}
}

// This function prints the content of the logging buffer to the STDOUT.
// TODO: Decide if it should be printed to console or wrote to a file or simmilar.
func printAndClearLoggerContent() {
	if loggingBuffer.Len() != 0 {
		if logFileObject == nil {
			loggingBuffer.WriteTo(os.Stdout)
			//fmt.Println(&loggingBuffer)
			loggingBuffer.Reset()
		} else {
			loggingBuffer.WriteTo(logFileObject)
		}
	}
}
