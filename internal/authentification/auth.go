package authentification

import (
	"crypto/sha256"
	"encoding/json"
	"final_project/internal/common"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(common.Response{Error: "ошибка аутентификации"})
			return
		}
		key := []byte("key")
		hashPassword := sha256.Sum256([]byte(common.Password))
		claims := jwt.MapClaims{"password": hashPassword}
		jwtUns := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		jwtSign, err := jwtUns.SignedString(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(common.Response{Error: "ошибка проверки подлинности токена аутентификации, попробуйте еще раз " + err.Error()})
			return
		}
		if tok.Value == jwtSign {
			next(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(common.Response{Error: "ошибка аутентификации, попробуйте еще раз "})
	}
}
