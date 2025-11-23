package handlers

import (
	"database/sql"
	"go-task-manager-final_project/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

// APIServer представляет собой структуру сервера API, содержащую подключение к базе данных.
type APIServer struct {
	DB *sql.DB
}

// Init настраивает роутинг для HTTP‑сервера.
// Параметры:
// r — роутер chi.Mux для регистрации эндпоинтов;
// db — подключение к базе данных SQL.
// Регистрирует обработчик для статических файлов и API‑эндпоинты, включая аутентифицированные маршруты для работы с задачами.
func Init(r *chi.Mux, db *sql.DB) {

	server := &APIServer{
		DB: db,
	}

	// Регистрируем обработчик API‑эндпоинта для вычисления следующей даты.
	// Метод: GET. Путь: http://localhost:7540/api/nextdate.
	r.Get("/api/nextdate", handleNextDay)

	// Регистрируем обработчик для аутентификации пользователя.
	// Метод: POST. Путь: http://localhost:7540/api/signin.
	r.Post("/api/signin", handleSignIn)

	// Регистрируем защищённый эндпоинт для получения списка задач.
	// Требуется аутентификация. Метод: GET. Путь: http://localhost:7540/api/tasks.
	r.Get("/api/tasks", middleware.Auth(server.tasksHandler))

	// Регистрируем защищённый эндпоинт для добавления новой задачи.
	// Требуется аутентификация. Метод: POST. Путь: http://localhost:7540/api/task.
	r.Post("/api/task", middleware.Auth(server.addTaskHandler))

	// Регистрируем защищённый эндпоинт для отметки задачи как выполненной.
	// Требуется аутентификация. Метод: POST. Путь: http://localhost:7540/api/task/done.
	r.Post("/api/task/done", middleware.Auth(server.doneTaskHandler))

	// Регистрируем защищённый эндпоинт для получения конкретной задачи.
	// Требуется аутентификация. Метод: GET. Путь: http://localhost:7540/api/task.
	r.Get("/api/task", middleware.Auth(server.getTaskHandler))

	// Регистрируем защищённый эндпоинт для обновления задачи.
	// Требуется аутентификация. Метод: PUT. Путь: http://localhost:7540/api/task.
	r.Put("/api/task", middleware.Auth(server.putTaskHandler))

	// Регистрируем защищённый эндпоинт для удаления задачи.
	// Требуется аутентификация. Метод: DELETE. Путь: http://localhost:7540/api/task.
	r.Delete("/api/task", middleware.Auth(server.deleteTaskHandler))

}
