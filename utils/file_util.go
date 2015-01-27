package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

func ReadNumbersFromFirstLine(fileObject *os.File, numberOfValues int) ([]int, error) {
	if fileObject == nil {
		return nil, errors.New("The file object was nil.")
	}
	if numberOfValues < 1 {
		return nil, errors.New("The number of values must be greater than 0.")
	}
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

func IncreaseNumbersFromFirstLine(fileObject *os.File, numberOfValues int) ([]int, error) {
	var buffer bytes.Buffer

	numbers, err := ReadNumbersFromFirstLine(fileObject, numberOfValues)
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

	for i := 0; i < len(numbers); i++ {
		if _, err := buffer.WriteString(string(numbers[i] + 48)); err != nil {
			return nil, err
		}
	}

	if _, err := fileObject.WriteAt(buffer.Bytes(), 0); err != nil {
		return nil, err
	} /*else {
		PrintMessage(fmt.Sprintf("Wrote %d bytes: %v.", wroteBytes, buffer.Bytes()))
	}*/
	return numbers, nil
}

func DecreaseNumbersFromFirstLine(fileObject *os.File, numberOfValues int) ([]int, error) {
	var buffer bytes.Buffer

	numbers, err := ReadNumbersFromFirstLine(fileObject, numberOfValues)
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

	for i := 0; i < len(numbers); i++ {
		if _, err := buffer.WriteString(string(numbers[i] + 48)); err != nil {
			return nil, err
		}
	}

	if _, err := fileObject.WriteAt(buffer.Bytes(), 0); err != nil {
		return nil, err
	} /*else {
		PrintMessage(fmt.Sprintf("Wrote %d bytes: %v.", wroteBytes, buffer.Bytes()))
	}*/
	return numbers, nil
}
