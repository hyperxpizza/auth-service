package utils

import (
	"errors"
	"fmt"
	"regexp"

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

type toValidate interface {
	GetUsername() string
	GetPassword1() string
	GetPassword2() string
}

func ValidateRegisterUser[V toValidate](user V) error {

	if err := ValidatePassword(user.GetPassword1(), user.GetPassword2()); err != nil {
		return err
	}

	if !usernameConventon(user.GetUsername()) {
		return fmt.Errorf(usernameNotValidError, user.GetUsername())
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

func ComparePasswords(hashedPassword, toCheck string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(toCheck))
}
