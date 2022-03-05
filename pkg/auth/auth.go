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
	tokenKey                  = "token-%d-%d-%s"
	redisDelError             = "something went wrong"
	tokenNotFoundError        = "token not found"
	tokenNotValidError        = "token is not valid"
	unexpectedSigingMethod    = "unexpected token signing method"
	redisConnectionError      = "cannot connect to redis"
	redisPONG                 = "PONG"
	redisOK                   = "OK"
	redisUnknownResponseError = "unknown redis response: %s"
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

var ctx = context.Background()

func NewAuthenticator(c *config.Config) (*Authenticator, error) {
	rdc := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		DB:   int(c.Redis.DB),
	})

	stat, err := rdc.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	if stat != redisPONG {
		return nil, errors.New(redisConnectionError)
	}

	return &Authenticator{cfg: c, rdc: rdc}, nil
}

func (a *Authenticator) GenerateTokenPairs(authServiceID, usersServiceID int64, username string) (string, string, error) {

	accessToken, accessTokenUID, err := a.generateToken(authServiceID, usersServiceID, a.cfg.AuthService.ExpirationTimeAccess, username, a.cfg.AuthService.AccessIssuer)
	if err != nil {
		return "", "", err
	}

	refreshToken, refreshTokenUID, err := a.generateToken(authServiceID, usersServiceID, a.cfg.AuthService.ExpirationTimeRefresh, username, a.cfg.AuthService.RefreshIssuer)
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

	a.SetTokensInRedis(authServiceID, usersServiceID, username, string(cacheJSON))

	return refreshToken, accessToken, nil
}

func (a *Authenticator) generateToken(authServiceID, usersServiceID, exp int64, username, issuer string) (string, string, error) {

	expTime := time.Now().Add(time.Minute * time.Duration(exp))
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

func (a *Authenticator) SetTokensInRedis(authServiceID, usersServiceID int64, username, cacheJSON string) error {
	res, err := a.rdc.Set(ctx, fmt.Sprintf(tokenKey, authServiceID, usersServiceID, username), cacheJSON, time.Hour*time.Duration(a.cfg.AuthService.AutoLogoff)).Result()
	if err != nil {
		return err
	}

	if res != redisOK {
		return fmt.Errorf(redisUnknownResponseError, res)
	}

	return nil
}

func (a *Authenticator) GetTokensFromRedis(authServiceID, usersServiceID int64, username string) (string, error) {
	cacheJSON, err := a.rdc.Get(ctx, fmt.Sprintf(tokenKey, authServiceID, usersServiceID, username)).Result()
	if err != nil {
		return "", err
	}

	return cacheJSON, nil
}

func (a *Authenticator) ValidateToken(tokenString string, isRefresh bool) error {
	claims, err := a.ParseToken(tokenString)
	if err != nil {
		return err
	}

	cacheJSON, err := a.GetTokensFromRedis(claims.AuthServiceID, claims.UsersServiceID, claims.Username)
	if err != nil {
		return err
	}

	cachedTokens := new(CachedTokens)
	err = json.Unmarshal([]byte(cacheJSON), &cachedTokens)
	if err != nil {
		return err
	}

	var tokenUID string
	if isRefresh {
		tokenUID = cachedTokens.RefreshUID
	} else {
		tokenUID = cachedTokens.AccessUID
	}

	if tokenUID != claims.Uid {
		return errors.New(tokenNotFoundError)
	}

	return nil
}

func (a *Authenticator) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(unexpectedSigingMethod)
		}

		return []byte(a.cfg.AuthService.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(tokenNotValidError)
}

func (a *Authenticator) DeleteToken(authServiceID, usersServiceID int64, username string) error {
	intCmd := a.rdc.Del(ctx, fmt.Sprintf(tokenKey, authServiceID, usersServiceID, username))
	if intCmd.Val() != 1 {
		return errors.New(redisDelError)
	}

	return nil
}

func (a *Authenticator) ExpireToken(authServiceID, usersServiceID int64, username string) {
	a.rdc.Expire(ctx, fmt.Sprintf(tokenKey, authServiceID, usersServiceID, username), time.Hour*time.Duration(a.cfg.AuthService.AutoLogoff))
}
