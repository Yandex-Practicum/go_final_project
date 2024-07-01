package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"example.com/todo/api"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var jwtKey = []byte("your_secret_key")

// Task представляет задачу
type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

// Response представляет ответ сервера
type Response struct {
	ID       string `json:"id,omitempty"`
	Error    string `json:"error,omitempty"`
	Tasks    []Task `json:"tasks,omitempty"`
	NextDate string `json:"nextDate,omitempty"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/task", auth(addTaskHandler(db)))
	http.HandleFunc("/api/task/done", auth(markTaskDoneHandler(db)))
	http.HandleFunc("/api/tasks", auth(getTasksHandler(db)))

	webDir := "web"
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Запуск сервера на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if creds.Password != expectedPassword {
		http.Error(w, `{"error": "Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		Username: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			tokenStr := cookie.Value

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

// addTaskHandler обработчик для добавления задачи
func addTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %s", err)
			http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
			return
		}
		log.Printf("Received body: %s", body)

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		err = json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Printf("JSON Decode Error: %s", err)
			http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, "Не указан заголовок задачи", http.StatusBadRequest)
			return
		}

		const layout = "20060102"
		now := time.Now()

		if task.Date == "" {
			task.Date = now.Format(layout)
		} else {
			parsedDate, err := time.Parse(layout, task.Date)
			if err != nil {
				http.Error(w, "Дата указана в неправильном формате", http.StatusBadRequest)
				return
			}
			if parsedDate.Before(now) {
				if task.Repeat == "" {
					task.Date = now.Format(layout)
				} else {
					nextDate, err := api.NextDate(now, task.Date, task.Repeat)
					if err != nil {
						http.Error(w, "Ошибка вычисления следующей даты: "+err.Error(), http.StatusBadRequest)
						return
					}
					task.Date = nextDate
				}
			}
		}

		res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
			task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, "Ошибка добавления задачи в базу данных", http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Ошибка получения ID новой задачи", http.StatusInternalServerError)
			return
		}

		response := Response{ID: fmt.Sprintf("%d", id)}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// initDB инициализирует базу данных
func initDB() (*sqlx.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	_, err = os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err := sqlx.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			title TEXT,
			comment TEXT,
			repeat TEXT
		)`)
		if err != nil {
			return nil, err
		}
		log.Println("Таблица scheduler создана или уже существует.")

		_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date)`)
		if err != nil {
			return nil, err
		}
		log.Println("Индекс idx_date создан или уже существует.")
	}

	return db, nil
}

// updateTaskHandler обработчик для обновления задачи
func updateTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
			return
		}

		if task.ID == "" {
			http.Error(w, "Не указан идентификатор задачи", http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, "Не указан заголовок задачи", http.StatusBadRequest)
			return
		}

		const layout = "20060102"
		now := time.Now()

		// Проверка и установка даты
		if task.Date == "" {
			task.Date = now.Format(layout)
		} else {
			parsedDate, err := time.Parse(layout, task.Date)
			if err != nil {
				http.Error(w, "Дата указана в неправильном формате", http.StatusBadRequest)
				return
			}
			if parsedDate.Before(now) {
				if task.Repeat == "" {
					task.Date = now.Format(layout)
				} else {
					nextDate, err := api.NextDate(now, task.Date, task.Repeat)
					if err != nil {
						http.Error(w, "Ошибка вычисления следующей даты: "+err.Error(), http.StatusBadRequest)
						return
					}
					task.Date = nextDate
				}
			}
		}

		// Обновление задачи в базе данных
		_, err = db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
			task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			http.Error(w, "Ошибка обновления задачи в базе данных", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{})
	}
}

// getTaskHandler обработчик для получения задачи по идентификатору
func getTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		var task Task
		err := db.QueryRowx("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).StructScan(&task)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "Ошибка при получении задачи"}`, http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

// deleteTaskHandler обработчик для удаления задачи
func deleteTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{})
	}
}

// markTaskDoneHandler обработчик для отметки задачи выполненной
func markTaskDoneHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		var task Task
		err := db.QueryRowx("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).StructScan(&task)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "Ошибка при получении задачи"}`, http.StatusInternalServerError)
			}
			return
		}

		if task.Repeat == "" {
			_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
			if err != nil {
				http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
				return
			}
		} else {
			now := time.Now()
			nextDate, err := api.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error": "Ошибка при вычислении следующей даты"}`, http.StatusInternalServerError)
				return
			}

			_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
			if err != nil {
				http.Error(w, `{"error": "Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{})
	}
}

// getTasksHandler обработчик для получения списка задач
func getTasksHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Queryx("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
		if err != nil {
			http.Error(w, "Ошибка выборки задач из базы данных", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			if err := rows.StructScan(&task); err != nil {
				http.Error(w, "Ошибка сканирования задачи", http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		if tasks == nil {
			tasks = []Task{}
		}

		response := Response{Tasks: tasks}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
