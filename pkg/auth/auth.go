package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hyperxpizza/auth-service/pkg/config"
)

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type Authenticator struct {
	cfg *config.Config
}

func NewAuthenticator(c *config.Config) *Authenticator {
	return &Authenticator{cfg: c}
}

func (a *Authenticator) GenerateToken(id int64, username string) (string, error) {
	//var token string

	expTime := time.Now().Add(time.Hour * time.Duration(a.cfg.AuthService.ExpirationTimeHours))

	claims := Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Audience:  a.cfg.AuthService.Audience,
			ExpiresAt: expTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    a.cfg.AuthService.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.cfg.AuthService.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *Authenticator) ValidateToken(tokenString string) (string, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token signing method is not valid: %v", token.Header["alg"])
		}

		return []byte(a.cfg.AuthService.JWTSecret), nil
	})
	if err != nil {
		return "", 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["id"]
		username := claims["username"]
		return username.(string), id.(int64), nil
	}

	return "", 0, err
}
