package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jzipfler/htw-ava/filehandler"
	"github.com/jzipfler/htw-ava/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	ERROR_HEADER = "-------------ERROR-------------"
	ERROR_FOOTER = "^^^^^^^^^^^^^ERROR^^^^^^^^^^^^^"
)

var (
	loggingBuffer bytes.Buffer
	logger        *log.Logger
	loggingPrefix string
	logFile       string
	id            int
	ipAddress     string
	port          int
	nodeListFile  string
	isController  bool
	allNodes      map[int]server.NetworkServer
	neighbor      map[int]server.NetworkServer
	thisNode      server.NetworkServer
)

// The ini function is called before the main function is started.
// It is used to do some stuff that should be done before the rest
// of the application starts and on what they may depend on.
func init() {
	flag.StringVar(&nodeListFile, "nodeList", "path/to/nodeList.txt", "A file where nodes are defined as follows: \"ID IP_ADDR:PORT\"")
	flag.IntVar(&id, "id", 1, "The if of the actual starting node.")
	flag.StringVar(&ipAddress, "ipAddress", "127.0.0.1", "The ip address of the actual starting node.")
	flag.IntVar(&port, "port", 15108, "The port of the actual starting node.")
	flag.StringVar(&loggingPrefix, "loggingPrefix", "LOGGING --> ", "This can be used to define which prefix the logger should use to print his messages.")
	flag.StringVar(&logFile, "logFile", "path/to/logfile.txt", "This parameter can be used to print the logging output to the given file.")
	flag.BoolVar(&isController, "isController", false, "Tell the node if he should act as controller or as independent node.")
}

// The main function is used when the programm is called / executed.
func main() {
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
	flag.Parse()

	if nodeListFile == "path/to/nodeList.txt" {
		log.Fatalf("The nodeListFile is required.\n%s\n", ERROR_FOOTER)
		os.Exit(1)
	}
	var readFromNodeListError error
	allNodes, readFromNodeListError = filehandler.CollectAllFromNodeListFile(nodeListFile)
	if readFromNodeListError != nil {
		log.Fatalf("%s\n%s\n", readFromNodeListError.Error(), ERROR_FOOTER)
	}

	go signalHandler() // Handle CTRL-C signals
	initializeLogger(loggingPrefix)

	if isController {
		startController()
	} else {
		startIndependentNode()
	}
}

func startIndependentNode() {
	quit := false
	printMessage("Start current instance as independent node.")
	thisNode = allNodes[id]
	printMessage("This node has the folowing settings: ")
	printMessage(thisNode)

	for !quit {
		quit = shouldRestartProgram()
	}
}

func startController() {
	quit := false
	printMessage("Start current instance as controller.")

	for !quit {
		quit = shouldRestartProgram()
	}
}

// Asks the user if he want to exit the program.
// Returns true if and only if the user types y or j. False otherwise.
func shouldRestartProgram() bool {
	var input string
	printMessage("Would you like to exit the program? (y/j/n)")
	fmt.Print("\nInput: ")
	if _, err := fmt.Scanln(&input); err == nil {
		switch input {
		case "y", "j":
			printMessage("Program exists.")
			return true
		case "n":
			printMessage(input)
			return false
		default:
			printMessage("Please only insert y/j for \"YES\" or n for \"NO\".\n" + ERROR_FOOTER)
			printMessage("Assume a \"n\" as input.")
			return false
		}
	} else {
		printMessage("Please only insert y/j for \"YES\" or n for \"NO\".\n" + ERROR_HEADER)
	}
	return false
}

// Creates a new logger which uses a buffer where he collects the messages.
func initializeLogger(preface string) {
	if logFile == "path/to/logfile.txt" {
		logger = log.New(os.Stdout, preface, log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = log.New(&loggingBuffer, preface, log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func printMessage(message interface{}) {
	if logger == nil {
		switch typeValue := message.(type) {
		case fmt.Stringer:
			log.Println(typeValue.String())
		case string:
			log.Println(message)
		case int:
			log.Println(strconv.Itoa(typeValue))
		}
	} else {
		switch typeValue := message.(type) {
		case fmt.Stringer:
			logger.Println(typeValue.String())
		case string:
			logger.Println(message)
		case int:
			logger.Println(strconv.Itoa(typeValue))
		}
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
