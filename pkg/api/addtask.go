package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/nextdate"
)

// Функция проверки и корректировки даты
func checkDate(task *db.Task) error {
	now := time.Now().In(time.UTC) // Приводим текущее время к UTC

	// Если дата не указана, устанавливаем текущую
	if task.Date == "" {
		task.Date = now.Format(nextdate.TimeFormat)
		log.Printf("Дата не указана, установлена текущая дата: %v", task.Date)
	}

	// Парсим дату задачи
	t, err := time.Parse(nextdate.TimeFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %v", err)
	}

	log.Printf("Парсинг даты задачи: %v", t)

	// Если дата задачи в прошлом, вычисляем следующую дату с учётом повторений
	if !afterNow(now, t) {
		log.Printf("Дата задачи (%v) не может быть меньше сегодняшней (%v)", t, now)

		// Если повторение не указано, устанавливаем текущую дату
		if task.Repeat == "" {
			task.Date = now.Format(nextdate.TimeFormat)
			log.Printf("Дата установлена на текущую: %v", task.Date)
		} else {
			// Если задано правило повторения, вычисляем следующую дату
			next, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка в правиле повторения: %v", err)
			}

			// Парсим вычисленную следующую дату
			nextParsed, err := time.Parse(nextdate.TimeFormat, next)
			if err != nil {
				return fmt.Errorf("ошибка разбора следующей даты: %v", err)
			}

			log.Printf("Вычислена следующая дата: %v", nextParsed)

			// Обработка повторения "d 1" (ежедневно)
			if task.Repeat == "d 1" {
				// Устанавливаем дату задачи на сегодняшнюю, если она в будущем
				task.Date = now.Format(nextdate.TimeFormat)
				log.Printf("Дата установлена на сегодняшнюю, так как повторение 'd 1': %v", task.Date)
			} else if task.Repeat == "y" {
				// Обработка повторения "y" (ежегодно)
				// Проверяем, если дата в прошлом, то устанавливаем дату на следующий год
				if !afterNow(now, nextParsed) {
					task.Date = nextParsed.AddDate(1, 0, 0).Format(nextdate.TimeFormat)
					log.Printf("Дата установлена на следующий год, так как повторение 'y': %v", task.Date)
				}
			} else {
				// Для других повторений, оставляем как есть
				task.Date = next
			}
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка чтения тела запроса: %v", err)}, http.StatusBadRequest)
		return
	}
	log.Printf("Полученные данные: %s", body)

	// Декодирование JSON
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка декодирования JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Поле 'title' обязательно"}, http.StatusBadRequest)
		return
	}

	// Проверка даты
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Добавление задачи в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка добавления задачи в базу данных: %v", err)}, http.StatusInternalServerError)
		return
	}

	// Ответ с id задачи
	writeJson(w, map[string]int64{"id": id}, http.StatusOK)
}

// утилита для json ответов
func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	// Логирование отправляемого ответа
	log.Printf("Response: %+v", data)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		http.Error(w, "Ошибка сериализации ответа", http.StatusInternalServerError)
	}
}

func afterNow(now, date time.Time) bool {
	truncatedNow := now.In(time.UTC).Truncate(24 * time.Hour)   // Переводим в UTC и обрезаем время
	truncatedDate := date.In(time.UTC).Truncate(24 * time.Hour) // Переводим в UTC и обрезаем время

	log.Printf("Сравнение дат: now = %v, date = %v", truncatedNow, truncatedDate)
	return truncatedDate.After(truncatedNow) || truncatedDate.Equal(truncatedNow)
}
