// ClientStarter
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jzipfler/htw-ava/client"
	"github.com/jzipfler/htw-ava/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	loggingBuffer bytes.Buffer
	logger        *log.Logger
	id            int
	ipAddress     string
	port          int
	nodeListFile  string
)

// The ini function is called before the main function is started.
// It is used to do some stuff that should be done before the rest
// of the application starts and on what they may depend on.
func init() {
	flag.StringVar(&nodeListFile, "nodeList", "path/to/nodeList.txt", "A file where nodes are defined as follows: \"ID IP_ADDR:PORT\"")
	flag.IntVar(&id, "id", 1, "The if of the actual starting node.")
	flag.StringVar(&ipAddress, "ipAddress", "127.0.0.1", "The ip address of the actual starting node.")
	flag.IntVar(&port, "port", 15108, "The port of the actual starting node.")
}

// The main function is used when the programm is called / executed.
func main() {
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
	flag.Parse()
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

	go signalHandler()

	for {
		fmt.Println("Test")
		time.Sleep(1 * time.Second)
	}
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

func signalHandler() {
	c := make(chan os.Signal, 1)
	// Signal für CTRL-C abfangen...
	signal.Notify(c, syscall.SIGINT)

	// foreach signal received
	for signal := range c {
		fmt.Println("\nSignal empfangen...")
		fmt.Println(signal.String())
		switch signal {
		case syscall.SIGINT:
			os.Exit(0)
		}
	}
}
