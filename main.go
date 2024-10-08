package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func createDB() {
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		log.Fatalf("Error when opening database: %v", err)
	}
	defer db.Close()

	commands := []string{
		`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8),
			title TEXT,
			comment TEXT,
			repeat CHAR (128)
					)`,
		"CREATE INDEX IF NOT EXISTS indexdate ON scheduler (date)",
	}

	for _, cmd := range commands {
		if _, err := db.Exec(cmd); err != nil {
			log.Fatalf("Error during command execution: %s, error: %v", cmd, err)
		}
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || date == "" || repeat == "" {
		http.Error(w, "Necessary parameters are missing", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {

		http.Error(w, "Incorrect time format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, nextDate)
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	fmt.Println()
	fmt.Println()
	fmt.Println("Accepted values: Now:", now.Format("20060102"), "Start:", date, "-", repeat, "	")

	if repeat == "" {
		fmt.Println("ERROR: There is no repetition rule!")
		return "", errors.New("no repetition rule")
	}

	rep := strings.Split(repeat, " ")
	fmt.Println("Parsing the repetition rule:", rep, " ")

	if len(rep) < 1 || (rep[0] != "y" && rep[0] != "d") {
		fmt.Println("ERROR: Unsupported repetition rule! ")
		return "repetition rule is in the wrong format", errors.New("repetition rule is in the wrong format")
	}

	timBase, err := time.Parse("20060102", date)
	if err != nil {
		fmt.Println("ERROR: Incorrect date!")
		return "", err
	}
	fmt.Println("The date has been successfully recognized:", timBase, " ")

	if rep[0] == "y" {

		fmt.Println("Define the repetition mode: year (y) ")
		timBase = timBase.AddDate(1, 0, 0)
		for timBase.Before(now) {
			timBase = timBase.AddDate(1, 0, 0)
			fmt.Println("Adding 1 year! New date:", timBase.Format("20060102"), " ")
		}
		result := timBase.Format("20060102")
		fmt.Println("Old date:", date, " ")
		fmt.Println("New Date:", result, " ")
		fmt.Println("Add: 1 year")
		fmt.Println("New Date:", result, " ")
		return result, nil
	}

	if rep[0] == "d" {

		fmt.Println("Define repeat mode: day (d) ")
		if len(rep) < 2 {
			fmt.Println("ERROR: Incorrect repeat mode specified! ")
			return "", errors.New("Repeat mode is incorrectly specified")
		}

		days, err := strconv.Atoi(rep[1])
		if err != nil {
			return "", err
		}

		if days > 400 {
			fmt.Println("ERROR: Number of days exceeds 400! ")
			return "", errors.New("postponement of the event for more than 400 days is not allowed")
		}

		fmt.Println("Number of days to add:", days, " ")
		if days == 1 && now.Format("20060102") == timBase.Format("20060102") {
			fmt.Println("")
			fmt.Println("", now.Format("20060102"), " = ", timBase.Format("20060102"), " ")
			fmt.Println("1 DAY!")
			fmt.Println("")
		} else {
			fmt.Println("TOTAL", days, "DAYS! OUT ")
			timBase = timBase.AddDate(0, 0, days)
			for timBase.Before(now) {
				timBase = timBase.AddDate(0, 0, days)
				fmt.Println("Adding", days, "days! New date:", timBase.Format("20060102"), " ")
			}
		}

		result := timBase.Format("20060102")
		fmt.Println("Current Date:", now.Format("20060102"), " ")
		fmt.Println("Start Date:", date, " ")
		fmt.Println("Add:", days, " days")
		fmt.Println("New Date:", result, " ")
		return result, nil
	}

	return "", errors.New("incorrect repetition rule")
}

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendJSONError(w, "JSON deserialization error: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Accepted Values:", "Start:", task.Date, "-", task.Repeat, "-", task.Title, "-", task.Comment, " ")

	if task.Title == "" {
		sendJSONError(w, "No task title specified", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			sendJSONError(w, "The date is in the wrong format, YYYYYYMMDD is expected YYYYMMDD", http.StatusBadRequest)
			return
		}

		if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format("20060102")
			} else {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					sendJSONError(w, "Error in the repetition rule: "+err.Error(), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			}
		}
	}

	if task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			sendJSONError(w, "The repetition rule is specified in the wrong format: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		sendJSONError(w, "Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)")
	if err != nil {
		sendJSONError(w, "An error in the preparation of the request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		sendJSONError(w, "Error when inserting a task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		sendJSONError(w, "Error while getting task ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	task.ID = strconv.FormatInt(id, 10)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"id": task.ID})
}

func getTasks(w http.ResponseWriter, _ *http.Request) {
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		http.Error(w, "server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC")
	if err != nil {
		http.Error(w, "Error when retrieving tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []map[string]string
	for rows.Next() {
		var id, date, title, comment, repeat string
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			http.Error(w, "Error while scanning a task: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, map[string]string{
			"id":      id,
			"date":    date,
			"title":   title,
			"comment": comment,
			"repeat":  repeat,
		})
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error in reading task data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}

func getTaskByID(w http.ResponseWriter, _ *http.Request, id string) {
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		http.Error(w, "Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var task struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		} else {
			http.Error(w, "Error when receiving a task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}

func markTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{"error": "No identifier specified"})
		return
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		http.Error(w, "Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var task struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		} else {
			http.Error(w, "Error when receiving a task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Error when deleting a task: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		sendJSONError(w, "Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	limit := 20

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	search := r.URL.Query().Get("search")
	dateParam := r.URL.Query().Get("date")

	var rows *sql.Rows

	if search != "" {

		searchPattern := "%" + search + "%"
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?"
		rows, err = db.Query(query, searchPattern, searchPattern, limit)
	} else if dateParam != "" {

		date := dateParam
		if len(dateParam) == 10 && strings.Contains(dateParam, ".") {

			t, err := time.Parse("02.01.2006", dateParam)
			if err != nil {
				sendJSONError(w, "Incorrect date format, expected YYYYYYMMDD or DD.MM.YYYYYY", http.StatusBadRequest)
				return
			}
			date = t.Format("20060102")
		}

		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?"
		rows, err = db.Query(query, date, limit)
	} else {

		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?"
		rows, err = db.Query(query, limit)
	}

	if err != nil {
		sendJSONError(w, "Error when retrieving tasks:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := make([]Task, 0)

	for rows.Next() {
		var task Task
		var id int64
		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			sendJSONError(w, "Error in reading the task: "+err.Error(), http.StatusInternalServerError)
			return
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		sendJSONError(w, "Error in task processing: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendJSONError(w, "JSON deserialization error:"+err.Error(), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		sendJSONError(w, "No task identifier specified", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		sendJSONError(w, "No task title specified", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			sendJSONError(w, "The date is in the wrong format, YYYYYYMMDD is expected YYYYMMDD", http.StatusBadRequest)
			return
		}

		if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format("20060102")
			} else {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					sendJSONError(w, "Error in the repetition rule: "+err.Error(), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			}
		}
	}

	if task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			sendJSONError(w, "The repetition rule is specified in the wrong format: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		sendJSONError(w, "Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingID string
	err = db.QueryRow("SELECT id FROM scheduler WHERE id = ?", task.ID).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendJSONError(w, "Task not found", http.StatusNotFound)
		} else {
			sendJSONError(w, "Error in checking the task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	stmt, err := db.Prepare("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?")
	if err != nil {
		sendJSONError(w, "An error in the preparation of the request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		sendJSONError(w, "Error when updating a task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTask(w, r)
	case http.MethodGet:
		if id := r.URL.Query().Get("id"); id != "" {
			getTaskByID(w, r, id)
		} else {

			sendJSONError(w, "No identifier specified", http.StatusBadRequest)
		}
	case http.MethodPut:
		updateTask(w, r)
	case http.MethodDelete:

	default:
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	if _, err := os.Stat("./scheduler.db"); os.IsNotExist(err) {
		createDB()
	}

	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", markTaskDone)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", nextDateHandler)
	log.Fatal(http.ListenAndServe(":7540", nil))
}
