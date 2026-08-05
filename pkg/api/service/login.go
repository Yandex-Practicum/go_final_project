package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	PasswordHash string `json:"passwordhash"`
}

var Secret string = "YandexPracticum"

func LoginAuther(w http.ResponseWriter, r *http.Request) ([]byte, error, int) {
	type passwordType struct {
		Password string `json:"password"`
	}
	var pass passwordType

	type TokenJWT struct {
		Token string `json:"token"`
	}

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	err = json.Unmarshal(bodyByte, &pass)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if pass.Password != AutherENV() {
		// SendErrorJSON(w, "Неверный пароль", http.StatusBadRequest)
		return nil, fmt.Errorf("Неверный пароль"), http.StatusUnauthorized
	} else {
		hashedPassword := sha256.Sum256([]byte(pass.Password))
		hashStringPassword := hex.EncodeToString(hashedPassword[:])

		var claims Claims
		claims.PasswordHash = hashStringPassword

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(Secret))
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}

		var tokenJWT TokenJWT

		tokenJWT.Token = tokenString

		jsonByte, err := json.Marshal(tokenJWT)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}

		return jsonByte, nil, http.StatusOK
	}

}
