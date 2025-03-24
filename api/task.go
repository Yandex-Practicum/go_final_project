package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Task представляет структуру задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var tasks []Task
	limit := 10 // Максимальное количество задач для возврата

	// Получаем все задачи из базы данных
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?;"
	rows, err := db.Query(query, limit) // Используем limit в запросе
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v\n", err)
		http.Error(w, `{"error": "Ошибка при получении задач: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Считываем результаты
	for rows.Next() {
		var task Task
		var dateStr string

		if err := rows.Scan(&task.ID, &dateStr, &task.Title, &task.Comment, &task.Repeat); err != nil {
			log.Printf("Ошибка при считывании задач: %v\n", err)
			http.Error(w, `{"error": "Ошибка при считывании задач: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		// Проверяем формат даты
		now := time.Now().Truncate(24 * time.Hour) // Убираем время, оставляем только дату
		if dateStr == "" {
			task.Date = now.Format("20060102") // Если дата пустая, устанавливаем текущую дату
		} else {
			// Преобразуем строку даты в time.Time
			date, err := time.Parse("20060102", dateStr) // Используем правильный формат
			if err != nil {
				log.Printf("Ошибка при парсинге даты: %v\n", err)
				http.Error(w, `{"error": "Неверный формат даты"}`, http.StatusBadRequest)
				return
			}
			task.Date = date.Format("20060102") // Форматируем дату в нужный формат
		}
		tasks = append(tasks, task)
	}

	// Чтобы избежать {"tasks":null} в ответе JSON
	if len(tasks) < 1 {
		tasks = make([]Task, 0)
	}

	// Формируем ответ
	response := map[string]interface{}{"tasks": tasks}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при формировании ответа: %v\n", err)
		http.Error(w, `{"error": "Ошибка при формировании ответа: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
}
