//The utils_test package is used for testing all functionalities in the
// utils package.
package utils_test

import (
	"github.com/jzipfler/htw/ava/utils"
	"os"
	"testing"
)

const (
	FILE_WITH_READ_PERMISSION          = "file_with_write.txt"
	FILE_WITHOUT_READ_WRITE_PERMISSION = "/tmp/file_without_read.txt"
)

// Use the both constant files and a not existing file to check if the
// CheckIfFileExists function works.
func TestIfFileExists(t *testing.T) {
	if err := utils.CheckIfFileExists(FILE_WITH_READ_PERMISSION); err != nil {
		t.Error("The file: \"" + FILE_WITH_READ_PERMISSION + "\" exists but the check fails.")
	}
	notExistingFile := "datei123ASD.txt"
	if err := utils.CheckIfFileExists(notExistingFile); err == nil {
		t.Error("The file: \"" + notExistingFile + "\" does not exists but no error occured.")
	}
}

// Use the both constant files and a not existing file to check if the
// CheckIfFileIsReadable function works.
func TestIfFileIsReadable(t *testing.T) {
	if err := utils.CheckIfFileIsReadable(FILE_WITH_READ_PERMISSION); err != nil {
		t.Error("The file: \"" + FILE_WITH_READ_PERMISSION + "\" should be readable but an error occured during the check.")
	}
	//Create the file.
	file, err := os.Create(FILE_WITHOUT_READ_WRITE_PERMISSION)
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that only the owner can read it.
	file.Chmod(0200)
	//Remove the file after checking. If a error occures and of not...
	if err := utils.CheckIfFileIsReadable(FILE_WITHOUT_READ_WRITE_PERMISSION); err == nil {
		os.Remove(file.Name())
		t.Error("The file: \"" + FILE_WITHOUT_READ_WRITE_PERMISSION + "\" is not readable but no error occured during the check.")
	}
	os.Remove(file.Name())
}

// Use the both constant files and a not existing file to check if the
// CheckIfFileIsWritable function works.
func TestIfFileIsWritable(t *testing.T) {
	//TODO: Add implementation
}
