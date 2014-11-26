package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jzipfler/htw-ava/exercise1"
	"github.com/jzipfler/htw-ava/filehandler"
	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

var (
	loggingPrefix             string
	logFile                   string
	id                        int
	ipAddress                 string
	port                      int
	messageContent            string
	rumor                     string
	nodeListFile              string
	graphvizFile              string
	isController              bool
	rumorExperimentMode       bool
	allNodes                  map[int]server.NetworkServer
	neighbors                 map[int]server.NetworkServer
	messageToAllNeighborsSend bool

	//targetId      int //TODO: TMP for client / server testing
)

// The ini function is called before the main function is started.
// It is used to do some stuff that should be done before the rest
// of the application starts and on what they may depend on.
func init() {
	flag.StringVar(&nodeListFile, "nodeList", "path/to/nodeList.txt", "A file where nodes are defined as follows: \"ID IP_ADDR:PORT\"")
	flag.StringVar(&graphvizFile, "graphvizFile", "path/to/graphviz.{txt,dot}", "A graphviz-dot file written with undirected edges, for example: \"graph G { 1 -- 2; }\"")
	flag.IntVar(&id, "id", 1, "The if of the actual starting node.")
	flag.StringVar(&ipAddress, "ipAddress", "127.0.0.1", "The ip address of the actual starting node.")
	flag.IntVar(&port, "port", 15100, "The port of the actual starting node.")
	flag.StringVar(&messageContent, "messageContent", "The earth is a disc.", "The message content is the string that is sent to the other nodes.")
	flag.StringVar(&rumor, "rumor", "The earth is a disc.", "The rumor that is sent to the other nodes.")
	flag.StringVar(&loggingPrefix, "loggingPrefix", "LOGGING --> ", "This can be used to define which prefix the logger should use to print his messages.")
	flag.StringVar(&logFile, "logFile", "path/to/logfile.txt", "This parameter can be used to print the logging output to the given file.")
	flag.BoolVar(&isController, "isController", false, "Tell the node if he should act as controller or as independent node.")
	flag.BoolVar(&rumorExperimentMode, "rumorExperiment", false, "The last part of the first exercise is a experiment that can be enabled with this parameter.")
}

// The main function is used when the programm is called / executed.
func main() {
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
	flag.Parse()

	if nodeListFile == "path/to/nodeList.txt" {
		log.Fatalf("The nodeListFile is required.\n%s\n", utils.ERROR_FOOTER)
	}
	var readFromNodeListError error
	allNodes, readFromNodeListError = filehandler.CollectAllFromNodeListFile(nodeListFile)
	if readFromNodeListError != nil {
		log.Fatalf("%s\n%s\n", readFromNodeListError.Error(), utils.ERROR_FOOTER)
	} else {
		if len(allNodes) == 1 {
			log.Fatalf("There is only one node in the nodeList. Ther must be at least 2.\n%s\n", utils.ERROR_FOOTER)
		}
	}
	if rumor != "The earth is a disc." {
		messageContent = rumor
	}

	go signalHandler() // Handle CTRL-C signals
	utils.InitializeLogger(logFile, fmt.Sprintf("%s(%d)", loggingPrefix, id))

	if isController {
		controllerNode := server.New()
		controllerNode.SetClientName("Controller")
		controllerNode.SetIpAddressAsString(ipAddress)
		controllerNode.SetPort(port)
		controllerNode.SetUsedProtocol("tcp")
		exercise1.StartController(controllerNode, allNodes, messageContent)
	} else {
		if graphvizFile == "path/to/graphviz.{txt,dot}" {
			neighbors = exercise1.ChooseThreeNeighbors(id, allNodes)
		} else {
			var err error
			neighbors, err = filehandler.CollectNeighborsFromGraphvizFile(graphvizFile, id, allNodes)
			if err != nil {
				utils.PrintMessage("An error occured during the reading of the graphviz file: " + err.Error())
				utils.PrintMessage("Choose three randam neighbors instead.")
				neighbors = exercise1.ChooseThreeNeighbors(id, allNodes)
			}
		}
		// Use this to set the number of used CPUs
		//runtime.GOMAXPROCS(runtime.NumCPU())
		utils.PrintMessage(fmt.Sprintf("The following %d neighbors are chosen: %v", len(neighbors), neighbors))
		exercise1.StartIndependentNode(id, allNodes, neighbors, rumorExperimentMode)
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
