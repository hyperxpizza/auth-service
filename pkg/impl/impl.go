package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/hyperxpizza/auth-service/pkg/auth"
	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/hyperxpizza/auth-service/pkg/database"
	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/hyperxpizza/auth-service/pkg/validator"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	UserNotFoundError = "User not found in the database"
)

type AuthServiceServer struct {
	cfg           *config.Config
	logger        logrus.FieldLogger
	authenticator *auth.Authenticator
	db            *database.Database
	pb.UnimplementedAuthServiceServer
}

func NewAuthServiceServer(pathToConfig string, logger logrus.FieldLogger) (*AuthServiceServer, error) {
	cfg, err := config.NewConfig(pathToConfig)
	if err != nil {
		return nil, err
	}

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}

	authenticator := auth.NewAuthenticator(cfg)

	return &AuthServiceServer{
		cfg:           cfg,
		logger:        logger,
		db:            db,
		authenticator: authenticator,
	}, nil
}

func (a *AuthServiceServer) Run() {
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, a)

	addr := fmt.Sprintf(":%d", a.cfg.AuthService.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		a.logger.Fatalf("net.Listen failed: %s", err.Error())
	}

	a.logger.Infof("auth service server running on %s:%d", a.cfg.AuthService.Host, a.cfg.AuthService.Port)

	if err := grpcServer.Serve(lis); err != nil {
		a.logger.Fatalf("failed to serve: %s", err.Error())
	}
}

func (a *AuthServiceServer) GenerateToken(ctx context.Context, req *pb.TokenRequest) (*pb.Token, error) {
	a.logger.Infof("generating token for: %s", req.Username)
	var tokenResponse pb.Token

	//check if user exists in the database
	user, err := a.db.GetUser(req.UsersServiceID, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.Infof("user with id: %d and username: %s was not found in the database", req.UsersServiceID, req.Username)
			return nil, status.Error(
				codes.NotFound,
				UserNotFoundError,
			)
		}
		a.logger.Infof("database GetUser function returned an error: %s", err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)

	}

	token, err := a.authenticator.GenerateToken(user.ID, req.UsersServiceID, req.Username)
	if err != nil {
		a.logger.Infof("generating jwt token for: %s with id: %d failed: %s", req.Username, req.UsersServiceID, err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	tokenResponse.Token = token

	return &tokenResponse, nil
}

func (a *AuthServiceServer) ValidateToken(ctx context.Context, token *pb.Token) (*pb.TokenData, error) {
	var tokenData pb.TokenData

	username, authServiceID, usersServiceID, err := a.authenticator.ValidateToken(token.Token)
	if err != nil {
		a.logger.Infof("validating jwt token failed: %s", err.Error())
		return nil, status.Error(
			codes.PermissionDenied,
			err.Error(),
		)
	}

	tokenData.AuthServiceID = authServiceID
	tokenData.UsersServiceID = usersServiceID
	tokenData.Username = username

	return &tokenData, nil
}

func (a *AuthServiceServer) AddUser(ctx context.Context, user *pb.AuthServiceUser) (*pb.ID, error) {
	var id pb.ID

	a.logger.Infof("adding user: %s into the database", user.Username)

	unmappedUser := unMapUser(user)
	err := validator.ValidateUser(unmappedUser)
	if err != nil {
		a.logger.Infof("user: %s is not valid: %s", user.Username, err.Error())
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	idInt, err := a.db.InsertUser(unmappedUser)
	if err != nil {
		a.logger.Infof("inserting user: %s into the database failed: %s", user.Username, err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	id.Id = idInt

	a.logger.Infof("user: %s id: %d has been successfully inserted into the database", user.Username, idInt)

	return &id, nil
}

func (a *AuthServiceServer) RemoveUser(ctx context.Context, id *pb.ID) (*emptypb.Empty, error) {

	a.logger.Infof("deleting user with id: %d", id.Id)

	err := a.db.DeleteUser(id.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.Infof("user with id: %d was not found in the database", id.Id)
			return nil, status.Error(
				codes.NotFound,
				UserNotFoundError,
			)
		}
		a.logger.Infof("deleting user with id: %d has failed: %s", id.Id, err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)

	}

	a.logger.Infof("deleted user with id: %d", id.Id)

	return &emptypb.Empty{}, nil
}

func (a *AuthServiceServer) UpdateUser(ctx context.Context, user *pb.AuthServiceUser) (*emptypb.Empty, error) {
	a.logger.Infof("updating user with id: %d", user.Id)

	mapppedUser := unMapUser(user)

	err := a.db.UpdateUser(mapppedUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.Infof("user with id: %d was not found in the database", user.Id)
			return nil, status.Error(
				codes.NotFound,
				err.Error(),
			)
		}

		a.logger.Infof("updating user with id: %d failed: %s", user.Id, err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	return &emptypb.Empty{}, nil
}
