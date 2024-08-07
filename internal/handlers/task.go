package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_final_project/internal/models"
)

type TaskUpdateDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func TaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTask(w, r, db)
		case http.MethodPost:
			handlePostTask(w, r, db)
		case http.MethodPut:
			handlePutTask(w, r, db)
		case http.MethodDelete:
			handleDeleteTask(w, r, db)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}

func handlePostTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var taskDTO models.Task
	err := json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		respondWithError(w, "Ошибка десериализации JSON")
		return
	}

	task, err := validateTask(&taskDTO)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		respondWithError(w, "Ошибка вставки в базу данных")
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		respondWithError(w, "Ошибка получения ID задачи")
		return
	}

	task.ID = id

	log.Printf("Задача добавлена: %+v\n", task)

	response := models.Response{ID: &task.ID}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func handleGetTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	stringId := r.URL.Query().Get("id")
	if len(stringId) == 0 {
		respondWithError(w, "Не указан идентификатор")
		return
	}
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		respondWithError(w, "Не указан идентификатор")
		return
	}

	query := `SELECT 
    			id,
    			date,
    			title,
    			comment,
    			repeat
			  FROM scheduler
			  WHERE id = ?`
	row := db.QueryRow(query, id)
	var task models.Task
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, "Задача не найдена")
			return
		}
		respondWithError(w, "Ошибка разбора задач из базы данных")
		return
	}
	response := GetTasksTask{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func handlePutTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var taskDTO TaskUpdateDTO
	err := json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		respondWithError(w, "Ошибка десериализации JSON")
		return
	}
	taskId, err := strconv.ParseInt(taskDTO.ID, 10, 64)
	if err != nil {
		respondWithError(w, "Неверный ID задачи")
		return
	}

	taskRequest := models.Task{
		ID:      taskId,
		Date:    taskDTO.Date,
		Title:   taskDTO.Title,
		Comment: taskDTO.Comment,
		Repeat:  taskDTO.Repeat,
	}

	task, err := validateTask(&taskRequest)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}

	query := `UPDATE scheduler
    		  SET
    			date = ?,
    			title = ?,
    			comment = ?,
    			repeat = ?
			  WHERE id = ?`
	updateResult, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		respondWithError(w, "Задача не найдена")
		return
	}

	rowsAffected, err := updateResult.RowsAffected()
	if err != nil || rowsAffected == 0 {
		respondWithError(w, "Задача не найдена")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	stringId := r.URL.Query().Get("id")
	if len(stringId) == 0 {
		respondWithError(w, "Не указан идентификатор")
		return
	}
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		respondWithError(w, "Не указан идентификатор")
		return
	}

	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	_, deleteErr := db.Exec(deleteQuery, id)
	if deleteErr != nil {
		respondWithError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
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

func respondWithError(w http.ResponseWriter, message string) {
	response := models.Response{Error: &message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}
