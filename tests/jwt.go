package main

import (
	"errors"
	"testing"

	"github.com/hyperxpizza/auth-service/pkg/auth"
	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestJWTToken(t *testing.T) {
	validateFlags := func() error {
		if *idOpt == 0 {
			return errors.New("ID flag not set")
		}

		if *usernameOpt == "" {
			return errors.New("username flag not set")
		}

		if *configPathOpt == "" {
			return errors.New("config path not set")
		}

		return nil
	}

	err := validateFlags()
	assert.NoError(t, err)

	cfg, err := config.NewConfig(*configPathOpt)
	assert.NoError(t, err)

	authenticator := auth.NewAuthenticator(cfg)
	token, err := authenticator.GenerateToken(*idOpt, *usernameOpt)
	assert.NoError(t, err)

	username, id, err := authenticator.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, *usernameOpt, username)
	assert.Equal(t, *idOpt, id)

}
