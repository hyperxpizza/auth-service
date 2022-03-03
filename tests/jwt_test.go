package main

import (
	"errors"
	"flag"
	"testing"

	"github.com/hyperxpizza/auth-service/pkg/auth"
	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/stretchr/testify/assert"
)

//go test -v ./tests/ --run TestJWTToken --id=1 --username=hyperxpizza --config=/home/hyperxpizza/dev/golang/auth-service/config.json
func TestJWTToken(t *testing.T) {

	flag.Parse()

	validateFlags := func() error {
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
	authenticator.GenerateTokenPairs(1, 1, "username")
}
