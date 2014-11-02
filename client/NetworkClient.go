// Client.go
package client

import (
	"errors"
	"net"
)

// A network type from me
// It should contain a name of the network entity, a protocol that is used
// ("udp" or "tcp") and the ip address.
type NetworkClient struct {
	clientName   string
	usedProtocol string
	ipAddress    net.IP
}

const (
	UDP = "udp"
	TCP = "tcp"
)

// New returns a NetworkClient type.
func New() NetworkClient {
	return *new(NetworkClient)
}

// Simple getter for the NetworkClient which returns the net.IP value from the type.
func (networkClient NetworkClient) IpAddress() net.IP {
	return networkClient.ipAddress
}

// Function to return the IP address as string.
// If the address is not set, nil is returned.
func (networkClient NetworkClient) IpAddressAsString() string {
	if networkClient.IpAddress() == nil {
		return ""
	}
	return networkClient.IpAddress().String()
}

// Setter for the value from IP address.
func (networkClient *NetworkClient) SetIpAddress(ipAddress net.IP) {
	networkClient.ipAddress = ipAddress
}

// This method uses a string to set the IP address
// and returns a error if the string cannot be parsed.
// If the string can be parsed, error is nil.
func (networkClient *NetworkClient) SetIpAddressAsString(ipAddress string) error {
	networkClient.ipAddress = net.ParseIP(ipAddress)
	if networkClient.ipAddress == nil {
		return errors.New("String konnte nicht als IP Adresse geparsed werden. Beispielformat: \"1.2.3.4\"")
	}
	return nil
}

// Defines a getter method to access the name of the NetworkClient type.
func (networkClient NetworkClient) ClientName() string {
	return networkClient.clientName
}

// This function is used to set the name of the NetworkClient type.
func (networkClient *NetworkClient) SetClientName(clientName string) {
	networkClient.clientName = clientName
}

// Returns the currently set protcol. For exmaple "udp" or "tcp".
func (networkClient NetworkClient) UsedProtocol() string {
	return networkClient.usedProtocol
}

// Set the usedProtocol field to udp or tcp. If a other protocol is given as
// argument, the function returns an error.
func (networkClient *NetworkClient) SetUsedProtocol(usedProtocol string) error {
	if usedProtocol != UDP && usedProtocol != TCP {
		return errors.New("The protocol has to be udp or tcp.")
	}
	networkClient.usedProtocol = usedProtocol
	return nil
}

// Implements a string method thas is part of the String interface to
// have easier access on print functions.
func (networkClient NetworkClient) String() string {
	name := networkClient.ClientName()
	if len(name) == 0 {
		name = "\"\""
	}
	return "Name:" + name + "; Address:" + networkClient.IpAddress().String()
}
