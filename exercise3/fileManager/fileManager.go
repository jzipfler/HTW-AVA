package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

var (
	filename    string
	managerName string
	logFile     string
	ipAddress   string
	port        int
	managedFile *os.File
	force       bool
)

func init() {
	flag.StringVar(&filename, "filename", "path/to/file.txt", "A file that is managed by this process.")
	flag.StringVar(&managerName, "name", "Manager A", "Define the name of this manager.")
	flag.StringVar(&logFile, "logFile", "path/to/logfile.txt", "This parameter can be used to print the logging output to the given file.")
	flag.StringVar(&ipAddress, "ipAddress", "127.0.0.1", "The ip address of the actual starting node.")
	flag.IntVar(&port, "port", 15100, "The port of the actual starting node.")
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

	utils.InitializeLogger(logFile, "")
	utils.PrintMessage(fmt.Sprintf("File \"%s\" is now managed by this process.", filename))

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

	for i := 0; i <= 100; i++ {
		if numbers, err := utils.IncreaseNumbersFromFirstLine(filename, 6); err != nil {
			log.Fatalln(err.Error())
		} else {
			fmt.Println(numbers)
		}
	}

	for i := 0; i <= 101; i++ {
		if numbers, err := utils.DecreaseNumbersFromFirstLine(filename, 6); err != nil {
			log.Println(err.Error())
		} else {
			fmt.Println(numbers)
		}
	}

	if err := utils.AppendStringToFile(filename, "Hier könnte Ihre Werbung stehen", false); err != nil {
		log.Fatalln(err)
	}
	if err := utils.AppendStringToFile(filename, " ::::: Oder vieles mehr!", true); err != nil {
		log.Fatalln(err)
	}
	if err := utils.AppendStringToFile(filename, "Das stimmt!", true); err != nil {
		log.Fatalln(err)
	}
	os.Exit(0)

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
