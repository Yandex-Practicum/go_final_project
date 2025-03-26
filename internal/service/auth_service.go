package service

import (
	"fmt"
	"hash/fnv"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go_final_project/internal/constants"
)

type AuthService struct {
	password string
	secret   string
}

func NewAuthService() *AuthService {
	pass := os.Getenv(constants.EnvPassword)
	secret := os.Getenv(constants.EnvSecret)

	return &AuthService{password: pass, secret: secret}
}

func (s *AuthService) IsAuthEnabled() bool {
	return len(s.password) > 0
}

func (s *AuthService) IsTokenValid(tokenString string) bool {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	return err == nil
}

func (s *AuthService) SignIn(password string) (string, error) {
	if strings.TrimSpace(password) != s.password {
		return "", constants.ErrInvalidPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash(password),
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", constants.ErrTokenCreate
	}
	return tokenString, nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
