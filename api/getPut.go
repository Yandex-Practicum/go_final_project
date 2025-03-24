package api

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"go_final_project/db"
)

// TaskHandler обрабатывает GET и PUT запросы для задач
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTask(w, r)
	case http.MethodPut:
		updateTask(w, r)
	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

// getTask обрабатывает GET-запрос для получения задачи по идентификатору
func getTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, что идентификатор является числом
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Запрашиваем задачу с ID: %d", id)

	// Получаем задачу из базы данных
	task, err := db.GetTaskByID(strconv.Itoa(id)) // функция реализована в db
	if err != nil {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Форматируем дату
	if task.Date != "" {
		// Преобразуем строку даты в time.Time
		date, err := time.Parse("20060102", task.Date) // Используем правильный формат
		if err != nil {
			log.Printf("Ошибка при парсинге даты: %v\n", err)
			http.Error(w, `{"error": "Неверный формат даты"}`, http.StatusBadRequest)
			return
		}
		task.Date = date.Format("20060102") // Форматируем дату в нужный формат
	} else {
		// Если дата пустая, устанавливаем текущую дату
		now := time.Now().Truncate(24 * time.Hour)
		task.Date = now.Format("20060102")
	}

	// Возвращаем задачу в формате JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, `{"error": "Ошибка формирования ответа"}`, http.StatusInternalServerError)
	}
}

// updateTask обрабатывает PUT-запрос для обновления задачи
func updateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error": "Некорректные данные"}`, http.StatusBadRequest)
		return
	}

	// Проверка на наличие идентификатора
	if task.ID == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Проверка, что идентификатор является числом
	_, err = strconv.Atoi(task.ID)
	if err != nil {
		http.Error(w, `{"error": "Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if task.Title == "" {
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	if _, err := time.Parse("20060102", task.Date); err != nil {
		http.Error(w, `{"error": "Дата представлена в неправильном формате"}`, http.StatusBadRequest)
		return
	}

	taskDate, err := time.Parse("20060102", task.Date) // Преобразуем строку в time.Time
	if err != nil {
		http.Error(w, `{"error": "Неверный формат даты"}`, http.StatusBadRequest)
		return
	}

	// Дата не должна быть в прошлом
	now := time.Now().Truncate(24 * time.Hour)
	if taskDate.Before(now) {
		http.Error(w, `{"error": "Дата не может быть в прошлом"}`, http.StatusBadRequest)
		return
	}

	// Проверка поля repeat
	if task.Repeat != "" && !isValidRepeat(task.Repeat) {
		http.Error(w, `{"error": "Некорректное значение для поля repeat"}`, http.StatusBadRequest)
		return
	}

	// Обновляем задачу в базе данных
	err = db.UpdateTask(task) // Реализовано в db.go
	if err != nil {
		http.Error(w, `{"error": "Не удалось обновить задачу"}`, http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON в случае успеха
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// isValidRepeat проверяет, является ли значение repeat допустимым
func isValidRepeat(repeat string) bool {
	// Регулярное выражение для проверки допустимых форматов
	re := regexp.MustCompile(`^(d \d{1,3}|y)$`)
	return re.MatchString(repeat)
}
