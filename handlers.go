package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func createTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(task.Title) == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		today := time.Now().Format("20060102")
		if strings.TrimSpace(task.Date) == "" {
			task.Date = today
		}

		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
			return
		}

		if parsedDate.Before(time.Now().Local().Truncate(24 * time.Hour)) {
			if strings.TrimSpace(task.Repeat) == "" {
				task.Date = today
			} else {
				nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Ошибка обработки правила повторения"}`, http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			}
		}

		res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка записи в БД"}`, http.StatusInternalServerError)
			return
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Ошибка получения ID"}`, http.StatusInternalServerError)
			return
		}

		resp := Response{ID: int(lastID)}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func getTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")

		var searchDate string
		if search != "" {
			parsedDate, err := time.Parse("02.01.2006", search)
			if err == nil {
				searchDate = parsedDate.Format("20060102")
			}
		}

		query := "SELECT id, date, title, comment, repeat FROM scheduler"
		args := []interface{}{}
		whereClauses := []string{}

		if searchDate != "" {
			whereClauses = append(whereClauses, "date = ?")
			args = append(args, searchDate)
		} else if search != "" { // Поиск по подстроке
			whereClauses = append(whereClauses, "(title LIKE ? OR comment LIKE ?)")
			searchParam := "%" + search + "%"
			args = append(args, searchParam, searchParam)
		}

		if len(whereClauses) > 0 {
			query += " WHERE " + strings.Join(whereClauses, " AND ")
		}

		query += " ORDER BY date ASC"

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, `{"error":"Ошибка получения списка задач"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		tasks := []Task{}
		for rows.Next() {
			var t Task
			if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
				http.Error(w, `{"error":"Ошибка чтения данных"}`, http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, t)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, `{"error":"Ошибка обработки результатов"}`, http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []Task{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks}); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
		}
	}
}

func getTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if strings.TrimSpace(id) == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		var task Task
		err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
			Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, `{"error":"Ошибка при поиске задачи"}`, http.StatusInternalServerError)
			log.Printf("Ошибка поиска задачи (ID: %s): %v", id, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func updateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(task.ID) == "" || strings.TrimSpace(task.Title) == "" {
			http.Error(w, `{"error":"Некорректные данные"}`, http.StatusBadRequest)
			return
		}

		if !isValidDate(task.Date) {
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
			return
		}

		if task.Repeat != "" && !isValidRepeat(task.Repeat) {
			http.Error(w, `{"error":"Некорректный repeat"}`, http.StatusBadRequest)
			return
		}

		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM scheduler WHERE id = ?)", task.ID).Scan(&exists)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при проверке задачи"}`, http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		_, err = db.Exec("UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?",
			task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}")) // Успешный ответ
	}
}

func isValidDate(date string) bool {
	if match, _ := regexp.MatchString(`^\d{8}$`, date); !match {
		return false
	}

	year := date[:4]
	month := date[4:6]
	day := date[6:8]

	y, _ := strconv.Atoi(year)
	m, _ := strconv.Atoi(month)
	d, _ := strconv.Atoi(day)

	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	return t.Year() == y && int(t.Month()) == m && t.Day() == d
}

func isValidRepeat(repeat string) bool {
	match, _ := regexp.MatchString(`^(d|w) \d+$`, repeat)
	return match
}
func doneTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"Не передан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		var date, repeat string
		err := db.QueryRow("SELECT date, repeat FROM scheduler WHERE id = ?", id).Scan(&date, &repeat)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, `{"error":"Ошибка запроса к БД"}`, http.StatusInternalServerError)
			return
		}

		if repeat == "" {
			_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
			if err != nil {
				http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{}"))
			return
		}

		now := time.Now()
		newDate, err := NextDate(now, date, repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка вычисления следующей даты"}`, http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}
}
func deleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"Не передан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}")) // Успешный ответ
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "Отсутствует важные параметры", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}
