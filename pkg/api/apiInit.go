package api

import "net/http"

// Функция для инициализации всех API маршрутов
func InitAPI() {
	http.HandleFunc("/api/nextdate", nextDayHandler) // Обработчик для вычисления следующей даты (без авторизации)
	http.HandleFunc("/api/signin", SignInHandler)    // Обработчик входа (без авторизации)

	http.HandleFunc("/api/task", AuthMiddleware(taskHandler))              // Обработчик задач с авторизацией
	http.HandleFunc("/api/tasks", AuthMiddleware(tasksHandler))            // Обработчик списка задач с авторизацией
	http.HandleFunc("/api/task/done", AuthMiddleware(doneTaskHandler))     // Обработчик выполнения задачи с авторизацией
	http.HandleFunc("/api/task/delete", AuthMiddleware(deleteTaskHandler)) // Обработчик удаления задачи с авторизацией
}
