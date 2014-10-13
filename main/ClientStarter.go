// ClientStarter
package main

import (
	"bytes"
	"fmt"
	"github.com/jzipfler/htw/ava/client"
	"github.com/jzipfler/htw/ava/server"
	"log"
)

var (
	loggingBuffer bytes.Buffer
	logger        *log.Logger
)

func main() {
	initializeLogger("LOG::: ")
	client := client.New()
	serverObject := server.New()
	logger.Printf("%T\n", client)
	logger.Printf("%T\n", serverObject)
	logger.Println("Client: " + client.String() + "\nServer: " + serverObject.String())
	error := client.SetIpAddressAsString("1.2.3.4")
	if error != nil {
		logger.Fatalln(error)
		return
	}
	client.SetClientName("First")
	logger.Println(client)

	printAndClearLoggerContent()
}

func initializeLogger(preface string) {
	logger = log.New(&loggingBuffer, preface, log.Lshortfile)
}

func printAndClearLoggerContent() {
	if loggingBuffer.Len() != 0 {
		fmt.Println(&loggingBuffer)
		loggingBuffer.Reset()
	}
}
