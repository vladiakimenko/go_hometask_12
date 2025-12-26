package service

import (
	"fmt"
	"unicode"

	"github.com/dlclark/regexp2"
)

var passwordRE = regexp2.MustCompile(`^(?=.*[A-Z])(?=.*\d)(?=.*[_^\w\s]).{8,80}$`, regexp2.Compiled)
var emailRE = regexp2.MustCompile(`^[a-zA-Z0-9._%+-]{1,64}@([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,63}$`, regexp2.Compiled)

const MinUsernameLength = 3
const MaxUsernameLength = 80

func ValidateEmail(email string) error {
	valid, _ := emailRE.MatchString(email) // прекомпилированный regexp не вернёт ошибку в рантайме, игнорим
	if !valid {
		return fmt.Errorf(
			"email must be in the format user@example.com; " +
				"local-part: 1–64 characters; " +
				"domain: 1–63 characters; " +
				"TLD: 2–63 characters; total: ≤254 characters",
		)
	}
	return nil
}

func ValidateUsername(username string) error {
	if MinUsernameLength > len(username) || len(username) > MaxUsernameLength {
		return fmt.Errorf(
			"username length must be between %d and %d characters",
			MinUsernameLength,
			MaxUsernameLength,
		)
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("username must be alphanumeric")
		}
	}
	return nil
}

func ValidatePassword(password string) error {
	valid, _ := passwordRE.MatchString(password)
	if !valid {
		return fmt.Errorf(
			"password must be at least 8 and at most 80 characters long " +
				"and include an uppercase letter, a number, and a special character",
		)
	}
	return nil
}
