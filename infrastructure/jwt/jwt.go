package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWT represents a JSON Web Token.
type JWT struct {
	SecretKey string // SecretKey is the secret key used for signing and verifying the token.
}

// New creates a new instance of JWT with the provided secret key.
func New(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

// CreateToken creates a new JWT token using the provided secret key.
// It returns the token string and any error encountered during the process.
func (j *JWT) CreateToken() (string, error) {
	tokenInstance := jwt.New(jwt.SigningMethodHS512)
	tokenString, err := tokenInstance.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", fmt.Errorf("tokenInstance.SignedString: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates the given JWT token string.
// It parses the token using the secret key stored in the JWT instance.
// If the token is valid, it returns nil. Otherwise, it returns an error.
func (j *JWT) ValidateToken(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return fmt.Errorf("jwt.Parse: %w", err)
	}

	return nil
}
