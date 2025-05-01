package server

import (
	"fmt"
	"github.com/AndreiCoder404/go_final_project/pkg/api"
	"github.com/AndreiCoder404/go_final_project/pkg/auth"
	"github.com/AndreiCoder404/go_final_project/pkg/db"
	"log"
	"net/http"
	"os"
	"strconv" // Добавляем импорт для преобразования строки в число
)

// recoveryMiddleware перехватывает паники и возвращает ошибку клиенту
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic in HTTP server: %v", err)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, `{"error":"internal server error"}`)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware добавляет заголовки CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// authMiddleware проверяет JWT-токен
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			valid, err := auth.ValidateJWT(cookie.Value, pass)
			if err != nil || !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Run запускает сервер
func Run() error {
	// Определяем путь к базе данных
	dbfile := "scheduler.db"
	if envDB := os.Getenv("TODO_DBFILE"); envDB != "" {
		dbfile = envDB
	}

	// Инициализируем базу данных
	err := db.Init(dbfile)
	if err != nil {
		return err
	}

	// Определяем порт
	port := 7540 // Порт по умолчанию
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		parsedPort, err := strconv.Atoi(envPort)
		if err != nil {
			return fmt.Errorf("invalid TODO_PORT value: %v", err)
		}
		if parsedPort < 1 || parsedPort > 65535 {
			return fmt.Errorf("TODO_PORT must be between 1 and 65535")
		}
		port = parsedPort
	}

	// Создаём мультиплексор
	mux := http.NewServeMux()

	// Регистрируем API-обработчики
	api.Init()

	// Регистрируем обработчик для статических файлов
	fileServer := http.FileServer(http.Dir("web"))
	mux.Handle("/", corsMiddleware(recoveryMiddleware(fileServer)))

	// Регистрируем обработчики API с middleware аутентификации
	mux.Handle("/api/task", authMiddleware(corsMiddleware(recoveryMiddleware(http.DefaultServeMux))))
	mux.Handle("/api/tasks", authMiddleware(corsMiddleware(recoveryMiddleware(http.DefaultServeMux))))
	mux.Handle("/api/task/done", authMiddleware(corsMiddleware(recoveryMiddleware(http.DefaultServeMux))))

	// Обработчик /api/signin без middleware аутентификации
	mux.HandleFunc("/api/signin", auth.SignInHandler)

	// Обработчик /api/nextdate без middleware аутентификации
	mux.Handle("/api/nextdate", corsMiddleware(recoveryMiddleware(http.DefaultServeMux)))

	// Запускаем сервер на указанном порту
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on port %d...", port)
	return http.ListenAndServe(addr, mux)
}
