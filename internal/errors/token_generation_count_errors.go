package custom_error

import "fmt"

// TokenGenerationCountError represents an exceeding token generation limit
type TokenGenerationCountError struct {
	Message string
}

func (e *TokenGenerationCountError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

func NewTokenGenerationError(message string) error {
	return &TokenGenerationCountError{
		Message: message,
	}
}
