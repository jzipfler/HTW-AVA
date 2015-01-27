package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

var (
	filename    string
	managerName string
	logFile     string
	ipAddress   string
	port        int
)

func init() {
	flag.StringVar(&filename, "filename", "path/to/file.txt", "A file that is managed by this process.")
	flag.StringVar(&managerName, "name", "Manager A", "Define the name of this manager.")
	flag.StringVar(&logFile, "logFile", "path/to/logfile.txt", "This parameter can be used to print the logging output to the given file.")
	flag.StringVar(&ipAddress, "ipAddress", "127.0.0.1", "The ip address of the actual starting node.")
	flag.IntVar(&port, "port", 15100, "The port of the actual starting node.")
}

func main() {

	if filename == "path/to/file.txt" {
		log.Printf("A filename is required.\n%s\n\n", utils.ERROR_FOOTER)
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	if exists := utils.CheckIfFileExists(filename); exists {
		if writable, err := utils.CheckIfFileIsReadableAndWritebale(filename); !writable {
			log.Fatalf("%s\n%s\n", err.Error(), utils.ERROR_FOOTER)
		}
	} else {
		os.Create(filename)
	}

	utils.InitializeLogger(logFile, "")
	utils.PrintMessage(fmt.Sprintf("File \"%s\" is now managed by this process.", filename))

	serverObject := server.New()

	serverObject.SetClientName(managerName)
	serverObject.SetIpAddressAsString(ipAddress)
	serverObject.SetPort(port)
	serverObject.SetUsedProtocol("tcp")

	if err := server.StartServer(serverObject, nil); err != nil {
		log.Fatalln("Could not start server. --> Exit.")
		os.Exit(1)
	}
	defer server.StopServer()
}
