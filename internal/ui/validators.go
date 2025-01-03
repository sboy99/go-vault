package ui

import (
	"errors"
	"strconv"
)

func portValidator(input string) error {
	// Validate port //
	_, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return errors.New("Invalid port")
	}
	return nil
}
