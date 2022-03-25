package main

import (
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/hyperxpizza/auth-service/pkg/database"
	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/hyperxpizza/auth-service/pkg/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
const (
	postgresImage = "postgres:alpine"
)

var bgContext context.Context

func init() {
	bgContext = context.Background()
}

func startPostgresContainer() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
*/
func sampleUser() *pb.AuthServiceUser {

	getPwdHash := func(pwd string) string {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(hashedPassword))

		return string(hashedPassword)
	}

	return &pb.AuthServiceUser{
		Username:              "pizza",
		PasswordHash:          getPwdHash("some-password"),
		Created:               timestamppb.Now(),
		Updated:               timestamppb.Now(),
		RelatedUsersServiceID: 1,
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

//go test -v ./tests --run TestInsertUser --config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.dev.json --delete=true
func TestInsertUser(t *testing.T) {
	flag.Parse()

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

	if *deleteOpt {
		err = db.DeleteUser(id)
		assert.NoError(t, err)
		fmt.Printf("user with id: %d deleted", id)
	}

}

//go test -v ./tests --run TestGetUser --config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.dev.json --delete=true
func TestGetUser(t *testing.T) {
	flag.Parse()

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

	user2, err := db.GetUser(user.RelatedUsersServiceID, user.Username)
	assert.NoError(t, err)

	assert.Equal(t, user.Username, user2.Username)

	if *deleteOpt {
		err = db.DeleteUser(id)
		assert.NoError(t, err)
		fmt.Printf("user with id: %d deleted", id)
	}

}

//go test -v ./tests --run TestUpdateUser --config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.dev.json --delete=true
func TestUpdateUser(t *testing.T) {
	flag.Parse()

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
	user2.Id = id
	user2.Username = "updatedUsername"

	err = db.UpdateUser(user2)
	assert.NoError(t, err)

	user3, err := db.GetUser(id, user2.Username)
	assert.NoError(t, err)
	assert.Equal(t, user2.Username, user3.Username)

	if *deleteOpt {
		err = db.DeleteUser(id)
		assert.NoError(t, err)
		fmt.Printf("user with id: %d deleted", id)
	}

}

//go test -v ./tests --run TestChangePassword --config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.dev.json --delete=true
func TestChangePassword(t *testing.T) {
	flag.Parse()

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

	newPwd, err := utils.GeneratePasswordHash("newPassword1#@")
	assert.NoError(t, err)

	err = db.ChangePassword(id, user.Username, newPwd)
	assert.NoError(t, err)

	user2, err := db.GetUser(id, user.Username)
	assert.NoError(t, err)

	assert.Equal(t, newPwd, user2.PasswordHash)

	if *deleteOpt {
		err = db.DeleteUser(id)
		assert.NoError(t, err)
		fmt.Printf("user with id: %d deleted", id)
	}
}
