package impl

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/hyperxpizza/auth-service/pkg/auth"
	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/hyperxpizza/auth-service/pkg/database"
	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	"github.com/hyperxpizza/auth-service/pkg/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	UserNotFoundError         = "user not found in the database"
	PasswordsNotMatchingError = "passwords do not match"
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

	authenticator, err := auth.NewAuthenticator(cfg)
	if err != nil {
		return nil, err
	}

	return &AuthServiceServer{
		cfg:           cfg,
		logger:        logger,
		db:            db,
		authenticator: authenticator,
	}, nil
}

func (a *AuthServiceServer) Run() error {

	cert, err := tls.LoadX509KeyPair(a.cfg.TLS.CertPath, a.cfg.TLS.KeyPath)
	if err != nil {
		a.logger.Error("cound not load X509 key pair")
		return err
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, a)

	addr := fmt.Sprintf(":%d", a.cfg.AuthService.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		a.logger.Error("net.Listen failed: %s", err.Error())
		return err
	}

	a.logger.Infof("auth service server running on %s:%d", a.cfg.AuthService.Host, a.cfg.AuthService.Port)

	if err := grpcServer.Serve(lis); err != nil {
		a.logger.Error("failed to serve: %s", err.Error())
		return err
	}

	return nil
}

func (a *AuthServiceServer) GenerateToken(ctx context.Context, req *pb.TokenRequest) (*pb.Tokens, error) {
	a.logger.Infof("generating token for: %s", req.Username)
	var tokenResponse pb.Tokens

	//check if user exists in the database
	user, err := a.db.GetUserByUsersServiceID(req.UsersServiceID, req.Username)
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

	refreshToken, accessToken, err := a.authenticator.GenerateTokenPairs(user.Id, req.UsersServiceID, req.Username)
	if err != nil {
		a.logger.Infof("generating jwt token for: %s with id: %d failed: %s", req.Username, req.UsersServiceID, err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	tokenResponse.RefreshToken = refreshToken
	tokenResponse.AccessToken = accessToken

	return &tokenResponse, nil
}

func (a *AuthServiceServer) ValidateToken(ctx context.Context, token *pb.AccessTokenData) (*pb.TokenData, error) {
	var tokenData pb.TokenData

	claims, err := a.authenticator.ParseToken(token.AccessToken)
	if err != nil {
		a.logger.Infof("parsing token failed: %s", err.Error())
		return nil, status.Error(
			codes.Unauthenticated,
			err.Error(),
		)
	}

	err = a.authenticator.ValidateToken(token.AccessToken, false)
	if err != nil {
		a.logger.Infof("validating token failed: %s", err.Error())
		return nil, status.Error(
			codes.Unauthenticated,
			err.Error(),
		)
	}

	tokenData.Username = claims.Username
	tokenData.AuthServiceID = claims.AuthServiceID
	tokenData.UsersServiceID = claims.UsersServiceID

	a.logger.Infof("token for: %s valid", claims.Username)
	go a.authenticator.ExpireToken(claims.AuthServiceID, claims.UsersServiceID, claims.Username)

	return &tokenData, nil
}

func (a *AuthServiceServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenData) (*pb.Tokens, error) {
	var tokens pb.Tokens

	claims, err := a.authenticator.ParseToken(req.RefreshToken)
	if err != nil {
		a.logger.Infof("parsing token failed: %s")
		return nil, status.Error(
			codes.Unauthenticated,
			err.Error(),
		)
	}

	err = a.authenticator.ValidateToken(req.RefreshToken, true)
	if err != nil {
		a.logger.Infof("validating token failed: %s", err.Error())
		return nil, status.Error(
			codes.Unauthenticated,
			err.Error(),
		)
	}

	refreshToken, accessToken, err := a.authenticator.GenerateTokenPairs(claims.AuthServiceID, claims.UsersServiceID, claims.Username)
	if err != nil {
		a.logger.Infof("generating jwt token for: %s with id: %d failed: %s", claims.Username, claims.UsersServiceID, err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken
	a.logger.Infof("refreshed token for: %s", claims.Username)

	return &tokens, nil
}

func (a *AuthServiceServer) DeleteTokens(ctx context.Context, req *pb.TokenData) (*emptypb.Empty, error) {
	err := a.authenticator.DeleteToken(req.AuthServiceID, req.UsersServiceID, req.Username)
	if err != nil {
		a.logger.Infof("deleting token for: %s failed: %s", req.Username, err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	a.logger.Infof("deleted token for: %s", req.Username)

	return &emptypb.Empty{}, nil
}

func (a *AuthServiceServer) AddUser(ctx context.Context, req *pb.AuthServiceUserRequest) (*pb.AuthServiceID, error) {
	var id pb.AuthServiceID

	a.logger.Infof("adding user: %s into the database", req.Username)

	err := utils.ValidateRegisterUser(req)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	passwordHash, err := utils.GeneratePasswordHash(req.Password1)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	user := pb.AuthServiceUser{
		Username:              req.Username,
		PasswordHash:          passwordHash,
		RelatedUsersServiceID: req.RelatedUsersServiceID,
	}

	idInt, err := a.db.InsertUser(&user)
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

func (a *AuthServiceServer) RemoveUser(ctx context.Context, id *pb.AuthServiceID) (*emptypb.Empty, error) {

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

func (a *AuthServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateAuthServiceUserRequest) (*emptypb.Empty, error) {
	a.logger.Infof("updating user with username: %d", req.Username)

	err := utils.ValidateRegisterUser(req)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	passwordHash, err := utils.GeneratePasswordHash(req.Password1)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	user := pb.AuthServiceUser{
		Id:           req.Id,
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	err = a.db.UpdateUser(&user)
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

func (a *AuthServiceServer) ChangePassword(ctx context.Context, req *pb.PasswordRequest) (*emptypb.Empty, error) {

	if req.Password1 != req.Password2 {
		return nil, status.Error(
			codes.InvalidArgument,
			PasswordsNotMatchingError,
		)
	}

	passwordHash, err := utils.GeneratePasswordHash(req.Password1)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	err = a.db.ChangePassword(req.AuthServiceID, req.Username, passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(
				codes.NotFound,
				UserNotFoundError,
			)
		} else {
			return nil, status.Error(
				codes.Internal,
				err.Error(),
			)
		}
	}

	return &emptypb.Empty{}, nil

}
