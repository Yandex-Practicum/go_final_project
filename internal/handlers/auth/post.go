package auth

import (
	"encoding/json"
	"hash/fnv"
	"net/http"
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
		utils.RespondWithError(w, utils.ErrInvalidJson)
		return
	}

	if strings.TrimSpace(signDTO.Password) != h.password {
		utils.RespondWithError(w, utils.ErrInvalidPassword)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash(signDTO.Password),
	})

	tokenString, err := token.SignedString([]byte(utils.AuthSecret))
	if err != nil {
		utils.RespondWithError(w, utils.ErrTokenCreate)
		return
	}

	response := SignResponseDTO{Token: tokenString}
	utils.Respond(w, response)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
