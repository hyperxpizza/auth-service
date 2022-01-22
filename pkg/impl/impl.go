package impl

import (
	"context"

	"github.com/hyperxpizza/auth-service/pkg/config"
	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/sirupsen/logrus"
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
	var tokenResponse pb.Token

	return &tokenResponse, nil
}

func (a AuthServiceServer) ValidateToken(ctx context.Context, token *pb.Token) (*pb.TokenData, error) {
	var tokenData pb.TokenData

	return &tokenData, nil
}
