package utils

import (
	"errors"
	"fmt"
	"regexp"

	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"golang.org/x/crypto/bcrypt"
)

const (
	usernameNotValidError                       = "username: %s not valid"
	passwordsNotMatchingError                   = "passwords are not matching"
	passwordNeedsAtLeastOneNumError             = "password needs at least one number"
	passwordNeedsAtLeastOneSmallCharError       = "password needs at least one small character"
	passwordNeedsAtLeastOneUpperCaseCharError   = "password needs at leat one uppercase character"
	passwordNeedsAtLeastOneSpecialCaseCharError = "password needs at least one special character"
	num                                         = `[0-9]{1}`
	az                                          = `[a-z]{1}`
	AZ                                          = `[A-Z]{1}`
	symbol                                      = `[!@#~$%^&*()+|_]{1}`
)

var usernameConventon = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`).MatchString

func ValidateRegisterUser(user *pb.InsertAuthServiceUserRequest) error {

	if err := ValidatePassword(user.Password1, user.Password2); err != nil {
		return err
	}

	if !usernameConventon(user.Username) {
		return fmt.Errorf(usernameNotValidError, user.Username)
	}

	return nil
}

func ValidatePassword(pwd1, pwd2 string) error {
	if pwd1 != pwd2 {
		return errors.New(passwordsNotMatchingError)
	}

	if ok, err := regexp.MatchString(num, pwd1); !ok || err != nil {
		return errors.New(passwordNeedsAtLeastOneNumError)
	}
	if ok, err := regexp.MatchString(az, pwd1); !ok || err != nil {
		return errors.New(passwordNeedsAtLeastOneSmallCharError)
	}
	if ok, err := regexp.MatchString(AZ, pwd1); !ok || err != nil {
		return errors.New(passwordNeedsAtLeastOneUpperCaseCharError)
	}
	if ok, err := regexp.MatchString(symbol, pwd1); !ok || err != nil {
		return errors.New(passwordNeedsAtLeastOneSpecialCaseCharError)
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
