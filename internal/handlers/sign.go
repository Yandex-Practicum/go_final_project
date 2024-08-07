package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"hash/fnv"
	"net/http"
	"os"
	"strings"
)

const envPassword = "TODO_PASSWORD"
const secret = "kjgsdfi632riusvd732fikd2--2!"

type SignDTO struct {
	Password string `json:"password"`
}

type SignResponseDTO struct {
	Token string `json:"token"`
}

func SignHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleSign(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}

func handleSign(w http.ResponseWriter, r *http.Request) {
	var signDTO SignDTO
	err := json.NewDecoder(r.Body).Decode(&signDTO)
	if err != nil {
		respondWithError(w, "Ошибка десериализации JSON")
		return
	}

	password := os.Getenv(envPassword)
	if strings.TrimSpace(signDTO.Password) != password {
		respondWithError(w, "Неверный пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash(signDTO.Password),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		respondWithError(w, "Ошибка генерации токена")
		return
	}

	response := SignResponseDTO{Token: tokenString}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func IsTokenValid(tokenString string) bool {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	return err == nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
