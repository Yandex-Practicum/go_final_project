package postHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *Store) AddTask(task Task) (string, err) {

	if t.Date == "" {
		t.Date = time.Now().Format("20060102")
	}
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты")
	}
	if t.Date < time.Now().Format("20060102") {
		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return "", fmt.Errorf("правило повторения указано в неправильном формате")
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format("20060102")
		}
	}

	if task.Title == "" {
		return "", fmt.Errorf("не указан заголовок")
	}

	// добавление задачи в БД

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", fmt.Errorf("задача не добавлена")
	}
	// возвращаем id задачи
	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("не удается получить id")
	}
	return fmt.Sprintf("%d", id), err
}

// POST-обработчик

func PostHandler(w http.ResponseWriter, req *http.Request) {
	var task Task

	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		return http.Error(w, "ошибка десериализации JSON", http.StatusBadRequest)
	}
	id, err := w.AddTask(task)
	if err != nil {
		return http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		http.Error(w, http.StatusBadRequest)
		return
	}
}
