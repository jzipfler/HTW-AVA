package exercise1

import (
	"fmt"
	"github.com/jzipfler/htw-ava/protobuf"
	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	localNode                 server.NetworkServer
	allNodes                  map[int]server.NetworkServer
	neighbors                 map[int]server.NetworkServer
	messageToAllNeighborsSent bool
)

// With this function an node that interacts independently gets started.
// He can be controlled with a controller.
func StartIndependentNode(localNodeId int, allAvailableNodes, neighborNodes map[int]server.NetworkServer) {
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
	utils.PrintMessage("This node has the folowing settings: ")
	utils.PrintMessage(localNode)

	protobufChannel := make(chan *protobuf.Nachricht)
	//A goroutine that receives the protobuf message and reacts to it.
	go handleReceivedProtobufMessage(localNode,protobufChannel)
	//go func() {
	//	for {
	//		message := <-protobufChannel
	//		utils.PrintMessage(fmt.Sprintf("Message received:\n\n%s\n\n", message.String()))
	//		utils.PrintMessage(message)
	//	}
	//}()
	//Listen to the TCP port
	if err := server.StartServer(localNode, nil); err != nil {
		log.Fatal("Error happened: " + err.Error())
	}
	defer server.StopServer()

	for {
		//ReceiveMessage blocks until a message comes in
		if conn, err := server.ReceiveMessage(); err == nil {
			//If err is nil then that means that data is available for us so we take up this data and pass it to a new goroutine
			go ReceiveAndParseIncomingProtobufMessageToChannel(conn, protobufChannel)
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
func handleReceivedProtobufMessage(localNode server.NetworkServer, receivingChannel chan *protobuf.Nachricht) {
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

func handleReceivedControlMessage(message *protobuf.Nachricht) {
	switch message.GetKontrollTyp() {
	case protobuf.Nachricht_INITIALISIEREN:
		if !messageToAllNeighborsSent {
			for key, value := range neighbors {
				SendProtobufApplicationMessage(localNode, value, key)
			}
		}
	case protobuf.Nachricht_BEENDEN:
		for key, value := range neighbors {
			SendProtobufControlMessage(localNode, value, key, utils.CONTROL_TYPE_EXIT)
		}
		utils.PrintMessage("")
		os.Exit(0)
	default:
		log.Fatalln("Read a unknown \"KontrollTyp\"")
	}
}

func handleReceivedApplicationMessage(message *protobuf.Nachricht) {
	if !messageToAllNeighborsSent {
		for key, value := range neighbors {
			SendProtobufApplicationMessage(localNode, value, key)
		}
	}

	// Because the SourceID is of type int32, I have to cast it here.
	sourceId := int(message.GetSourceID())
	// Check if the node that sends the message is in the neighbors map.
	// If not, add him and send him a response.
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
