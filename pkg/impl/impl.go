package impl

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hyperxpizza/auth-service/pkg/auth"
	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/hyperxpizza/auth-service/pkg/database"
	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (a AuthServiceServer) GenerateToken(ctx context.Context, data *pb.TokenData) (*pb.Token, error) {
	a.logger.Infof("generating token for: %s", data.Username)
	var tokenResponse pb.Token

	//check if user exists in the database
	_, err := a.db.GetUser(data.Id, data.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(
				codes.NotFound,
				"User not found in the database",
			)
		}
	}

	token, err := a.authenticator.GenerateToken(data.Id, data.Username)
	if err != nil {
		a.logger.Warnf("generating jwt token for: %s with id: %d failed: %s", data.Username, data.Id, err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	tokenResponse.Token = token

	return &tokenResponse, nil
}

func (a AuthServiceServer) ValidateToken(ctx context.Context, token *pb.Token) (*pb.TokenData, error) {
	var tokenData pb.TokenData

	username, id, err := a.authenticator.ValidateToken(token.Token)
	if err != nil {
		a.logger.Warnf("validating jwt token failed: %s", err.Error())
		return nil, status.Error(
			codes.PermissionDenied,
			err.Error(),
		)
	}

	tokenData.Id = id
	tokenData.Username = username

	return &tokenData, nil
}

func (a AuthServiceServer) AddUser(ctx context.Context, user *pb.User) (*pb.ID, error) {
	var id pb.ID

	return &id, nil
}

func (a AuthServiceServer) RemoveUser(ctx context.Context, id *pb.ID) (emptypb.Empty, error) {
	return emptypb.Empty{}, nil
}

func (a AuthServiceServer) UpdateUser(ctx context.Context, user *pb.User) (emptypb.Empty, error) {
	return emptypb.Empty{}, nil
}
