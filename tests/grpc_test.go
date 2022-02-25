package main

import (
	"context"

	"google.golang.org/grpc/test/bufconn"
)

const (
	buffer = 1024 * 1024
	target = "bufnet"
)

var lis *bufconn.Listener
var ctx = context.Background()

/*
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

func samplePbUser() pb.AuthServiceUser {

	getPwdHash := func(pwd string) string {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
		if err != nil {
			panic(err)
		}

		return string(hashedPassword)
	}

	return pb.AuthServiceUser{
		Id:                    1,
		Username:              "pizza",
		PasswordHash:          getPwdHash("some-password"),
		Created:               timestamppb.Now(),
		Updated:               timestamppb.Now(),
		RelatedUsersServiceID: 1,
	}
}

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

	user := samplePbUser()

	//insert user into the database
	id, err := client.AddUser(ctx, &user)
	assert.NoError(t, err)

	req := pb.TokenRequest{
		UsersServiceID: id.Id,
		Username:       user.Username,
	}

	token, err := client.GenerateToken(ctx, &req)
	assert.NoError(t, err)

	data, err := client.ValidateToken(ctx, token)
	assert.NoError(t, err)

	assert.Equal(t, data.AuthServiceID, id.Id)
	assert.Equal(t, data.Username, user.Username)

	if *deleteOpt {
		req := pb.ID{
			Id: id.Id,
		}
		_, err := client.RemoveUser(ctx, &req)
		assert.NoError(t, err)
	}
}
*/
