package server

import (
	"errors"
	"github.com/jzipfler/htw-ava/client"
	"strconv"
)

const (
	MINIMUM_PORT = 0
	MAXIMUM_PORT = 65535
)

// A server includes a client.NetworkClient type
// and adds a port where a client can connect to
type NetworkServer struct {
	client.NetworkClient
	port int
}

// New returns a NetworkServer type.
func New() NetworkServer {
	return NetworkServer{client.New(), 0}
}

func (networkServer *NetworkServer) SetPort(port int) error {
	if port < MINIMUM_PORT || port > MAXIMUM_PORT {
		return errors.New("Port must be between " + strconv.Itoa(MINIMUM_PORT) + " and " + strconv.Itoa(MAXIMUM_PORT))
	}
	networkServer.port = port
	return nil
}

// Returns the port as its integer value.
func (networkServer NetworkServer) Port() int {
	return networkServer.port
}

// IpAndPortAsString returns the content of ipAddress and port in the following
// format: "IP_ADDRESS:PORT".
// (The " are used to visualize the beginning and end of the format.)
func (networkServer NetworkServer) IpAndPortAsString() string {
	if networkServer.IpAddress() == nil {
		return ""
	}
	return networkServer.IpAddressAsString() + ":" + strconv.Itoa(networkServer.Port())
}

// String gives a string object which contains a representation of a NetworkServer type.
func (networkServer NetworkServer) String() string {
	return networkServer.IpAddressAsString() + ":" + strconv.Itoa(networkServer.Port())
}
