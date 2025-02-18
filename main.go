package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // Импортируем драйвер sqlite3
)

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func initDB() error {
	// Получаем текущий рабочий каталог
	appPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Определяем полный путь к файлу базы данных
	dbFile := filepath.Join(appPath, "scheduler.db")

	// Проверяем, существует ли файл базы данных
	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("failed to check database file: %v", err)
		}
	}

	// Открываем или создаем базу данных
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	if install {
		// Если база данных не существует, создаем таблицу и индекс
		_, err = db.Exec(`
            CREATE TABLE scheduler (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                date TEXT NOT NULL,
                title TEXT NOT NULL,
                comment TEXT,
                repeat TEXT CHECK (LENGTH(repeat) <= 128)
            );
            CREATE INDEX idx_date ON scheduler(date);
        `)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
		log.Println("Database and table created successfully.")
	} else {
		log.Println("Database already exists.")
	}

	return nil
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	d, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}

	switch {
	case repeat == "":
		return "", fmt.Errorf("repeat rule is empty")
	case repeat == "y":
		for d.Before(now) || d.Equal(now) {
			d = d.AddDate(1, 0, 0)
		}
		return d.Format("20060102"), nil
	case strings.HasPrefix(repeat, "d "):
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("invalid repeat rule: %v", repeat)
		}
		for d.Before(now) || d.Equal(now) {
			d = d.AddDate(0, 0, days)
		}
		return d.Format("20060102"), nil
	case strings.HasPrefix(repeat, "w "):
		var daysOfWeek []int
		_, err := fmt.Sscanf(repeat, "w %v", &daysOfWeek)
		if err != nil {
			return "", fmt.Errorf("invalid repeat rule: %v", repeat)
		}
		for {
			d = d.AddDate(0, 0, 1)
			for _, day := range daysOfWeek {
				if int(d.Weekday()) == day {
					if d.After(now) {
						return d.Format("20060102"), nil
					}
				}
			}
		}
	case strings.HasPrefix(repeat, "m "):
		var daysOfMonth []int
		var months []int
		_, err := fmt.Sscanf(repeat, "m %v %v", &daysOfMonth, &months)
		if err != nil {
			return "", fmt.Errorf("invalid repeat rule: %v", repeat)
		}
		for {
			d = d.AddDate(0, 1, 0)
			for _, day := range daysOfMonth {
				if d.Day() == day {
					if len(months) == 0 || contains(months, int(d.Month())) {
						if d.After(now) {
							return d.Format("20060102"), nil
						}
					}
				}
			}
		}
	default:
		return "", fmt.Errorf("unsupported repeat rule: %v", repeat)
	}
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Некорректный формат now", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, nextDate)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error decoding JSON: %v"}`, err), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error": "Invalid date format"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(time.Now()) {
			if task.Repeat == "" {
				task.Date = time.Now().Format("20060102")
			} else {
				nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error": "Invalid repeat rule: %v"}`, err), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			}
		}
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to open database: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to insert task: %v"}`, err), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to retrieve last insert ID: %v"}`, err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"id": id}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", addTaskHandler) // Добавляем маршрут для обработчика

	port := 7540
	addr := ":" + strconv.Itoa(port)
	log.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
