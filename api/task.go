package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// Task представляет структуру задачи
type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// GetTasksHandler обрабатывает GET-запросы для получения списка задач
func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		http.Error(w, `{"error": "Ошибка при получении задач"}`, http.StatusInternalServerError)
		log.Println(`{"Ошибка при выполнении запроса:"}`, err)
		return
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, `{"error": "Ошибка при обработке задач"}`, http.StatusInternalServerError)
			log.Println(`{"Ошибка при сканировании строки:"}`, err)
			return
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, `{"error": "Ошибка при чтении результатов"}`, http.StatusInternalServerError)
		return
	}

	// Обработка пустого результата
	if len(tasks) == 0 {
		response := map[string]interface{}{
			"tasks": []Task{}, // Возвращаем пустой массив задач
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error": "Ошибка при формировании ответа"}`, http.StatusInternalServerError)
			log.Println(`{"Ошибка при кодировании JSON:"}`, err)
			return
		}
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
	}

	log.Println("Формируемый ответ:", response) // Логирование ответа

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Ошибка при формировании ответа"}`, http.StatusInternalServerError)
		log.Println(`{"Ошибка при кодировании JSON:"}`, err)
		return
	}
}
