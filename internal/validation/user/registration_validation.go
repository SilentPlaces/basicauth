package validation

import (
	"errors"
	"fmt"
	"github.com/SilentPlaces/basicauth.git/internal/config"
	"regexp"
	"unicode"
)

// ValidateEmail checks if the provided email matches a basic regex pattern.
func ValidateEmail(email string) error {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidatePassword checks whether the provided password meets the requirements of the given strength level.
func ValidatePassword(password string, strength *config.RegistrationPasswordConfig) error {

	if len(password) < strength.MinLength {
		return errors.New(fmt.Sprintf("password must be at least %d characters long", strength.MinLength))
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if strength.RequireUpper && !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if strength.RequireLower && !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if strength.RequireNumber && !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if strength.RequireSpecial && !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
