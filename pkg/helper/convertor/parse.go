package convertor

import (
	"fmt"
	"strconv"
)

// ParseInt convert value to integer, In case of error it returns 0
// and log message
func ParseInt(paramName, value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid mysql %s value: %w", paramName, err)
	}
	return i, nil
}

// ParseBool converts a string to a bool.
func ParseBool(paramName, value string) (bool, error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid %s value: %w", paramName, err)
	}
	return b, nil
}
