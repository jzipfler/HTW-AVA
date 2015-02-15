package filehandler

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

func CollectNeighborsFromGraphvizFile(filename string, ownId int, allNodes map[int]server.NetworkServer) (neighbors map[int]server.NetworkServer, err error) {
	if ownId < 0 {
		return nil, errors.New("The id must not be negative.")
	}
	if allNodes == nil {
		return nil, errors.New("The allNodes parameter must not be nil")
	}
	if len(allNodes) <= 1 {
		return nil, errors.New("There are not enough nodes in the allNodes map.")
	}
	if readable, _ := utils.CheckIfFileIsReadable(filename); !readable {
		return nil, errors.New("The given file is not readable.")
	}
	graphvizFileObject, _ := os.Open(filename)
	defer graphvizFileObject.Close()
	scanner := bufio.NewScanner(graphvizFileObject)
	neighbors = make(map[int]server.NetworkServer)
	for scanner.Scan() {
		var firstScannedId, secondScannedId int
		line := scanner.Text()
		if line == "" {
			//log.Printf("Leere Zeile gelesen.\n")
			continue
		}
		if strings.HasPrefix(line, "//") {
			//log.Printf("Kommentar gelesen: \"%s\"\n", line)
			continue
		}
		if strings.Contains(line, "graph") {
			//log.Printf("Graph definition read: \"%s\"\n", line)
			continue
		}
		if strings.Contains(line, "{") || strings.Contains(line, "}") {
			//log.Printf("Line contained brackets. Skip.")
			continue
		}
		splittedNodeIdArray := strings.Split(line, "--")
		//log.Print("Here the read line: ")
		//log.Println(splittedNodeIdArray)
		var err error
		if firstScannedId, err = strconv.Atoi(strings.Trim(splittedNodeIdArray[0], "\t; ")); err != nil {
			log.Printf("Could not one of the parts to a number: \"%s\"\n%s", splittedNodeIdArray, utils.ERROR_FOOTER)
			continue
		}
		if secondScannedId, err = strconv.Atoi(strings.Trim(splittedNodeIdArray[1], "\t; ")); err != nil {
			log.Printf("Could not one of the parts to a number: \"%s\"\n%s", splittedNodeIdArray, utils.ERROR_FOOTER)
			continue
		}
		if firstScannedId == ownId {
			neighbors[secondScannedId] = allNodes[secondScannedId]
		}
		if secondScannedId == ownId {
			neighbors[firstScannedId] = allNodes[firstScannedId]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(neighbors) == 0 {
		return nil, errors.New("No neighbors found.")
	}
	return
}
