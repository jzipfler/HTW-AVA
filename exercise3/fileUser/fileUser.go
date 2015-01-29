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

	"github.com/golang/protobuf/proto"
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
	id          int
	managedFile *os.File
	force       bool
	managerA    string
	managerB    string

	serverObject   server.NetworkServer
	managerAObject server.NetworkServer
	managerBObject server.NetworkServer
)

const (
	GET      = iota // 0
	RELEASE         // 1
	RENOUNCE        // 2
)

const (
	MANAGER_A = iota
	MANAGER_B
)

func init() {
	flag.StringVar(&logFile, "logFile", "path/to/logfile.txt", "This parameter can be used to print the logging output to the given file.")
	flag.StringVar(&ipAddress, "ipAddress", "127.0.0.1", "The ip address of the actual starting node.")
	flag.IntVar(&port, "port", 15100, "The port of the actual starting node.")
	flag.StringVar(&managerA, "managerA", "127.0.0.1:15100", "The ip address and port of manager A.")
	flag.StringVar(&managerB, "managerB", "127.0.0.1:15100", "The ip address and port of manager B.")
	flag.IntVar(&id, "id", 1337, "With this option, a optional id can be specified. If not, the id becomes the process id of this program.")
}

func main() {

	var containsAddress, containsPort, containsManagerA, containsManagerB, containsId bool
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
		if strings.Contains(argument, "-id") {
			containsId = true
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

	//Store the processId to decide later if it is an "even" or "uneven" process, or use the given id.
	if containsId && id > 0 {
		processId = id
	} else {
		processId = os.Getpid()
		log.Printf("Process id \"%d\" will be used.", processId)
	}

	if logFile == "path/to/logfile.txt" {
		logFile = ""
	}

	utils.InitializeLogger(logFile, fmt.Sprintf("%d > ", processId))

	var err error
	managerAObject, err = parseIpColonPortToNetworkServer(managerA)
	if err != nil {
		log.Fatalf("%s\n%s\n\n", err.Error(), utils.ERROR_FOOTER)
	}
	managerAObject.SetClientName("ManagerA")
	utils.PrintMessage(fmt.Sprint("ManagerA: ", managerA))
	managerBObject, err = parseIpColonPortToNetworkServer(managerB)
	if err != nil {
		log.Fatalf("%s\n%s\n\n", err.Error(), utils.ERROR_FOOTER)
	}
	managerBObject.SetClientName("ManagerB")
	utils.PrintMessage(fmt.Sprint("ManagerB: ", managerB))

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

func parseIpColonPortToNetworkServer(managerInformation string) (server.NetworkServer, error) {
	serverObject := server.New()
	managerInformation = strings.Trim(managerInformation, " \t")
	if managerInformation == "" {
		return serverObject, errors.New("The information about the entiry was empty.")
	}
	if !strings.Contains(managerInformation, ":") {
		return serverObject, errors.New("The information must use the format \"IPADDRESS:PORT\", but no \":\" was found.")
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
	utils.PrintMessage("Incoming FilemanagerResponse")
	//Close the connection when the function exits
	defer conn.Close()
	//Create a data buffer of type byte slice with capacity of 4096
	data := make([]byte, 4096)
	//Read the data waiting on the connection and put it in the data buffer
	n, err := conn.Read(data)
	if err != nil {
		return nil, err
	}
	utils.PrintMessage("Decoding Protobuf FilemanagerResponse")
	//Create an struct pointer of type ProtobufTest.TestMessage struct
	protodata := new(protobuf.FilemanagerResponse)
	//Convert all the data retrieved into the ProtobufTest.TestMessage struct type
	err = proto.Unmarshal(data[0:n], protodata)
	if err != nil {
		return nil, err
	}
	utils.PrintMessage("FilemanagerResponse decoded.")
	return protodata, nil
}

func sendFilemanagerRequest(destinationFileManager server.NetworkServer, reaction int) error {
	if destinationFileManager.IpAndPortAsString() == "" {
		return errors.New(fmt.Sprintf("The target server information has no ip address or port.\n%s\n", utils.ERROR_FOOTER))
	}
	utils.PrintMessage(fmt.Sprintf("Encode protobuf FilemanagerRequest message for node with IP:PORT : %s.", destinationFileManager.IpAndPortAsString()))
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
	utils.PrintMessage(fmt.Sprintf("FilemanagerRequest message from %s to %s sent:\n\n%s\n\n", serverObject.String(), destinationFileManager.IpAndPortAsString(), protobufMessage.String()))
	utils.PrintMessage("Sent " + strconv.Itoa(n) + " bytes")
	return nil
}

func workerFunctionForEvenProcesses() {
	//Get write access on A then on B
	//Increase A and decrease B
	receivedMessageFromManagerA, err := waitForAccessFromManagerA(false)
	if err != nil {
		log.Fatalln(err)
	}
	receivedMessageFromManagerB, err := waitForAccessFromManagerB(true)
	if err != nil {
		if err.Error() == "Released first resource, need to restart access process." {
			return
		} else {
			log.Fatalln(err)
		}
	}
	utils.IncreaseNumbersFromFirstLine(receivedMessageFromManagerA.GetFilename(), 6)
	utils.AppendStringToFile(receivedMessageFromManagerB.GetFilename(), strconv.Itoa(processId), true)
	utils.DecreaseNumbersFromFirstLine(receivedMessageFromManagerB.GetFilename(), 6)
	utils.AppendStringToFile(receivedMessageFromManagerA.GetFilename(), strconv.Itoa(processId), true)
	releaseResourceFromManagerA()
	releaseResourceFromManagerB()
}

func workerFunctionForUnevenProcesses() {
	//Get write access on B then on A
	//Increase B and decrease A
	receivedMessageFromManagerB, err := waitForAccessFromManagerB(false)
	if err != nil {
		log.Fatalln(err)
	}
	receivedMessageFromManagerA, err := waitForAccessFromManagerA(true)
	if err != nil {
		if err.Error() == "Released first resource, need to restart access process." {
			return
		} else {
			log.Fatalln(err)
		}
	}
	utils.IncreaseNumbersFromFirstLine(receivedMessageFromManagerB.GetFilename(), 6)
	utils.AppendStringToFile(receivedMessageFromManagerB.GetFilename(), strconv.Itoa(processId), true)
	utils.DecreaseNumbersFromFirstLine(receivedMessageFromManagerA.GetFilename(), 6)
	utils.AppendStringToFile(receivedMessageFromManagerA.GetFilename(), strconv.Itoa(processId), true)
	releaseResourceFromManagerB()
	releaseResourceFromManagerA()
}

func waitForAccessFromManagerA(handleDeny bool) (*protobuf.FilemanagerResponse, error) {
	for {
		if err := sendFilemanagerRequest(managerAObject, GET); err != nil {
			return nil, err
		}
		receivedMessageFromManagerA, err := receiveAndParseFilemanagerResponse()
		if err != nil {
			return nil, err
		}
		switch receivedMessageFromManagerA.GetRequestReaction() {
		case protobuf.FilemanagerResponse_ACCESS_GRANTED:
			utils.PrintMessage("Access granted from manager A")
			return receivedMessageFromManagerA, nil
		case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
			fallthrough
		case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
			log.Fatalln("Received wrong answer from manager B.")
		case protobuf.FilemanagerResponse_ACCESS_DENIED:
			fallthrough
		default:
			utils.PrintMessage("Access denied from manager A")
		}
		switch receivedMessageFromManagerA.GetProcessThatUsesResource() {
		case "":
			//Value was not set (try again?!)
			time.Sleep(3 * time.Second)
			continue
		case serverObject.IpAndPortAsString():
			//We already have access
			return receivedMessageFromManagerA, nil
		default:
			//An other process uses the resource (maybe deadlock)
		}
		if handleDeny {
			if resourceReleased := tryRecoveringPossibleDeadlock(receivedMessageFromManagerA.GetProcessThatUsesResource(), MANAGER_A); resourceReleased {
				return nil, errors.New("Released first resource, need to restart access process.")
			} else {
				continue
			}
		} else {
			time.Sleep(3 * time.Second)
			continue
		}
	}
	return nil, errors.New("This error should never happen")
}

func waitForAccessFromManagerB(handleDeny bool) (*protobuf.FilemanagerResponse, error) {
	for {
		if err := sendFilemanagerRequest(managerBObject, GET); err != nil {
			return nil, err
		}
		receivedMessageFromManagerB, err := receiveAndParseFilemanagerResponse()
		if err != nil {
			return nil, err
		}
		switch receivedMessageFromManagerB.GetRequestReaction() {
		case protobuf.FilemanagerResponse_ACCESS_GRANTED:
			utils.PrintMessage("Access granted from manager A")
			return receivedMessageFromManagerB, nil
		case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
			fallthrough
		case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
			log.Fatalln("Received wrong answer from manager A.")
		case protobuf.FilemanagerResponse_ACCESS_DENIED:
			fallthrough
		default:
			utils.PrintMessage("Access denied from manager B")
		}
		switch receivedMessageFromManagerB.GetProcessThatUsesResource() {
		case "":
			//Value was not set (try again?!)
			time.Sleep(3 * time.Second)
			continue
		case serverObject.IpAndPortAsString():
			//We already have access
			return receivedMessageFromManagerB, nil
		default:
			//A other process uses the resource (maybe deadlock)
		}
		if handleDeny {
			if resourceReleased := tryRecoveringPossibleDeadlock(receivedMessageFromManagerB.GetProcessThatUsesResource(), MANAGER_B); resourceReleased {
				return nil, errors.New("Released first resource, need to restart access process.")
			} else {
				continue
			}
		} else {
			time.Sleep(3 * time.Second)
			continue
		}
	}
	return nil, errors.New("This error should never happen")
}

func releaseResourceFromManagerA() error {
	for {
		if err := sendFilemanagerRequest(managerAObject, RELEASE); err != nil {
			return err
		}
		receivedMessageFromManagerA, err := receiveAndParseFilemanagerResponse()
		if err != nil {
			return err
		}
		switch receivedMessageFromManagerA.GetRequestReaction() {
		case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
			utils.PrintMessage("Resource from manager A successfully released.")
			return nil
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
	return errors.New("This error should never happen")
}

func releaseResourceFromManagerB() error {
	for {
		if err := sendFilemanagerRequest(managerBObject, RELEASE); err != nil {
			return err
		}
		receivedMessageFromManagerB, err := receiveAndParseFilemanagerResponse()
		if err != nil {
			return err
		}
		switch receivedMessageFromManagerB.GetRequestReaction() {
		case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
			utils.PrintMessage("Resource from manager B successfully released.")
			return nil
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
	return errors.New("This error should never happen")
}

func tryRecoveringPossibleDeadlock(processThatUsesResource string, whichManager int) bool {
	utils.PrintMessage(fmt.Sprintf("This process got already one resource the other is blocked by %s.", processThatUsesResource))
	blockingProcess, err := parseIpColonPortToNetworkServer(processThatUsesResource)
	if err != nil {
		log.Fatalln(err)
	}
	if err := sendAccessControlMessage(blockingProcess); err != nil {
		log.Fatalln(err)
	}
	receivedAccessControlMessage, err := receiveAccessControlMessage()
	if err != nil {
		log.Fatalln(err)
	}
	//The highest process id will release the resource
	if receivedAccessControlMessage.GetSourceID() < int32(processId) {
		if whichManager == MANAGER_A {
			if err := sendFilemanagerRequest(managerAObject, RENOUNCE); err != nil {
				log.Fatalln(err)
			}
		} else if whichManager == MANAGER_B {
			if err := sendFilemanagerRequest(managerBObject, RENOUNCE); err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalln("Wrong manager for deadlock recovery given.")
		}
		receivedResponse, err := receiveAndParseFilemanagerResponse()
		if err != nil {
			log.Fatalln(err)
		}
		if receivedResponse.GetRequestReaction() == protobuf.FilemanagerResponse_RESOURCE_RELEASED {
			return true
		} else {
			log.Fatalln("Received not a RESOURCE_RELEASED message for recovering deadlock.")
		}
	}
	return false
}

func receiveAccessControlMessage() (*protobuf.AccessControl, error) {
	//ReceiveMessage blocks until a message comes in
	conn, err := server.ReceiveMessage()
	if err != nil {
		return nil, err
	}
	utils.PrintMessage("Incoming AccessControl message")
	//Close the connection when the function exits
	defer conn.Close()
	//Create a data buffer of type byte slice with capacity of 4096
	data := make([]byte, 4096)
	//Read the data waiting on the connection and put it in the data buffer
	n, err := conn.Read(data)
	if err != nil {
		return nil, err
	}
	utils.PrintMessage("Decoding AccessControl message")
	//Create an struct pointer of type ProtobufTest.TestMessage struct
	protodata := new(protobuf.AccessControl)
	//Convert all the data retrieved into the ProtobufTest.TestMessage struct type
	err = proto.Unmarshal(data[0:n], protodata)
	if err != nil {
		return nil, err
	}
	utils.PrintMessage("Message decoded.")
	return protodata, nil
}

func sendAccessControlMessage(destinationFileManager server.NetworkServer) error {
	if destinationFileManager.IpAndPortAsString() == "" {
		return errors.New(fmt.Sprintf("The target server information has no ip address or port.\n%s\n", utils.ERROR_FOOTER))
	}
	utils.PrintMessage(fmt.Sprintf("Encode protobuf AccessControl message for node with IP:PORT : %s.", destinationFileManager.IpAndPortAsString()))
	protobufMessage := new(protobuf.AccessControl)
	protobufMessage.SourceIP = proto.String(serverObject.IpAddressAsString())
	protobufMessage.SourcePort = proto.Int(serverObject.Port())
	protobufMessage.SourceID = proto.Int(processId)
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
	utils.PrintMessage(fmt.Sprintf("AccessControl message from %s to %s sent:\n\n%s\n\n", serverObject.String(), destinationFileManager.IpAndPortAsString(), protobufMessage.String()))
	utils.PrintMessage("Sent " + strconv.Itoa(n) + " bytes")
	return nil
}
