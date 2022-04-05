package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"testing"

	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/hyperxpizza/auth-service/pkg/impl"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	buffer = 1024 * 1024
	target = "bufnet"
)

var lis *bufconn.Listener
var ctx = context.Background()

func mockGrpcServer(configPath string, secure bool) error {
	lis = bufconn.Listen(buffer)
	server := grpc.NewServer()

	logger := logrus.New()
	if level, err := logrus.ParseLevel(*loglevelOpt); err == nil {
		logger.Level = level
	}

	authServiceServer, err := impl.NewAuthServiceServer(configPath, logger)
	if err != nil {
		return err
	}

	pb.RegisterAuthServiceServer(server, authServiceServer)

	if err := server.Serve(lis); err != nil {
		logger.Fatal(err)
		return err
	}

	return nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func sampleUserRequest() pb.AuthServiceUserRequest {
	return pb.AuthServiceUserRequest{
		Username:              "pizza",
		Password1:             "Password!3",
		Password2:             "Password!3",
		RelatedUsersServiceID: 1,
	}
}

// go test -v ./tests --run TestGenerateToken -config=/home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.json
func TestGenerateToken(t *testing.T) {

	flag.Parse()

	validateFlags := func() error {
		if *configPathOpt == "" {
			return errors.New("config path not set")
		}

		return nil
	}

	err := validateFlags()
	assert.NoError(t, err)
	go mockGrpcServer(*configPathOpt, false)

	connection, err := grpc.DialContext(ctx, target, grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)

	defer connection.Close()

	client := pb.NewAuthServiceClient(connection)

	user := sampleUserRequest()

	//insert user into the database
	id, err := client.AddUser(ctx, &user)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	req := pb.TokenRequest{
		Username:       user.Username,
		UsersServiceID: user.RelatedUsersServiceID,
	}

	tokens, err := client.GenerateToken(ctx, &req)
	assert.NoError(t, err)

	accToken := pb.AccessTokenData{
		AccessToken: tokens.AccessToken,
	}

	data, err := client.ValidateToken(ctx, &accToken)
	assert.NoError(t, err)

	assert.Equal(t, user.Username, data.Username)
	assert.Equal(t, user.RelatedUsersServiceID, data.UsersServiceID)
	assert.Equal(t, id.Id, data.AuthServiceID)

	refToken := pb.RefreshTokenData{
		RefreshToken: tokens.RefreshToken,
	}

	_, err = client.RefreshToken(ctx, &refToken)
	assert.NoError(t, err)

	if *deleteOpt {
		_, err := client.RemoveUser(ctx, id)
		assert.NoError(t, err)
	}

}
