package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hyperxpizza/auth-service/pkg/config"
)

const (
	TokenKey      = "token-%d-%d-%s"
	RedisDelError = "Something went wrong"
)

type Claims struct {
	AuthServiceID  int64  `json:"authServiceID"`
	UsersServiceID int64  `json:"usersServiceID"`
	Uid            string `json:"uid"`
	Username       string `json:"username"`
	jwt.StandardClaims
}

type CachedTokens struct {
	AccessUID  string `json:"access"`
	RefreshUID string `json:"refresh"`
}

type Authenticator struct {
	cfg *config.Config
	rdc *redis.Client
}

func NewAuthenticator(c *config.Config) *Authenticator {
	rdc := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Network:  c.Redis.Network,
		Password: c.Redis.Password,
		DB:       int(c.Redis.DB),
	})

	return &Authenticator{cfg: c, rdc: rdc}
}

func (a *Authenticator) GenerateTokenPairs(authServiceID, usersServiceID, exp int64, username, issuer string) (string, string, error) {

	accessToken, accessTokenUID, err := a.generateToken(authServiceID, usersServiceID, exp, username, a.cfg.AuthService.AccessIssuer)
	if err != nil {
		return "", "", err
	}

	refreshToken, refreshTokenUID, err := a.generateToken(authServiceID, usersServiceID, exp, username, a.cfg.AuthService.RefreshIssuer)
	if err != nil {
		return "", "", err
	}

	cacheJSON, err := json.Marshal(CachedTokens{
		AccessUID:  accessTokenUID,
		RefreshUID: refreshTokenUID,
	})
	if err != nil {
		return "", "", err
	}

	a.rdc.Set(context.Background(), fmt.Sprintf(TokenKey, authServiceID, usersServiceID, username), string(cacheJSON), time.Hour*time.Duration(a.cfg.AuthService.AutoLogoff))

	return refreshToken, accessToken, nil
}

func (a *Authenticator) generateToken(authServiceID, usersServiceID, exp int64, username, issuer string) (string, string, error) {

	expTime := time.Now().Add(time.Hour * time.Duration(exp))
	uid := uuid.New().String()

	claims := Claims{
		AuthServiceID:  authServiceID,
		UsersServiceID: usersServiceID,
		Username:       username,
		Uid:            uid,
		StandardClaims: jwt.StandardClaims{
			Audience:  a.cfg.AuthService.Audience,
			ExpiresAt: expTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.cfg.AuthService.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return tokenString, uid, nil
}

func (a *Authenticator) ValidateToken(tokenString string) (username string, authServiceID int64, usersServiceID int64, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token signing method is not valid: %v", token.Header["alg"])
		}

		return []byte(a.cfg.AuthService.JWTSecret), nil
	})
	if err != nil {
		return "", 0, 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		authSerivceID := claims["authServiceID"]
		authSerivceIDFloat := authSerivceID.(float64)

		usersServiceID := claims["usersServiceID"]
		usersServiceIDFloat := usersServiceID.(float64)

		username := claims["username"]
		return username.(string), int64(authSerivceIDFloat), int64(usersServiceIDFloat), nil
	}

	return "", 0, 0, err
}

func (a *Authenticator) DeleteToken(authServiceID, usersServiceID int64, username string) error {
	intCmd := a.rdc.Del(context.Background(), fmt.Sprintf(TokenKey, authServiceID, usersServiceID, username))
	if intCmd.Val() != 1 {
		return errors.New(RedisDelError)
	}

	return nil
}
