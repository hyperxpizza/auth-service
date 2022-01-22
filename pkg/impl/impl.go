package impl

import (
	"context"

	"github.com/hyperxpizza/auth-service/pkg/config"
	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/sirupsen/logrus"
)

type AuthServiceServer struct {
	cfg    config.Config
	logger logrus.FieldLogger
	pb.UnimplementedAuthServiceServer
}

func (a AuthServiceServer) GenerateToken(ctx context.Context, data *pb.TokenData) (*pb.Token, error) {
	var tokenResponse pb.Token

	return &tokenResponse, nil
}

func (a AuthServiceServer) ValidateToken(ctx context.Context, token *pb.Token) (*pb.TokenData, error) {
	var tokenData pb.TokenData

	return &tokenData, nil
}
