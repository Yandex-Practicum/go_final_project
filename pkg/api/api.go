package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var pass = os.Getenv("TODO_PASSWORD")

// Init Регистрация маршрутов
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", TaskDoneHandler)
	http.HandleFunc("/api/signin", signinHandler)

}

const Layout = "20060102"

// Обработчик NextDate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(Layout, nowStr)
		if err != nil {
			http.Error(w, "Invalid format", http.StatusBadRequest)
			return
		}
	}

	//Без даты - не пускаем
	if dstart == "" {
		http.Error(w, "The date field cannot be empty.", http.StatusBadRequest)
		return
	}

	//Без правила повторения - тоже
	if repeat == "" {
		http.Error(w, "The repeat field cannot be empty.", http.StatusBadRequest)
		return
	}

	//Вычисление следующей даты
	result, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(result))
}

func signinHandler(w http.ResponseWriter, r *http.Request) {

	//Если пустой - лови 403
	if pass == "" {
		http.Error(w, "Authentication disabled", http.StatusForbidden)
		return
	}

	var req struct {
		Password string `json:"password"`
	}

	//Декодируем
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//Неправильный пароль - держи 401
	if req.Password != pass {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(map[string]string{"error": "Wrong password"})
		if err != nil {
			return
		}
		return
	}

	//Правильный пароль - собираем токен Джейсона
	jwtKey := []byte(pass)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pwd_hash": req.Password,
		"exp":      jwt.NewNumericDate(jwt.TimeFunc().Add(8 * 3600e9)),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	tokenS := map[string]any{"token": tokenString}

	if err := json.NewEncoder(w).Encode(tokenS); err != nil {
		log.Printf("Encode response: %v", err)
	}
}
