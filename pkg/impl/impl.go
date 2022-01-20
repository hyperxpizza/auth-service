package impl

import (
	"context"

	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
)

type AuthServiceServer struct {
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
