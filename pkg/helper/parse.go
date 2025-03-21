package helpers

import (
	"fmt"
	"strconv"
)

func ParseInt(paramName, value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid mysql %s value: %w", paramName, err)
	}
	return i, nil
}
