package exercise2

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jzipfler/htw-ava/protobuf"
	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

var (
	localCompanyNode          CompanyNode
	localCustomerNode         CustomerNode
	localNode                 server.NetworkServer
	customerNodeMap           map[int]bool
	allNodes                  map[int]server.NetworkServer
	neighbors                 map[int]server.NetworkServer
	messageToAllNeighborsSent bool
	localeId                  int
)

// With this function an node that interacts independently gets started.
// He can be controlled with a controller.
func StartIndependentNode(localNodeId int, customerNode bool, allAvailableNodes, neighborNodes map[int]server.NetworkServer) {
	if allAvailableNodes == nil {
		utils.PrintMessage(fmt.Sprintf("To start the node, there must be a node map which is currently nil.\n%s\n", utils.ERROR_FOOTER))
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
	localeId = localNodeId
	allNodes = allAvailableNodes
	neighbors = neighborNodes
	localNode = allAvailableNodes[localeId]
	customerNodeMap = make(map[int]bool, len(allNodes))
	if customerNode {
		customerNodeMap[localeId] = true
		localCustomerNode = NewCustomerNodeWithServerObject(localNode)
		localCustomerNode.SetCustomerId(localeId)
		localCustomerNode.SetFriends(neighborNodes)
		utils.PrintMessage("This node has the folowing settings: ")
		utils.PrintMessage(localCustomerNode)

		//Set the companyNode's id to -1, because it is easier to identify the "unset" node.
		localCompanyNode.SetCompanyId(-1)
	} else {
		customerNodeMap[localeId] = false
		localCompanyNode = NewCompanyNodeWithServerObject(localNode)
		localCompanyNode.SetCompanyId(localeId)
		localCompanyNode.SetProduct(fmt.Sprintf("Product%d", localeId))
		localCompanyNode.SetRegularCustomers(neighborNodes)
		localCompanyNode.InitAdvertisingBudget()
		utils.PrintMessage("This node has the folowing settings: ")
		utils.PrintMessage(localCompanyNode)

		//Set the customerNode's id to -1, because it is easier to identify the "unset" node.
		localCustomerNode.SetCustomerId(-1)
	}
	messageToAllNeighborsSent = false

	protobufChannel := make(chan *protobuf.MessageTwo)
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
			//Method is placed in messageHandler.go
			go ReceiveAndParseIncomingProtobufMessageToChannel(conn, protobufChannel)
			//ReceiveAndParseIncomingProtobufMessageToChannel(conn, protobufChannel)
			//protodata := ReceiveAndParseInfomingProtoufMessage(conn)
			//utils.PrintMessage(fmt.Sprintf("Message on %s received:\n\n%s\n\n", localNode.String(), protodata.String()))
			//handleReceivedProtobufMessage(protodata)
		}
	}
}

func isCustomerInitialized() bool {
	if localCustomerNode.CustomerId() != -1 {
		return true
	}
	return false
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
func handleReceivedProtobufMessageWithChannel(localNode server.NetworkServer, receivingChannel chan *protobuf.MessageTwo) {
	for {
		// This call blocks until a new message is available.
		message := <-receivingChannel
		handleReceivedProtobufMessage(localNode, message)
	}
}

// This method gets a protobuf message and decides if it is a control or a
// application message and gives it to the related function.
func handleReceivedProtobufMessage(localNode server.NetworkServer, protoMessage *protobuf.MessageTwo) {
	utils.PrintMessage(fmt.Sprintf("Message on %s received:\n\n%s\n\n", localNode.String(), protoMessage.String()))
	switch protoMessage.GetMessageType() {
	case protobuf.MessageTwo_CONTROLMESSAGE:
		utils.PrintMessage("Message is of type CONTROLMESSAGE.")
		handleReceivedControlMessage(protoMessage)
	case protobuf.MessageTwo_APPLICATIONMESSAGE:
		utils.PrintMessage("Message is of type APPLICATIONMESSAGE.")
		handleReceivedApplicationMessage(protoMessage)
	default:
		log.Fatalln("Read a unknown \"NachrichtenTyp\"")
	}
}

func handleReceivedControlMessage(message *protobuf.MessageTwo) {
	switch message.GetControlType() {
	case protobuf.MessageTwo_INITIALIZE:
		if !messageToAllNeighborsSent {
			for _, value := range neighbors {
				SendProtobufApplicationMessage(localNode, value, localeId, message.GetMessageContent(), isCustomerInitialized())
			}
			messageToAllNeighborsSent = true
		}
	case protobuf.MessageTwo_QUIT:
		for id, destinationNode := range neighbors {
			SendProtobufControlMessage(localNode, destinationNode, id, utils.CONTROL_TYPE_EXIT, message.GetMessageContent(), isCustomerInitialized())
		}
		utils.PrintMessage("Received a QUIT message, so program will be exited.")
		os.Exit(0)
	default:
		log.Fatalln("Read a unknown \"ControlType\"")
	}
}

func handleReceivedApplicationMessage(message *protobuf.MessageTwo) {
	if !messageToAllNeighborsSent {
		for _, value := range neighbors {
			SendProtobufApplicationMessage(localNode, value, localeId, message.GetMessageContent(), isCustomerInitialized())
		}
		messageToAllNeighborsSent = true
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
