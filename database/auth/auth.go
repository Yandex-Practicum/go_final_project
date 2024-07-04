package auth

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	pass string = os.Getenv("TODO_PASSWORD")
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		if len(pass) > 0 {
			err := getAndVerifyToken(r)
			if err != nil {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

// verifyToken проверяет токен на подлинность, возвращает true если токен корректен
func verifyToken(signedToken string) bool {
	passByte := []byte(pass)
	passwordChecksum := sha256.Sum256(passByte)

	jwtToken, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
		return passByte, nil
	})
	if err != nil {
		log.Printf("Failed to parse token: %s\n", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	passRaw, ok := claims["password"]
	if !ok {
		return false
	}
	// костыль чтобы лениво преобразовать jwt.Claims password и password из .env к одному типу :')
	claimsPassChSm := fmt.Sprintf("%v", passRaw)
	envPassChSm := fmt.Sprintf("%v", passwordChecksum)

	return claimsPassChSm == envPassChSm

}

// getAndVerifyToken проверяет cookie на наличие токена авторизации, и проверяет его подлинность.
// Возвращает ошибку, если токен не найден, или токен не прошёл проверку.
func getAndVerifyToken(r *http.Request) error {
	token, err := r.Cookie("token")

	if err != nil {
		return err
	}
	if verifyToken(token.Value) {
		return nil
	}
	return fmt.Errorf("ошибка авторизации")
}
