package server

import (
	"errors"
	"fmt"
	"log"
	"net"
)

var (
	serverConnection net.Listener
	externalLogger   *log.Logger
)

const (
	ERROR = "ERROR::"
)

// Starts the server on the ip address and port given by the NetworkServer type.
// If the logger value is nil, then the server uses the general logging mechanism
// to report some information.
// If a error happened, then the related one is returned otherwise it is nil.
func StartServer(serverObject NetworkServer, logger *log.Logger) error {
	if logger != nil {
		externalLogger = logger
	}
	var err error
	serverConnection, err = net.Listen(serverObject.UsedProtocol(), serverObject.IpAndPortAsString())
	if err != nil {
		serverErrorPrint(err.Error())
		return err
	}
	serverPrint(fmt.Sprintf("Listen on %s://%s", serverObject.UsedProtocol(), serverObject.IpAndPortAsString()))
	return nil
}

// ReceiveMessages is used to react on a received message from the network.
// A net.Conn type is returned if something is available.
func ReceiveMessage() (net.Conn, error) {
	if serverConnection == nil {
		return nil, errors.New("First start server.") //TODO: Besser machen
	}
	connection, err := serverConnection.Accept()
	if err != nil {
		return nil, err
	}
	return connection, nil
}

// Stops the server
func StopServer() {
	serverConnection.Close()
	serverConnection = nil
	serverPrint("Server closed.")
}

func serverPrint(message string) {
	if externalLogger == nil {
		log.Println(message)
	} else {
		externalLogger.Println(message)

	}
}

func serverErrorPrint(message string) {
	if externalLogger == nil {
		log.Println(ERROR + message)
	} else {
		externalLogger.Println(ERROR + message)

	}
}
