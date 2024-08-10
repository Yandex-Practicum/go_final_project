package tasks

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_final_project/internal/handlers/task"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

const (
	FilterTypeNone   = iota
	FilterTypeDate   = iota
	FilterTypeSearch = iota
)

type GetTasksResponseDTO struct {
	Tasks []task.GetTaskResponseDTO `json:"tasks"`
}

func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	filterType := FilterTypeNone
	filterValue := ""
	if len(search) > 0 {
		searchDate, err := time.Parse("01.02.2006", search)
		if err == nil {
			filterType = FilterTypeDate
			filterValue = searchDate.Format(utils.ParseDateFormat)
		} else {
			filterType = FilterTypeSearch
			filterValue = fmt.Sprintf("%%%s%%", search)
		}
	}

	var tasks []*models.Task
	var err error
	switch filterType {
	case FilterTypeDate:
		tasks, err = h.repository.GetAllTasksFilterByDate(filterValue)
	case FilterTypeSearch:
		tasks, err = h.repository.GetAllTasksFilterByTitleOrComment(filterValue)
	default:
		tasks, err = h.repository.GetAllTasks()
	}

	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	response := GetTasksResponseDTO{Tasks: make([]task.GetTaskResponseDTO, 0)}
	for _, t := range tasks {
		response.Tasks = append(response.Tasks, task.GetTaskResponseDTO{
			ID:      strconv.FormatInt(t.ID, 10),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		})
	}

	utils.Respond(w, response)
}
