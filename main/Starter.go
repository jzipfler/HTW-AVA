package main

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"flag"
	"fmt"
	"github.com/jzipfler/htw-ava/filehandler"
	"github.com/jzipfler/htw-ava/protobuf"
	"github.com/jzipfler/htw-ava/server"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"
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
	loggingBuffer             bytes.Buffer
	logger                    *log.Logger
	loggingPrefix             string
	logFile                   string
	id                        int
	ipAddress                 string
	port                      int
	nodeListFile              string
	graphvizFile              string
	isController              bool
	allNodes                  map[int]server.NetworkServer
	neighbors                 map[int]server.NetworkServer
	thisNode                  server.NetworkServer
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
	flag.IntVar(&port, "port", 15108, "The port of the actual starting node.")
	flag.StringVar(&loggingPrefix, "loggingPrefix", "LOGGING --> ", "This can be used to define which prefix the logger should use to print his messages.")
	flag.StringVar(&logFile, "logFile", "path/to/logfile.txt", "This parameter can be used to print the logging output to the given file.")
	flag.BoolVar(&isController, "isController", false, "Tell the node if he should act as controller or as independent node.")

	//TODO: TMP for client / server testing
	//flag.IntVar(&targetId, "targetId", 1, "The id from the server the message should be sent to. Must be != the id.")
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

// With this function an node that interacts independently gets started.
// He can be controlled with a controller.
func startIndependentNode() {
	printMessage("Start current instance as independent node.")
	thisNode = allNodes[id]
	printMessage("This node has the folowing settings: ")
	printMessage(thisNode)

	if graphvizFile == "path/to/graphviz.{txt,dot}" {
		neighbors = chooseThreeNeighbors()
	} else {
		var err error
		neighbors, err = filehandler.CollectNeighborsFromGraphvizFile(graphvizFile)
		if err != nil {
			printMessage("An error occured during the reading of the graphviz file: " + err.Error())
			printMessage("Choose three randam neighbors instead.")
			neighbors = chooseThreeNeighbors()
		}
	}

	printMessage(fmt.Sprintf("The following %d neighbors are chosen: %v", len(neighbors), neighbors))

	protobufChannel := make(chan *protobuf.Nachricht)
	//A goroutine that receives the protobuf message and reacts to it.
	go handleReceivedProtobufMessage(protobufChannel)
	//go func() {
	//	for {
	//		message := <-protobufChannel
	//		printMessage(fmt.Sprintf("Message received:\n\n%s\n\n", message.String()))
	//		printMessage(message)
	//	}
	//}()
	//Listen to the TCP port
	if err := server.StartServer(thisNode, nil); err != nil {
		log.Fatal("Error happened: " + err.Error())
	}
	defer server.StopServer()

	for {
		//ReceiveMessage blocks until a message comes in
		if conn, err := server.ReceiveMessage(); err == nil {
			//If err is nil then that means that data is available for us so we take up this data and pass it to a new goroutine
			go receiveAndParseIncomingProtobufMessageToChannel(conn, protobufChannel)
		}
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
		if err := sendProtobufControlMessage(targetId, controlAction); err != nil {
			printMessage(fmt.Sprintf("The following error occured while trying to send a control message: %s\n%s\n", err.Error(), ERROR_FOOTER))
		} else {
			printMessage(fmt.Sprintf("Message to node with id %d, successfully sent.", targetId))
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
		// Add only the nodes with the id which exists.
		if value, ok := allNodes[randomNumber]; ok {
			// And check here if the node already exists in the neighbors map.
			if _, ok := neighbors[randomNumber]; !ok {
				neighbors[randomNumber] = value
			}
		}
	}
	return
}

// This function sends a application message (ANWENDUNGSNACHRICHT) to the neighbor with
// the given targetId. If the id does not exists, it just returns and does nothing.
func sendProtobufApplicationMessage(targetId int) error {
	printMessage("Encode protobuf message.")
	protobufMessage := new(protobuf.Nachricht)
	protobufMessage.SourceIP = proto.String(thisNode.IpAddressAsString())
	protobufMessage.SourcePort = proto.Int(thisNode.Port())
	protobufMessage.SourceID = proto.Int(id)
	nachrichtenTyp := protobuf.Nachricht_NachrichtenTyp(protobuf.Nachricht_ANWENDUNGSNACHRICHT)
	protobufMessage.NachrichtenTyp = &nachrichtenTyp
	protobufMessage.NachrichtenInhalt = proto.String("Nachrichteninhalt")
	protobufMessage.ZeitStempel = proto.String(time.Now().UTC().String())
	//Protobuf message filled with data. Now marshal it.
	data, err := proto.Marshal(protobufMessage)
	if err != nil {
		return err
	}
	printMessage(fmt.Sprintf("Application message sent:\n\n%s\n\n", protobufMessage.String()))
	// Use the neighbors map for the dial method. And abort when the targetId is not available in the map.
	if _, ok := neighbors[targetId]; !ok {
		return errors.New(fmt.Sprintf("The node with the ID %d is not a neighbor of this one. Abort sending.\n%s\n", targetId, ERROR_FOOTER))
	}
	conn, err := net.Dial(neighbors[targetId].UsedProtocol(), neighbors[targetId].IpAndPortAsString())
	if err != nil {
		return err
	}
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	printMessage("Sent " + strconv.Itoa(n) + " bytes")
	return nil
}

// This function sends a control message (KONTROLLNACHRICHT) to the node with
// the given targetId. If the id does not exists, it just returns and does nothing.
func sendProtobufControlMessage(targetId, controlType int) error {
	printMessage(fmt.Sprintf("Encode protobuf control message for node : %d.", targetId))
	protobufMessage := new(protobuf.Nachricht)
	protobufMessage.SourceIP = proto.String(thisNode.IpAddressAsString())
	protobufMessage.SourcePort = proto.Int(thisNode.Port())
	protobufMessage.SourceID = proto.Int(id)
	nachrichtenTyp := protobuf.Nachricht_NachrichtenTyp(protobuf.Nachricht_KONTROLLNACHRICHT)
	protobufMessage.NachrichtenTyp = &nachrichtenTyp
	var kontrollTyp protobuf.Nachricht_KontrollTyp
	switch controlType {
	case CONTROL_TYPE_INIT:
		kontrollTyp = protobuf.Nachricht_KontrollTyp(protobuf.Nachricht_INITIALISIEREN)
	case CONTROL_TYPE_EXIT:
		kontrollTyp = protobuf.Nachricht_KontrollTyp(protobuf.Nachricht_BEENDEN)
	default:
		printMessage("No valid controlType given. Assume EXIT.")
		kontrollTyp = protobuf.Nachricht_KontrollTyp(protobuf.Nachricht_BEENDEN)
	}
	protobufMessage.KontrollTyp = &kontrollTyp
	protobufMessage.NachrichtenInhalt = proto.String("Nachrichteninhalt")
	protobufMessage.ZeitStempel = proto.String(time.Now().UTC().String())
	//Protobuf message filled with data. Now marshal it.
	data, err := proto.Marshal(protobufMessage)
	if err != nil {
		return err
	}
	printMessage(fmt.Sprintf("Control message sent:\n\n%s\n\n", protobufMessage.String()))
	// Use the allNodes map for the dial method. And abort when the targetId is not available in the map.
	if _, ok := allNodes[targetId]; !ok {
		return errors.New(fmt.Sprintf("The node with the ID %d is not in the node list. Abort sending.\n%s\n", targetId, ERROR_FOOTER))
	}
	conn, err := net.Dial(allNodes[targetId].UsedProtocol(), allNodes[targetId].IpAndPortAsString())
	if err != nil {
		return err
	}
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	printMessage("Sent " + strconv.Itoa(n) + " bytes")
	return nil
}

// This function uses a established connection to parse the data there to the
// protobuf message. The result gets assigned to the channel.
func receiveAndParseIncomingProtobufMessageToChannel(conn net.Conn, c chan *protobuf.Nachricht) {
	printMessage("Incoming message")
	//Close the connection when the function exits
	defer conn.Close()
	//Create a data buffer of type byte slice with capacity of 4096
	data := make([]byte, 4096)
	//Read the data waiting on the connection and put it in the data buffer
	n, err := conn.Read(data)
	if err != nil {
		log.Fatal("Error happened: " + err.Error())
	}
	fmt.Println("Decoding Protobuf message")
	//Create an struct pointer of type ProtobufTest.TestMessage struct
	protodata := new(protobuf.Nachricht)
	//Convert all the data retrieved into the ProtobufTest.TestMessage struct type
	err = proto.Unmarshal(data[0:n], protodata)
	if err != nil {
		log.Fatal("Error happened: " + err.Error())
	}
	//Push the protobuf message into a channel
	c <- protodata
}

// This function waits for a message that is sent to the channel and
// splits the handling of the message depending on the NachrichtenTyp (message type)
func handleReceivedProtobufMessage(receivingChannel chan *protobuf.Nachricht) {
	for {
		// This call blocks until a new message is available.
		message := <-receivingChannel
		printMessage(fmt.Sprintf("Message received:\n\n%s\n\n", message.String()))
		switch message.GetNachrichtenTyp() {
		case protobuf.Nachricht_KONTROLLNACHRICHT:
			printMessage("Message is of type KONTROLLNACHRICHT.")
			handleReceivedControlMessage(message)
		case protobuf.Nachricht_ANWENDUNGSNACHRICHT:
			printMessage("Message is of type ANWENDUNGSNACHRICHT.")
			handleReceivedApplicationMessage(message)
		default:
			log.Fatalln("Read a unknown \"NachrichtenTyp\"")
		}
	}
}

func handleReceivedControlMessage(message *protobuf.Nachricht) {
	switch message.GetKontrollTyp() {
	case protobuf.Nachricht_INITIALISIEREN:
		if !messageToAllNeighborsSend {
			for key := range neighbors {
				sendProtobufApplicationMessage(key)
			}
		}
	case protobuf.Nachricht_BEENDEN:
		for key := range neighbors {
			sendProtobufControlMessage(key, CONTROL_TYPE_EXIT)
		}
		printMessage("")
		os.Exit(0)
	default:
		log.Fatalln("Read a unknown \"KontrollTyp\"")
	}
}

func handleReceivedApplicationMessage(message *protobuf.Nachricht) {
	if !messageToAllNeighborsSend {
		for key := range neighbors {
			sendProtobufApplicationMessage(key)
		}
	}

	// Because the SourceID is of type int32, I have to cast it here.
	sourceId := int(message.GetSourceID())
	// Check if the node that sends the message is in the neighbors map.
	// If not, add him and send him a response.
	if _, ok := neighbors[sourceId]; !ok {
		networkServerObject := server.New()
		networkServerObject.SetClientName(strconv.Itoa(sourceId))
		networkServerObject.SetIpAddressAsString(message.GetSourceIP())
		networkServerObject.SetPort(sourceId)
		networkServerObject.SetUsedProtocol("tcp") //TODO: Maybe a different approach...
		neighbors[int(message.GetSourceID())] = networkServerObject
		//sendProtobufApplicationMessage(sourceId)
	}
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
