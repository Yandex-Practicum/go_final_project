package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"TODOGo/database"
	"TODOGo/rules"
	"TODOGo/schemas"
)

// NextDateHandler обрабатывает запросы для вычисления следующей даты
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	// Парсим строку "now" в тип time.Time
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		log.Printf("Ошибка парсинга now: %v", err)
		http.Error(w, `{"error":"некорректное now"}`, http.StatusBadRequest)
		return
	}
	// Вычисляем следующую дату с помощью функции из пакета rules
	nextDate, err := rules.NextDate(now, dateStr, repeat)
	if err != nil {
		log.Printf("Ошибка вычисления следующей даты: %v", err)
		http.Error(w, `{"error":"нкорректное вычисление даты"}`, http.StatusBadRequest)
		return
	}
	// Возвращаем следующую дату в ответе
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(nextDate))
}

// AddTaskHandler обрабатывает запросы на добавление задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	// Декодирование JSON-запроса
	decoder := json.NewDecoder(r.Body)
	var task schemas.Table
	if err := decoder.Decode(&task); err != nil {
		log.Printf("Ошибка десериализации JSON: %v", err)
		http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Получен запрос на добавление задачи: %+v", task)

	// Проверка обязательного поля title
	if task.Title == "" {
		log.Println("Ошибка: не указан заголовок задачи")
		http.Error(w, `{"error":"не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	today := now.Format("20060102")
	log.Printf("Текущая дата перед установкой: %s", today)

	// Установка даты задачи на 'сегодня', если она не указана
	if task.Date == "" || task.Date == "today" {
		task.Date = today
		log.Printf("Дата задачи установлена на 'сегодня': %s", task.Date)
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			log.Printf("Ошибка парсинга даты задачи: %v", err)
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}

		log.Printf("Проверка даты задачи: %s", task.Date)

		// Проверка, если дата задачи меньше текущей
		if parsedDate.Before(time.Now()) {
			log.Printf("Ошибка: дата задачи (%s) меньше текущей даты (%s) и задача не повторяется", task.Date, today)
			if task.Repeat == "" {
				task.Date = today
			} else {
				nextDate, err := rules.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
					return
				}
				if task.Date != today {
					task.Date = nextDate
				}
			}
		}
	}
	// Проверка на повторяющиеся задачи
	if task.Repeat != "" {
		_, err := rules.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
			return
		}
	}
	// Добавление задачи в базу данных
	id, err := database.AddTaskToDB(task)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи в базу данных: %v", err)
		http.Error(w, `{"error":"ошибка при добавлении задачи в базу данных"}`, http.StatusInternalServerError)
		return
	}
	// Формируем успешный ответ
	response := map[string]interface{}{
		"id":      id,
		"message": "Задача успешно добавлена",
		"task":    task,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Ошибка при формировании ответа: %v", err)
		http.Error(w, `{"error":"ошибка при формировании ответа"}`, http.StatusInternalServerError)
		return
	}
	// Устанавливаем заголовки и отправляем ответ
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(responseJSON)

	log.Printf("Задача добавлена с ID: %d", id)
}

// GetAllTasksHandler обрабатывает запросы на получение всех задач
func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search") // Получаем параметр поиска из URL
	const limit = 50                      // Устанавливаем лимит на количество задач

	tasks, err := database.SearchTasksInDB(search, limit)
	if err != nil {
		http.Error(w, `{"error":"ошибка при получении задач"}`, http.StatusInternalServerError)
		return
	}
	// Создаем массив для ответных задач
	var tasksForResponse []map[string]string
	for _, task := range tasks {
		tasksForResponse = append(tasksForResponse, map[string]string{
			"id":      task.ID,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		})
	}
	// Если задач нет, инициализируем пустой массив
	if tasksForResponse == nil {
		tasksForResponse = []map[string]string{}
	}
	response := map[string]interface{}{
		"tasks": tasksForResponse, // Добавляем задачи в ответ
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error":"Ошибка форматирование JSON-ответа"}`, http.StatusInternalServerError)
		return
	}

}

// GetTaskHandler обрабатывает запросы на получение одной задачи по ID
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id") // Получаем ID задачи из URL
	if idStr == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr) // Преобразуем ID из строки в целое число
	if err != nil {
		http.Error(w, `{"error":"Неверный идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := database.FetchTaskByID(id) // Получаем задачу из базы данных по ID
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":      task.ID,
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error":"Ошибка форматирование JSON-ответа"}`, http.StatusInternalServerError)
		return
	}
}

// PutTaskHandler обновляет задачу
func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен PUT-запрос для обновления задачи")

	if r.Method != http.MethodPut {
		log.Println("Метод не поддерживается")
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}
	// Создаем переменную для задачи
	var task schemas.Table
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		log.Println("Не указан идентификатор задачи")
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		log.Println("Не указан заголовок задачи")
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	today := now.Format("20060102")
	if task.Date == "" || task.Date == "today" {
		task.Date = today
	} else {
		parsedDate, err := time.Parse("20060102", task.Date) // Парсим дату задачи
		if err != nil {
			log.Println("Некорректный формат даты:", err)
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Repeat = today // Устанавливаем дату на сегодня
			} else {
				nextDate, err := rules.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					log.Println("Некорректное правило повторения:", err)
					http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
					return
				}
				if task.Date != today {
					task.Date = nextDate // Устанавливаем следующую дату
				}
			}
		}
	}
	if task.Repeat != "" {
		if _, err := rules.NextDate(now, task.Date, task.Repeat); err != nil {
			log.Println("Некорректное правило повторения:", err)
			http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
			return
		}
	}
	// Обновление задачи в базе данных
	if err := database.UpdateTask(task); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Задача не найдена с ID:", task.ID)
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			log.Println("Ошибка обновления задачи:", err)
			http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
		}
		return
	}
	log.Println("Задача успешно обновлена:", task.ID)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}

// DoneTaskHandler отметка о выполнении задачи
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		idStr = r.URL.Query().Get("id") // Если ID не найден в параметрах, получаем его из запроса
	}
	if idStr == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr) // Преобразуем ID из строки в целое число
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"Некорректный идентификатор задачи"}`, http.StatusBadRequest)
		return
	}
	task, err := database.FetchTaskByID(id) // Получаем задачу из базы данных по ID
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Ошибка получения задачи из БД"}`, http.StatusInternalServerError)
		}
		return
	}
	if task.Repeat == "" {
		if err := database.DeleteTaskByID(id); err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := rules.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка расчета следующей даты повторения"}`, http.StatusBadRequest)
			return
		}
		if err := database.UpdateTaskDate(uint64(id), nextDate); err != nil {
			http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}

// DeleteTaskHandler обрабатывает запросы на удаление задачи
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	// Получаем ID задачи из параметров запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}
	// Преобразуем ID из строки в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"Некорректный идентификатор задачи"}`, http.StatusInternalServerError)
		return
	}
	// Удаляем задачу из базы данных по ID
	err = database.DeleteTaskByID(id)
	if err != nil {
		http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
	}
	// Устанавливаем заголовок Content-Type для ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}
