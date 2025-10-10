package api

import "net/http"

// Init регистрирует все API обработчики
func Init() {
	// Публичные endpoints (без аутентификации)
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/signin", SingInHandler)

	// Защищенные endpoints (с аутентификацией)
	http.HandleFunc("/api/task", AuthMiddleware(TaskHandler))
	http.HandleFunc("/api/tasks", AuthMiddleware(TasksHandler))
	http.HandleFunc("/api/task/done", AuthMiddleware(DoneTaskHandler))
}
