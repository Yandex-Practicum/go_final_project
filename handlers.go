package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		http.Error(w, "Ошибка получения задач", http.StatusInternalServerError)
		log.Printf("Ошибка выполнения запроса: %v", err)
		return
	}
	defer rows.Close()

	tasks := []map[string]string{}
	for rows.Next() {
		var id, date, title, comment, repeat string
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			http.Error(w, "Ошибка обработки данных", http.StatusInternalServerError)
			log.Printf("Ошибка чтения строки: %v", err)
			return
		}
		tasks = append(tasks, map[string]string{
			"id":      id,
			"date":    date,
			"title":   title,
			"comment": comment,
			"repeat":  repeat,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}

func EditTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса - PUT, потому что мы редактируем задачу
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Создаем структуру для получения данных из тела запроса
	var task struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Date    string `json:"date"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Ошибка при разборе данных", http.StatusBadRequest)
		log.Printf("Ошибка при разборе данных: %v", err)
		return
	}

	if task.ID == 0 || task.Title == "" || task.Date == "" {
		http.Error(w, "ID, title и date обязательны", http.StatusBadRequest)
		return
	}

	// Обновляем данные задачи в базе данных
	_, err := db.Exec(`
		UPDATE scheduler 
		SET title = ?, date = ?, comment = ?, repeat = ?
		WHERE id = ?`,
		task.Title, task.Date, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		http.Error(w, "Ошибка обновления задачи", http.StatusInternalServerError)
		log.Printf("Ошибка обновления задачи: %v", err)
		return
	}

	// Возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Задача успешно обновлена",
		"id":      strconv.Itoa(task.ID),
	}
	json.NewEncoder(w).Encode(response)
}

func MarkTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса - POST или PUT, так как мы обновляем задачу
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID задачи из URL параметров
	taskIDStr := r.URL.Query().Get("id")
	if taskIDStr == "" {
		http.Error(w, "ID задачи обязателен", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Неверный формат ID задачи", http.StatusBadRequest)
		return
	}

	// Обновляем статус задачи
	_, err = db.Exec(`
		UPDATE scheduler 
		SET status = 'done' 
		WHERE id = ?`,
		taskID,
	)
	if err != nil {
		http.Error(w, "Ошибка обновления статуса задачи", http.StatusInternalServerError)
		log.Printf("Ошибка выполнения запроса: %v", err)
		return
	}

	// Возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Задача помечена как выполненная",
		"id":      strconv.Itoa(taskID),
	}
	json.NewEncoder(w).Encode(response)
}
