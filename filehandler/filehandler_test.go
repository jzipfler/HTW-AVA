// This package is used to test the functionalities of the filehandler package.
package filehandler_test

import (
	"github.com/jzipfler/htw-ava/filehandler"
	"os"
	"path"
	"testing"
)

const (
	FILE_WITH_COMMENTS       = "file_with_comments.txt"
	FILE_WITH_HOSTNAMES      = "file_with_hostnames.txt"
	FILE_WITH_IPV4_ADDRESSES = "file_with_ipv4.txt"
	EMPTY_FILE               = "empty_file.txt"

	COMMENT_LINE                   = "#This is a comment\n"
	NODE_ENTRY_WITH_IPV4           = "1 127.0.0.1:15101\n"
	NODE_ENTRY_WITH_HOSTNAME       = "1 localhost:15102\n"
	NODE_ENTRY_WITH_WRONG_HOSTNAME = "1 host123name123:15100\n"
)

// Tests the CollectAllFromNodeListFile in HandleNodeListFile.go.
func TestNodeListFilehandler(t *testing.T) {
	// First test the comment function
	pathToFileWithComments := path.Join(os.TempDir(), FILE_WITH_COMMENTS)
	fileWithComments, err := os.Create(pathToFileWithComments)
	if err != nil {
		t.Fatal("Could not create the comment file.")
	} else {
		defer fileWithComments.Close()
		defer os.Remove(pathToFileWithComments)
		if _, err := fileWithComments.WriteString(COMMENT_LINE); err != nil {
			t.Fatal("Could not write to the comment file.")
		}
		fileWithComments.Sync()
		if _, err := filehandler.CollectAllFromNodeListFile(pathToFileWithComments); err == nil {
			t.Error("There is no error at the fileWithComments: that should not happen because the nodeList is empty.")
		} else {
			if err.Error() != "No nodes present... ABORT" {
				t.Error("The fileWithComments should return a error with the string \"No nodes present... ABORT\".")
			} else {
				t.Log("Everithing fine on comment test")
			}
		}
	}
	// Then test the hostname resolving
	pathToFileWithHostnames := path.Join(os.TempDir(), FILE_WITH_HOSTNAMES)
	fileWithHostnames, err := os.Create(pathToFileWithHostnames)
	if err != nil {
		t.Fatal("Could not create the hostname file.")
	} else {
		defer fileWithHostnames.Close()
		defer os.Remove(pathToFileWithHostnames)
		if _, err := fileWithHostnames.WriteString(NODE_ENTRY_WITH_HOSTNAME); err != nil {
			t.Fatal("Could not write to the hostname file.")
		}
		if nodes, err := filehandler.CollectAllFromNodeListFile(pathToFileWithHostnames); err != nil {
			t.Error("There is a error on the hostname which should not happen: " + err.Error())
		} else {
			if len(nodes) != 1 {
				t.Error("The lengh of the nodes should be 1 for the hostname test.")
			} else {
				t.Log("Everything fine for the hostname test.")
			}
		}
	}
	// The recognition of version 4 ip addresses
	pathToFileWithIpV4Addresses := path.Join(os.TempDir(), FILE_WITH_IPV4_ADDRESSES)
	fileWithIpV4Addresses, err := os.Create(pathToFileWithIpV4Addresses)
	if err != nil {
		t.Fatal("Could not create the ipv4 file.")
	} else {
		defer fileWithIpV4Addresses.Close()
		defer os.Remove(pathToFileWithIpV4Addresses)
		if _, err := fileWithIpV4Addresses.WriteString(NODE_ENTRY_WITH_IPV4); err != nil {
			t.Fatal("Could not write to the ipv4 file.")
		}
		if nodes, err := filehandler.CollectAllFromNodeListFile(pathToFileWithIpV4Addresses); err != nil {
			t.Error("There is a error on the ipv4 file which should not happen: " + err.Error())
		} else {
			if len(nodes) != 1 {
				t.Error("The lengh of the nodes should be 1 for the ipv4 test.")
			} else {
				t.Log("Everything fine for the ipv4 test.")
			}
		}
	}
	// The behaviour of a empty file.
	pathToEmptyFile := path.Join(os.TempDir(), EMPTY_FILE)
	emptyFile, err := os.Create(pathToEmptyFile)
	if err != nil {
		t.Fatal("Could not create the empty file.")
	} else {
		defer emptyFile.Close()
		defer os.Remove(pathToEmptyFile)
		if _, err := filehandler.CollectAllFromNodeListFile(pathToEmptyFile); err == nil {
			t.Error("No error is reported for the empty file. That should not happen.")
		} else {
			t.Log("Everything fine for the empty file test. Error message: " + err.Error())
		}
	}
	//TODO:
	// The behaviour of a not existing file.
	// What happens with wrong port ranges/definitions?
	// What happens with from formatted lines?
}

func TestGraphvizFilehandler(t *testing.T) {
	if _, err := filehandler.CollectNeighborsFromGraphvizFile("filename"); err == nil {
		t.Error("The error should not be nil because the function is not implemented yet.")
	}
}
