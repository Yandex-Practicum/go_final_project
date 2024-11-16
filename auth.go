package main

import (
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type errornew struct {
	Error string `json:"error"`
}

type password struct {
	Password string `json:"password"`
}

type token struct {
	Token string `json:"token"`
}

func sign(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errornew{Error: "по URI api/sign доступен только POST-запрос"})
		return
	}
	pswrd := &password{}
	json.NewDecoder(r.Body).Decode(pswrd)
	if pswrd.Password != os.Getenv("TODO_PASSWORD") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errornew{Error: "неверый пароль"})
		return
	}

	key := []byte("key")
	hashPassword := sha256.Sum256([]byte(pswrd.Password))
	claims := jwt.MapClaims{"password": hashPassword}
	jwtUns := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSign, err := jwtUns.SignedString(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errornew{Error: "ошибка создания токена аутентификации, попробуйте еще раз " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(token{Token: jwtSign})

}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errornew{Error: "ошибка аутентификации"})
			return
		}
		key := []byte("key")
		hashPassword := sha256.Sum256([]byte(os.Getenv("TODO_PASSWORD")))
		claims := jwt.MapClaims{"password": hashPassword}
		jwtUns := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		jwtSign, err := jwtUns.SignedString(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errornew{Error: "ошибка проверки подлинности токена аутентификации, попробуйте еще раз " + err.Error()})
			return
		}
		if tok.Value == jwtSign {
			next(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errornew{Error: "ошибка аутентификации, попробуйте еще раз "})
	}
}
