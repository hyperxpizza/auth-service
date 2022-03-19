package utils

import (
	"errors"
	"regexp"

	"github.com/hyperxpizza/auth-service/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var usernameConventon = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`).MatchString

func ValidateUser(user models.User) error {
	if !usernameConventon(user.Username) {
		return errors.New("username not valid")
	}

	return nil
}

func GeneratePasswordHash(pwd string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
