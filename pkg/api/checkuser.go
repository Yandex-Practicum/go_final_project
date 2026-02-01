package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/JKasus/go_final_project/pkg/config"
	"github.com/JKasus/go_final_project/pkg/entities"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func checkUser(w http.ResponseWriter, r *http.Request) {
	var userData entities.UserData
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err = errors.New("Error reading body: " + err.Error())
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &userData)
	if err != nil {
		err = errors.New("Error unmarshalling body: " + err.Error())
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	cfg, err := config.NewConfig()
	if err != nil {
		err = errors.New("Error loading config: " + err.Error())
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	if cfg.Password == userData.Password {
		pass := []byte(userData.Password)
		jwtToken := jwt.New(jwt.SigningMethodHS256)
		signedToken, err := jwtToken.SignedString(pass)
		if err != nil {
			err = errors.New("Error signing token: " + err.Error())
			writeJSON(w, http.StatusInternalServerError, err)
			return
		}

		userData.Token = &signedToken
		writeJSON(w, http.StatusOK, &userData)
	} else {
		err = errors.New("Error verifying token: " + err.Error())
		writeJSON(w, http.StatusUnauthorized, err)
		return
	}
}
