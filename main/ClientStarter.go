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
	//initializeLogger("LOG::: ")
	client := client.New()
	serverObject := server.New()
	printMessage("Client: " + client.String() + "\n\t\tServer: " + serverObject.String())
	error := client.SetIpAddressAsString("1.2.3.4")
	if error != nil {
		printMessage(error.Error())
		return
	}
	client.SetClientName("First")
	printMessage(client.String())

	doServerStuff(serverObject)

	printAndClearLoggerContent()
}

func doServerStuff(serverObject server.NetworkServer) {
	serverObject.SetClientName("Server1")
	serverObject.SetIpAddressAsString("127.0.0.1")
	serverObject.SetPort(15108)
	serverObject.SetUsedProtocol("tcp")
	printMessage(serverObject.String())
	err := server.StartServer(serverObject, logger)
	if err != nil {
		return
	}
	defer server.StopServer()
	//for {
	//	connection := server.ReceiveMessage()
	//	printMessage(connection.LocalAddr().String())
	//}
}

// Creates a new logger which uses a buffer where he collects the messages.
func initializeLogger(preface string) {
	logger = log.New(&loggingBuffer, preface, log.Lshortfile)
}

func printMessage(message string) {
	if logger == nil {
		log.Println(message)
	} else {
		logger.Println(message)
	}
}

func printAndClearLoggerContent() {
	if loggingBuffer.Len() != 0 {
		fmt.Println(&loggingBuffer)
		loggingBuffer.Reset()
	}
}
