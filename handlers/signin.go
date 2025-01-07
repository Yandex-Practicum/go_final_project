package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var request struct {
	Password string `json:"password"`
}

var tokenResponse struct {
	Token string `json:"token"`
}

// SigninHandler is a handler for "/api/signin" endpoint.
// It expects a POST request with a JSON object with a single field:
// - password: a string representing the password to be verified
// If the password is invalid, it returns an error with HTTP status code 401.
// If the password is valid, it returns a JSON object with the following field:
// - token: a string representing the generated JWT token
// It returns the following HTTP status codes:
// - 200 OK: the password was successfully verified
// - 400 Bad Request: the request body is invalid
// - 401 Unauthorized: the password is invalid
// - 500 Internal Server Error: an error occurred while generating the token
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, `{"error":"can't read body"}`, http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		http.Error(w, `{"error":"can't unmarshal body"}`, http.StatusBadRequest)
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")
	if storedPassword == "" {
		http.Error(w, `{"error":"password is empty"}`, http.StatusInternalServerError)
		return
	}

	if request.Password != storedPassword {
		http.Error(w, `{"error":"invalid password"}`, http.StatusUnauthorized)
		return
	}

	tokenResponse.Token, err = generateTokenJWT(storedPassword)
	if err != nil {
		http.Error(w, `{"error":"can't generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(tokenResponse); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
}

// generateTokenJWT generates a JWT token using the given password.
//
// The generated token will contain the password hash and an expiration time
// set to 8 hours from now. If there is an error while generating the token,
// it will be returned in the second return value.
func generateTokenJWT(password string) (string, error) {
	claims := jwt.MapClaims{
		"password_hash": fmt.Sprintf("%x", hashPassword(password)),
		"exp":           time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("TODO_PASSWORD"))
	return token.SignedString(secretKey)
}

// hashPassword converts the given password string into a byte slice.
//
// The function takes a string parameter representing the password and returns
// a byte slice. This conversion is necessary for handling password data in
// a format suitable for cryptographic operations.
func hashPassword(password string) []byte {
	return []byte(password)
}

// isValidTokenJWT verifies the validity of a given JWT token using the provided password.
//
// The function takes two string parameters: a JWT token and a password.
// It parses the token using the password as the key and checks its validity.
// If the token is valid, it compares the password hash stored in the token's claims
// with the hash of the provided password. If these hashes match, the function returns true,
// indicating that the token is valid. Otherwise, it returns false.
//
// Returns:
// - true if the token is valid and the password hash matches the stored hash.
// - false if the token is invalid, parsing fails, or the password hash does not match.
func isValidTokenJWT(token, password string) bool {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(password), nil
	})
	if err != nil || !parsedToken.Valid {
		return false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	storedHash := claims["password_hash"].(string)
	return storedHash == fmt.Sprintf("%x", hashPassword(password))
}

// Authorization is a middleware function that checks for the presence of a valid JWT
// token in the request. If the token is valid, it will call the next handler in the
// chain. If the token is invalid or missing, it will return a 401 Unauthorized response.
//
// The function expects the TODO_PASSWORD environment variable to be set with the
// password to be verified. If the password is not set, the function will not check
// for the token.
//
// The function will look for a cookie named "token" in the request. If the cookie
// is not present, it will not check for the token. If the cookie is present, it will
// verify the token using the provided password. If the token is valid, it will call
// the next handler in the chain. If the token is invalid, it will return a 401
// Unauthorized response.
func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			var tokenJWT string

			cookie, err := r.Cookie("token")
			if err != nil {
				tokenJWT = cookie.Value
			}

			if tokenJWT == "" || !isValidTokenJWT(tokenJWT, password) {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
