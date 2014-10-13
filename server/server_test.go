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
	if error := serverObject.SetPort(-1); error == nil {
		t.Error("Port has not to be lower than " + strconv.Itoa(server.MINIMUM_PORT) + ".")
	} else {
		t.Log(error)
	}
	if error := serverObject.SetPort(151088); error == nil {
		t.Error("Port has not to be higher than" + strconv.Itoa(server.MAXIMUM_PORT) + ".")
	}
}
