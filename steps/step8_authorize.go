package steps

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type P struct {
	Password string `json:"password,omitempty"`
}

func auth(w http.ResponseWriter, r *http.Request) {
	envPassword, exists := os.LookupEnv("TODO_PASSWORD")

	if len(envPassword) == 0 || !exists {
		writeError(w, Token{Error: "не определён пароль в переменной окружения"})
		return
	}

	r.Method = http.MethodPost
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, Token{Error: err.Error()})
		return
	}

	p := P{}

	if err = json.Unmarshal(buf.Bytes(), &p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, Task{Error: err.Error()})
		return
	}

	password := p.Password
	tok := Token{}

	if password == envPassword {
		var secretKey = []byte("secret")
		hash := sha256.Sum256([]byte(password))
		claims := jwt.MapClaims{
			"hash": hash}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(secretKey)

		if err != nil {
			writeError(w, Token{Error: err.Error()})
			return
		}

		tok.Token = signedToken
		writeError(w, Token{Token: tok.Token})
		return

	}

	writeError(w, Token{Error: "неверный пароль"})
}

func authTask(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, exists := os.LookupEnv("TODO_PASSWORD")

		if !exists {
			writeError(w, Token{Error: "не определён пароль в переменной окружения"})
			return
		}

		password := os.Getenv("TODO_PASSWORD")

		if len(password) > 0 {
			var jwtCookie string
			cookie, err := r.Cookie("token")

			if err == nil {
				jwtCookie = cookie.Value
			}

			var secretKey = []byte("secret")
			jwtToken, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (any, error) {
				return (secretKey), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				writeError(w, Token{Error: err.Error()})
				return
			}

			if !jwtToken.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				writeError(w, Token{Error: "ошибка аутентификации"})
				return
			}

			payload, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				writeError(w, Token{Error: "ошибка проверки JWT-токена"})
				return
			}

			hashRaw := payload["hash"]

			hashOK, ok := hashRaw.([]interface{})
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				writeError(w, Token{Error: "ошибка проверки JWT-токена"})
				return
			}

			hashPassFromToken := []byte(fmt.Sprint(hashOK))

			hashIn := sha256.Sum256([]byte(password))
			slice := hashIn[:]
			hashPass := []byte(fmt.Sprint(slice[:]))

			if !bytes.Equal(hashPassFromToken, hashPass) {
				w.WriteHeader(http.StatusUnauthorized)
				writeError(w, Token{Error: "ошибка аутентификации"})
				return
			}

		}

		next.ServeHTTP(w, r)
	})
}
