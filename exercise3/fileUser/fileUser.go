package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	serverObject := server.New()

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
		work = func() {
			//Get write access on A then on B
			//Increase A and decrease B
			go sendFilemanagerRequest(server.New())
			filenameA, err := receiveAndParseFilemanagerResponse()
			if err != nil {
				utils.PrintMessage(err)
			}
			go sendFilemanagerRequest(server.New())
			filenameB, err := receiveAndParseFilemanagerResponse()
			if err != nil {
				utils.PrintMessage(err)
			}
			utils.IncreaseNumbersFromFirstLine(filenameA, 6)
			utils.AppendStringToFile(filenameA, string(processId), true)
			utils.DecreaseNumbersFromFirstLine(filenameB, 6)
			utils.AppendStringToFile(filenameB, string(processId), true)
		}
	} else {
		work = func() {
			//Get write access on B then on A
			//Increase B and decrease A
			go sendFilemanagerRequest(server.New())
			filenameB, err := receiveAndParseFilemanagerResponse()
			if err != nil {
				utils.PrintMessage(err)
			}
			go sendFilemanagerRequest(server.New())
			filenameA, err := receiveAndParseFilemanagerResponse()
			if err != nil {
				utils.PrintMessage(err)
			}
			utils.IncreaseNumbersFromFirstLine(filenameB, 6)
			utils.AppendStringToFile(filenameB, string(processId), true)
			utils.DecreaseNumbersFromFirstLine(filenameA, 6)
			utils.AppendStringToFile(filenameA, string(processId), true)
		}
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

func handleReceivedProtobufMessageWithChannel(localNode server.NetworkServer, receivingChannel chan *protobuf.FilemanagerRequest) {
	for {
		// This call blocks until a new message is available.
		message := <-receivingChannel
		utils.PrintMessage(fmt.Sprintf("Message on %s received:\n\n%s\n\n", localNode.String(), message.String()))
	}
}

func receiveAndParseFilemanagerResponse() (string, error) {
	conn, err := server.ReceiveMessage()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return "test.txt", nil
}

func sendFilemanagerRequest(fileManager server.NetworkServer) {

}

func appendProcessId(filename string) {

}
