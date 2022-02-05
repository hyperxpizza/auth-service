package validator

import (
	"errors"
	"regexp"

	"github.com/hyperxpizza/auth-service/pkg/models"
)

var usernameConventon = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`).MatchString

func ValidateUser(user models.User) error {
	if !usernameConventon(user.Username) {
		return errors.New("username not valid")
	}

	return nil
}
