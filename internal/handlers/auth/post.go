package auth

import (
	"encoding/json"
	"hash/fnv"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go_final_project/internal/utils"
)

type SignDTO struct {
	Password string `json:"password"`
}

type SignResponseDTO struct {
	Token string `json:"token"`
}

func (h *Handler) handleSign(w http.ResponseWriter, r *http.Request) {
	var signDTO SignDTO
	err := json.NewDecoder(r.Body).Decode(&signDTO)
	if err != nil {
		utils.RespondWithError(w, "Ошибка десериализации JSON")
		return
	}

	password := os.Getenv(envPassword)
	if strings.TrimSpace(signDTO.Password) != password {
		utils.RespondWithError(w, "Неверный пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash(signDTO.Password),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		utils.RespondWithError(w, "Ошибка генерации токена")
		return
	}

	response := SignResponseDTO{Token: tokenString}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
