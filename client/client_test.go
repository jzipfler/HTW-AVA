// The client_test package is used to give some unit tests
// which are used to provide basic functionalities for the
// client package.
package client_test

import (
	"github.com/jzipfler/htw/ava/client"
	"testing"
)

// Tests the SetIpAddressAsString function.
func TestIpAddressInitialisation(t *testing.T) {
	clientObject := client.New()
	if error := clientObject.SetIpAddressAsString("1"); error == nil {
		t.Error("The string represenation of a IP address has to be like 1.2.3.4")
	} else {
		t.Log("ErrorOuput of the function: " + error.Error())
	}
	if error := clientObject.SetIpAddressAsString("1.2.3.4"); error != nil {
		t.Error("A correct IP address was given as string type.")
	}
}
