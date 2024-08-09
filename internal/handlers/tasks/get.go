package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_final_project/internal/handlers/task"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

const getLimit = 50

type GetTasksResponseDTO struct {
	Tasks []task.GetTaskResponseDTO `json:"tasks"`
}

func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	filterDate := ""
	if len(search) > 0 {
		searchDate, err := time.Parse("01.02.2006", search)
		if err == nil {
			filterDate = searchDate.Format("20060201")
		} else {
			search = fmt.Sprintf("%%%s%%", search)
		}
	}

	var rows *sql.Rows
	var selectErr error
	if len(filterDate) > 0 {
		query := `SELECT 
    			id,
    			date,
    			title,
    			comment,
    			repeat
			  FROM scheduler
			  WHERE date = ?
			  ORDER BY date
			  LIMIT ?`
		rows, selectErr = h.db.Query(query, filterDate, getLimit)
	} else if len(search) > 0 {
		query := `SELECT 
    			id,
    			date,
    			title,
    			comment,
    			repeat
			  FROM scheduler
			  WHERE title LIKE ? OR comment LIKE ?
			  ORDER BY date
			  LIMIT ?`
		rows, selectErr = h.db.Query(query, search, search, getLimit)
	} else {
		query := `SELECT 
    			id,
    			date,
    			title,
    			comment,
    			repeat
			  FROM scheduler
			  ORDER BY date
			  LIMIT ?`
		rows, selectErr = h.db.Query(query, getLimit)
	}

	if selectErr != nil {
		utils.RespondWithError(w, "Ошибка чтения из базы данных")
		return
	}

	response := GetTasksResponseDTO{Tasks: make([]task.GetTaskResponseDTO, 0)}
	for rows.Next() {
		var selectTask models.Task
		err := rows.Scan(&selectTask.ID, &selectTask.Date, &selectTask.Title, &selectTask.Comment, &selectTask.Repeat)
		if err != nil {
			utils.RespondWithError(w, "Ошибка разбора задач из базы данных")
			return
		}
		response.Tasks = append(response.Tasks, task.GetTaskResponseDTO{
			ID:      strconv.FormatInt(selectTask.ID, 10),
			Date:    selectTask.Date,
			Title:   selectTask.Title,
			Comment: selectTask.Comment,
			Repeat:  selectTask.Repeat,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}
