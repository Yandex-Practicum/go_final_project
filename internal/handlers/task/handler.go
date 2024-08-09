package task

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"go_final_project/internal/models"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetTask(w, r)
		case http.MethodPost:
			h.handlePostTask(w, r)
		case http.MethodPut:
			h.handlePutTask(w, r)
		case http.MethodDelete:
			h.handleDeleteTask(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}

func validateTask(task *models.Task) (*models.Task, error) {
	if task.Title == "" {
		return nil, errors.New("Не указан заголовок задачи")
	}

	now := time.Now()
	today := now.Format("20060102")
	if len(strings.TrimSpace(task.Date)) > 0 {
		taskDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			return nil, errors.New("Не верно указана дата задачи")
		}
		if taskDate.Format("20060102") < today {
			if len(strings.TrimSpace(task.Repeat)) == 0 {
				task.Date = today
			} else {
				nextDate, err := models.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return nil, errors.New("Не верно указана дата задачи и повтор")
				}
				task.Date = nextDate
			}
		}
	} else {
		task.Date = today
	}
	return task, nil
}
