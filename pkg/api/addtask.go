package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"todo-list/pkg/db"
)

type TaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := processDate(&req); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Repeat != "" {
		if err := validateRepeat(req.Repeat); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	id, err := db.AddTask(task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"id": formatID(id)}
	writeJSONResponse(w, response, http.StatusOK)
}

func processDate(req *TaskRequest) error {
	// Получаем текущую дату без времени
	now := time.Now()
	nowDate := now.Format(DateFormat)

	// Парсим текущую дату для сравнения (без времени)
	nowTrunc, _ := time.Parse(DateFormat, nowDate)

	// Если дата не указана, используем сегодняшнюю
	if req.Date == "" {
		req.Date = nowDate
		return nil
	}

	// Проверяем формат даты
	date, err := time.Parse(DateFormat, req.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты: ожидается YYYYMMDD")
	}

	// Сравниваем только даты (без времени)
	// Если дата в прошлом (меньше сегодняшней)
	if date.Before(nowTrunc) {
		// Если правило повторения не указано, используем сегодняшнюю дату
		if req.Repeat == "" {
			req.Date = nowDate
			return nil
		}

		// Если правило повторения указано, вычисляем следующую дату
		// Используем now для вычислений, но передаём дату из запроса
		nextDate, err := NextDate(now, req.Date, req.Repeat)
		if err != nil {
			return fmt.Errorf("ошибка вычисления следующей даты: %v", err)
		}
		req.Date = nextDate
	}
	// Если дата сегодня или в будущем, оставляем как есть

	return nil
}

func validateRepeat(repeat string) error {
	if repeat == "" {
		return fmt.Errorf("правило повторения не может быть пустым")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return fmt.Errorf("неверный формат правила повторения")
	}

	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return fmt.Errorf("формат правила d неверный, ожидался: d <число>")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("неверное количество дней: %s", parts[1])
		}
		if days < 1 || days > 400 {
			return fmt.Errorf("количество дней должно быть от 1 до 400, получено: %d", days)
		}
		return nil

	case "y":
		if len(parts) != 1 {
			return fmt.Errorf("формат правила y неверный, лишние параметры")
		}
		return nil

	default:
		return fmt.Errorf("неподдерживаемый формат правила: %s (поддерживаются: d, y)", rule)
	}
}
