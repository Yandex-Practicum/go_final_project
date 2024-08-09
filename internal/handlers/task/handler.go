package task

import (
	"database/sql"
	"go_final_project/internal/utils"
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
			http.Error(w, utils.ErrUnsupportedMethod, http.StatusMethodNotAllowed)
		}
	}
}

func validateTask(task *models.Task) (*models.Task, error) {
	if task.Title == "" {
		return nil, utils.ErrInvalidTaskTitle
	}

	now := time.Now()
	today := now.Format(utils.ParseDateFormat)
	if len(strings.TrimSpace(task.Date)) > 0 {
		taskDate, err := time.Parse(utils.ParseDateFormat, task.Date)
		if err != nil {
			return nil, utils.ErrInvalidTaskDate
		}
		if taskDate.Format(utils.ParseDateFormat) < today {
			if len(strings.TrimSpace(task.Repeat)) == 0 {
				task.Date = today
			} else {
				nextDate, err := models.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return nil, utils.ErrInvalidTaskRepeat
				}
				task.Date = nextDate
			}
		}
	} else {
		task.Date = today
	}
	return task, nil
}
