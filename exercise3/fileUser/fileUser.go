package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jzipfler/htw-ava/client"
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
	slowRunning bool
	useTCP      bool
	managerA    string
	managerB    string

	serverObject        server.NetworkServer
	managerAObject      server.NetworkServer
	managerBObject      server.NetworkServer
	localNodeUdpAddress *net.UDPAddr

	tokenServer        server.NetworkServer
	nonBlockingManager server.NetworkServer
	waitForIpAndPort   string
	waitForId          int
	blocking           bool
	gotOneResource     bool
)

const (
	SECONDS_UNTIL_NEXT_TRY = 3
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
	flag.IntVar(&port, "port", 15100, "The port of the actual starting node. (Portnumber must be even)")
	flag.StringVar(&managerA, "managerA", "127.0.0.1:15100", "The ip address and port of manager A. (Portnumber must be even)")
	flag.StringVar(&managerB, "managerB", "127.0.0.1:15100", "The ip address and port of manager B. (Portnumber must be even)")
	flag.IntVar(&id, "id", 1337, "With this option, a optional id can be specified. If not, the id becomes the process id of this program.")
	flag.BoolVar(&slowRunning, "slow", false, "If slow is set to true, the fileUser will restart each request after a time interval from [0,5) seconds randomly, otherwise between 1ms and 1s.")
	flag.BoolVar(&useTCP, "useTCP", false, "If this value is set to true, the application uses TCP to communicate.")
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

	if port%2 != 0 {
		log.Printf("The port number must be even.\n%s\n\n", utils.ERROR_FOOTER)
		os.Exit(1)
	}

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

	//Initialize random generator
	rand.Seed(time.Now().UTC().UnixNano())

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
	if useTCP {
		serverObject.SetUsedProtocol(client.TCP)
	} else {
		serverObject.SetUsedProtocol(client.UDP)
	}

	utils.PrintMessage(fmt.Sprintf("Server with the following settings created: %s", serverObject))

	if serverObject.UsedProtocol() == client.TCP {
		if err := server.StartServer(serverObject, nil); err != nil {
			log.Fatalln("Could not start server. --> Exit.")
		}
		defer server.StopServer()
	} else if serverObject.UsedProtocol() == client.UDP {
		//Do nothing except to get the UDPAdress because the UDP call can gather packages directly
		//without call the Listen and then the Accept function. (like in TCP)
		var err error
		localNodeUdpAddress, err = net.ResolveUDPAddr(serverObject.UsedProtocol(), serverObject.IpAndPortAsString())
		if err != nil {
			log.Fatalln("Error happened: Can not convert the local address information to a UdpAdressObject.")
		}
		utils.PrintMessage(fmt.Sprintf("Created UDP information for node %d: %v", processId, localNodeUdpAddress))
	} else {
		log.Fatalln("Error happened: The given protocol to start the server on the independend node, was neigther tcp nor udp.")
	}

	var work func()

	if processId%2 == 0 {
		work = workerFunctionForEvenProcesses
	} else {
		work = workerFunctionForUnevenProcesses
	}

	go handleTokenMessages()

	for {
		work()
		if slowRunning {
			time.Sleep(time.Duration(rand.Float32()*5) * time.Second)
		} else {
			time.Sleep(time.Duration(rand.Float32()*100) * time.Millisecond)
		}
	}
}

func parseIpColonPortToNetworkServer(managerInformation string) (server.NetworkServer, error) {
	serverObject := server.New()
	managerInformation = strings.Trim(managerInformation, " \t")
	if managerInformation == "" {
		return serverObject, errors.New("Can not parse IP:PORT because the incoming string was empty.")
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
	if useTCP {
		serverObject.SetUsedProtocol(client.TCP)
	} else {
		serverObject.SetUsedProtocol(client.UDP)
	}
	return serverObject, nil
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
	conn, err := net.Dial(serverObject.UsedProtocol(), destinationFileManager.IpAndPortAsString())
	if err != nil {
		return err
	}
	defer conn.Close()

	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	utils.PrintMessage(fmt.Sprintf("FilemanagerRequest message from %s to %s sent:\n\n%s\n", serverObject.String(), destinationFileManager.IpAndPortAsString(), protobufMessage.String()))
	utils.PrintMessage("Sent " + strconv.Itoa(n) + " bytes")

	return nil
}

func receiveFilemanagerResponses() *protobuf.FilemanagerResponse {
	var conn net.Conn
	var err error
	if serverObject.UsedProtocol() == client.TCP {
		//ReceiveMessage blocks until a message comes in
		conn, err = server.ReceiveMessage()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		conn, err = net.ListenUDP(serverObject.UsedProtocol(), localNodeUdpAddress)
		if err != nil {
			log.Fatalln(err)
		}
	}

	defer conn.Close()
	data := make([]byte, 4096)
	n, err := conn.Read(data)
	utils.PrintMessage("Incoming message")
	if err != nil {
		log.Fatalln(err)
	}
	protoFilemanagerResponseMessage := new(protobuf.FilemanagerResponse)
	err = proto.Unmarshal(data[0:n], protoFilemanagerResponseMessage)
	if err != nil {
		log.Fatalln(err)
	}
	utils.PrintMessage(fmt.Sprintf("FilemanagerResponse decoded.\n\n%s\n", protoFilemanagerResponseMessage))
	return protoFilemanagerResponseMessage
}

func workerFunctionForEvenProcesses() {
	//Get write access on A then on B
	//Increase A and decrease B
	receivedMessageFromManagerA, err := waitForAccessFromManagerA()
	if err != nil {
		log.Fatalln(err)
	}
	gotOneResource = true
	receivedMessageFromManagerB, err := waitForAccessFromManagerB()
	if err != nil {
		log.Fatalln(err)
	}
	if err == nil && receivedMessageFromManagerB == nil {
		blocking = false
		gotOneResource = false
		return
	}
	utils.IncreaseNumbersFromFirstLine(receivedMessageFromManagerA.GetFilename(), 6)
	utils.AppendStringToFile(receivedMessageFromManagerB.GetFilename(), strconv.Itoa(processId), true)
	utils.DecreaseNumbersFromFirstLine(receivedMessageFromManagerB.GetFilename(), 6)
	utils.AppendStringToFile(receivedMessageFromManagerA.GetFilename(), strconv.Itoa(processId), true)
	if err := releaseResourceFromManager(MANAGER_A); err != nil {
		log.Fatalln(err)
	}
	if err := releaseResourceFromManager(MANAGER_B); err != nil {
		log.Fatalln(err)
	}
	blocking = false
	gotOneResource = false
}

func workerFunctionForUnevenProcesses() {
	//Get write access on B then on A
	//Increase B and decrease A
	receivedMessageFromManagerB, err := waitForAccessFromManagerB()
	if err != nil {
		log.Fatalln(err)
	}
	gotOneResource = true
	receivedMessageFromManagerA, err := waitForAccessFromManagerA()
	if err != nil {
		log.Fatalln(err)
	}
	if err == nil && receivedMessageFromManagerA == nil {
		blocking = false
		gotOneResource = false
		return
	}
	utils.IncreaseNumbersFromFirstLine(receivedMessageFromManagerB.GetFilename(), 6)
	utils.AppendStringToFile(receivedMessageFromManagerB.GetFilename(), strconv.Itoa(processId), true)
	utils.DecreaseNumbersFromFirstLine(receivedMessageFromManagerA.GetFilename(), 6)
	utils.AppendStringToFile(receivedMessageFromManagerA.GetFilename(), strconv.Itoa(processId), true)
	if err := releaseResourceFromManager(MANAGER_B); err != nil {
		log.Fatalln(err)
	}
	if err := releaseResourceFromManager(MANAGER_A); err != nil {
		log.Fatalln(err)
	}
	blocking = false
	gotOneResource = false
}

func waitForAccessFromManagerA() (*protobuf.FilemanagerResponse, error) {
	for {
		utils.PrintMessage("Begin waitForAccessFromManagerA()")
		if err := sendFilemanagerRequest(managerAObject, GET); err != nil {
			return nil, err
		}
		var receivedMessageFromManagerA *protobuf.FilemanagerResponse
		if receivedMessageFromManagerA = receiveFilemanagerResponses(); receivedMessageFromManagerA == nil {
			continue
		}
		switch receivedMessageFromManagerA.GetRequestReaction() {
		case protobuf.FilemanagerResponse_ACCESS_GRANTED:
			utils.PrintMessage("Access granted from manager A")
			if gotOneResource {
				blocking = false
				nonBlockingManager = server.New()
			}
			return receivedMessageFromManagerA, nil
		case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
			utils.PrintMessage("Received RESOURCE_RELEASED from manager A. Releasing...")
			gotOneResource = false
			blocking = false
			nonBlockingManager = server.New()
			waitForId = 0
			waitForIpAndPort = ""
			time.Sleep(time.Duration(SECONDS_UNTIL_NEXT_TRY*100*rand.Float32()) * time.Millisecond)
			return nil, nil
		case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
			utils.PrintMessage("Received RESOURCE_NOT_RELEASED from manager A.")
			if receivedMessageFromManagerA.GetProcessIdThatUsesResource() == int32(processId) {
				utils.PrintMessage("...but we have already access to that resource, so return it to go on.")
				return receivedMessageFromManagerA, nil
			} else {
				utils.PrintMessage("...but we do not own this resource, so continue try to getting it.")
				time.Sleep(time.Duration(SECONDS_UNTIL_NEXT_TRY*100*rand.Float32()) * time.Millisecond)
				continue
			}
		case protobuf.FilemanagerResponse_ACCESS_DENIED:
			fallthrough
		default:
			utils.PrintMessage("Access denied from manager A.")
			if gotOneResource {
				utils.PrintMessage("Got already one resource, waiting for the next one.")
				targetServerObject, err := parseIpColonPortToNetworkServer(receivedMessageFromManagerA.GetProcessIpAndPortThatUsesResource())
				if err != nil {
					log.Fatalln(err)
				}
				utils.PrintMessage("Send token to WAIT-FOR node.")
				nonBlockingManager = managerBObject
				blocking = true
				waitForIpAndPort = receivedMessageFromManagerA.GetProcessIpAndPortThatUsesResource()
				waitForId = int(receivedMessageFromManagerA.GetProcessIdThatUsesResource())
				if err := sendGoldmanToken(targetServerObject, nil); err != nil {
					log.Fatalln(err)
				}
			}
			time.Sleep(time.Duration(SECONDS_UNTIL_NEXT_TRY*100*rand.Float32()) * time.Millisecond)
			continue
		}
	}
	return nil, errors.New("This error should never happen")
}

func waitForAccessFromManagerB() (*protobuf.FilemanagerResponse, error) {
	for {
		utils.PrintMessage("Begin waitForAccessFromManagerB()")
		if err := sendFilemanagerRequest(managerBObject, GET); err != nil {
			return nil, err
		}
		var receivedMessageFromManagerB *protobuf.FilemanagerResponse
		if receivedMessageFromManagerB = receiveFilemanagerResponses(); receivedMessageFromManagerB == nil {
			time.Sleep(time.Duration(SECONDS_UNTIL_NEXT_TRY*100*rand.Float32()) * time.Millisecond)
			continue
		}
		switch receivedMessageFromManagerB.GetRequestReaction() {
		case protobuf.FilemanagerResponse_ACCESS_GRANTED:
			utils.PrintMessage("Access granted from manager B")
			if gotOneResource {
				blocking = false
				nonBlockingManager = server.New()
			}
			return receivedMessageFromManagerB, nil
		case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
			utils.PrintMessage("Received RESOURCE_RELEASED from manager B. Releasing...")
			gotOneResource = false
			blocking = false
			nonBlockingManager = server.New()
			waitForId = 0
			waitForIpAndPort = ""
			time.Sleep(time.Duration(SECONDS_UNTIL_NEXT_TRY*100*rand.Float32()) * time.Millisecond)
			return nil, nil
		case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
			utils.PrintMessage("Received RESOURCE_NOT_RELEASED from manager B.")
			if receivedMessageFromManagerB.GetProcessIdThatUsesResource() == int32(processId) {
				utils.PrintMessage("...but we have already access to that resource, so return it to go on.")
				return receivedMessageFromManagerB, nil
			} else {
				utils.PrintMessage("...but we do not own this resource, so continue try to getting it.")
				time.Sleep(time.Duration(SECONDS_UNTIL_NEXT_TRY*100*rand.Float32()) * time.Millisecond)
				continue
			}
		case protobuf.FilemanagerResponse_ACCESS_DENIED:
			fallthrough
		default:
			utils.PrintMessage("Access denied from manager B.")
			if gotOneResource {
				utils.PrintMessage("Got already one resource, waiting for the next one.")
				targetServerObject, err := parseIpColonPortToNetworkServer(receivedMessageFromManagerB.GetProcessIpAndPortThatUsesResource())
				if err != nil {
					log.Fatalln(err)
				}
				utils.PrintMessage("Send token to WAIT-FOR node.")
				nonBlockingManager = managerAObject
				blocking = true
				waitForIpAndPort = receivedMessageFromManagerB.GetProcessIpAndPortThatUsesResource()
				waitForId = int(receivedMessageFromManagerB.GetProcessIdThatUsesResource())
				if err := sendGoldmanToken(targetServerObject, nil); err != nil {
					log.Fatalln(err)
				}
			}
			time.Sleep(time.Duration(SECONDS_UNTIL_NEXT_TRY*100*rand.Float32()) * time.Millisecond)
			continue
		}
	}
	return nil, errors.New("This error should never happen")
}

func releaseResourceFromManager(managerToRecover int) error {
	var receivedMessageFromManager *protobuf.FilemanagerResponse
	for {
		if managerToRecover == MANAGER_A {
			if err := sendFilemanagerRequest(managerAObject, RELEASE); err != nil {
				return err
			}
		} else if managerToRecover == MANAGER_B {
			if err := sendFilemanagerRequest(managerBObject, RELEASE); err != nil {
				return err
			}
		} else {
			log.Fatalln("Wrong manager identifier given on releaseResource")
		}
		if receivedMessageFromManager = receiveFilemanagerResponses(); receivedMessageFromManager == nil {
			continue
		}
		switch receivedMessageFromManager.GetRequestReaction() {
		case protobuf.FilemanagerResponse_RESOURCE_RELEASED:
			if managerToRecover == MANAGER_A {
				utils.PrintMessage("Resource from manager A successfully released.")
			} else if managerToRecover == MANAGER_B {
				utils.PrintMessage("Resource from manager B successfully released.")
			} else {
				log.Fatalln("Got an unknown state at releaseResourceFromManager().")
			}
			return nil
		case protobuf.FilemanagerResponse_ACCESS_GRANTED:
			fallthrough
		case protobuf.FilemanagerResponse_RESOURCE_NOT_RELEASED:
			fallthrough
		case protobuf.FilemanagerResponse_ACCESS_DENIED:
			fallthrough
		default:
			utils.PrintMessage("Trying to release the resource but don't get RESOURCE_RELEASED. Check if we own it.")
			if receivedMessageFromManager.GetProcessIdThatUsesResource() != int32(processId) {
				utils.PrintMessage("Nope, we do not own it, let us abort this and try to get access again.")
				return nil
			}
			utils.PrintMessage("Yep, we own it, lets wait for the OK.")
			continue
		}
	}
	return errors.New("This error should never happen")
}

func sendGoldmanToken(destinationNode server.NetworkServer, blockingProcesses []int32) error {
	if destinationNode.IpAndPortAsString() == "" {
		return errors.New(fmt.Sprintf("The target server information has no ip address or port.\n%s\n", utils.ERROR_FOOTER))
	}
	if destinationNode.Port()%2 == 0 {
		destinationNode.SetPort(destinationNode.Port() + 1)
	}
	if blockingProcesses == nil {
		blockingProcesses = make([]int32, 0)
	}
	utils.PrintMessage(fmt.Sprintf("Encode protobuf Token message for node with IP:PORT : %s.", destinationNode.IpAndPortAsString()))
	protobufMessage := new(protobuf.GoldmanToken)
	protobufMessage.BlockingProcesses = append(blockingProcesses, int32(processId))
	protobufMessage.SourceIP = proto.String(tokenServer.IpAddressAsString())
	protobufMessage.SourcePort = proto.Int(tokenServer.Port())
	//Protobuf message filled with data. Now marshal it.
	data, err := proto.Marshal(protobufMessage)
	if err != nil {
		return err
	}
	conn, err := net.Dial(destinationNode.UsedProtocol(), destinationNode.IpAndPortAsString())
	if err != nil {
		return err
	}
	defer conn.Close()
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	utils.PrintMessage(fmt.Sprintf("Token message from %s to %s sent:\n\n%s\n", tokenServer.String(), destinationNode.IpAndPortAsString(), protobufMessage.String()))
	utils.PrintMessage("Sent " + strconv.Itoa(n) + " bytes")
	return nil
}

func receiveGoldmanToken(tokenListener net.Listener) *protobuf.GoldmanToken {
	var conn net.Conn

	if tokenListener != nil && useTCP {
		var err error
		conn, err = tokenListener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		tokenServerUDPAddress, err := net.ResolveUDPAddr(tokenServer.UsedProtocol(), tokenServer.IpAndPortAsString())
		conn, err = net.ListenUDP(tokenServer.UsedProtocol(), tokenServerUDPAddress)
		if err != nil {
			log.Fatalln(err)
		}
	}

	defer conn.Close()
	data := make([]byte, 4096)
	n, err := conn.Read(data)
	utils.PrintMessage("Incoming message")
	if err != nil {
		log.Fatalln(err)
	}
	token := new(protobuf.GoldmanToken)
	err = proto.Unmarshal(data[0:n], token)
	if err != nil {
		log.Fatalln(err)
	}
	utils.PrintMessage(fmt.Sprintf("Token decoded.\n\n%s\n", token))
	return token
}

func handleTokenMessages() {
	var tokenListener net.Listener
	var err error

	tokenServer = server.New()
	tokenServer.SetClientName(string(processId))
	tokenServer.SetIpAddressAsString(ipAddress)
	tokenServer.SetPort(port + 1)
	if useTCP {
		tokenServer.SetUsedProtocol(client.TCP)
		tokenListener, err = net.Listen(tokenServer.UsedProtocol(), tokenServer.IpAndPortAsString())
		if err != nil {
			log.Fatalln(err)
		}
		defer tokenListener.Close()
	} else {
		tokenServer.SetUsedProtocol(client.UDP)
		tokenListener = nil
	}

	for {

		token := receiveGoldmanToken(tokenListener)
		if token.GetSourcePort()%2 == 0 {
			token.SourcePort = proto.Int32(token.GetSourcePort() + 1)
		}
		//Check if the array of blocking processes contains the id of this process
		if blocking {
			thisClientCausesDeadlock := false
			for _, value := range token.GetBlockingProcesses() {
				if value == int32(processId) {
					thisClientCausesDeadlock = true
				}
			}
			if thisClientCausesDeadlock {
				utils.PrintMessage("Deadlock, release resource!")
				if err := sendFilemanagerRequest(nonBlockingManager, RENOUNCE); err != nil {
					log.Fatalln(err)
				}
				time.Sleep(time.Duration(SECONDS_UNTIL_NEXT_TRY*100*rand.Float32()) * time.Millisecond)
			} else {
				time.Sleep(time.Duration(SECONDS_UNTIL_NEXT_TRY*100*rand.Float32()) * time.Millisecond)
				targetServerObject, err := parseIpColonPortToNetworkServer(fmt.Sprintf("%s:%d", token.GetSourceIP(), token.GetSourcePort()))
				if err != nil {
					log.Fatalln(err)
				}
				utils.PrintMessage("Send token to WAIT-FOR node.")
				if err := sendGoldmanToken(targetServerObject, token.GetBlockingProcesses()); err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
}
