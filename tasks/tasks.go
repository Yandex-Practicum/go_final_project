package tasks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"main.go/parsedate"
)

type TaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type DBTask struct {
	ID      int    `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

type JSONTask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	ID int64 `json:"id"`
}

func respondWithError(w http.ResponseWriter, head int, message string) {
	w.WriteHeader(head)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func respondWithSuccess(w http.ResponseWriter, head int, id int64) {
	w.WriteHeader(head)
	_ = json.NewEncoder(w).Encode(SuccessResponse{ID: id})
}

func ValidateAndProcessTaskRequest(req *TaskRequest, now time.Time) (time.Time, error) {
	if req.Title == "" {
		return time.Time{}, errors.New("Missed header")
	}

	var finalDate time.Time

	if req.Date == "" || req.Date == "today" || req.Date == now.Format("20060102") {
		finalDate = now
	} else {
		parsedDate, err := time.ParseInLocation("20060102", req.Date, time.Local)
		if err != nil {
			return time.Time{}, errors.New("Invalid date format")
		}
		finalDate = parsedDate

		if finalDate.Before(now) {
			if req.Repeat == "" {
				finalDate = now
			} else {
				next, err := parsedate.NextDate(now, finalDate.Format("20060102"), req.Repeat)
				if err != nil {
					return time.Time{}, err
				}
				finalDate, _ = time.ParseInLocation("20060102", next, time.Local)
			}
		}
	}
	if req.Repeat != "" && req.Date != "today" && req.Date != now.Format("20060102") {
		if _, err := parsedate.NextDate(now, finalDate.Format("20060102"), req.Repeat); err != nil {
			return time.Time{}, err
		}
	}
	return finalDate, nil
}

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method denied")
			return
		}

		var req TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Wrong request format")
			return
		}

		now := time.Now().Local().Truncate(24 * time.Hour)
		finalDate, err := ValidateAndProcessTaskRequest(&req, now)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := db.Exec(
			`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
			finalDate.Format("20060102"),
			req.Title,
			req.Comment,
			req.Repeat,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		id, _ := res.LastInsertId()
		respondWithSuccess(w, http.StatusCreated, id)

	}
}

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		search := r.URL.Query().Get("search")

		var tasks []DBTask
		var err error

		query := "SELECT id, date, title, comment, repeat FROM scheduler"
		args := []interface{}{}
		whereAdded := false
		limit := 50

		if search != "" {

			if date, err := time.Parse("02.01.2006", search); err == nil {
				if whereAdded {
					query += " AND date = ?"
				} else {
					query += " WHERE date = ?"
					whereAdded = true
				}
				args = append(args, date.Format("20060102"))
			} else {

				search = strings.ReplaceAll(search, "%", "\\%")
				search = strings.ReplaceAll(search, "_", "\\_")
				searchTerm := "%" + search + "%"
				if whereAdded {
					query += " AND (title LIKE ? OR comment LIKE ?)"
				} else {
					query += " WHERE (title LIKE ? OR comment LIKE ?)"
					whereAdded = true
				}
				args = append(args, searchTerm, searchTerm)
			}
		}

		query += " ORDER BY date LIMIT ?"
		args = append(args, limit)

		rows, err := db.QueryContext(r.Context(), query, args...)
		if err != nil {
			log.Printf("Error request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var task DBTask
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				log.Printf("Row scan error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
				return
			}

			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Rows error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}

		jsonTasks := make([]JSONTask, 0, len(tasks))
		for _, task := range tasks {
			jsonTasks = append(jsonTasks, JSONTask{
				ID:      strconv.Itoa(task.ID),
				Date:    task.Date,
				Title:   task.Title,
				Comment: task.Comment,
				Repeat:  task.Repeat,
			})
		}

		response := struct {
			Tasks []JSONTask `json:"tasks"`
		}{
			Tasks: jsonTasks,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("JSON encode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

}

func GetTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Missed id"})
			return
		}

		var task DBTask

		err := db.QueryRowContext(r.Context(),
			"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
			Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Task not found"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}
		json.NewEncoder(w).Encode(JSONTask{
			ID:      strconv.Itoa(task.ID),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})

	}
}

func UpdateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}
		if req.ID == "" {
			respondWithError(w, http.StatusBadRequest, "Missed ID")
			return
		}

		now := time.Now().Local().Truncate(24 * time.Hour)
		finalDate, err := ValidateAndProcessTaskRequest(&req, now)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		res, err := db.ExecContext(r.Context(),
			"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
			finalDate.Format("20060102"),
			req.Title,
			req.Comment,
			req.Repeat,
			req.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			respondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func DoneMarkHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Missed id"})
			return
		}

		var task struct {
			Date   string
			Repeat string
		}
		err := db.QueryRowContext(r.Context(),
			"SELECT date, repeat FROM scheduler WHERE id = ?", id).
			Scan(&task.Date, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Task not found"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Server error"})
			return
		}
		now := time.Now().UTC()

		if task.Repeat != "" {
			parsedDate, _ := time.Parse("20060102", task.Date)
			next, err := parsedate.NextDate(now, parsedDate.Format("20060102"), task.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
				return
			}

			_, err = db.ExecContext(r.Context(),
				"UPDATE scheduler SET date = ? WHERE id = ?",
				next,
				id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Server error"})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(struct{}{})
		} else {
			_, err := db.ExecContext(r.Context(),
				"DELETE FROM scheduler WHERE id = ?", id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Server error"})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(struct{}{})
		}
	}
}

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Missed id"})
			return
		}
		result, err := db.ExecContext(r.Context(),
			"DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Server error"})
			return
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Task not found"})
			return
		}
		json.NewEncoder(w).Encode(struct{}{})
	}
}
