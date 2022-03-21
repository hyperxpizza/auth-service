package main

import (
	"testing"

	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/hyperxpizza/auth-service/pkg/utils"
	"github.com/stretchr/testify/assert"
)

//go test -v ./tests --run TestUserValidation
func TestUserValidation(t *testing.T) {
	t.Run("Test validate register request", func(t *testing.T) {
		t.Parallel()
		req := pb.AuthServiceUserRequest{
			Username:              "pizza",
			Password1:             "PasswOrd1@",
			Password2:             "PasswOrd1@",
			RelatedUsersServiceID: 1,
		}
		err := utils.ValidateRegisterUser(&req)
		assert.NoError(t, err)
	})

	t.Run("Test validate update request", func(t *testing.T) {
		t.Parallel()
		req := pb.UpdateAuthServiceUserRequest{
			Id:        1,
			Username:  "pizza",
			Password1: "PasswOrd1@",
			Password2: "PasswOrd1@",
		}
		err := utils.ValidateRegisterUser(&req)
		assert.NoError(t, err)
	})

}
