package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(id, exp int64, username, issuer, aud, secret string) (string, error) {
	//var token string

	expTime := time.Now().Add(time.Hour * time.Duration(exp))

	claims := Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Audience:  aud,
			ExpiresAt: expTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) {

}
