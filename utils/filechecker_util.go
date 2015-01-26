package utils

import (
	"errors"
	"os"
)

// This method takes the string and checks if the string points to a file
// that exists.
func CheckIfFileExists(pathToFile string) bool {
	if _, err := os.Stat(pathToFile); err == nil {
		//Not needed anymore
		//if file.IsDir() {
		//	return true
		//}
		return true
	}
	return false
}

// This method checks if the file that is given is readable.
func CheckIfFileIsReadable(pathToFile string) (bool, error) {
	if exists := CheckIfFileExists(pathToFile); !exists {
		return false, errors.New("The file does not exists.")
	}
	if _, err := os.Open(pathToFile); os.IsPermission(err) {
		return false, errors.New("The file " + pathToFile + " is not readable.")
	}
	return true, nil
}

func CheckIfFileIsWritebale(pathToFile string) (bool, error) {
	if _, err := os.OpenFile(pathToFile, os.O_WRONLY, 0777); os.IsPermission(err) {
		return false, errors.New("The file " + pathToFile + " is not writable.")
	}
	return true, nil
}

func CheckIfFileIsReadableAndWritebale(pathToFile string) (bool, error) {
	if readable, err := CheckIfFileIsReadable(pathToFile); !readable {
		return false, err
	}
	if _, err := os.OpenFile(pathToFile, os.O_RDWR, 0777); os.IsPermission(err) {
		return false, errors.New("The file " + pathToFile + " is not writable.")
	}
	return true, nil
}
