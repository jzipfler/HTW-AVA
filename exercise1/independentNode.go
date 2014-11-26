package exercise1

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jzipfler/htw-ava/protobuf"
	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

const (
	// The number of nodes from which the rumor must be heared before the node belives in it.
	BELIVE_IN_RUMOR_THRESHOLD = 2
)

var (
	localNode                 server.NetworkServer
	allNodes                  map[int]server.NetworkServer
	neighbors                 map[int]server.NetworkServer
	messageToAllNeighborsSent bool
	rumorExperiment           bool
	rumors                    []int
)

// With this function an node that interacts independently gets started.
// He can be controlled with a controller.
func StartIndependentNode(localNodeId int, allAvailableNodes, neighborNodes map[int]server.NetworkServer, rumorExperimentMode bool) {
	if allAvailableNodes == nil {
		utils.PrintMessage(fmt.Sprintf("To start the controller, there must be a node map which is currently nil.\n%s\n", utils.ERROR_FOOTER))
		os.Exit(1)
	}
	if _, ok := allAvailableNodes[localNodeId]; !ok {
		utils.PrintMessage(fmt.Sprintf("The given id exists not in the node map.\n%s\n", utils.ERROR_FOOTER))
		os.Exit(1)
	}
	if neighborNodes == nil {
		utils.PrintMessage(fmt.Sprintf("No neighbors given. Use the ChooseThreeNeighbors function to get some.\n%s\n", utils.ERROR_FOOTER))
		neighborNodes = ChooseThreeNeighbors(localNodeId, allAvailableNodes)
	}
	utils.PrintMessage("Start current instance as independent node.")
	allNodes = allAvailableNodes
	neighbors = neighborNodes
	localNode = allAvailableNodes[localNodeId]
	rumors = make([]int, len(allNodes), len(allNodes))
	messageToAllNeighborsSent = false
	rumorExperiment = rumorExperimentMode
	utils.PrintMessage("This node has the folowing settings: ")
	utils.PrintMessage(localNode)

	protobufChannel := make(chan *protobuf.Nachricht)
	//A goroutine that receives the protobuf message and reacts to it.
	go handleReceivedProtobufMessageWithChannel(localNode, protobufChannel)
	if err := server.StartServer(localNode, nil); err != nil {
		log.Fatal("Error happened: " + err.Error())
	}
	defer server.StopServer()

	for {
		//ReceiveMessage blocks until a message comes in
		if conn, err := server.ReceiveMessage(); err == nil {
			//If err is nil then that means that data is available for us so we take up this data and pass it to a new goroutine
			go ReceiveAndParseIncomingProtobufMessageToChannel(conn, protobufChannel)
			//ReceiveAndParseIncomingProtobufMessageToChannel(conn, protobufChannel)
			//protodata := ReceiveAndParseInfomingProtoufMessage(conn)
			//utils.PrintMessage(fmt.Sprintf("Message on %s received:\n\n%s\n\n", localNode.String(), protodata.String()))
			//handleReceivedProtobufMessage(protodata)
		}
	}
}

// The chooseThreeNeighbors function uses the allAvailableNodes map to return
// another map that contains 3 nodes at the most.
// It calls os.Exit(1) if only one node is available in the allAvailableNodes map.
func ChooseThreeNeighbors(localNodeId int, allAvailableNodes map[int]server.NetworkServer) (neighbors map[int]server.NetworkServer) {
	neighbors = make(map[int]server.NetworkServer, 3)
	// If there are only 1, 2 or 3 possible neighbors...take them.
	switch len(allAvailableNodes) {
	case 1:
		utils.PrintMessage(fmt.Sprintf("There is only one node in the nodeList. Ther must be at least 2.\n%s\n", utils.ERROR_FOOTER))
		os.Exit(1)
	case 2, 3, 4:
		for key, value := range allAvailableNodes {
			if key != localNodeId {
				neighbors[key] = value
			}
		}
		return
	}
	// Because of
	var highestId int
	for key := range allAvailableNodes {
		if highestId < key {
			highestId = key
		}
	}
	randomObject := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(neighbors) != 3 {
		var randomNumber int
		randomNumber = randomObject.Intn(highestId + 1)
		if randomNumber == localNodeId {
			continue
		}
		// Add only the nodes with the id which exists.
		if value, ok := allAvailableNodes[randomNumber]; ok {
			// And check here if the node already exists in the neighbors map.
			if _, ok := neighbors[randomNumber]; !ok {
				neighbors[randomNumber] = value
				// Now remove the added node from the map.
				delete(allAvailableNodes, randomNumber)
			}
		}
	}
	return
}

// This function waits for a message that is sent to the channel and
// splits the handling of the message depending on the NachrichtenTyp (message type)
func handleReceivedProtobufMessageWithChannel(localNode server.NetworkServer, receivingChannel chan *protobuf.Nachricht) {
	for {
		// This call blocks until a new message is available.
		message := <-receivingChannel
		utils.PrintMessage(fmt.Sprintf("Message on %s received:\n\n%s\n\n", localNode.String(), message.String()))
		switch message.GetNachrichtenTyp() {
		case protobuf.Nachricht_KONTROLLNACHRICHT:
			utils.PrintMessage("Message is of type KONTROLLNACHRICHT.")
			handleReceivedControlMessage(message)
		case protobuf.Nachricht_ANWENDUNGSNACHRICHT:
			utils.PrintMessage("Message is of type ANWENDUNGSNACHRICHT.")
			handleReceivedApplicationMessage(message)
		default:
			log.Fatalln("Read a unknown \"NachrichtenTyp\"")
		}
	}
}

// This method gets a protobuf message and decides if it is a control or a
// application message and gives it to the related function.
func handleReceivedProtobufMessage(protoMessage *protobuf.Nachricht) {
	switch protoMessage.GetNachrichtenTyp() {
	case protobuf.Nachricht_KONTROLLNACHRICHT:
		utils.PrintMessage("Message is of type KONTROLLNACHRICHT.")
		handleReceivedControlMessage(protoMessage)
	case protobuf.Nachricht_ANWENDUNGSNACHRICHT:
		utils.PrintMessage("Message is of type ANWENDUNGSNACHRICHT.")
		handleReceivedApplicationMessage(protoMessage)
	default:
		log.Fatalln("Read a unknown \"NachrichtenTyp\"")
	}
}

func handleReceivedControlMessage(message *protobuf.Nachricht) {
	switch message.GetKontrollTyp() {
	case protobuf.Nachricht_INITIALISIEREN:
		if !messageToAllNeighborsSent {
			for key, value := range neighbors {
				SendProtobufApplicationMessage(localNode, value, key, message.GetNachrichtenInhalt())
			}
			messageToAllNeighborsSent = true
		}
	case protobuf.Nachricht_BEENDEN:
		for id, destinationNode := range neighbors {
			SendProtobufControlMessage(localNode, destinationNode, id, utils.CONTROL_TYPE_EXIT, message.GetNachrichtenInhalt())
		}
		utils.PrintMessage("Received a EXIT message, so program will be exited.")
		os.Exit(0)
	default:
		log.Fatalln("Read a unknown \"KontrollTyp\"")
	}
}

func handleReceivedApplicationMessage(message *protobuf.Nachricht) {
	/*
	 *	Check if the last part of exercise one or if the first part should be
	 *	applied. The first part simply sends the own ID to all neighbors once.
	 * 	The last part is where a rumor should be be telled to a defined number
	 * 	of neighbors. For example (d-2) neighbors, where d is the degree of
	 *	the node.
	 */
	if rumorExperiment {
		// TODO: Place for the last part of the exercise (RUMORS)
		rumors[int(message.GetSourceID())-1]++
		utils.PrintMessage(fmt.Sprintln("Current rumors counted: ", rumors[int(message.GetSourceID())-1]))
		if rumors[int(message.GetSourceID())-1] == BELIVE_IN_RUMOR_THRESHOLD {
			filename := localNode.ClientName() + "_belives.txt"
			if exists := utils.CheckIfFileExists(filename); !exists {
				stringBuffer := bytes.NewBufferString(message.GetNachrichtenInhalt())
				ioutil.WriteFile(filename, stringBuffer.Bytes(), 0644)
			}
		}
		/*
		 *	The three variations to
		 *	0) allen Nachbarn
		 *	i) 2 Nachbarn
		 *	ii) d-2 Nachbarn, wobei d der Grad des Knotens ist
		 *	iii) (d-1)/2 Nachbarn, wobei d der Grad des Knotens ist
		 */
		sentToNeighborsThreshold := len(neighbors)
		for key, value := range neighbors {
			alreadySent := 0
			// Send it only that often as defined in the comment above.
			if alreadySent >= sentToNeighborsThreshold {
				break
			}
			SendProtobufApplicationMessage(localNode, value, key, message.GetNachrichtenInhalt())
			alreadySent++
		}
	} else {
		if !messageToAllNeighborsSent {
			for key, value := range neighbors {
				SendProtobufApplicationMessage(localNode, value, key, message.GetNachrichtenInhalt())
			}
			messageToAllNeighborsSent = true
		}
	}
	// Because the SourceID is of type int32, I have to cast it here.
	sourceId := int(message.GetSourceID())
	// Check if the node that sends the message is in the neighbors map.
	// If not, add him.
	// Optional: Send him a response that he is added as neighbor.
	if _, ok := neighbors[sourceId]; !ok {
		networkServerObject := server.New()
		networkServerObject.SetClientName(strconv.Itoa(sourceId))
		networkServerObject.SetIpAddressAsString(message.GetSourceIP())
		networkServerObject.SetPort(sourceId)
		networkServerObject.SetUsedProtocol("tcp") //TODO: Maybe a different approach...
		neighbors[int(message.GetSourceID())] = networkServerObject
		//sendProtobufApplicationMessage(sourceId)
	}
}
