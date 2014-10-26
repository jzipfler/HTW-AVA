//The utils_test package is used for testing all functionalities in the
// utils package.
package utils_test

import (
	"github.com/jzipfler/htw-ava/utils"
	"os"
	"path"
	"testing"
)

const (
	FILE_WITH_READ_PERMISSION          = "file_with_write.txt"
	FILE_WITHOUT_READ_WRITE_PERMISSION = "file_without_read.txt"
)

// Use the both constant files and a not existing file to check if the
// CheckIfFileExists function works.
func TestIfFileExists(t *testing.T) {
	//Create the file.
	readWriteFilePath := path.Join(os.TempDir(), FILE_WITH_READ_PERMISSION)
	readWriteFile, err := os.Create(readWriteFilePath)
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that only the owner can read and write it.
	readWriteFile.Chmod(0600)
	if err := utils.CheckIfFileExists(readWriteFilePath); err != nil {
		os.Remove(readWriteFile.Name())
		t.Error("The file: \"" + readWriteFilePath + "\" exists but the check fails.")
	}
	os.Remove(readWriteFile.Name())
	notExistingFile := "datei123ASD.txt"
	if err := utils.CheckIfFileExists(notExistingFile); err == nil {
		t.Error("The file: \"" + notExistingFile + "\" does not exists but no error occured.")
	}
}

// Use the both constant files and a not existing file to check if the
// CheckIfFileIsReadable function works.
func TestIfFileIsReadable(t *testing.T) {
	//Create the file.
	readWriteFilePath := path.Join(os.TempDir(), FILE_WITH_READ_PERMISSION)
	readWriteFile, err := os.Create(readWriteFilePath)
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that only the owner can read and write it.
	readWriteFile.Chmod(0600)
	if err := utils.CheckIfFileIsReadable(readWriteFilePath); err != nil {
		os.Remove(readWriteFile.Name())
		t.Error("The file: \"" + readWriteFilePath + "\" should be readable but an error occured during the check.")
	}
	os.Remove(readWriteFile.Name())
	//Create the file.
	withoutReadWriteFilePath := path.Join(os.TempDir(), FILE_WITHOUT_READ_WRITE_PERMISSION)
	file, err := os.Create(withoutReadWriteFilePath)
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that only the owner can read it.
	file.Chmod(0200)
	//Remove the file after checking. If a error occures and of not...
	if err := utils.CheckIfFileIsReadable(withoutReadWriteFilePath); err == nil {
		os.Remove(file.Name())
		t.Error("The file: \"" + withoutReadWriteFilePath + "\" is not readable but no error occured during the check.")
	}
	os.Remove(file.Name())
}

// Use the both constant files and a not existing file to check if the
// CheckIfFileIsWritable function works.
func TestIfFileIsWritable(t *testing.T) {
	//TODO: Add implementation
}
