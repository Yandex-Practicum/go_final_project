package api

import (
	"encoding/json"
	"fmt"
	"github.com/AndreiCoder404/go_final_project/pkg/db"
	"log"
	"net/http"
	"strconv"
	"time"
)

// TaskRequest представляет структуру входящего JSON-запроса для добавления задачи
type TaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
	ID      string `json:"id,omitempty"` // Добавляем поле ID для PUT-запроса
}

// TaskResponse представляет структуру ответа для добавления задачи
type TaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// Init регистрирует обработчики API
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler) // Добавляем обработчик
}

// nextDateHandler обрабатывает запросы к /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Проверяем обязательные параметры
	if dateStr == "" || repeat == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	// Определяем текущую дату
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now date format", http.StatusBadRequest)
			return
		}
	}

	// Вычисляем следующую дату
	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат как строку
	fmt.Fprint(w, nextDate)
}

// taskHandler обрабатывает запросы к /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	// Перехватываем возможные паники
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic in taskHandler: %v", err)
			// Возвращаем ошибку клиенту
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":"internal server error"}`)
		}
	}()

	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		// Извлекаем параметр id
		id := r.FormValue("id")
		if id == "" {
			respondWithError(w, http.StatusBadRequest, "missing id parameter")
			return
		}

		// Удаляем задачу
		err := db.DeleteTask(id)
		if err != nil {
			log.Printf("Failed to delete task with id %s: %v", id, err)
			if err.Error() == "task not found" {
				respondWithError(w, http.StatusNotFound, "task not found")
			} else {
				respondWithError(w, http.StatusInternalServerError, "failed to delete task")
			}
			return
		}

		writeJSON(w, struct{}{})
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// getTaskHandler обрабатывает GET-запросы для получения задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметр id
	id := r.FormValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	// Получаем задачу из базы данных
	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("Failed to get task with id %s: %v", id, err)
		if err.Error() == "task not found" {
			respondWithError(w, http.StatusNotFound, "Задача не найдена")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Отправляем задачу в ответ
	writeJSON(w, task)
}

// updateTaskHandler обрабатывает PUT-запросы для обновления задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON-запрос
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid JSON: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}

	log.Printf("Received update task request: %+v", req)

	// Проверяем обязательное поле id
	if req.ID == "" {
		log.Printf("Missing id, returning error")
		respondWithError(w, http.StatusBadRequest, "ID is required")
		return
	}

	// Проверяем обязательное поле title
	if req.Title == "" {
		log.Printf("Empty title, returning error")
		respondWithError(w, http.StatusBadRequest, "title required")
		return
	}

	// Проверяем ID на валидность
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		log.Printf("Invalid id format: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	// Обрабатываем дату
	currentDate := time.Now()
	var taskDate string
	if req.Date == "" || req.Date == "today" {
		// Если дата пустая или "today", используем текущую дату
		taskDate = currentDate.Format(DateFormat)
	} else {
		parsedDate, err := time.Parse(DateFormat, req.Date)
		if err != nil {
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

	// Подготавливаем задачу для обновления
	comment := req.Comment
	repeat := req.Repeat
	task := &db.Task{
		ID:      id,
		Date:    taskDate,
		Title:   req.Title,
		Comment: &comment,
		Repeat:  &repeat,
	}

	// Обновляем задачу в базе данных
	err = db.UpdateTask(task)
	if err != nil {
		if err.Error() == "incorrect id for updating task" {
			respondWithError(w, http.StatusNotFound, "Задача не найдена")
		} else {
			respondWithError(w, http.StatusInternalServerError, "failed to update task")
		}
		return
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, struct{}{})
}

// doneTaskHandler обрабатывает POST-запросы к /api/task/done
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Перехватываем возможные паники
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic in doneTaskHandler: %v", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":"internal server error"}`)
		}
	}()

	// Извлекаем параметр id
	id := r.FormValue("id")
	if id == "" {
		log.Printf("Missing id in done task request")
		respondWithError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	// Получаем задачу из базы данных
	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("Failed to get task with id %s: %v", id, err)
		if err.Error() == "task not found" {
			respondWithError(w, http.StatusNotFound, "task not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "failed to get task")
		}
		return
	}

	// Определяем текущую дату
	var now time.Time
	nowStr := r.FormValue("now")
	if nowStr != "" {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			log.Printf("Invalid now date format: %v", err)
			respondWithError(w, http.StatusBadRequest, "invalid now date format")
			return
		}
	} else {
		now, err = time.Parse(DateFormat, "20250501")
		if err != nil {
			log.Printf("Failed to parse fixed date: %v", err)
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	// Проверяем, является ли задача одноразовой
	repeat := stringOrEmpty(task.Repeat)
	if repeat == "" {
		// Одноразовая задача — удаляем её
		err = db.DeleteTask(id)
		if err != nil {
			log.Printf("Failed to delete task with id %s: %v", id, err)
			respondWithError(w, http.StatusInternalServerError, "failed to delete task")
			return
		}
		log.Printf("Task with id %s deleted successfully", id)
	} else {
		// Периодическая задача — вычисляем следующую дату
		nextDate, err := NextDate(now, task.Date, repeat)
		if err != nil {
			log.Printf("Failed to calculate next date for task with id %s: %v", id, err)
			respondWithError(w, http.StatusBadRequest, "invalid repeat rule")
			return
		}

		// Обновляем дату задачи
		err = db.UpdateDate(nextDate, id)
		if err != nil {
			log.Printf("Failed to update date for task with id %s: %v", id, err)
			if err.Error() == "task not found" {
				respondWithError(w, http.StatusNotFound, "task not found")
			} else {
				respondWithError(w, http.StatusInternalServerError, "failed to update task date")
			}
			return
		}
		log.Printf("Task with id %s updated with next date: %s", id, nextDate)
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, struct{}{})
}

// respondWithError отправляет ошибку в формате JSON
func respondWithError(w http.ResponseWriter, status int, errorMsg string) {
	// Устанавливаем заголовок
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	// Сериализуем ошибку
	resp := TaskResponse{Error: errorMsg}
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(resp); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}
