package impl

import (
	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/hyperxpizza/auth-service/pkg/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UnMapUser(user *pb.User) models.User {
	return models.User{
		ID:           user.Id,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Created:      user.Created.AsTime(),
		Updated:      user.Updated.AsTime(),
	}
}

func MapUser(user models.User) *pb.User {
	return &pb.User{
		Id:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Created:      timestamppb.New(user.Created),
		Updated:      timestamppb.New(user.Updated),
	}
}
