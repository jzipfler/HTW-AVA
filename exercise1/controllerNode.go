package exercise1

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

// The controller is used to control the independent nodes.
// He can initialize or shutdown the nodes.
func StartController(localNode server.NetworkServer, allNodes map[int]server.NetworkServer, messageContent string) {
	if allNodes == nil {
		utils.PrintMessage(fmt.Sprintf("To start the controller, there must be a node map which is currently nil.\n%s\n", utils.ERROR_FOOTER))
		os.Exit(1)
	}
	quit := false
	utils.PrintMessage("Start current instance as controller.")

	for !quit {
		var input string
		utils.PrintMessage("Printing main menu.")
		printMainMenu(allNodes)
		fmt.Print("\nEnter ID of the node you would like to send a message.\nInput: ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			utils.PrintMessage(fmt.Sprintf("Error while reading the input. Quit program...\n%s\n", utils.ERROR_FOOTER))
			os.Exit(1)
		}
		targetId, err := strconv.Atoi(input)
		if err != nil {
			utils.PrintMessage(fmt.Sprintf("No number read. You have to enter the id of the node where you want to sent a message.\n%s\n", utils.ERROR_FOOTER))
			continue
		}
		if targetId == utils.CONTROLLER_MENU_NOTHING {
			quit = askForProgramRestart()
			continue
		}
		printControlMessageActionMenu()
		fmt.Print("\nEnter value of the control action you would like to sent.\nInput: ")
		_, err = fmt.Scanln(&input)
		if err != nil {
			utils.PrintMessage(fmt.Sprintf("Error while reading the input. Quit program...\n%s\n", utils.ERROR_FOOTER))
			os.Exit(1)
		}
		controlAction, err := strconv.Atoi(input)
		if controlAction == utils.CONTROLLER_MENU_NOTHING {
			quit = askForProgramRestart()
			continue
		}
		if err := SendProtobufControlMessage(localNode, allNodes[targetId], targetId, controlAction, messageContent); err != nil {
			utils.PrintMessage(fmt.Sprintf("The following error occured while trying to send a control message: %s\n%s\n", err.Error(), utils.ERROR_FOOTER))
		} else {
			utils.PrintMessage(fmt.Sprintf("Message to node with id %d, successfully sent.", targetId))
		}
		quit = askForProgramRestart()
	}
}

// Asks the user if he want to exit the program.
// Returns true if and only if the user types y or j. False otherwise.
func askForProgramRestart() bool {
	var input string
	utils.PrintMessage("Would you like to exit the program? (y/j/n)")
	fmt.Print("\nInput: ")
	if _, err := fmt.Scanln(&input); err == nil {
		switch input {
		case "y", "j":
			utils.PrintMessage("Program exists.")
			return true
		case "n":
			utils.PrintMessage(input)
			return false
		default:
			utils.PrintMessage("Please only insert y/j for \"YES\" or n for \"NO\".\n" + utils.ERROR_FOOTER)
			utils.PrintMessage("Assume a \"n\" as input.")
			return false
		}
	} else {
		utils.PrintMessage("Please only insert y/j for \"YES\" or n for \"NO\".\n" + utils.ERROR_HEADER)
	}
	return false
}

func printMainMenu(allNodes map[int]server.NetworkServer) {
	fmt.Println("")
	tabwriterObject := new(tabwriter.Writer)
	defer tabwriterObject.Flush()
	// Format in tab-separated columns with a tab stop of 4.
	tabwriterObject.Init(os.Stdout, 0, 4, 0, '\t', 0)
	fmt.Fprintln(tabwriterObject, "ID\tIP Address\tPort\tProtocol")
	fmt.Fprintln(tabwriterObject, utils.MENU_SEPERATOR)
	for key, value := range allNodes {
		fmt.Fprintf(tabwriterObject, "%d\t%s\t%d\t%s\n", key, value.IpAddressAsString(), value.Port(), value.UsedProtocol())
	}
	fmt.Fprintf(tabwriterObject, "\n%d\tAbort\n", utils.CONTROLLER_MENU_NOTHING)
}

func printControlMessageActionMenu() {
	fmt.Println("")
	tabwriterObject := new(tabwriter.Writer)
	defer tabwriterObject.Flush()
	// Format in tab-separated columns with a tab stop of 4.
	tabwriterObject.Init(os.Stdout, 0, 4, 0, '\t', 0)
	fmt.Fprintln(tabwriterObject, "Value\tControl message action")
	fmt.Fprintln(tabwriterObject, utils.MENU_SEPERATOR)
	fmt.Fprintf(tabwriterObject, "%d\tInitialize\n", utils.CONTROL_TYPE_INIT)
	fmt.Fprintf(tabwriterObject, "%d\tShutdown\n", utils.CONTROL_TYPE_EXIT)
	fmt.Fprintf(tabwriterObject, "\n%d\tAbort\n", utils.CONTROLLER_MENU_NOTHING)
}
