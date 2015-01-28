package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

func ReadNumbersFromFirstLine(filename string, numberOfValues int) ([]int, error) {
	if filename == "" {
		return nil, errors.New("The filename was empty.")
	}
	if numberOfValues < 1 {
		return nil, errors.New("The number of values must be greater than 0.")
	}
	if readable, err := CheckIfFileIsReadable(filename); !readable {
		return nil, err
	}
	fileObject, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer fileObject.Close()
	buffer := make([]byte, numberOfValues)
	if _, err := fileObject.ReadAt(buffer[:], 0); err != nil {
		return nil, err
	} /*else {
		PrintMessage(fmt.Sprintf("Read %d bytes: %v", numberOfReadBytes, buffer))
	}*/
	var numbers []int
	for _, value := range buffer {
		var number int
		if _, err := fmt.Sscanf(string(value), "%d", &number); err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}
	return numbers, nil
}

func IncreaseNumbersFromFirstLine(filename string, numberOfValues int) ([]int, error) {
	if rwPossible, err := CheckIfFileIsReadableAndWritebale(filename); !rwPossible {
		return nil, err
	}
	numbers, err := ReadNumbersFromFirstLine(filename, numberOfValues)
	if err != nil {
		return nil, err
	}
	for i := len(numbers) - 1; i >= 0; i-- {
		if i == 0 && numbers[i] == 9 {
			return nil, errors.New(fmt.Sprintf("The maximum value for this %d digit number is reached", len(numbers)))
		}
		if numbers[i] == 9 {
			numbers[i] = 0
		} else {
			numbers[i]++
			break
		}
	}
	var buffer bytes.Buffer
	for i := 0; i < len(numbers); i++ {
		if _, err := buffer.WriteString(string(numbers[i] + 48)); err != nil {
			return nil, err
		}
	}
	fileObject, err := os.OpenFile(filename, os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer fileObject.Close()
	if _, err := fileObject.WriteAt(buffer.Bytes(), 0); err != nil {
		return nil, err
	} /*else {
		PrintMessage(fmt.Sprintf("Wrote %d bytes: %v.", wroteBytes, buffer.Bytes()))
	}*/
	return numbers, nil
}

func DecreaseNumbersFromFirstLine(filename string, numberOfValues int) ([]int, error) {
	if rwPossible, err := CheckIfFileIsReadableAndWritebale(filename); !rwPossible {
		return nil, err
	}
	numbers, err := ReadNumbersFromFirstLine(filename, numberOfValues)
	if err != nil {
		return nil, err
	}
	minimum := true
	for i := 0; i < len(numbers); i++ {
		if numbers[i] != 0 {
			minimum = false
		}
	}
	if minimum {
		return nil, errors.New(fmt.Sprintf("The minimum value for this %d digit number is reached", len(numbers)))
	}
	for i := len(numbers) - 1; i >= 0; i-- {
		if numbers[i] == 0 {
			numbers[i] = 9
		} else {
			numbers[i]--
			break
		}
	}
	var buffer bytes.Buffer
	for i := 0; i < len(numbers); i++ {
		if _, err := buffer.WriteString(string(numbers[i] + 48)); err != nil {
			return nil, err
		}
	}
	fileObject, err := os.OpenFile(filename, os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer fileObject.Close()
	if _, err := fileObject.WriteAt(buffer.Bytes(), 0); err != nil {
		return nil, err
	} /*else {
		PrintMessage(fmt.Sprintf("Wrote %d bytes: %v.", wroteBytes, buffer.Bytes()))
	}*/
	return numbers, nil
}

func AppendStringToFile(filename, stringToAppend string, appendNewLine bool) error {
	if filename == "" {
		return errors.New("The filename can not be a empty string.")
	}
	if readable, err := CheckIfFileIsWritebale(filename); !readable {
		return err
	}
	fileObject, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer fileObject.Close()
	if appendNewLine {
		if _, err := fileObject.WriteString(stringToAppend + "\n"); err != nil {
			return err
		}
	} else {
		if _, err := fileObject.WriteString(stringToAppend); err != nil {
			return err
		}
	}
	return nil
}
