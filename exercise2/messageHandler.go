package exercise2

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"code.google.com/p/goprotobuf/proto"
	"github.com/jzipfler/htw-ava/protobuf"
	"github.com/jzipfler/htw-ava/server"
	"github.com/jzipfler/htw-ava/utils"
)

// This function sends a application message (ANWENDUNGSNACHRICHT) to the neighbor with
// the given targetId. If the id does not exists, it just returns and does nothing.
func SendProtobufApplicationMessage(sourceServer, destinationServer server.NetworkServer, destnationServerId int, messageContent string) error {
	if destinationServer.IpAddressAsString() == "" {
		return errors.New(fmt.Sprintf("The target server has no ip address or port.\n%s\n", destinationServer.IpAndPortAsString(), utils.ERROR_FOOTER))
	}
	utils.PrintMessage(fmt.Sprintf("Encode protobuf application message for node with IP:PORT : %s.", destinationServer.IpAndPortAsString()))
	protobufMessage := new(protobuf.Nachricht)
	protobufMessage.SourceIP = proto.String(sourceServer.IpAddressAsString())
	protobufMessage.SourcePort = proto.Int(sourceServer.Port())
	protobufMessage.SourceID = proto.Int(destnationServerId)
	nachrichtenTyp := protobuf.Nachricht_NachrichtenTyp(protobuf.Nachricht_ANWENDUNGSNACHRICHT)
	protobufMessage.NachrichtenTyp = &nachrichtenTyp
	protobufMessage.NachrichtenInhalt = proto.String(messageContent)
	protobufMessage.ZeitStempel = proto.String(time.Now().UTC().String())
	//Protobuf message filled with data. Now marshal it.
	data, err := proto.Marshal(protobufMessage)
	if err != nil {
		return err
	}
	conn, err := net.Dial(destinationServer.UsedProtocol(), destinationServer.IpAndPortAsString())
	if err != nil {
		return err
	}
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	utils.PrintMessage(fmt.Sprintf("Application message from %s to %s sent:\n\n%s\n\n", sourceServer.String(), destinationServer.String(), protobufMessage.String()))
	utils.PrintMessage("Sent " + strconv.Itoa(n) + " bytes")
	return nil
}

// This function sends a control message (KONTROLLNACHRICHT) to the node with
// the given targetId. If the id does not exists, it just returns and does nothing.
func SendProtobufControlMessage(sourceServer, destinationServer server.NetworkServer, destnationServerId, controlType int, messageContent string) error {
	if destinationServer.IpAddressAsString() == "" {
		return errors.New(fmt.Sprintf("The target server has no ip address or port.\n%s\n", destinationServer.IpAndPortAsString(), utils.ERROR_FOOTER))
	}
	utils.PrintMessage(fmt.Sprintf("Encode protobuf control message for node with IP:PORT : %s.", destinationServer.IpAndPortAsString()))
	protobufMessage := new(protobuf.Nachricht)
	protobufMessage.SourceIP = proto.String(sourceServer.IpAddressAsString())
	protobufMessage.SourcePort = proto.Int(sourceServer.Port())
	protobufMessage.SourceID = proto.Int(destnationServerId)
	nachrichtenTyp := protobuf.Nachricht_NachrichtenTyp(protobuf.Nachricht_KONTROLLNACHRICHT)
	protobufMessage.NachrichtenTyp = &nachrichtenTyp
	var kontrollTyp protobuf.Nachricht_KontrollTyp
	switch controlType {
	case utils.CONTROL_TYPE_INIT:
		kontrollTyp = protobuf.Nachricht_KontrollTyp(protobuf.Nachricht_INITIALISIEREN)
	case utils.CONTROL_TYPE_EXIT:
		kontrollTyp = protobuf.Nachricht_KontrollTyp(protobuf.Nachricht_BEENDEN)
	default:
		utils.PrintMessage("No valid controlType given. Assume EXIT.")
		kontrollTyp = protobuf.Nachricht_KontrollTyp(protobuf.Nachricht_BEENDEN)
	}
	protobufMessage.KontrollTyp = &kontrollTyp
	protobufMessage.NachrichtenInhalt = proto.String(messageContent)
	protobufMessage.ZeitStempel = proto.String(time.Now().UTC().String())
	//Protobuf message filled with data. Now marshal it.
	data, err := proto.Marshal(protobufMessage)
	if err != nil {
		return err
	}
	utils.PrintMessage(fmt.Sprintf("Control message from %s to %s sent:\n\n%s\n\n", sourceServer.String(), destinationServer.String(), protobufMessage.String()))
	conn, err := net.Dial(destinationServer.UsedProtocol(), destinationServer.IpAndPortAsString())
	if err != nil {
		return err
	}
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	utils.PrintMessage("Sent " + strconv.Itoa(n) + " bytes")
	return nil
}

// This function uses a established connection to parse the data there to the
// protobuf message and returns it.
func ReceiveAndParseInfomingProtoufMessage(conn net.Conn) *protobuf.Nachricht {
	utils.PrintMessage("Incoming message")
	//Close the connection when the function exits
	defer conn.Close()
	//Create a data buffer of type byte slice with capacity of 4096
	data := make([]byte, 4096)
	//Read the data waiting on the connection and put it in the data buffer
	n, err := conn.Read(data)
	if err != nil {
		log.Fatal("Error happened: " + err.Error())
	}
	utils.PrintMessage("Decoding Protobuf message")
	//Create an struct pointer of type ProtobufTest.TestMessage struct
	protodata := new(protobuf.Nachricht)
	//Convert all the data retrieved into the ProtobufTest.TestMessage struct type
	err = proto.Unmarshal(data[0:n], protodata)
	if err != nil {
		log.Fatal("Error happened: " + err.Error())
	}
	utils.PrintMessage("Message decoded.")
	return protodata
}

// This function uses a established connection to parse the data there to the
// protobuf message. The result gets assigned to the channel instead of
// returning it.
func ReceiveAndParseIncomingProtobufMessageToChannel(conn net.Conn, c chan *protobuf.Nachricht) {
	protodata := ReceiveAndParseInfomingProtoufMessage(conn)
	utils.PrintMessage("Sending decoded message to channel.")
	//Push the protobuf message into a channel
	c <- protodata
}
