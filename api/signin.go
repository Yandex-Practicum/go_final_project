package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func PostSigninHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var err error
	var body map[string]string
	targetPassword := os.Getenv("TODO_PASSWORD")
	var password string
	var signedToken string

	write := func() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var resp []byte
		if err != nil {
			writeErr(err, w)
			return
		} else {
			tokenResp := map[string]string{
				"token": signedToken,
			}
			resp, err = json.Marshal(tokenResp)
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(resp)
			if err != nil {
				log.Println(err)
			}
			return
		}

	}

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		write()
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &body); err != nil {
		write()
		return
	}
	if len(body["password"]) == 0 {
		err = fmt.Errorf("пустая строка вместо password")
		write()
		return
	} else {
		password = body["password"]
	}
	if password != targetPassword {
		err = fmt.Errorf("неправильный пароль")
		write()
		return
	}

	claims := jwt.MapClaims{
		"password": sha256.Sum256([]byte(targetPassword)),
		"Exp":      1550946689,
	}

	// создаём jwt и указываем алгоритм хеширования
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// получаем подписанный токен
	signedToken, err = jwtToken.SignedString([]byte(targetPassword))
	if err != nil {
		write()
		return
	}

	write()

}
