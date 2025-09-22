package task

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	_ "modernc.org/sqlite"
)

var DateFormat = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask обрабатывает POST-запросы для добавления задачи 4ый шаг
func AddTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)
	if req.Method != http.MethodPost {
		http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
		return
	}
	var task Task

	// Декодируем JSON из запроса
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&task); err != nil {
		log.Printf("Ошибка десериализации JSON")
		http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		log.Printf("Не указан заголовок задачи")
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверяем и парсим дату
	date, err := ParseDate(task.Date)
	if err != nil {
		log.Printf("Неправильный формат времени")
		http.Error(w, `{"error": "Неправильный формат времени"}`, http.StatusBadRequest)
		return
	}

	if task.Repeat != "" && !regexp.MustCompile(`^(d \d{1,3}|y)$`).MatchString(task.Repeat) {
		log.Printf("Неправильно указано правило повторения")
		http.Error(w, `{"error": "Неправильно указано правило повторения"}`, http.StatusBadRequest)
		return
	}

	var newDateStr string
	now := time.Now().Local().Truncate(24 * time.Hour)
	// Если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число
	if date.Before(now) && task.Repeat == "" {
		date = now
	}
	// Если дата задачи меньше сегодняшней, вычисляем новую дату с учетом повторения
	if date.Before(now) && task.Repeat != "" {
		// Если повторение указано, вычисляем следующую дату с учетом правила
		newDateStr, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Не удается вычислить новую дату")
			http.Error(w, `{"error": "Не удается вычислить новую дату"}`, http.StatusBadRequest)
			return
		}
	} else {
		newDateStr = date.Format(DateFormat)
	}

	task.Date = newDateStr

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("insert failed")
		http.Error(w, `{"error": "insert failed"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("scan id failed")
		http.Error(w, `{"error": "scan id failed"}`, http.StatusInternalServerError)
		return
	}

	// Формируем ответ

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]interface{}{"id": id}
	json.NewEncoder(w).Encode(response)

}

// GatTasks обрабатывает GET-запросы для добавления задачи 5ый шаг
func GetTasks(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)
	if req.Method != http.MethodGet {
		http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
		return
	}

	var tasks []Task
	limit := 10
	// Получаем список задач с ограничением в limit  штук
	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?;`, limit)
	if err != nil {
		log.Printf("select failed")
		http.Error(w, `{"error": "select failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		t := Task{}

		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Printf("scan failed")
			http.Error(w, `{"error": "scan failed"}`, http.StatusInternalServerError)
			return
		}
		if t.Date == "" {
			t.Date = time.Now().Format(DateFormat)
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows unpacking failed")
		http.Error(w, `{"error": "rows unpacking failed"}`, http.StatusInternalServerError)
		return
	}
	// Если список задач пуст, то получаем пустой список
	if len(tasks) < 1 {
		tasks = make([]Task, 0)
	}

	// Формируем ответ
	response := map[string]interface{}{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при формировании ответа: %v\n", err)
		http.Error(w, `{"error": "Ошибка при формировании ответа: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
}

func GetTask(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)
	if req.Method != http.MethodGet {
		log.Printf("Метод не разрешен ")
		http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
		return
	}
	id := req.URL.Query().Get("id")

	if id == "" {
		log.Printf("Не указан идентификатор задачи")
		http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	var t Task
	// Формируем запрос к БД
	row := db.QueryRow("SELECT  id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err == sql.ErrNoRows {
		log.Printf("Задача (ID: %s): %v не найдена", id, err)
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Ошибка поиска задачи (ID: %s): %v", id, err)
		http.Error(w, `{"error": "Ошибка поиска задачи в БД"}`, http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Printf("Ошибка при формировании ответа: %v\n", err)
		http.Error(w, `{"error": "Ошибка при формировании ответа: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
}
