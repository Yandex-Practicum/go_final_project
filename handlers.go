package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const dateFormat = "20060102"

var appPassword = os.Getenv("TODO_PASSWORD")

func compareDates(t1, t2 time.Time) (time.Time, time.Time) {
	// Обнуляем время, оставляем только дату
	date1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	date2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())
	return date1, date2
}

// Обработчик для /api/signin
func signInHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var creds struct {
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}
	if appPassword == "" || creds.Password != appPassword {
		http.Error(w, `{"error":"Неверный пароль"}`, http.StatusUnauthorized)
		return
	}
	tokenString, err := generateToken()
	if err != nil {
		http.Error(w, `{"error":"Ошибка генерации токена"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Обработчик для /api/task
func taskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r, db)
	case http.MethodGet:
		getTaskHandler(w, r, db)
	case http.MethodPut:
		editTaskHandler(w, r, db)
	case http.MethodDelete:
		deleteTaskHandler(w, r, db)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработчик для /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, "invalid 'now' date format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, nextDate)
}

// Обработчик добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
	} else {
		_, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}
	}

	now := time.Now()
	tDate, _ := time.Parse(dateFormat, task.Date)
	tDate, now_fix := compareDates(tDate, now)
	if tDate.Before(now_fix) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка добавления задачи в БД: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка получения ID задачи: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	task.ID = strconv.Itoa(int(id))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// Обработчик получения задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, id)
	var task Task
	var taskID int
	err := row.Scan(&taskID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка поиска задачи: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	task.ID = strconv.Itoa(taskID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// Обработчик обновления задачи
func editTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
	} else {
		_, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}
	}

	now := time.Now()
	tDate, _ := time.Parse(dateFormat, task.Date)
	if tDate.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка обновления задачи в БД: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка при получении количества обновленных строк: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{}`)
}

// Обработчик удаления задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка, существует ли задача с данным id
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM scheduler WHERE id = ?)`, id).Scan(&exists)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка проверки существования задачи в БД: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, `{"error":"Задача с таким ID не существует"}`, http.StatusNotFound)
		return
	}

	// Удаление задачи
	query := `DELETE FROM scheduler WHERE id = ?`
	_, err = db.Exec(query, id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка удаления задачи из БД: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{}`)
}

// Обработчик для /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	search := r.URL.Query().Get("search")
	var tasks []Task
	var query string
	var args []interface{}

	if search != "" {
		if searchDate, err := time.Parse("02.01.2006", search); err == nil {
			query = "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT 50"
			args = append(args, searchDate.Format(dateFormat))
		} else {
			query = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 50"
			searchPattern := "%" + search + "%"
			args = append(args, searchPattern, searchPattern)
		}
	} else {
		query = "SELECT * FROM scheduler ORDER BY date LIMIT 50"
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка выборки задач из БД: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		var id int
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка чтения задачи из БД: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		task.ID = strconv.Itoa(id)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка при обработке задач: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
	}

	if len(tasks) == 0 {
		response["tasks"] = []Task{}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Обработчик для /api/task/done
func taskDoneHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	query := "SELECT date, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, id)

	var date string
	var repeat string
	err := row.Scan(&date, &repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка поиска задачи: %s"}`, err.Error()), http.StatusInternalServerError)
		}
		return
	}

	if repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка удаления задачи из БД: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{}`)
		return
	}

	now := time.Now()
	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка расчета следующей даты: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	query = `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = db.Exec(query, nextDate, id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка обновления даты: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{}`)
}
