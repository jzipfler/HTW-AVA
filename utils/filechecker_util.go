package utils

import (
	"errors"
	"os"
)

// This method takes the string and checks if the string points to a file
// that exists.
func CheckIfFileExists(pathToFile string) error {
	if file, err := os.Stat(pathToFile); err == nil {
		if file.IsDir() {
			return errors.New("The given path belongs to a folder.")
		}
		return nil
	}
	return errors.New("The given file does not exist.")
}

// This method checks if the file that is given is readable.
func CheckIfFileIsReadable(pathToFile string) error {
	if err := CheckIfFileExists(pathToFile); err != nil {
		return err
	}
	if _, err := os.Open(pathToFile); !os.IsPermission(err) {
		return nil
	}
	return errors.New("The user does not have permissions to read the given file")
}

func CheckIfFileIsWritebale(pathToFile string) error {
	panic("Writable test is not implemented yet.")
}
