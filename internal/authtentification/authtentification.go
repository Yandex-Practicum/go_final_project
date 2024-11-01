package authtentification

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Pass struct {
	Password string `json:"password"`
}

type Claims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("secret")

var password string

func init() {
	environment := os.Getenv("TODO_PASSWORD")
	if len(environment) > 0 {
		password = environment
	}
}

func Authentification(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(password) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
				return
			}

			jwtToken := cookie.Value

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
				return
			}

			calcHash := sha256.Sum256([]byte(jwtKey))
			expectedHash := hex.EncodeToString(calcHash[:])
			if claims.Hash != expectedHash {
				http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func HandleSiginingIn(w http.ResponseWriter, r *http.Request) {
	var pass Pass
	err := json.NewDecoder(r.Body).Decode(&pass)
	if err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	if password == "" || pass.Password != password {
		http.Error(w, `{"error":"Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(8 * time.Hour)
	calcHash := sha256.Sum256([]byte(jwtKey))
	calcHashString := hex.EncodeToString(calcHash[:])
	claims := &Claims{
		Hash: calcHashString,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "ошибка подписи токена", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
