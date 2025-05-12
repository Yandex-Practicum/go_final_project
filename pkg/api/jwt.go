package api

import "github.com/golang-jwt/jwt/v5"

const (
	secret = "vErY VeRy sEcReT InfOrMaTiOn"
)

func generateJWT(password string) (string, error) {
	claims := jwt.MapClaims{
		"user_password": password,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
