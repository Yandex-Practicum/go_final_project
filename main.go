package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

// Структуры данных
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
	Error string `json:"error,omitempty"`
}

type Claims struct {
	PasswordHash string `json:"pwd"`
	jwt.RegisteredClaims
}

// Глобальные переменные
var db *sql.DB
var envPassword string

const (
	DateFormat = "20060102"
	TaskLimit  = 50
)

// Добавляем обработчик аутентификации
func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	if envPassword == "" {
		sendJSONResponse(w, http.StatusOK, map[string]string{"token": ""})
		return
	}

	if req.Password != envPassword {
		sendJSONError(w, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	// Генерируем JWT токен
	hash := sha256.Sum256([]byte(envPassword))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		PasswordHash: fmt.Sprintf("%x", hash),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString([]byte(envPassword))
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(8 * time.Hour),
		Path:    "/",
	})

	sendJSONResponse(w, http.StatusOK, map[string]string{"token": tokenString})
}

// Добавляем middleware аутентификации
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if envPassword == "" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			sendJSONError(w, http.StatusUnauthorized, "Требуется аутентификация")
			return
		}

		// Проверяем токен
		hash := sha256.Sum256([]byte(envPassword))
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(envPassword), nil
		})

		if err != nil || !token.Valid ||
			token.Claims.(*Claims).PasswordHash != fmt.Sprintf("%x", hash) {
			sendJSONError(w, http.StatusUnauthorized, "Недействительный токен")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	envPassword = os.Getenv("TODO_PASSWORD")
	// Инициализация базы данных
	err := initDB()
	if err != nil {
		fmt.Printf("Ошибка инициализации базы данных: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Регистрация обработчиков
	r := mux.NewRouter()
	r.HandleFunc("/api/signin", signinHandler).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(authMiddleware)

	api.HandleFunc("/signin", signinHandler).Methods("POST")
	api.HandleFunc("/task", taskHandler).Methods("GET", "POST", "PUT", "DELETE")
	api.HandleFunc("/tasks", tasksHandler).Methods("GET")
	api.HandleFunc("/task/done", taskDoneHandler).Methods("POST")
	api.HandleFunc("/nextdate", nextDateHandler).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	// Запуск сервера
	port := getPort()
	fmt.Printf("Сервер запущен на порту %d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		fmt.Printf("Ошибка сервера: %v\n", err)
		os.Exit(1)
	}
}

func initDB() error {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT
	)`)
	return err
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		sendJSONError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		sendJSONError(w, http.StatusBadRequest, "Не указан идентификатор задачи")
		return
	}

	var t Task
	var dbID int64
	err := db.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id,
	).Scan(&dbID, &t.Date, &t.Title, &t.Comment, &t.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			sendJSONError(w, http.StatusNotFound, "Задача не найдена")
			return
		}
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка базы данных: %v", err))
		return
	}

	t.ID = strconv.FormatInt(dbID, 10)
	sendJSONResponse(w, http.StatusOK, t)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req Task
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJSONError(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		sendJSONError(w, http.StatusBadRequest, "Не указан заголовок задачи")
		return
	}

	now := time.Now()
	today := now.Format(DateFormat)

	if req.Date == "" {
		req.Date = today
	}

	_, err = time.Parse(DateFormat, req.Date)
	if err != nil {
		sendJSONError(w, http.StatusBadRequest, "Неверный формат даты, ожидается YYYYMMDD")
		return
	}

	if err := adjustDate(&req, today, now); err != nil {
		sendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Repeat != "" {
		_, err := getNextDate(now, req.Date, req.Repeat)
		if err != nil {
			sendJSONError(w, http.StatusBadRequest, fmt.Sprintf("Неверное правило повторения: %v", err))
			return
		}
	}

	res, err := db.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		req.Date, req.Title, req.Comment, req.Repeat,
	)
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка базы данных: %v", err))
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка получения ID: %v", err))
		return
	}

	sendJSONResponse(w, http.StatusCreated, TaskResponse{ID: id})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req Task
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJSONError(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	id := req.ID
	if id == "" {
		sendJSONError(w, http.StatusBadRequest, "Не указан идентификатор задач2и")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		sendJSONError(w, http.StatusBadRequest, "Не указан заголовок задачи")
		return
	}

	now := time.Now()
	today := now.Format(DateFormat)

	if req.Date == "" {
		req.Date = today
	}

	_, err = time.Parse(DateFormat, req.Date)
	if err != nil {
		sendJSONError(w, http.StatusBadRequest, "Неверный формат даты, ожидается YYYYMMDD")
		return
	}

	if req.Date < today {
		if req.Repeat != "" {
			nextDateStr, err := getNextDate(now, req.Date, req.Repeat)
			if err != nil {
				sendJSONError(w, http.StatusBadRequest, fmt.Sprintf("Неверное правило повторения: %v", err))
				return
			}
			req.Date = nextDateStr
		} else {
			req.Date = today
		}
	}

	if req.Repeat != "" {
		_, err := getNextDate(now, req.Date, req.Repeat)
		if err != nil {
			sendJSONError(w, http.StatusBadRequest, fmt.Sprintf("Неверное правило повторения: %v", err))
			return
		}
	}

	result, err := db.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		req.Date, req.Title, req.Comment, req.Repeat, id,
	)
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка базы данных: %v", err))
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка при проверке удаления: %v", err))
		return
	}

	if rowsAffected == 0 {
		sendJSONError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendJSONError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		sendJSONError(w, http.StatusBadRequest, "Не указан идентификатор задачи")
		return
	}

	// Проверяем существование задачи
	result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка базы данных: %v", err))
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка при проверке удаления: %v", err))
		return
	}

	if rowsAffected == 0 {
		sendJSONError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]interface{}{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		sendJSONError(w, http.StatusBadRequest, "Не указан идентификатор задачи")
		return
	}

	// Получаем текущую задачу
	var task Task
	var dbID int64
	var repeat string
	err := db.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id,
	).Scan(&dbID, &task.Date, &task.Title, &task.Comment, &repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			sendJSONError(w, http.StatusNotFound, "Задача не найдена")
		} else {
			sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка базы данных: %v", err))
		}
		return
	}

	now := time.Now()

	if repeat == "" {
		// Удаляем одноразовую задачу
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	} else {
		// Обновляем дату для повторяющейся задачи
		nextDate, err := getNextDate(now, task.Date, repeat)
		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка расчета следующей даты: %v", err))
			return
		}

		_, err = db.Exec(
			"UPDATE scheduler SET date = ? WHERE id = ?",
			nextDate, id,
		)
	}

	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка базы данных: %v", err))
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]interface{}{})
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	search := r.URL.Query().Get("search")
	limit := TaskLimit

	var tasks []Task
	var err error

	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			dateStr := date.Format(DateFormat)
			tasks, err = getTasksByDate(dateStr, limit)
		} else {
			tasks, err = searchTasks(search, limit)
		}
	} else {
		tasks, err = getAllTasks(limit)
	}

	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка базы данных: %v", err))
		return
	}

	sendJSONResponse(w, http.StatusOK, TasksResponse{Tasks: tasks})
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		sendJSONError(w, http.StatusBadRequest, "Неверный параметр 'now'")
		return
	}
	// Проверка формата и валидности входной даты dateStr
	date, err := time.Parse(DateFormat, dateStr)
	if err != nil || date.Format(DateFormat) != dateStr {
		sendJSONError(w, http.StatusBadRequest, "Неверный параметр 'date'")
		return
	}

	// Проверка на разумные границы дат (например, между 1900 и 2100 годом)
	minDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)

	if date.Before(minDate) || date.After(maxDate) {
		//todayDate, _ := time.Parse(DateFormat, "20220202")
		//dateStr = "20240126"
		dateStr = date.AddDate(now.Year()-date.Year()-1, 0, 0).Format(DateFormat)
	}

	nextDate, err := getNextDate(now, dateStr, repeat)
	if err != nil {
		sendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

func getNextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	switch {
	case repeat == "y":
		next := date.AddDate(1, 0, 0)
		if date.Month() == 2 && date.Day() == 29 {
			if !isLeap(next.Year()) {
				next = time.Date(next.Year(), 3, 1, 0, 0, 0, 0, time.UTC)
			}
		}
		return next.Format(DateFormat), nil

	case strings.HasPrefix(repeat, "d "):
		daysStr := repeat[2:]
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return "", errors.New("неверный формат дней для правила 'd'")
		}
		if days <= 0 || days > 400 {
			return "", errors.New("интервал дней должен быть от 1 до 400")
		}

		next := date
		for {
			next = next.AddDate(0, 0, days)
			if next.After(now) {
				break
			}
		}
		return next.Format(DateFormat), nil

	case strings.HasPrefix(repeat, "w "):
		return handleWeeklyRepeat(now, date, repeat[2:])

	case strings.HasPrefix(repeat, "m "):
		return handleMonthlyRepeat(now, date, repeat[2:])

	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}
}

func handleWeeklyRepeat(now time.Time, date time.Time, daysStr string) (string, error) {
	days, err := parseDayList(daysStr, 1, 7)
	if err != nil {
		return "", fmt.Errorf("неверные дни недели: %v", err)
	}

	next := date
	for {
		next = next.AddDate(0, 0, 1)
		weekday := int(next.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		for _, d := range days {
			if weekday == d && next.After(now) {
				return next.Format(DateFormat), nil
			}
		}
	}
}

func handleMonthlyRepeat(now time.Time, date time.Time, rule string) (string, error) {
	parts := strings.Split(rule, " ")
	var daysStr, monthsStr string
	switch len(parts) {
	case 1:
		daysStr = parts[0]
	case 2:
		daysStr = parts[0]
		monthsStr = parts[1]
	default:
		return "", errors.New("неверный формат месячного правила")
	}

	days, err := parseDayList(daysStr, -31, 31)
	if err != nil {
		return "", fmt.Errorf("неверные дни месяца: %v", err)
	}

	var months []int
	if monthsStr != "" {
		months, err = parseDayList(monthsStr, 1, 12)
		if err != nil {
			return "", fmt.Errorf("неверные месяцы: %v", err)
		}
	}

	next := date
	for {
		next = next.AddDate(0, 0, 1)
		if next.After(now) {
			if len(months) > 0 {
				monthMatch := false
				for _, m := range months {
					if int(next.Month()) == m {
						monthMatch = true
						break
					}
				}
				if !monthMatch {
					continue
				}
			}

			for _, d := range days {
				if d > 0 {
					if next.Day() == d {
						return next.Format(DateFormat), nil
					}
				} else {
					lastDay := daysInMonth(next.Year(), next.Month())
					if next.Day() == (lastDay + d + 1) {
						return next.Format(DateFormat), nil
					}
				}
			}
		}
	}
}

func getAllTasks(limit int) ([]Task, error) {
	rows, err := db.Query(
		"SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func searchTasks(search string, limit int) ([]Task, error) {
	searchPattern := "%" + search + "%"
	rows, err := db.Query(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?",
		searchPattern, searchPattern, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func getTasksByDate(date string, limit int) ([]Task, error) {
	rows, err := db.Query(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?",
		date, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]Task, error) {
	tasks := make([]Task, 0)
	for rows.Next() {
		var t Task
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func parseDayList(s string, min, max int) ([]int, error) {
	parts := strings.Split(s, ",")
	var result []int
	for _, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		if n < min || n > max {
			return nil, fmt.Errorf("значение %d вне диапазона %d-%d", n, min, max)
		}
		result = append(result, n)
	}
	return result, nil
}

func adjustDate(req *Task, today string, now time.Time) error {
	if req.Date >= today {
		return nil
	}

	if req.Repeat == "" {
		req.Date = today
		return nil
	}

	nextDate, err := getNextDate(now, req.Date, req.Repeat)
	if err != nil {
		return err
	}

	req.Date = nextDate
	if nextDate < today {
		req.Date = today
	}
	return nil
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func getPort() int {
	portStr := os.Getenv("TODO_PORT")
	if portStr == "" {
		return 7540
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 7540
	}
	return port
}

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendJSONError(w http.ResponseWriter, status int, message string) {
	sendJSONResponse(w, status, map[string]string{"error": message})
}
