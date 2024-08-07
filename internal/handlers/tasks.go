package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/internal/models"
	"net/http"
	"strconv"
	"time"
)

const getLimit = 50

type GetTasksTask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type GetTasksResponse struct {
	Tasks []GetTasksTask `json:"tasks"`
}

func GetTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTasks(w, r, db)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}

func handleGetTasks(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
		rows, selectErr = db.Query(query, filterDate, getLimit)
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
		rows, selectErr = db.Query(query, search, search, getLimit)
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
		rows, selectErr = db.Query(query, getLimit)
	}

	if selectErr != nil {
		respondWithError(w, "Ошибка чтения из базы данных")
		return
	}

	response := GetTasksResponse{Tasks: make([]GetTasksTask, 0)}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			respondWithError(w, "Ошибка разбора задач из базы данных")
			return
		}
		response.Tasks = append(response.Tasks, GetTasksTask{
			ID:      strconv.FormatInt(task.ID, 10),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}
