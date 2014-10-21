// The client_test package is used to give some unit tests
// which are used to provide basic functionalities for the
// client package.
package client_test

import (
	"fmt"
	"github.com/jzipfler/htw/ava/client"
	"testing"
)

// Tests the SetIpAddressAsString function.
func TestIpAddressParsing(t *testing.T) {
	clientObject := client.New()
	correctIP := "127.0.0.1"
	wrongIP := "12345"
	ipAddressOutOfRange := "620.10.42.201"
	characterIP := "Zehn.Zehn.Zehn.Hundert"
	if error := clientObject.SetIpAddressAsString(correctIP); error != nil {
		t.Error("A ip address of " + correctIP + " must be correct.")
	}
	if error := clientObject.SetIpAddressAsString(wrongIP); error == nil {
		t.Error("A ip address of " + wrongIP + " should not be parsed, because it is a simple number.")
	}
	if error := clientObject.SetIpAddressAsString(ipAddressOutOfRange); error == nil {
		t.Error("A ip address of " + ipAddressOutOfRange + " should not be parsed, because it is out of the ip address range.")
	}
	if error := clientObject.SetIpAddressAsString(characterIP); error == nil {
		t.Error("A ip address of " + characterIP + " should not be parsed, because it is a orinary string (word).")
	}
}

// Tests the SetUsedProtocol function.
func TestUsedProtocol(t *testing.T) {
	clientObject := client.New()
	if error := clientObject.SetUsedProtocol(client.TCP); error != nil {
		t.Error("The protocol " + client.TCP + " must be accepted.")
	}
	if error := clientObject.SetUsedProtocol(client.UDP); error != nil {
		t.Error("The protocol " + client.UDP + " must be accepted.")
	}
	wrongProtocol := "protocol"
	if error := clientObject.SetUsedProtocol(wrongProtocol); error == nil {
		t.Error("The protocol " + wrongProtocol + " should be not supported.")
	}
}

// Tests if the NetworkClient type implements the fmt.Stringer interface.
func TestStringer(t *testing.T) {
	clientObject := client.New()
	if output := fmt.Sprintln(clientObject); output == "" {
		t.Error("The NetworkClient seems not to implement the fmt.Stringer interface.")
	}
}
