//The utils_test package is used for testing all functionalities in the
// utils package.
package utils_test

import (
	"os"
	"path"
	"testing"

	"github.com/jzipfler/htw-ava/utils"
)

const (
	FILE_WITH_READ_PERMISSION          = "file_with_read.txt"
	FILE_WITH_WRITE_PERMISSION         = "file_with_write.txt"
	FILE_WITH_READ_WRITE_PERMISSION    = "file_with_read_write.txt"
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
	if exists := utils.CheckIfFileExists(readWriteFilePath); !exists {
		os.Remove(readWriteFile.Name())
		t.Error("The file: \"" + readWriteFilePath + "\" exists but the check fails.")
	}
	os.Remove(readWriteFile.Name())
	notExistingFile := "datei123ASD.txt"
	if exists := utils.CheckIfFileExists(notExistingFile); exists {
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
	if readable, _ := utils.CheckIfFileIsReadable(readWriteFilePath); !readable {
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
	if readable, _ := utils.CheckIfFileIsReadable(withoutReadWriteFilePath); readable {
		os.Remove(file.Name())
		t.Error("The file: \"" + withoutReadWriteFilePath + "\" is not readable but no error occured during the check.")
	}
	os.Remove(file.Name())
}

// Use the constant files and a not existing file to check if the
// CheckIfFileIsWritable function works.
func TestIfFileIsWritable(t *testing.T) {
	//Create the file.
	withoutReadWriteFilePath := path.Join(os.TempDir(), FILE_WITHOUT_READ_WRITE_PERMISSION)
	file, err := os.Create(withoutReadWriteFilePath)
	defer os.Remove(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that nobody can read or write it.
	file.Chmod(0000)
	//Remove the file after checking. If a error occures and of not...
	if writable, _ := utils.CheckIfFileIsWritebale(withoutReadWriteFilePath); writable {
		t.Error("The file: \"" + withoutReadWriteFilePath + "\" is not writable but no error occured during the check.")
	}
	//Create the file.
	readFilePath := path.Join(os.TempDir(), FILE_WITH_READ_PERMISSION)
	readFile, err := os.Create(readFilePath)
	defer os.Remove(readFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that only the owner can read it.
	readFile.Chmod(0400)
	if writable, _ := utils.CheckIfFileIsWritebale(readFilePath); writable {
		t.Error("The file: \"" + readFilePath + "\" should not be writable but no error occured during the check.")
	}
	//Create the file.
	writeFilePath := path.Join(os.TempDir(), FILE_WITH_WRITE_PERMISSION)
	writeFile, err := os.Create(writeFilePath)
	defer os.Remove(writeFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that only the owner can read and write it.
	writeFile.Chmod(0200)
	if writable, _ := utils.CheckIfFileIsWritebale(writeFilePath); !writable {
		t.Error("The file: \"" + writeFilePath + "\" should be writable but an error occured during the check.")
	}
}

// Use the both constant files and a not existing file to check if the
// CheckIfFileIsWritable function works.
func TestIfFileIsReadableAndWritable(t *testing.T) {
	/*
		THE EXISTING FILE
	*/
	withoutReadWriteFilePath := path.Join(os.TempDir(), FILE_WITHOUT_READ_WRITE_PERMISSION)
	file, err := os.Create(withoutReadWriteFilePath)
	defer os.Remove(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that nobody can read or write it.
	file.Chmod(0000)
	//Remove the file after checking. If a error occures and of not...
	if readAndWritable, _ := utils.CheckIfFileIsReadableAndWritebale(withoutReadWriteFilePath); readAndWritable {
		t.Error("The file: \"" + withoutReadWriteFilePath + "\" is not read and writable but no error occured during the check.")
	}
	/*
		THE READ FILE
	*/
	readFilePath := path.Join(os.TempDir(), FILE_WITH_READ_PERMISSION)
	readFile, err := os.Create(readFilePath)
	defer os.Remove(readFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that only the owner can read it.
	readFile.Chmod(0400)
	if readAndWritable, _ := utils.CheckIfFileIsReadableAndWritebale(readFilePath); readAndWritable {
		t.Error("The file: \"" + readFilePath + "\" should not be read and writable but no error occured during the check.")
	}
	/*
		THE WRITE FILE
	*/
	writeFilePath := path.Join(os.TempDir(), FILE_WITH_WRITE_PERMISSION)
	writeFile, err := os.Create(writeFilePath)
	defer os.Remove(writeFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that only the owner can read and write it.
	writeFile.Chmod(0200)
	if readAndWritable, _ := utils.CheckIfFileIsReadableAndWritebale(writeFilePath); readAndWritable {
		t.Error("The file: \"" + writeFilePath + "\" should not be read and writable but an error occured during the check.")
	}
	/*
		THE READ WRITE FILE
	*/
	readWriteFilePath := path.Join(os.TempDir(), FILE_WITH_READ_WRITE_PERMISSION)
	readWriteFile, err := os.Create(readWriteFilePath)
	defer os.Remove(readWriteFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	//Change the permissions that only the owner can read and write it.
	readWriteFile.Chmod(0600)
	if readAndWritable, _ := utils.CheckIfFileIsReadableAndWritebale(readWriteFilePath); !readAndWritable {
		t.Error("The file: \"" + readWriteFilePath + "\" should be read and writable but an error occured during the check.")
	}
}

// Test the creation of a graphviz file.
func TestGraphvizGeneration(t *testing.T) {
	if err := utils.GenerateGraphvizFile("filename", 5, 6); err == nil {
		t.Error("The error should not be nil because the function is not implemented yet.")
	}
}
