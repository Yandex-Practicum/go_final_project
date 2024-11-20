package authentification

import (
	"crypto/sha256"
	"encoding/json"
	"final_project/internal/common"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type password struct {
	Password string `json:"password"`
}

type token struct {
	Token string `json:"token"`
}

func Sign(w http.ResponseWriter, r *http.Request) {
	pswrd := &password{}
	json.NewDecoder(r.Body).Decode(pswrd)
	if pswrd.Password != common.Password {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "неверый пароль"})
		return
	}

	key := []byte("key")
	hashPassword := sha256.Sum256([]byte(pswrd.Password))
	claims := jwt.MapClaims{"password": hashPassword}
	jwtUns := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSign, err := jwtUns.SignedString(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Response{Error: "ошибка создания токена аутентификации, попробуйте еще раз " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(token{Token: jwtSign})

}
