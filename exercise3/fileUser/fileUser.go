package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/goprotobuf/proto"

	"github.com/jzipfler/htw-ava/protobuf"
	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

var (
	filename    string
	processId   int
	logFile     string
	ipAddress   string
	port        int
	managedFile *os.File
	force       bool
	managerA    string
	managerB    string

	serverObject server.NetworkServer
)

const (
	GET      = iota // 0
	RELEASE         // 1
	RENOUNCE        // 2
)

func init() {
	flag.StringVar(&logFile, "logFile", "path/to/logfile.txt", "This parameter can be used to print the logging output to the given file.")
	flag.StringVar(&ipAddress, "ipAddress", "127.0.0.1", "The ip address of the actual starting node.")
	flag.IntVar(&port, "port", 15100, "The port of the actual starting node.")
	flag.StringVar(&managerA, "managerA", "127.0.0.1:15100", "The ip address and port of manager A.")
	flag.StringVar(&managerB, "managerB", "127.0.0.1:15100", "The ip address and port of manager B.")
}

func main() {

	var containsAddress, containsPort, containsManagerA, containsManagerB bool
	for _, argument := range os.Args {
		if strings.Contains(argument, "-ipAddress") {
			containsAddress = true
		}
		if strings.Contains(argument, "-port") {
			containsPort = true
		}
		if strings.Contains(argument, "-managerA") {
			containsManagerA = true
		}
		if strings.Contains(argument, "-managerB") {
			containsManagerB = true
		}
	}
	if !containsAddress {
		log.Printf("A IP address is required.\n%s\n\n", utils.ERROR_FOOTER)
		flag.Usage()
		os.Exit(0)
	}
	if !containsPort {
		log.Printf("A port number is required.\n%s\n\n", utils.ERROR_FOOTER)
		flag.Usage()
		os.Exit(0)
	}
	if !containsManagerA || !containsManagerB {
		log.Printf("The information for manager A and B are required.\n%s\n\n", utils.ERROR_FOOTER)
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	if logFile == "path/to/logfile.txt" {
		logFile = ""
	}

	utils.InitializeLogger(logFile, "")

	managerAObject, err := parseManagerInformation(managerA)
	if err != nil {
		log.Fatalf("%s\n%s\n\n", err.Error(), utils.ERROR_FOOTER)
	}
	managerAObject.SetClientName("ManagerA")
	managerBObject, err := parseManagerInformation(managerB)
	if err != nil {
		log.Fatalf("%s\n%s\n\n", err.Error(), utils.ERROR_FOOTER)
	}
	managerBObject.SetClientName("ManagerB")

	processId = os.Getpid()

	serverObject = server.New()

	serverObject.SetClientName(string(processId))
	serverObject.SetIpAddressAsString(ipAddress)
	serverObject.SetPort(port)
	serverObject.SetUsedProtocol("tcp")

	utils.PrintMessage(fmt.Sprintf("Server with the following settings created: %s", serverObject))

	if err := server.StartServer(serverObject, nil); err != nil {
		log.Fatalln("Could not start server. --> Exit.")
	}
	defer server.StopServer()

	var work func()

	if processId%2 == 0 {
		work = workerFunctionForEvenProcesses
	} else {
		work = workerFunctionForUnevenProcesses
	}

	for {
		work()
	}
}

func parseManagerInformation(managerInformation string) (server.NetworkServer, error) {
	serverObject := server.New()
	managerInformation = strings.Trim(managerInformation, " \t")
	if managerInformation == "" {
		return serverObject, errors.New("The information about the manager was empty.")
	}
	if !strings.Contains(managerInformation, ":") {
		return serverObject, errors.New("The managerInformation must use the format \"IPADDRESS:PORT\", but no \":\" was found.")
	}
	ipAndPortArray := strings.Split(managerInformation, ":")
	port, err := strconv.Atoi(ipAndPortArray[1])
	if err != nil {
		return serverObject, err
	}
	serverObject.SetIpAddressAsString(ipAndPortArray[0])
	serverObject.SetPort(port)
	serverObject.SetUsedProtocol("tcp")
	return serverObject, nil
}

func receiveAndParseFilemanagerResponse() (*protobuf.FilemanagerResponse, error) {
	//ReceiveMessage blocks until a message comes in
	conn, err := server.ReceiveMessage()
	if err != nil {
		return nil, err
	}
	utils.PrintMessage("Incoming message")
	//Close the connection when the function exits
	defer conn.Close()
	//Create a data buffer of type byte slice with capacity of 4096
	data := make([]byte, 4096)
	//Read the data waiting on the connection and put it in the data buffer
	n, err := conn.Read(data)
	if err != nil {
		return nil, err
	}
	utils.PrintMessage("Decoding Protobuf message")
	//Create an struct pointer of type ProtobufTest.TestMessage struct
	protodata := new(protobuf.FilemanagerResponse)
	//Convert all the data retrieved into the ProtobufTest.TestMessage struct type
	err = proto.Unmarshal(data[0:n], protodata)
	if err != nil {
		return nil, err
	}
	utils.PrintMessage("Message decoded.")
	return protodata, nil
}

func sendFilemanagerRequest(destinationFileManager server.NetworkServer, reaction int) error {
	if destinationFileManager.IpAndPortAsString() == "" {
		return errors.New(fmt.Sprintf("The target server information has no ip address or port.\n%s\n", destinationFileManager.IpAndPortAsString(), utils.ERROR_FOOTER))
	}
	utils.PrintMessage(fmt.Sprintf("Encode protobuf application message for node with IP:PORT : %s.", destinationFileManager.IpAndPortAsString()))
	protobufMessage := new(protobuf.FilemanagerRequest)
	protobufMessage.SourceIP = proto.String(serverObject.IpAddressAsString())
	protobufMessage.SourcePort = proto.Int(serverObject.Port())
	protobufMessage.SourceID = proto.Int(processId)
	var accessOperation protobuf.FilemanagerRequest_AccessOperation
	switch reaction {
	case GET:
		accessOperation = protobuf.FilemanagerRequest_GET
	case RELEASE:
		accessOperation = protobuf.FilemanagerRequest_RELEASE
	case RENOUNCE:
		accessOperation = protobuf.FilemanagerRequest_RENOUNCE
	default:
		return errors.New("The given reaction do not matches any case.")
	}
	protobufMessage.AccessOperation = &accessOperation
	//Protobuf message filled with data. Now marshal it.
	data, err := proto.Marshal(protobufMessage)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", destinationFileManager.IpAndPortAsString())
	if err != nil {
		return err
	}
	defer conn.Close()
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	utils.PrintMessage(fmt.Sprintf("Application message from %s to %s sent:\n\n%s\n\n", serverObject.String(), destinationFileManager.IpAndPortAsString(), protobufMessage.String()))
	utils.PrintMessage("Sent " + strconv.Itoa(n) + " bytes")
	return nil
}

func workerFunctionForEvenProcesses() {
	//Get write access on A then on B
	//Increase A and decrease B
	if err := sendFilemanagerRequest(server.New(), GET); err != nil {
		utils.PrintMessage(err)
	}
	receivedMessageFromManagerA, err := receiveAndParseFilemanagerResponse()
	if err != nil {
		utils.PrintMessage(err)
	}
	switch receivedMessageFromManagerA.GetRequestReaction() {
	case protobuf.FilemanagerResponse_ACCESS_GRANTED:
		utils.PrintMessage("Access granted from manager A")
	case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
		fallthrough
	case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
		log.Fatalln("Received wrong answer from the server.")
	case protobuf.FilemanagerResponse_ACCESS_DENIED:
		fallthrough
	default:
		utils.PrintMessage("Access denied from manager A")
		time.Sleep(1 * time.Second)
		return
	}
	if err := sendFilemanagerRequest(server.New(), GET); err != nil {
		utils.PrintMessage(err)
	}
	receivedMessageFromManagerB, err := receiveAndParseFilemanagerResponse()
	if err != nil {
		utils.PrintMessage(err)
	}
	switch receivedMessageFromManagerB.GetRequestReaction() {
	case protobuf.FilemanagerResponse_ACCESS_GRANTED:
		utils.PrintMessage("Access granted from manager A")
	case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
		fallthrough
	case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
		log.Fatalln("Received wrong answer from the server.")
	case protobuf.FilemanagerResponse_ACCESS_DENIED:
		fallthrough
	default:
		utils.PrintMessage("Access denied from manager A")
		time.Sleep(1 * time.Second)
		return
	}
	utils.IncreaseNumbersFromFirstLine(*receivedMessageFromManagerA.Filename, 6)
	utils.AppendStringToFile(*receivedMessageFromManagerB.Filename, string(processId), true)
	utils.DecreaseNumbersFromFirstLine(*receivedMessageFromManagerB.Filename, 6)
	utils.AppendStringToFile(*receivedMessageFromManagerA.Filename, string(processId), true)
	if err := sendFilemanagerRequest(server.New(), RELEASE); err != nil {
		utils.PrintMessage(err)
	}
	receivedMessageFromManagerB, err = receiveAndParseFilemanagerResponse()
	if err != nil {
		utils.PrintMessage(err)
	}
	switch receivedMessageFromManagerA.GetRequestReaction() {
	case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
		utils.PrintMessage("Resource from manager A successfully released.")
	case protobuf.FilemanagerResponse_ACCESS_GRANTED:
		fallthrough
	case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
		fallthrough
	case protobuf.FilemanagerResponse_ACCESS_DENIED:
		fallthrough
	default:
		log.Fatalln("Received wrong answer from the server.")
	}
	if err := sendFilemanagerRequest(server.New(), RELEASE); err != nil {
		utils.PrintMessage(err)
	}
	receivedMessageFromManagerB, err = receiveAndParseFilemanagerResponse()
	if err != nil {
		utils.PrintMessage(err)
	}
	switch receivedMessageFromManagerA.GetRequestReaction() {
	case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
		utils.PrintMessage("Resource from manager B successfully released.")
	case protobuf.FilemanagerResponse_ACCESS_GRANTED:
		fallthrough
	case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
		fallthrough
	case protobuf.FilemanagerResponse_ACCESS_DENIED:
		fallthrough
	default:
		log.Fatalln("Received wrong answer from the server.")
	}
}

func workerFunctionForUnevenProcesses() {
	//Get write access on B then on A
	//Increase B and decrease A
	if err := sendFilemanagerRequest(server.New(), GET); err != nil {
		utils.PrintMessage(err)
	}
	receivedMessageFromManagerB, err := receiveAndParseFilemanagerResponse()
	if err != nil {
		utils.PrintMessage(err)
	}
	switch receivedMessageFromManagerB.GetRequestReaction() {
	case protobuf.FilemanagerResponse_ACCESS_GRANTED:
		utils.PrintMessage("Access granted from manager A")
	case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
		fallthrough
	case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
		log.Fatalln("Received wrong answer from the server.")
	case protobuf.FilemanagerResponse_ACCESS_DENIED:
		fallthrough
	default:
		utils.PrintMessage("Access denied from manager A")
		time.Sleep(1 * time.Second)
		return
	}
	if err := sendFilemanagerRequest(server.New(), GET); err != nil {
		utils.PrintMessage(err)
	}
	receivedMessageFromManagerA, err := receiveAndParseFilemanagerResponse()
	if err != nil {
		utils.PrintMessage(err)
	}
	switch receivedMessageFromManagerA.GetRequestReaction() {
	case protobuf.FilemanagerResponse_ACCESS_GRANTED:
		utils.PrintMessage("Access granted from manager A")
	case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
		fallthrough
	case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
		log.Fatalln("Received wrong answer from the server.")
	case protobuf.FilemanagerResponse_ACCESS_DENIED:
		fallthrough
	default:
		utils.PrintMessage("Access denied from manager A")
		time.Sleep(1 * time.Second)
		return
	}
	utils.IncreaseNumbersFromFirstLine(*receivedMessageFromManagerB.Filename, 6)
	utils.AppendStringToFile(*receivedMessageFromManagerB.Filename, string(processId), true)
	utils.DecreaseNumbersFromFirstLine(*receivedMessageFromManagerA.Filename, 6)
	utils.AppendStringToFile(*receivedMessageFromManagerA.Filename, string(processId), true)
	if err := sendFilemanagerRequest(server.New(), RELEASE); err != nil {
		utils.PrintMessage(err)
	}
	receivedMessageFromManagerB, err = receiveAndParseFilemanagerResponse()
	if err != nil {
		utils.PrintMessage(err)
	}
	switch receivedMessageFromManagerB.GetRequestReaction() {
	case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
		utils.PrintMessage("Resource from manager B successfully released.")
	case protobuf.FilemanagerResponse_ACCESS_GRANTED:
		fallthrough
	case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
		fallthrough
	case protobuf.FilemanagerResponse_ACCESS_DENIED:
		fallthrough
	default:
		log.Fatalln("Received wrong answer from the server.")
	}
	if err := sendFilemanagerRequest(server.New(), RELEASE); err != nil {
		utils.PrintMessage(err)
	}
	receivedMessageFromManagerA, err = receiveAndParseFilemanagerResponse()
	if err != nil {
		utils.PrintMessage(err)
	}
	switch receivedMessageFromManagerA.GetRequestReaction() {
	case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
		utils.PrintMessage("Resource from manager A successfully released.")
	case protobuf.FilemanagerResponse_ACCESS_GRANTED:
		fallthrough
	case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
		fallthrough
	case protobuf.FilemanagerResponse_ACCESS_DENIED:
		fallthrough
	default:
		log.Fatalln("Received wrong answer from the server.")
	}
}
