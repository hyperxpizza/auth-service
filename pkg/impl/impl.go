package impl

import (
	"context"

	"github.com/hyperxpizza/auth-service/pkg/config"
	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthServiceServer struct {
	cfg    *config.Config
	logger logrus.FieldLogger
	pb.UnimplementedAuthServiceServer
}

func NewAuthServiceServer(pathToConfig string, logger logrus.FieldLogger) (*AuthServiceServer, error) {
	cfg, err := config.NewConfig(pathToConfig)
	if err != nil {
		return nil, err
	}

	return &AuthServiceServer{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (a AuthServiceServer) GenerateToken(ctx context.Context, data *pb.TokenData) (*pb.Token, error) {
	a.logger.Infof("generating token for: %s", data.Username)
	var tokenResponse pb.Token

	return &tokenResponse, nil
}

func (a AuthServiceServer) ValidateToken(ctx context.Context, token *pb.Token) (*pb.TokenData, error) {
	var tokenData pb.TokenData

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
