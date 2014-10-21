// The server_test package is used to give some unit tests
// which are used to provide basic functionalities for the
// client package.
package server_test

import (
	"github.com/jzipfler/htw/ava/server"
	"strconv"
	"testing"
)

// Tests if the portranges are corret implemented
func TestServerPortRange(t *testing.T) {
	serverObject := server.New()
	toLow := -1
	if error := serverObject.SetPort(toLow); error == nil {
		t.Error("A port of " + strconv.Itoa(toLow) + " is not allowed to exists, because the minimal port number is " + strconv.Itoa(server.MINIMUM_PORT) + " .")
	}
	correctPortSize := 15108
	if error := serverObject.SetPort(correctPortSize); error != nil {
		t.Error("A port of " + strconv.Itoa(correctPortSize) + " should not return an error, because it is in a correct port range.")
	}
	toHigh := 151088
	if error := serverObject.SetPort(toHigh); error == nil {
		t.Error("A port of " + strconv.Itoa(toHigh) + " is not allowed to exists, because the maximum port number is " + strconv.Itoa(server.MAXIMUM_PORT) + " .")
	}
}
