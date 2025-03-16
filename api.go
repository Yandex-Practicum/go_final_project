package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// taskHandler обрабатывает запросы на добавление, получение или обновление задачи
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTask(w, r) // Вызов функции для добавления задачи
	case http.MethodGet:
		getTask(w, r) // Вызов функции для получения задачи по ID
	case http.MethodPut:
		updateTask(w, r) // Вызов функции для обновления задачи
	case http.MethodDelete:
		deleteTaskHandler(w, r) // Вызов функции для удаления задачи
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

const dateLayout = "20060102" // Шаблон формата даты

const limit = 50 // Максимальное количество задач для возврата

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		return err
	}
	return nil
}

// updateTask обрабатывает PUT-запросы на обновление задачи
func updateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if task.Title == "" {
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	taskDate, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		http.Error(w, `{"error": "Дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
		return
	}

	// Проверка идентификатора
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		http.Error(w, `{"error": "Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Проверка формата поля repeat
	if !isValidRepeatFormat(task.Repeat) {
		http.Error(w, `{"error": "Некорректный формат поля repeat"}`, http.StatusBadRequest)
		return
	}

	// Запрос для обновления задачи
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	result, err := db.Exec(query, taskDate.Format(dateLayout), task.Title, task.Comment, task.Repeat, id)
	if err != nil {
		http.Error(w, `{"error": "Ошибка обновления задачи"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, `{"error": "Ошибка получения количества обновленных строк"}`, http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON в случае успеха
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// addTask обрабатывает POST-запросы на добавление задачи
func addTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	var taskReq Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&taskReq); err != nil {
		http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if taskReq.Title == "" {
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	now := time.Now()
	var taskDate time.Time

	if taskReq.Date == "" {
		taskDate = now
	} else {
		var err error
		taskDate, err = time.Parse(dateLayout, taskReq.Date)
		if err != nil {
			http.Error(w, `{"error": "Дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
			return
		}

		// Проверка на существование даты
		if !isValidDate(taskReq.Date) {
			http.Error(w, `{"error": "Указана несуществующая дата"}`, http.StatusBadRequest)
			return
		}
	}

	// Проверка, что задача должна быть запланирована после now
	if taskDate.Before(now) {
		if taskReq.Repeat == "" {
			taskDate = now
		} else {
			// Проверка формата поля repeat
			if !isValidRepeatFormat(taskReq.Repeat) {
				http.Error(w, `{"error": "Некорректный формат поля repeat"}`, http.StatusBadRequest)
				return
			}

			// Вычисляем следующую дату выполнения с помощью функции NextDate
			nextDateStr, err := NextDate(now, taskReq.Date, taskReq.Repeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			taskDate, _ = time.Parse(dateLayout, nextDateStr)
		}
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, taskDate.Format(dateLayout), taskReq.Title, taskReq.Comment, taskReq.Repeat)
	if err != nil {
		http.Error(w, `{"error": "Ошибка добавления задачи в базу данных"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error": "Ошибка получения ID созданной задачи"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON с ID созданной задачи
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]int64{"id": id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Ошибка формирования ответа"}`, http.StatusInternalServerError)
	}
}

// getTask обрабатывает GET-запросы на получение задачи по идентификатору
func getTask(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметр id из URL
	idStr := r.URL.Query().Get("id")

	// Проверяем, указан ли идентификатор
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Преобразуем идентификатор в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Запрос для получения задачи по идентификатору
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, id)

	var task struct {
		ID      int
		Date    string
		Title   string
		Comment string
		Repeat  string
	}

	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		} else {
			log.Printf("Ошибка выполнения запроса к базе данных: %v", err) // Логируем ошибку
			http.Error(w, `{"error": "Ошибка выполнения запроса к базе данных"}`, http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем задачу в формате JSON
	response := map[string]string{
		"id":      strconv.Itoa(task.ID),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

// tasksHandler обрабатывает запросы на получение списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var query string
	var rows *sql.Rows

	// Проверка на поиск по дате
	if _, err := time.Parse("02.01.2006", search); err == nil {
		// Если строка поиска соответствует формату даты
		date, _ := time.Parse("02.01.2006", search)
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = db.Query(query, date.Format(dateLayout), limit)

		if err != nil {
			http.Error(w, `{"error": "Ошибка выполнения запроса к базе данных"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
	} else {
		// Поиск по заголовку и комментарию
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		searchPattern := "%" + search + "%"
		rows, err = db.Query(query, searchPattern, searchPattern, limit)

		if err != nil {
			http.Error(w, `{"error": "Ошибка выполнения запроса к базе данных"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
	}

	// Получаем задачи
	tasks := []Task{} // Используем структуру Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, `{"error": "Ошибка чтения данных"}`, http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task) // Добавляем задачу в срез
	}

	// Возвращаем список задач
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}

// nextDateHandler обрабатывает запросы на получение следующей даты
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	// Парсинг текущей даты
	now, err := time.Parse(dateLayout, nowStr)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат даты now"}`, http.StatusBadRequest)
		return
	}

	// Вызываем функцию NextDate с now и dateStr
	nextDateStr, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем только значение nextDateStr в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Устанавливаем статус 200 OK
	w.Write([]byte(nextDateStr)) // Возвращаем значение как строку в JSON
}

// doneTaskHandler обрабатывает запросы на пометку задачи как выполненной
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Получаем текущую задачу
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, id)

	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "Ошибка выполнения запроса к базе данных"}`, http.StatusInternalServerError)
		}
		return
	}

	// Проверяем, является ли задача периодической
	if task.Repeat == "" {
		// Удаляем задачу, если она одноразовая
		deleteQuery := `DELETE FROM scheduler WHERE id = ?`
		_, err = db.Exec(deleteQuery, id)
		if err != nil {
			http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// Рассчитываем следующую дату выполнения
		now := time.Now()
		nextDateStr, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		updateQuery := `UPDATE scheduler SET date = ? WHERE id = ?`
		_, err = db.Exec(updateQuery, nextDateStr, id)
		if err != nil {
			http.Error(w, `{"error": "Ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	// Возвращаем пустой JSON в случае успеха
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// deleteTaskHandler обрабатывает запросы на удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Удаляем задачу
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	result, err := db.Exec(deleteQuery, id)
	if err != nil {
		http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, `{"error": "Ошибка получения количества удаленных строк"}`, http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON в случае успеха
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}")) // Возвращаем пустой JSON
}

// isValidRepeatFormat проверяет корректность формата поля repeat
func isValidRepeatFormat(repeat string) bool {
	// Регулярное выражение для проверки формата repeat
	var re = regexp.MustCompile(`^(d \d{1,3}|y)?$`)
	return re.MatchString(repeat)
}

// isValidDate проверяет, существует ли указанная дата
func isValidDate(dateStr string) bool {
	_, err := time.Parse(dateLayout, dateStr)
	return err == nil
}
