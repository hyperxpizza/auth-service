package main

import (
	"errors"
	"flag"
	"testing"
	"time"

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

	authenticator, err := auth.NewAuthenticator(cfg)
	assert.NoError(t, err)
	authenticator.GenerateTokenPairs(1, 1, *usernameOpt)
}

// go test -v ./tests --run TestTimeValidation --config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.json
func TestTimeValidation(t *testing.T) {
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

	cfg.AuthService.ExpirationTimeAccess = 2
	cfg.AuthService.ExpirationTimeRefresh = 3

	authenticator, err := auth.NewAuthenticator(cfg)
	assert.NoError(t, err)

	refToken, accToken, err := authenticator.GenerateTokenPairs(*authServiceIDOpt, *usersServiceIDOpt, *usernameOpt)
	assert.NoError(t, err)

	t.Run("Test Acces Token Validation", func(t *testing.T) {
		t.Parallel()
		err := authenticator.ValidateToken(accToken, false)
		assert.NoError(t, err)
		time.Sleep(time.Minute * time.Duration(cfg.AuthService.ExpirationTimeAccess+1))
		err = authenticator.ValidateToken(accToken, false)
		assert.Error(t, err)
	})

	t.Run("Test Refresh Token Validation", func(t *testing.T) {
		t.Parallel()
		err := authenticator.ValidateToken(refToken, true)
		assert.NoError(t, err)
		time.Sleep(time.Minute * time.Duration(cfg.AuthService.ExpirationTimeRefresh+1))
		err = authenticator.ValidateToken(refToken, true)
		assert.Error(t, err)
	})
}

//go test -v ./tests --run TestLogout --config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.json
func TestLogout(t *testing.T) {
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

	cfg.AuthService.ExpirationTimeAccess = 2
	cfg.AuthService.ExpirationTimeRefresh = 3

	authenticator, err := auth.NewAuthenticator(cfg)
	assert.NoError(t, err)

	refToken, accToken, err := authenticator.GenerateTokenPairs(*authServiceIDOpt, *usersServiceIDOpt, *usernameOpt)
	assert.NoError(t, err)

	err = authenticator.DeleteToken(*authServiceIDOpt, *usersServiceIDOpt, *usernameOpt)
	assert.NoError(t, err)

	t.Run("Test logout refresh token", func(t *testing.T) {
		t.Parallel()
		err = authenticator.ValidateToken(refToken, true)
		assert.Error(t, err)
	})

	t.Run("Test logout access token", func(t *testing.T) {
		t.Parallel()
		err = authenticator.ValidateToken(accToken, false)
		assert.Error(t, err)
	})
}
