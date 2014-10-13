// Client.go
package client

import (
	"errors"
	"net"
)

// A network type from me
type NetworkClient struct {
	clientName string
	ipAddress  net.IP
}

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
		return errors.New("String konnte nicht als IP Adresse geparsed werden.")
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

// Implements a string method thas is part of the String interface to
// have easier access on print functions.
func (networkClient NetworkClient) String() string {
	name := networkClient.ClientName()
	if len(name) == 0 {
		name = "\"\""
	}
	return "Name:" + name + "; Address:" + networkClient.IpAddress().String()
}
