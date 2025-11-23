package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go-task-manager-final_project/internal/api"
	"go-task-manager-final_project/internal/db"
	"go-task-manager-final_project/internal/services"
)

// Функция проверяет и корректирует дату задачи.
// Параметры:
// task - указатель на структуру задачи, поле Date которой подлежит проверке и корректировке.
// Возвращает: ошибку, если дата некорректна или возникла проблема при обработке.
func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата не указана или равна "today", устанавливаем текущую дату в формате services.DateFormat
	if task.Date == "" || task.Date == "today" {
		task.Date = now.Format(services.DateFormat)
	}

	// Преобразуем строку с датой в объект time.Time по формату services.DateFormat
	t, err := time.Parse(services.DateFormat, task.Date)
	if err != nil {
		return err
	}

	// Проверяем, не превышает ли дата текущую (t > now)
	if services.AfterNow(now, t) {
		if task.Repeat == "" {
			// Если повторение не задано, устанавливаем текущую дату
			task.Date = now.Format(services.DateFormat)
		} else {
			// Если задано повторение, вычисляем следующую допустимую дату выполнения
			next, err := services.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			// Обновляем дату задачи на вычисленную следующую дату
			task.Date = next
		}
	}

	return nil
}

// Метод обработчика HTTP-запроса для добавления новой задачи.
// Параметры:
// w - интерфейс для записи HTTP-ответа.
// r - HTTP-запрос с данными новой задачи.
func (s *APIServer) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем значение заголовка Content-Type из запроса
	contentType := r.Header.Get("Content-Type")

	// Проверяем, что Content-Type начинается с "application/json" (с учётом регистра)
	if !strings.EqualFold(strings.TrimSpace(contentType), "application/json") {
		api.WriteJSON(w, http.StatusUnsupportedMediaType, map[string]string{
			"error": "content type must be application/json",
		})
		return
	}

	var task db.Task

	// Декодируем JSON из тела запроса в структуру задачи
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		r.Body.Close()
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON payload",
		})
		// Завершаем обработку из‑за некорректного JSON
		return
	}
	// Закрываем тело после успешного декодирования
	r.Body.Close()

	// Проверяем, что поле Title не пустое (обязательное поле)
	if task.Title == "" {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title cannot be empty",
		})
		// Завершаем обработку, так как Title обязателен
		return
	}

	// Проверяем и корректируем дату задачи согласно бизнес‑логике
	if err := checkDate(&task); err != nil {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		// Завершаем обработку при ошибке валидации даты
		return
	}

	// Сохраняем задачу в базу данных через функцию AddTask
	id, err := db.AddTask(s.DB, &task)
	if err != nil {
		log.Printf("failed to save task: %v, task data: %+v", err, task)
		api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to save task",
		})
		// Завершаем обработку при ошибке сохранения
		return
	}

	// Формируем успешный ответ:
	// - id: идентификатор созданной задачи
	// - location: URL для доступа к задаче
	// - message: текстовое подтверждение создания
	api.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       id,
		"location": fmt.Sprintf("/tasks/%d", id),
		"message":  "Task created successfully",
	})
}
