package filehandler

import (
	"bufio"
	"errors"
	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	ERROR_HEADER = "-------------ERROR-------------"
	ERROR_FOOTER = "^^^^^^^^^^^^^ERROR^^^^^^^^^^^^^"
)

// This function tries to read the content of the file with the given name
// as parameter. If it is not readable or nothing can be parsed, the function
// returns nil for the map and sets the value for the error type.
func CollectAllFromNodeListFile(nodeListFile string) (allNodes map[int]server.NetworkServer, err error) {
	if err := utils.CheckIfFileIsReadable(nodeListFile); err != nil {
		return nil, err
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
			log.Printf("Leere Zeile gelesen.\n")
			continue
		}
		if strings.HasPrefix(line, "#") {
			log.Printf("Kommentar gelesen: \"%s\"\n", line)
			continue
		}
		idAndIpPortArray := strings.Split(line, " ")
		if len(idAndIpPortArray) == 0 {
			log.Printf("Could not split the line with a space: \"%s\".\n%s\n", line, ERROR_FOOTER)
			continue
		}
		scanId, err := strconv.Atoi(idAndIpPortArray[0])
		if err != nil {
			log.Printf("Could not parse the first part of the line to a number : \"%s\".\n%s\n", idAndIpPortArray[0], ERROR_FOOTER)
			continue
		} else {
			scanServerObject.SetClientName(idAndIpPortArray[0])
		}
		ipAndPortArray := strings.Split(idAndIpPortArray[1], ":")
		if len(ipAndPortArray) == 0 {
			log.Printf("Could not split the ip address and port with a colon: \"%s\".\n%s\n", idAndIpPortArray[1], ERROR_FOOTER)
			continue
		}
		if splitIpArray, err := net.LookupIP(ipAndPortArray[0]); err != nil {
			if splitIpHostnameArray, err := net.LookupHost(ipAndPortArray[0]); err != nil {
				log.Printf("Could not lookup this ip/host: \"%s\".\n%s\n", ipAndPortArray[0], ERROR_FOOTER)
				continue
			} else {
				scanServerObject.SetIpAddressAsString(splitIpHostnameArray[0])
			}
		} else {
			if len(splitIpArray) == 0 {
				log.Printf("No ip found: \"%s\".\n%s\n", ipAndPortArray[0], ERROR_FOOTER)
				continue
			}
			scanServerObject.SetIpAddress(splitIpArray[0])
		}
		if splitPort, err := strconv.Atoi(ipAndPortArray[1]); err != nil {
			log.Printf("Could not parse the port: \"%s\".\n%s\n", ipAndPortArray[1], ERROR_FOOTER)
			continue
		} else {
			scanServerObject.SetPort(splitPort)
		}
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
