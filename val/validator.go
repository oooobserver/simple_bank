package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(input string, minLength, maxLength int) error {
	l := len(input)
	if l < minLength || l > maxLength {
		return fmt.Errorf("must contain %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUserName(input string) error {
	// Name length should between 3-100
	err := ValidateString(input, 3, 100)
	if err != nil {
		return err
	}

	if !isValidUsername(input) {
		return fmt.Errorf("name should only contains character, digits, or underscore")
	}

	return nil
}

func ValidFullName(input string) error {
	// Name length should between 3-100
	err := ValidateString(input, 3, 100)
	if err != nil {
		return err
	}

	if !isValidFullName(input) {
		return fmt.Errorf("name should only contains letters")
	}

	return nil
}

func ValidatePasswaord(input string) error {
	return ValidateString(input, 3, 100)
}

func ValidateEmail(input string) error {
	if err := ValidateString(input, 1, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(input); err != nil {
		return fmt.Errorf("is not valid email address")
	}
	return nil
}
