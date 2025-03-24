package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	ID    string `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type Password struct {
	Pass string `json:"password"`
}

var db *sql.DB

var now = time.Now()

func createTable() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date INTEGER NOT NULL DEFAULT 20060102,
	title TEXT NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat TEXT CHECK (LENGTH(repeat) <= 128) NOT NULL DEFAULT "");`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ошибка загрузки файла .env")
	}
}

// dateParse возвращает близжайщее время задачи
func dateParse(now time.Time, dateStr string, repeat string) (string, error) {
	date, err := time.Parse(timeFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("ошибка при парсинге времени date: %v", err)
	}

	taskDays, err := NextDate(now, date.Format(timeFormat), repeat)
	if err != nil {
		return "", fmt.Errorf("ошибка в функции NextDate: %v", err)
	}

	slices.Sort(taskDays)
	return taskDays[0], nil
}

func sendErrorResponse(res http.ResponseWriter, errorMessage string) {
	response := Response{Error: errorMessage}
	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(res).Encode(response); err != nil {
		http.Error(res, "не удалось обработать ошибку", http.StatusInternalServerError)
	}
}

func sendJSONResponse(res http.ResponseWriter, response interface{}) {
	resp, err := json.Marshal(response)
	if err != nil {
		sendErrorResponse(res, fmt.Sprintf("ошибка при сериализации JSON: %v", err))
		return
	}

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(resp)
}

func addTaskRules(t *Task) (int, error) {
	if t.Date == "" {
		t.Date = time.Now().Format(timeFormat)
	}
	_, err := time.Parse(timeFormat, t.Date)
	if err != nil {
		return 0, fmt.Errorf("дата представлена в формате, отличном от 20060102")
	}

	if t.Date < time.Now().Format(timeFormat) {
		if t.Repeat == "" {
			t.Date = time.Now().Format(timeFormat)
		} else {
			taskDay, err := dateParse(now, t.Date, t.Repeat)
			if err != nil {
				return 0, fmt.Errorf("ошибка при парсинге даты: %v", err)
			}
			t.Date = taskDay
		}
	}

	dateInt, err := strconv.Atoi(t.Date)
	if err != nil {
		return 0, fmt.Errorf("не удалось конвертировать дату в число")
	}

	if t.Title == "" {
		return 0, fmt.Errorf("не указан заголовок задачи")
	}

	return dateInt, nil
}

func nextDateHandler(res http.ResponseWriter, req *http.Request) {
	nowStr := req.FormValue("now")
	dateStr := req.FormValue("date")
	repeatStr := req.FormValue("repeat")

	t, err := time.Parse(timeFormat, nowStr)
	if err != nil {
		sendErrorResponse(res, fmt.Sprintf("ошибка при парсинге времени now: %v", err))
		return
	}

	taskDay, err := dateParse(t, dateStr, repeatStr)
	if err != nil {
		sendErrorResponse(res, fmt.Sprintf("ошибка при парсинге даты: %v", err))
		return
	}

	_, _ = res.Write([]byte(taskDay))
}

func postTask(res http.ResponseWriter, req *http.Request) {
	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		sendErrorResponse(res, fmt.Sprintf("ошибка десериализации JSON: %v", err))
		return
	}

	dateInt, err := addTaskRules(&task)
	if err != nil {
		sendErrorResponse(res, fmt.Sprintf("ошибка при добавлении задачи: %v", err))
		return
	}

	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		dateInt, task.Title, task.Comment, task.Repeat)
	if err != nil {
		sendErrorResponse(res, "не удалось вставить задачу в таблицу")
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		sendErrorResponse(res, "ошибка получения последнего ID")
		return
	}

	response := Response{ID: fmt.Sprintf("%d", id)}

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(res).Encode(response)
}

func getTasks(res http.ResponseWriter, req *http.Request) {
	var (
		tasks []Task
		query string
		args  []interface{}
	)

	limit := 50
	searchStr := req.URL.Query().Get("search")

	if searchStr == "" {
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit"
		args = append(args, limit)
	} else {
		t, err := time.Parse("02.01.2006", searchStr)
		if err == nil {
			query = "SELECT * FROM scheduler WHERE date = :date LIMIT :limit"
			args = append(args, t.Format(timeFormat), limit)
		} else {
			query = "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
			args = append(args, "%"+searchStr+"%", limit)
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		sendErrorResponse(res, "не удалось получить задачу")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			sendErrorResponse(res, fmt.Sprintf("ошибка при извлечении данных: %v", err))
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		sendErrorResponse(res, fmt.Sprintf("ошибка при итерации: %v", err))
		return
	}

	if tasks == nil {
		response := TasksResponse{Tasks: []Task{}}
		sendJSONResponse(res, response)
		return
	}

	response := TasksResponse{Tasks: tasks}
	sendJSONResponse(res, response)
}

func getTaskId(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	var task Task

	rows, err := db.Query("SELECT *FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		sendErrorResponse(res, "задача не найдена")
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			sendErrorResponse(res, fmt.Sprintf("ошибка при извлечении данных: %v", err))
			return
		}
	}

	if task == (Task{}) {
		sendErrorResponse(res, "не указан идентификатор")
		return
	}

	sendJSONResponse(res, task)
}

func putTask(res http.ResponseWriter, req *http.Request) {
	var (
		taskCurrent Task
		taskUpdate  Task
	)

	err := json.NewDecoder(req.Body).Decode(&taskUpdate)
	if err != nil {
		sendErrorResponse(res, fmt.Sprintf("ошибка десериализации JSON: %v", err))
		return
	}

	err = db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", taskUpdate.ID)).
		Scan(&taskCurrent.ID, &taskCurrent.Date, &taskCurrent.Title, &taskCurrent.Comment, &taskCurrent.Repeat)
	if err != nil {
		sendErrorResponse(res, "ошибка при извлечении данных")
		return
	}

	dateInt, err := addTaskRules(&taskUpdate)
	if err != nil {
		sendErrorResponse(res, fmt.Sprintf("ошибка при добавлении задачи: %v", err))
		return
	}

	if taskUpdate.Title == "" {
		sendErrorResponse(res, "не указан заголовок задачи")
		return
	}

	if taskUpdate.Date != "" {
		taskCurrent.Date = taskUpdate.Date
	}

	if taskUpdate.Title != "" {
		taskCurrent.Title = taskUpdate.Title
	}

	if taskUpdate.Comment != "" {
		taskCurrent.Comment = taskUpdate.Comment
	}

	if taskUpdate.Repeat != "" {
		taskCurrent.Repeat = taskUpdate.Repeat
	}

	_, err = db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment , repeat = :repeat WHERE id = :id",
		sql.Named("date", dateInt),
		sql.Named("title", taskCurrent.Title),
		sql.Named("comment", taskCurrent.Comment),
		sql.Named("repeat", taskCurrent.Repeat),
		sql.Named("id", taskCurrent.ID))

	if err != nil {
		sendErrorResponse(res, "ошибка при обновлении данных")
		return
	}

	if taskCurrent == (Task{}) {
		sendErrorResponse(res, "не указан идентификатор")
		return
	}

	sendJSONResponse(res, taskCurrent)
}

func taskDone(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	if id == "" {
		sendErrorResponse(res, "не указан идентификатор")
		return
	}

	var taskCurrent Task

	err := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id)).
		Scan(&taskCurrent.ID, &taskCurrent.Date, &taskCurrent.Title, &taskCurrent.Comment, &taskCurrent.Repeat)
	if err != nil {
		sendErrorResponse(res, "задача не найдена")
		return
	}

	if taskCurrent.Repeat == "" {
		deleteTask(res, req)
		return
	}

	date, err := dateParse(now, taskCurrent.Date, taskCurrent.Repeat)
	if err != nil {
		sendErrorResponse(res, "ошибка при получении даты")
		return
	}

	dateInt, err := strconv.Atoi(date)
	if err != nil {
		sendErrorResponse(res, fmt.Sprintf("ошибка при преобразовании: %v", err))
	}

	_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
		sql.Named("date", dateInt),
		sql.Named("id", id))
	if err != nil {
		sendErrorResponse(res, "ошибка при обновлении данных")
		return
	}
	taskCurrent.Date = date
	sendJSONResponse(res, struct{}{})
}

func deleteTask(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	result, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		sendErrorResponse(res, "ошибка при удалении задачи")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendErrorResponse(res, "ошибка при получении количества удалённых строк")
		return
	}

	if rowsAffected == 0 {
		sendErrorResponse(res, "задача с таким id не найдена")
		return
	}

	sendJSONResponse(res, struct{}{})
}

func signIn(res http.ResponseWriter, req *http.Request) {
	var p Password

	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) > 0 {
		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			sendErrorResponse(res, fmt.Sprintf("ошибка десериализации JSON: %v", err))
			return
		}

		if p.Pass == pass {
			secret := []byte("secret_key")

			hashedPass := sha256.Sum256([]byte(p.Pass))
			hashString := hex.EncodeToString(hashedPass[:])

			claims := jwt.MapClaims{
				"hash": hashString,
			}

			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			signedToken, err := jwtToken.SignedString(secret)
			if err != nil {
				http.Error(res, "ошибка получения подписанного токена", http.StatusBadRequest)
			}
			response := Response{Token: signedToken}
			sendJSONResponse(res, response)
		} else {
			sendErrorResponse(res, "невереный пароль")
		}

	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtCookie string // JWT-токен из куки
			// получаем куку
			cookie, err := req.Cookie("token")
			if err == nil {
				jwtCookie = cookie.Value
			}

			//валидация и проверка JWT-токена
			secret := []byte("secret_key")
			jwtToken, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil {
				http.Error(res, fmt.Sprintf("Ошибка при парсинге токена: %v", err), http.StatusUnauthorized)
				return
			}

			if !jwtToken.Valid {
				sendErrorResponse(res, fmt.Sprint("токен не валиден"))
				return
			}

			result, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				sendErrorResponse(res, fmt.Sprint("не удалось выполнить приведение типа к  jwt.MapClaims"))
				return
			}

			hashRow := result["hash"]

			hash, ok := hashRow.(string)
			if !ok {
				sendErrorResponse(res, fmt.Sprint("не удалось выполнить приведение типа к string"))
				return
			}

			hashedPass := sha256.Sum256([]byte(pass))
			expectedHash := hex.EncodeToString(hashedPass[:])

			if hash != expectedHash {
				sendErrorResponse(res, fmt.Sprint("ошибка валидации хэша"))
				return
			}

			// if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
			// 	hashValue, exists := claims["hash"]

			// 	hashedPass := sha256.Sum256([]byte(pass))
			// 	expectedHash := hex.EncodeToString(hashedPass[:])

			// 	if !exists || hashValue != expectedHash {
			// 		http.Error(res, "Ошибка валадации хэша", http.StatusUnauthorized)
			// 		return
			// 	}
			// } else {
			// 	// возвращаем ошибку авторизации 401
			// 	http.Error(res, "Неверный токен", http.StatusUnauthorized)
			// 	return
			// }
		}
		next(res, req)
	})
}

func main() {
	// appPath, err := os.Executable()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	dbFile := filepath.Join(dbPath)
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if install {
		createTable()
	}

	r := chi.NewRouter()

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("web"))))
	r.Get("/api/nextdate", nextDateHandler)
	r.Post("/api/task", auth(postTask))
	r.Get("/api/tasks", auth(getTasks))
	r.Get("/api/task", auth(getTaskId))
	r.Put("/api/task", auth(putTask))
	r.Post("/api/task/done", auth(taskDone))
	r.Delete("/api/task", auth(deleteTask))
	r.Post("/api/signin", signIn)

	// r.Get("/api/nextdate", nextDateHandler)
	// r.Post("/api/task", postTask)
	// r.Get("/api/tasks", getTasks)
	// r.Get("/api/task", getTaskId)
	// r.Put("/api/task", putTask)
	// r.Post("/api/task/done", taskDone)
	// r.Delete("/api/task", deleteTask)
	// r.Post("/api/signin", signIn)

	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "0.0.0.0:7540"
	}

	adress := "0.0.0.0:" + port

	if err := http.ListenAndServe(adress, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
