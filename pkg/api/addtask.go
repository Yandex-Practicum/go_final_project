package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AndreiCoder404/go_final_project/pkg/db"
)

// stringOrEmpty возвращает пустую строку, если указатель nil
func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// addTaskHandler обрабатывает POST-запросы для добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Перехватываем возможные паники
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic in addTaskHandler: %v", err)
		}
	}()

	// Парсим JSON-запрос
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid JSON: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}

	// Проверяем обязательное поле title
	if req.Title == "" {
		log.Printf("Empty title, returning error")
		respondWithError(w, http.StatusBadRequest, "title required")
		return
	}

	// Обрабатываем дату
	currentDate := time.Now()
	var taskDate string
	if req.Date == "" || req.Date == "today" {
		// Если дата пустая или "today", используем текущую дату
		taskDate = currentDate.Format(DateFormat)
	} else {
		// Парсим дату
		parsedDate, err := time.Parse(DateFormat, req.Date)
		if err != nil {
			// Если дата некорректная, возвращаем ошибку
			log.Printf("Invalid date format: %v", err)
			respondWithError(w, http.StatusBadRequest, "invalid date format")
			return
		}

		// Если дата меньше текущей, заменяем на текущую или вычисляем следующую
		currentDateTruncated := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
		parsedDateTruncated := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
		if parsedDateTruncated.Before(currentDateTruncated) {
			if req.Repeat == "" {
				taskDate = currentDate.Format(DateFormat)
			} else {
				// Вычисляем следующую дату с помощью NextDate
				nextDate, err := NextDate(currentDate, req.Date, req.Repeat)
				if err != nil {
					log.Printf("Invalid repeat rule when calculating next date: %v", err)
					respondWithError(w, http.StatusBadRequest, "invalid repeat rule")
					return
				}
				taskDate = nextDate
			}
		} else {
			taskDate = req.Date
		}
	}

	// Проверяем правило повторения, если оно указано
	if req.Repeat != "" {
		// Используем текущую дату для проверки repeat
		_, err := NextDate(currentDate, taskDate, req.Repeat)
		if err != nil {
			log.Printf("Invalid repeat rule: %v", err)
			respondWithError(w, http.StatusBadRequest, "invalid repeat rule")
			return
		}
	}

	// Подготавливаем задачу для сохранения
	comment := req.Comment
	repeat := req.Repeat
	task := db.Task{
		Date:  taskDate,
		Title: req.Title,
		// Всегда указатель на строку
		Comment: &comment,
		Repeat:  &repeat,
	}

	// Сохраняем задачу в базе данных
	id, err := db.AddTask(task)
	if err != nil {
		log.Printf("Failed to save task: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to save task")
		return
	}

	// Формируем ответ
	resp := TaskResponse{
		ID: fmt.Sprintf("%d", id),
	}

	// Отправляем ответ
	writeJSON(w, resp)
}

// writeJSON сериализует данные в JSON и отправляет их в ответ
func writeJSON(w http.ResponseWriter, data any) {
	// Устанавливаем заголовок
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Сериализуем данные
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false) // Отключаем экранирование HTML
	if err := encoder.Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
