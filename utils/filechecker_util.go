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
	if _, err := os.Open(pathToFile); !os.IsPermission(err) {
		return true, nil
	}
	return false, nil
}

func CheckIfFileIsWritebale(pathToFile string) error {
	panic("Writable test is not implemented yet.")
}
