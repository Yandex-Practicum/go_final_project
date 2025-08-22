package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_final_project/pkg/db"
	"go_final_project/pkg/nextdate"
	"log"
	"net/http"
	"strconv"
	"time"
)

// updateTaskHandler обрабатывает PUT-запрос на обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID      interface{} `json:"id"` // Меняем на interface{}
		Date    string      `json:"date"`
		Title   string      `json:"title"`
		Comment string      `json:"comment"`
		Repeat  string      `json:"repeat"`
	}

	// Чтение тела запроса в буфер
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		writeJson(w, map[string]string{"error": "Error reading request body"}, http.StatusBadRequest)
		return
	}

	// Логирование тела запроса
	log.Printf("Received raw request body: %s", buf.String())

	// Декодирование тела запроса с использованием json.Unmarshal
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		log.Printf("Error decoding request body: %v", err) // Логируем ошибку декодирования
		errEncode := json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		if errEncode != nil {
			log.Printf("Error encoding error response: %v", errEncode)
		}
		return
	}

	// Логируем полученные данные для отладки
	log.Printf("Полученные данные: %+v", req)

	// Преобразование ID в строку, если оно числовое
	idStr := ""
	switch v := req.ID.(type) {
	case string:
		idStr = v
	case float64: // Когда ID передан как число
		idStr = strconv.FormatFloat(v, 'f', 0, 64)
	default:
		log.Printf("Invalid ID type: %T", v)
		writeJson(w, map[string]string{"error": "Invalid ID type"}, http.StatusBadRequest)
		return
	}

	// Логирование ID для отладки
	log.Printf("Парсинг ID: %s", idStr)

	// Парсинг ID как целого числа
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		log.Printf("Invalid ID: %v", idStr) // Логируем ошибку парсинга ID
		writeJson(w, map[string]string{"error": "Invalid ID"}, http.StatusBadRequest)
		return
	}

	// Проверка на обязательные поля
	if req.Date == "" || req.Title == "" {
		log.Printf("Missing fields: Date=%s, Title=%s", req.Date, req.Title) // Логируем недостающие поля
		writeJson(w, map[string]string{"error": "Missing required fields"}, http.StatusBadRequest)
		return
	}

	// Парсинг даты в формате "20060102"
	dateParsed, err := time.Parse("20060102", req.Date)
	if err != nil {
		log.Printf("Invalid date format: %s", req.Date) // Логируем ошибку парсинга даты
		writeJson(w, map[string]string{"error": "Invalid date format"}, http.StatusBadRequest)
		return
	}

	// Получаем текущую дату без времени
	today := time.Now().Truncate(24 * time.Hour)

	// Сравнение только по дате (без учёта времени)
	if dateParsed.Before(today) {
		log.Printf("Date cannot be in the past: %s", req.Date) // Логируем ошибку, если дата в прошлом
		writeJson(w, map[string]string{"error": "Date cannot be in the past"}, http.StatusBadRequest)
		return
	}

	// Проверка на корректность правила повторений
	if req.Repeat != "" {
		_, err := nextdate.NextDate(time.Now(), req.Date, req.Repeat)
		if err != nil {
			log.Printf("Invalid repeat rule: %s", req.Repeat) // Логируем ошибку повторения
			writeJson(w, map[string]string{"error": "Invalid repeat rule"}, http.StatusBadRequest)
			return
		}
	}

	// Обновление задачи в базе данных
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err = db.GetDB().Exec(query, req.Date, req.Title, req.Comment, req.Repeat, id)
	if err != nil {
		log.Printf("Failed to update task: %v", err) // Логируем ошибку при обновлении в базе данных
		writeJson(w, map[string]string{"error": fmt.Sprintf("Failed to update task: %v", err)}, http.StatusInternalServerError)
		return
	}

	// Возврат успешного ответа
	writeJson(w, map[string]string{"status": "success", "id": idStr}, http.StatusOK)
}
