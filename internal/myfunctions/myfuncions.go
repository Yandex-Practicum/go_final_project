package myfunctions

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"todo_restapi/internal/dto"
	"todo_restapi/pkg/constants"
)

func ValidateJWT(request *http.Request) error {

	cookie, err := request.Cookie("token")
	if err != nil {
		return errors.New("token not found")
	}

	tokenString := cookie.Value

	secret, exists := os.LookupEnv("TODO_SECRET")
	if !exists {
		return errors.New("TODO_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return errors.New("token expired")
		}
	} else {
		return errors.New("missing exp claim")
	}

	hashFromToken, ok := claims["pwd"].(string)
	if !ok {
		return errors.New("missing password hash in token")
	}

	storedPassword, exists := os.LookupEnv("TODO_PASSWORD")
	if !exists {
		return errors.New("TODO_PASSWORD is not set")
	}

	hash := sha256.New()
	hash.Write([]byte(storedPassword))
	hashPassword := hex.EncodeToString(hash.Sum(nil))

	if hashFromToken != hashPassword {
		return errors.New("invalid token: password hash mismatch")
	}

	return nil
}

func PwdValidateGenerateJWT(password string) (string, error) {

	secret, exists := os.LookupEnv("TODO_SECRET")
	if !exists {
		return "", errors.New("TODO_SECRET is not set")
	}

	storedPassword, exists := os.LookupEnv("TODO_PASSWORD")
	if !exists {
		return "", errors.New("TODO_PASSWORD is not set")
	}

	if password != storedPassword {
		return "", errors.New("invalid password")
	}

	hash := sha256.New()
	hash.Write([]byte(password))
	hashPassword := hex.EncodeToString(hash.Sum(nil))

	payload := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 8).Unix(),
		"pwd": hashPassword,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("cannot sign JWT: %w", err)
	}

	return signedToken, nil
}

func IsDate(searchQuery string) (string, error) {

	isTime, err := time.Parse("02.01.2006", searchQuery)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	return isTime.Format(constants.DateFormat), nil
}

func ValidateTaskRequest(newTask *models.Task, now string) error {

	if newTask.Title == "" {
		return errors.New("title is empty")
	}

	if newTask.Date == "" {
		newTask.Date = now
	}

	_, err := time.Parse(constants.DateFormat, newTask.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	if newTask.Date < now {
		if newTask.Repeat == "" {
			newTask.Date = now
		} else {
			dateCalculation, err := NextDate(time.Now(), newTask.Date, newTask.Repeat)
			if err != nil {
				return fmt.Errorf("NextDate: function error: %w", err)
			}
			newTask.Date = dateCalculation
		}
	}
	return nil
}

func WriteJSONError(write http.ResponseWriter, statusCode int, errMsg string) {

	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(statusCode)

	response := map[string]string{"error": errMsg}

	if err := json.NewEncoder(write).Encode(response); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

func parseNumbers(input string) []int {

	stringNums := strings.Split(input, ",")
	output := make([]int, 0, len(stringNums))

	for _, num := range stringNums {
		intNum, err := strconv.Atoi(num)
		if err != nil {
			fmt.Println("parseNumbers: string to int conversion error")
		}
		output = append(output, intNum)
	}
	return output
}

func parseRepeat(repeat string) (string, []int, []int) {

	var firstRepeatPattern []int
	var secondRepeatPattern []int

	repeatParse := strings.Fields(repeat)

	if len(repeatParse) > 1 {
		firstRepeatPattern = parseNumbers(repeatParse[1])
	}

	if len(repeatParse) > 2 {
		secondRepeatPattern = parseNumbers(repeatParse[2])
	}

	repeatType := repeatParse[0]

	return repeatType, firstRepeatPattern, secondRepeatPattern
}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("repeat cannot be empty")
	}

	dateParse, err := time.Parse(constants.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("date parse error: %w", err)
	}

	repeatType, firstRepeatPattern, _ := parseRepeat(repeat)

	switch repeatType {

	case "d":

		if len(firstRepeatPattern) == 0 {
			return "", errors.New("\"d\" parameter is empty")
		}

		if firstRepeatPattern[0] > 400 {
			return "", errors.New("invalid \"d\" value (400 is max)")
		}

		dateParse = dateParse.AddDate(0, 0, firstRepeatPattern[0])
		for dateParse.Before(now) {
			dateParse = dateParse.AddDate(0, 0, firstRepeatPattern[0])
		}

		return dateParse.Format(constants.DateFormat), nil

	case "y":

		dateParse = dateParse.AddDate(1, 0, 0)
		for dateParse.Before(now) {
			dateParse = dateParse.AddDate(1, 0, 0)
		}

		return dateParse.Format(constants.DateFormat), nil

	default:
		return "", errors.New("invalid repeat value")
	}
}
