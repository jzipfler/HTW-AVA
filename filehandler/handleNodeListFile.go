package filehandler

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

// This function tries to read the content of the file with the given name
// as parameter. If it is not readable or nothing can be parsed, the function
// returns nil for the map and sets the value for the error type.
func CollectAllFromNodeListFile(nodeListFile string) (allNodes map[int]server.NetworkServer, err error) {
	if readable, _ := utils.CheckIfFileIsReadable(nodeListFile); !readable {
		return nil, errors.New("The given file is not readable.")
	}
	nodeListFileObject, _ := os.Open(nodeListFile)
	defer nodeListFileObject.Close()
	allNodes = make(map[int]server.NetworkServer, 10)
	scanner := bufio.NewScanner(nodeListFileObject)
	for scanner.Scan() {
		var scanId int
		var scanServerObject server.NetworkServer
		line := scanner.Text()
		if line == "" {
			//log.Printf("Leere Zeile gelesen.\n")
			continue
		}
		if strings.HasPrefix(line, "#") {
			//log.Printf("Kommentar gelesen: \"%s\"\n", line)
			continue
		}
		idAndIpPortArray := strings.Split(line, " ")
		if len(idAndIpPortArray) == 0 {
			log.Printf("Could not split the line with a space: \"%s\".\n%s\n", line, utils.ERROR_FOOTER)
			continue
		}
		scanId, err := strconv.Atoi(idAndIpPortArray[0])
		if err != nil {
			log.Printf("Could not parse the first part of the line to a number : \"%s\".\n%s\n", idAndIpPortArray[0], utils.ERROR_FOOTER)
			continue
		} else {
			scanServerObject.SetClientName(idAndIpPortArray[0])
		}
		ipAndPortArray := strings.Split(idAndIpPortArray[1], ":")
		if len(ipAndPortArray) == 0 {
			log.Printf("Could not split the ip address and port with a colon: \"%s\".\n%s\n", idAndIpPortArray[1], utils.ERROR_FOOTER)
			continue
		}
		//Check if the given part is a ip address or a host.
		if splitIpArray, err := net.LookupIP(ipAndPortArray[0]); err != nil {
			log.Printf("Could not lookup this ip/host: \"%s\".\n%s\n", ipAndPortArray[0], utils.ERROR_FOOTER)
			continue
		} else {
			if len(splitIpArray) == 0 {
				log.Printf("No ip found: \"%s\".\n%s\n", ipAndPortArray[0], utils.ERROR_FOOTER)
				continue
			}
			//If we have some ip addresses, lets check if they are from version 4.
			ipv4Found := false
			for _, value := range splitIpArray {
				if value.To4() != nil {
					ipv4Found = true
					scanServerObject.SetIpAddress(value)
					break
				}
			}
			if !ipv4Found {
				log.Printf("No ipv4 found: \"%s\".\n%s\n", ipAndPortArray[0], utils.ERROR_FOOTER)
				continue
			}
		}
		if splitPort, err := strconv.Atoi(ipAndPortArray[1]); err != nil {
			log.Printf("Could not parse the port: \"%s\".\n%s\n", ipAndPortArray[1], utils.ERROR_FOOTER)
			continue
		} else {
			scanServerObject.SetPort(splitPort)
		}
		//TODO: Maybe the usedProtocol should be set in a other way.
		scanServerObject.SetUsedProtocol("tcp")
		allNodes[scanId] = scanServerObject
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(allNodes) == 0 {
		return nil, errors.New("No nodes present... ABORT")
	}
	return allNodes, nil
}
