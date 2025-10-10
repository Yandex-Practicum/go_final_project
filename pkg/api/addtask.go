package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

// AddTaskRequest структура для входящего запроса
type AddTaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTaskResponse структура для ответа
type AddTaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// addTaskHandler обработчик для добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, AddTaskResponse{Error: "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	var req AddTaskRequest

	// Декодируем JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, AddTaskResponse{Error: fmt.Sprintf("Ошибка декодирования JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	fmt.Printf("Получена задача: Date='%s', Title='%s', Comment='%s', Repeat='%s'\n",
		req.Date, req.Title, req.Comment, req.Repeat)

	// Валидируем заголовок
	if req.Title == "" {
		fmt.Println("Не указан заголовок задачи")
		writeJSON(w, AddTaskResponse{Error: "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Подготавливаем задачу
	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	// Обрабатываем дату и правило повторения
	if err := processTaskDate(task); err != nil {
		fmt.Printf("Ошибка даты : %v\n", err)
		writeJSON(w, AddTaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	fmt.Printf("Задача после обработки: Date='%s', Title='%s', Comment='%s', Repeat='%s'\n",
		task.Date, task.Title, task.Comment, task.Repeat)

	// Добавляем задачу в БД
	id, err := db.AddTask(task)
	if err != nil {
		fmt.Printf("Ошибка базы данных: %v", err)
		writeJSON(w, AddTaskResponse{Error: fmt.Sprintf("Ошибка базы данных: %v", err)}, http.StatusInternalServerError)
		return
	}

	fmt.Printf("Задача успешно добавлена с ID: %d\n", id)

	// Возвращаем успешный ответ
	writeJSON(w, AddTaskResponse{ID: fmt.Sprintf("%d", id)}, http.StatusOK)
}

// processTaskDate обрабатывает дату задачи и правило повторения
func processTaskDate(task *db.Task) error {
	now := time.Now()

	// Если дата не указана, используем сегодня
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	// Парсим дату для проверки формата
	parsedDate, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %v", err)
	}

	// Если дата в прошлом
	if !afterNow(parsedDate, now) {
		if task.Repeat == "" {
			// Без повторения - используем сегодня
			task.Date = now.Format(dateFormat)
		} else {
			// С повторением - вычисляем следующую дату
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("неверный формат повторения: %v", err)
			}
			task.Date = nextDate
		}
	} else if task.Repeat != "" {
		// Если дата в будущем, но есть правило - проверяем его корректность
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("неверный формат повторения: %v", err)
		}
	}

	return nil
}

// writeJSON вспомогательная функция для отправки JSON ответов
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	// Устанавливаем заголовки ОДИН раз
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode) // ← ТОЛЬКО ОДИН ВЫЗОВ

	if data == nil {
		w.Write([]byte("{}"))
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		// Если ошибка, пишем простой JSON ошибки
		errorJson := []byte(`{"error":"Ошибка формирования JSON"}`)
		w.Write(errorJson)
		return
	}

	w.Write(jsonData)
}
