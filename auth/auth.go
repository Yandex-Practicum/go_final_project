package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type PasswordRequest struct {
	Password string `json:"password"`
}

type Claims struct {
	jwt.StandardClaims
}

func CreateJWT() (string, error) {
	secret := []byte(os.Getenv("SECRET"))

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("error signing jwt secret: %w", err)
	}
	return signedToken, nil
}

func ValidateJWT(tokenStr string) (bool, error) {
	secret := os.Getenv("SECRET")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected sign in method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return false, errors.New("invalid token")
	}
	return true, nil
}
