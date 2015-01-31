package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/jzipfler/htw-ava/protobuf"
	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

var (
	filename        string
	managerName     string
	logFile         string
	ipAddress       string
	port            int
	managedFile     *os.File
	force           bool
	fileInUse       bool
	usedById        int32
	usedByIpAndPort string

	serverObject server.NetworkServer

	releaseServer server.NetworkServer
)

const (
	ACCESS_GRANTED        = iota // 0
	ACCESS_DENIED                // 1
	RESOURCE_RELEASED            // 2
	RESOURCE_NOT_RELEASED        // 3
)

func init() {
	flag.StringVar(&filename, "filename", "path/to/file.txt", "A file that is managed by this process.")
	flag.StringVar(&managerName, "name", "Manager A", "Define the name of this manager.")
	flag.StringVar(&logFile, "logFile", "path/to/logfile.txt", "This parameter can be used to print the logging output to the given file.")
	flag.StringVar(&ipAddress, "ipAddress", "127.0.0.1", "The ip address of the actual starting node.")
	flag.IntVar(&port, "port", 15100, "The port of the actual starting node. (Portnumber must be even)")
	flag.BoolVar(&force, "force", false, "If force is enabled, the programm removes a existing management file and creates a new one without asking.")
}

func main() {

	var containsAddress, containsPort, containsFilename bool
	for _, argument := range os.Args {
		if strings.Contains(argument, "-ipAddress") {
			containsAddress = true
		}
		if strings.Contains(argument, "-port") {
			containsPort = true
		}
		if strings.Contains(argument, "-filename") {
			containsFilename = true
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
	if !containsFilename {
		log.Printf("A filename is required.\n%s\n\n", utils.ERROR_FOOTER)
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	if port%2 != 0 {
		log.Printf("The port number must be even.\n%s\n\n", utils.ERROR_FOOTER)
		os.Exit(1)
	}

	utils.InitializeLogger(logFile, "")
	utils.PrintMessage(fmt.Sprintf("File \"%s\" is now managed by this process.", filename))

	var err error
	filename, err = filepath.Abs(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if exists := utils.CheckIfFileExists(filename); exists {
		if !force {
			if deleteIt := askForToDeleteFile(); !deleteIt {
				fmt.Println("Do not delete the file and exit the program.")
				utils.PrintMessage(fmt.Sprintf("The file \"%s\" already exists and should not be deleted.", filename))
				os.Exit(0)
			}
		}
		if err := os.Remove(filename); err != nil {
			log.Fatalf("%s\n%s\n", err.Error(), utils.ERROR_FOOTER)
		}
		utils.PrintMessage(fmt.Sprintf("Removed the file \"%s\"", filename))
	}

	managedFile, err := os.Create(filename)
	utils.PrintMessage(fmt.Sprintf("Created the file \"%s\"", filename))
	if err != nil {
		log.Fatalf("%s\n%s\n", err.Error(), utils.ERROR_FOOTER)
	}

	managedFile.WriteString("000000\n")
	utils.PrintMessage("Wrote 000000 to the file.")
	managedFile.Close()

	serverObject = server.New()
	serverObject.SetClientName(managerName)
	serverObject.SetIpAddressAsString(ipAddress)
	serverObject.SetPort(port)
	serverObject.SetUsedProtocol("tcp")

	if err := server.StartServer(serverObject, nil); err != nil {
		log.Fatalln("Could not start server. --> Exit.")
		os.Exit(1)
	}
	defer server.StopServer()

	go handleMessagesOnUnevenPort()

	for {
		receivedMessage := receiveAndParseFilemanagerRequest()
		var reaction int
		switch receivedMessage.GetAccessOperation() {
		case protobuf.FilemanagerRequest_GET:
			if fileInUse {
				if usedById == receivedMessage.GetSourceID() {
					utils.PrintMessage(fmt.Sprintf("Somebody asks about permission where he already have access to. --> Access granted."))
					reaction = ACCESS_GRANTED
				} else {
					utils.PrintMessage(fmt.Sprintf("Access denied, file in use by %d-%s", usedById, usedByIpAndPort))
					reaction = ACCESS_DENIED
				}
			} else {
				fileInUse = true
				usedByIpAndPort = fmt.Sprintf("%s:%d", receivedMessage.GetSourceIP(), int(receivedMessage.GetSourcePort()))
				usedById = receivedMessage.GetSourceID()
				utils.PrintMessage(fmt.Sprintf("Access granted, file is now in use by %d-%s", usedById, usedByIpAndPort))
				reaction = ACCESS_GRANTED
			}
		case protobuf.FilemanagerRequest_RELEASE, protobuf.FilemanagerRequest_RENOUNCE:
			utils.PrintMessage(fmt.Sprintf("Release/Renounce requested, file is in use by %d-%s", usedById, usedByIpAndPort))
			if fileInUse && usedById == receivedMessage.GetSourceID() {
				fileInUse = false
				usedById = 0
				usedByIpAndPort = ""
				reaction = RESOURCE_RELEASED
				utils.PrintMessage("File successfully released/renounced.")
			} else {
				reaction = RESOURCE_NOT_RELEASED
				utils.PrintMessage(fmt.Sprintf("Resource not released. ID %d != %d!", usedById, receivedMessage.GetSourceID()))
			}
		}
		if err := sendFilemanagerResponse(receivedMessage.GetSourceIP(), int(receivedMessage.GetSourcePort()), reaction); err != nil {
			log.Fatalln(err)
		}
	}
}

func askForToDeleteFile() bool {
	var input string
	fmt.Printf("Would you like to delete the file \"%s\"? (y/j/n)", filename)
	fmt.Print("\nInput: ")
	if _, err := fmt.Scanln(&input); err == nil {
		switch input {
		case "y", "j":
			fmt.Println("File gets deleted.")
			return true
		case "n":
			fmt.Println(input)
			return false
		default:
			fmt.Println("Please only insert y/j for \"YES\" or n for \"NO\".\n" + utils.ERROR_FOOTER)
			fmt.Println("Assume a \"n\" as input.")
			return false
		}
	} else {
		fmt.Println("Please only insert y/j for \"YES\" or n for \"NO\".\n" + utils.ERROR_HEADER)
	}
	return false
}

func receiveAndParseFilemanagerRequest() *protobuf.FilemanagerRequest {
	//ReceiveMessage blocks until a message comes in
	conn, err := server.ReceiveMessage()
	if err != nil {
		utils.PrintMessage(err)
	}
	utils.PrintMessage("Incoming message")
	//Close the connection when the function exits
	defer conn.Close()
	//Create a data buffer of type byte slice with capacity of 4096
	data := make([]byte, 4096)
	//Read the data waiting on the connection and put it in the data buffer
	n, err := conn.Read(data)
	if err != nil {
		log.Fatal("Error happened: " + err.Error())
	}
	utils.PrintMessage("Decoding Protobuf message")
	//Create an struct pointer of type ProtobufTest.TestMessage struct
	protodata := new(protobuf.FilemanagerRequest)
	//Convert all the data retrieved into the ProtobufTest.TestMessage struct type
	err = proto.Unmarshal(data[0:n], protodata)
	if err != nil {
		log.Fatal("Error happened: " + err.Error())
	}
	utils.PrintMessage(fmt.Sprintf("FilemanagerRequest decoded.\n\n%s\n\n", protodata))
	return protodata
}

func sendFilemanagerResponse(destinationIp string, destinationPort, reaction int) error {
	if destinationIp == "" || port == 0 {
		return errors.New(fmt.Sprintf("The target server information has no ip address or port.\n%s:%d\n", destinationIp, destinationPort, utils.ERROR_FOOTER))
	}
	ipAddressAndPort := destinationIp + ":" + strconv.Itoa(destinationPort)
	utils.PrintMessage(fmt.Sprintf("Encode protobuf application message for node with IP:PORT : %s.", ipAddressAndPort))
	protobufMessage := new(protobuf.FilemanagerResponse)
	protobufMessage.SourceIP = proto.String(serverObject.IpAddressAsString())
	protobufMessage.SourcePort = proto.Int(serverObject.Port())
	var requestReaction protobuf.FilemanagerResponse_RequestReaction
	switch reaction {
	case ACCESS_GRANTED:
		requestReaction = protobuf.FilemanagerResponse_RequestReaction(ACCESS_GRANTED)
		protobufMessage.Filename = proto.String(filename)
	case ACCESS_DENIED:
		requestReaction = protobuf.FilemanagerResponse_RequestReaction(ACCESS_DENIED)
	case RESOURCE_RELEASED:
		requestReaction = protobuf.FilemanagerResponse_RequestReaction(RESOURCE_RELEASED)
	case RESOURCE_NOT_RELEASED:
		requestReaction = protobuf.FilemanagerResponse_RequestReaction(RESOURCE_NOT_RELEASED)
	default:
		utils.PrintMessage("No valid reaction given. Assume DENIE.")
		requestReaction = protobuf.FilemanagerResponse_RequestReaction(ACCESS_DENIED)
	}
	protobufMessage.RequestReaction = &requestReaction
	if usedByIpAndPort != "" && usedById != 0 {
		protobufMessage.ProcessIpAndPortThatUsesResource = proto.String(usedByIpAndPort)
		protobufMessage.ProcessIdThatUsesResource = proto.Int32(usedById)
	}
	//Protobuf message filled with data. Now marshal it.
	data, err := proto.Marshal(protobufMessage)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", ipAddressAndPort)
	if err != nil {
		return err
	}
	defer conn.Close()
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	utils.PrintMessage(fmt.Sprintf("Application message from %s to %s sent:\n\n%s\n\n", serverObject.String(), ipAddressAndPort, protobufMessage.String()))
	utils.PrintMessage("Sent " + strconv.Itoa(n) + " bytes")
	return nil
}

func handleMessagesOnUnevenPort() {
	releaseServer = server.New()
	releaseServer.SetClientName(managerName)
	releaseServer.SetIpAddressAsString(ipAddress)
	releaseServer.SetPort(port + 1)
	releaseServer.SetUsedProtocol("tcp")

	releaseListener, err := net.Listen(releaseServer.UsedProtocol(), releaseServer.IpAndPortAsString())
	if err != nil {
		log.Fatalln(err)
	}
	defer releaseListener.Close()

	for {
		var reaction int
		conn, err := releaseListener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		utils.PrintMessage("Incoming message (Used to be a DEADLOCK-RENOUNCE message)")
		defer conn.Close()
		data := make([]byte, 4096)
		n, err := conn.Read(data)
		if err != nil {
			log.Fatalln(err)
		}
		request := new(protobuf.FilemanagerRequest)
		err = proto.Unmarshal(data[0:n], request)
		if err != nil {
			log.Fatalln(err)
		}
		utils.PrintMessage(fmt.Sprintf("Request decoded.\n\n%s\n\n", request))
		if request.GetAccessOperation() == protobuf.FilemanagerRequest_GET {
			sendFilemanagerResponse(request.GetSourceIP(), int(request.GetSourcePort()-1), RESOURCE_NOT_RELEASED)
		}
		//Check if the array of blocking processes contains the id of this process
		utils.PrintMessage("Message is a RELEASE or RENOUNCE message, check if this process id blocks the file.")
		if fileInUse && usedById == request.GetSourceID() {
			utils.PrintMessage("YES, this process id blocks the file, RELEASE it!")
			fileInUse = false
			usedById = 0
			usedByIpAndPort = ""
			reaction = RESOURCE_RELEASED
			utils.PrintMessage("File successfully released/renounced.")
		} else {
			utils.PrintMessage("NO, this process id does not block the file, DO NOT RELEASE it!")
			reaction = RESOURCE_NOT_RELEASED
		}
		sendFilemanagerResponse(request.GetSourceIP(), int(request.GetSourcePort()-1), reaction)
	}
}
