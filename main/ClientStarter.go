// ClientStarter
package main

import (
	"fmt"
	"github.com/jzipfler/htw/ava/client"
	"github.com/jzipfler/htw/ava/server"
)

func main() {
	client := client.New()
	serverObject := server.New()
	fmt.Printf("%T\n", client)
	fmt.Printf("%T\n", serverObject)
	fmt.Println("Client: " + client.String() + "\nServer: " + serverObject.String())
	error := client.SetIpAddressAsString("1.2.3.4")
	if error != nil {
		fmt.Println(error)
		return
	}
	client.SetClientName("First")
	fmt.Println(client)
}
