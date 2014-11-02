package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jzipfler/htw-ava/filehandler"
	"github.com/jzipfler/htw-ava/server"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/tabwriter"
)

const (
	ERROR_HEADER            = "------------------ERROR------------------"
	ERROR_FOOTER            = "^^^^^^^^^^^^^^^^^^ERROR^^^^^^^^^^^^^^^^^^"
	MENU_SEPERATOR          = "-----------------------------------------"
	CONTROL_TYPE_INIT       = 1
	CONTROL_TYPE_EXIT       = 2
	CONTROLLER_MENU_NOTHING = 0
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
	neighbors     map[int]server.NetworkServer
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
	}
	var readFromNodeListError error
	allNodes, readFromNodeListError = filehandler.CollectAllFromNodeListFile(nodeListFile)
	if readFromNodeListError != nil {
		log.Fatalf("%s\n%s\n", readFromNodeListError.Error(), ERROR_FOOTER)
	} else {
		if len(allNodes) == 1 {
			log.Fatalf("There is only one node in the nodeList. Ther must be at least 2.\n%s\n", ERROR_FOOTER)
		}
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

// The controller is used to control the independent nodes.
// He can initialize or shutdown the nodes.
func startController() {
	quit := false
	printMessage("Start current instance as controller.")

	for !quit {
		var input string
		printMessage("Printing main menu.")
		printMainMenu()
		fmt.Print("\nEnter ID of the node you would like to send a message.\nInput: ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			printMessage(fmt.Sprintf("Error while reading the input. Quit program...\n%s\n", ERROR_FOOTER))
			os.Exit(1)
		}
		targetId, err := strconv.Atoi(input)
		if err != nil {
			printMessage(fmt.Sprintf("No number read. You have to enter the id of the node where you want to sent a message.\n%s\n", ERROR_FOOTER))
			continue
		}
		if targetId == CONTROLLER_MENU_NOTHING {
			quit = askForProgramRestart()
			continue
		}
		printControlMessageActionMenu()
		fmt.Print("\nEnter value of the control action you would like to sent.\nInput: ")
		_, err = fmt.Scanln(&input)
		if err != nil {
			printMessage(fmt.Sprintf("Error while reading the input. Quit program...\n%s\n", ERROR_FOOTER))
			os.Exit(1)
		}
		controlAction, err := strconv.Atoi(input)
		if controlAction == CONTROLLER_MENU_NOTHING {
			quit = askForProgramRestart()
			continue
		}
		quit = askForProgramRestart()
	}
}

// The chooseThreeNeighbors function uses the allNodes map to return
// another map that contains 3 nodes at the most.
// It calls os.Exit(1) if only one node is available in the allNodes map.
func chooseThreeNeighbors() (neighbors map[int]server.NetworkServer) {
	neighbors = make(map[int]server.NetworkServer, 3)
	// If there are only 1, 2 or 3 possible neighbors...take them.
	switch len(allNodes) {
	case 1:
		printMessage(fmt.Sprintf("There is only one node in the nodeList. Ther must be at least 2.\n%s\n", ERROR_FOOTER))
		os.Exit(1)
	case 2, 3, 4:
		for key, value := range allNodes {
			if key != id {
				neighbors[key] = value
			}
		}
		return
	}
	randomObject := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(neighbors) != 3 {
		var randomNumber int
		randomNumber = randomObject.Intn(len(allNodes))
		if randomNumber == id {
			continue
		}
		if value, ok := allNodes[randomNumber]; ok {
			neighbors[randomNumber] = value
		}
	}
	return
}
// Asks the user if he want to exit the program.
// Returns true if and only if the user types y or j. False otherwise.
func askForProgramRestart() bool {
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

func printMainMenu() {
	fmt.Println("")
	tabwriterObject := new(tabwriter.Writer)
	defer tabwriterObject.Flush()
	// Format in tab-separated columns with a tab stop of 4.
	tabwriterObject.Init(os.Stdout, 0, 4, 0, '\t', 0)
	fmt.Fprintln(tabwriterObject, "ID\tIP Address\tPort\tProtocol")
	fmt.Fprintln(tabwriterObject, MENU_SEPERATOR)
	for key, value := range allNodes {
		fmt.Fprintf(tabwriterObject, "%d\t%s\t%d\t%s\n", key, value.IpAddressAsString(), value.Port(), value.UsedProtocol())
	}
	fmt.Fprintf(tabwriterObject, "\n%d\tAbort\n", CONTROLLER_MENU_NOTHING)
}

func printControlMessageActionMenu() {
	fmt.Println("")
	tabwriterObject := new(tabwriter.Writer)
	defer tabwriterObject.Flush()
	// Format in tab-separated columns with a tab stop of 4.
	tabwriterObject.Init(os.Stdout, 0, 4, 0, '\t', 0)
	fmt.Fprintln(tabwriterObject, "Value\tControl message action")
	fmt.Fprintln(tabwriterObject, MENU_SEPERATOR)
	fmt.Fprintf(tabwriterObject, "%d\tInitialize\n", CONTROL_TYPE_INIT)
	fmt.Fprintf(tabwriterObject, "%d\tShutdown\n", CONTROL_TYPE_EXIT)
	fmt.Fprintf(tabwriterObject, "\n%d\tAbort\n", CONTROLLER_MENU_NOTHING)
}

// Creates a new logger which uses a buffer where he collects the messages.
func initializeLogger(preface string) {
	if logFile == "path/to/logfile.txt" {
		logger = log.New(os.Stdout, preface, log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = log.New(&loggingBuffer, preface, log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// This function is used as a wrapper for the logging functinality.
// If the logger is not initialized, it calls the log methods from the
// log package.
func printMessage(message interface{}) {
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
func printAndClearLoggerContent() {
	if loggingBuffer.Len() != 0 {
		fmt.Println(&loggingBuffer)
		loggingBuffer.Reset()
	}
}

// The sinalHanlder function waits for a signal and catches it.
// If it is a SIGINT, the program exists with return code 0.
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
