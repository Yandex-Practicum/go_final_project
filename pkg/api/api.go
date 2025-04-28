package api

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func Init() {
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", auth(nextDateHandler))
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// auth выполняет аутентификацию запросов с помощью JWT
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
			tokenStr := cookie.Value
			// Метод проверки:
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(pass), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
