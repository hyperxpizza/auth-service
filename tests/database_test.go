package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/hyperxpizza/auth-service/pkg/database"
	"github.com/hyperxpizza/auth-service/pkg/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func sampleUser() models.User {

	getPwdHash := func(pwd string) string {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
		if err != nil {
			panic(err)
		}

		return string(hashedPassword)
	}

	return models.User{
		ID:           1,
		Username:     "pizza",
		PasswordHash: getPwdHash("some-password"),
		Created:      time.Now(),
		Updated:      time.Now(),
	}
}

func newDB(cfgpath string) (*database.Database, error) {
	cfg, err := config.NewConfig(cfgpath)
	if err != nil {
		return nil, err
	}

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}

//go test -v ./tests --run TestInsertUser --config=/home/hyperxpizza/dev/golang/auth-service/config.json --delete=true
func TestInsertUser(t *testing.T) {
	validateFlags := func() error {
		if *configPathOpt == "" {
			return errors.New("config path not set")
		}

		return nil
	}

	err := validateFlags()
	assert.NoError(t, err)

	db, err := newDB(*configPathOpt)
	assert.NoError(t, err)

	defer db.Close()

	id, err := db.InsertUser(sampleUser())
	assert.NoError(t, err)

	fmt.Printf("inserted id: %d\n", id)

	if *delete {
		err = db.DeleteUser(id)
		assert.NoError(t, err)
		fmt.Printf("user with id: %d deleted", id)
	}

}

//go test -v ./tests --run TestGetUser --config=/home/hyperxpizza/dev/golang/auth-service/config.json --delete=true
func TestGetUser(t *testing.T) {

	validateFlags := func() error {
		if *configPathOpt == "" {
			return errors.New("config path not set")
		}

		return nil
	}

	err := validateFlags()
	assert.NoError(t, err)

	db, err := newDB(*configPathOpt)
	assert.NoError(t, err)

	defer db.Close()

	user := sampleUser()

	id, err := db.InsertUser(user)
	assert.NoError(t, err)

	fmt.Printf("inserted id: %d\n", id)

	user2, err := db.GetUser(id, user.Username)
	assert.NoError(t, err)

	assert.Equal(t, user.Username, user2.Username)

	if *delete {
		err = db.DeleteUser(id)
		assert.NoError(t, err)
		fmt.Printf("user with id: %d deleted", id)
	}

}

//go test -v ./tests --run TestUpdateUser --config=/home/hyperxpizza/dev/golang/auth-service/config.json --delete=true
func TestUpdateUser(t *testing.T) {
	validateFlags := func() error {
		if *configPathOpt == "" {
			return errors.New("config path not set")
		}

		return nil
	}

	err := validateFlags()
	assert.NoError(t, err)

	db, err := newDB(*configPathOpt)
	assert.NoError(t, err)

	defer db.Close()

	user := sampleUser()

	id, err := db.InsertUser(user)
	assert.NoError(t, err)

	fmt.Printf("inserted id: %d\n", id)

	user2 := sampleUser()
	user2.ID = id
	user2.Username = "updatedUsername"

	err = db.UpdateUser(user2)
	assert.NoError(t, err)

	user3, err := db.GetUser(id, user2.Username)
	assert.NoError(t, err)
	assert.Equal(t, user2.Username, user3.Username)

	if *delete {
		err = db.DeleteUser(id)
		assert.NoError(t, err)
		fmt.Printf("user with id: %d deleted", id)
	}

}
