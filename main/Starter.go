// ClientStarter
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	collectAllNodesFromFile()
	go signalHandler() // Handle CTRL-C signals
	initializeLogger(loggingPrefix)

	if isController {
		startController()
	} else {
		startIndependentNode()
	}
}

func collectAllNodesFromFile() {
	if nodeListFile == "path/to/nodeList.txt" {
		fmt.Printf("The nodeListFile is required.\n%s\n", ERROR_FOOTER)
		os.Exit(1)
	}
	if err := utils.CheckIfFileIsReadable(nodeListFile); err != nil {
		fmt.Printf("%s\n%s\n", err.Error(), ERROR_FOOTER)
		os.Exit(1)
	}
	nodeListFileObject, _ := os.Open(nodeListFile)
	defer nodeListFileObject.Close()
	allNodes = make(map[int]server.NetworkServer, 10)
	scanner := bufio.NewScanner(nodeListFileObject)
	for scanner.Scan() {
		var scanId int
		var scanServerObject server.NetworkServer
		line := scanner.Text()
		if line == "" {
			log.Printf("Leere Zeile gelesen.\n")
			continue
		}
		if strings.HasPrefix(line, "#") {
			log.Printf("Kommentar gelesen: \"%s\"\n", line)
			continue
		}
		idAndIpPortArray := strings.Split(line, " ")
		if len(idAndIpPortArray) == 0 {
			log.Printf("Could not split the line with a space: \"%s\".\n%s\n", line, ERROR_FOOTER)
			continue
		}
		scanId, err := strconv.Atoi(idAndIpPortArray[0])
		if err != nil {
			log.Printf("Could not parse the first part of the line to a number : \"%s\".\n%s\n", idAndIpPortArray[0], ERROR_FOOTER)
			continue
		} else {
			scanServerObject.SetClientName(idAndIpPortArray[0])
		}
		ipAndPortArray := strings.Split(idAndIpPortArray[1], ":")
		if len(ipAndPortArray) == 0 {
			log.Printf("Could not split the ip address and port with a colon: \"%s\".\n%s\n", idAndIpPortArray[1], ERROR_FOOTER)
			continue
		}
		if splitIpArray, err := net.LookupIP(ipAndPortArray[0]); err != nil {
			if splitIpHostnameArray, err := net.LookupHost(ipAndPortArray[0]); err != nil {
				log.Printf("Could not lookup this ip/host: \"%s\".\n%s\n", ipAndPortArray[0], ERROR_FOOTER)
				continue
			} else {
				scanServerObject.SetIpAddressAsString(splitIpHostnameArray[0])
			}
		} else {
			if len(splitIpArray) == 0 {
				log.Printf("No ip found: \"%s\".\n%s\n", ipAndPortArray[0], ERROR_FOOTER)
				continue
			}
			scanServerObject.SetIpAddress(splitIpArray[0])
		}
		if splitPort, err := strconv.Atoi(ipAndPortArray[1]); err != nil {
			log.Printf("Could not parse the port: \"%s\".\n%s\n", ipAndPortArray[1], ERROR_FOOTER)
			continue
		} else {
			scanServerObject.SetPort(splitPort)
		}
		allNodes[scanId] = scanServerObject
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	if len(allNodes) == 0 {
		log.Fatal("No nodes present... ABORT")
	}
}

func startIndependentNode() {
	printMessage("Start current instance as independent node.")
	printMessage("This node has the folowing settings: ")
	thisNode = allNodes[id]
	printMessage(thisNode)
}

func startController() {
	printMessage("Start current instance as controller.")
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
