package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
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

type TaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// updateTask обрабатывает PUT-запросы на обновление задачи
func updateTask(w http.ResponseWriter, r *http.Request) {
	var task TaskRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if task.Title == "" {
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, "Дата представлена в формате, отличном от 20060102", http.StatusBadRequest)
		return
	}

	// Проверка идентификатора
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		http.Error(w, `{"error": "Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "scheduler.db") // Убедитесь, что используете правильный путь
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Запрос для обновления задачи
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	result, err := db.Exec(query, taskDate.Format("20060102"), task.Title, task.Comment, task.Repeat, id)
	if err != nil {
		http.Error(w, "Ошибка обновления задачи", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Ошибка получения количества обновленных строк", http.StatusInternalServerError)
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
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var taskReq TaskRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&taskReq); err != nil {
		http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if taskReq.Title == "" {
		http.Error(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Парсинг даты
	now := time.Now()
	var taskDate time.Time

	if taskReq.Date == "" {
		taskDate = now
	} else {
		var err error
		taskDate, err = time.Parse("20060102", taskReq.Date)
		if err != nil {
			http.Error(w, "Дата представлена в формате, отличном от 20060102", http.StatusBadRequest)
			return
		}
	}

	// Проверка, что задача должна быть запланирована после now
	if taskDate.Before(now) {
		if taskReq.Repeat == "" {
			taskDate = now
		} else {
			// Вычисляем следующую дату выполнения с помощью функции NextDate
			nextDateStr, err := NextDate(now, taskReq.Date, taskReq.Repeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			taskDate, _ = time.Parse("20060102", nextDateStr)
		}
	}

	// Добавление задачи в базу данных
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, taskDate.Format("20060102"), taskReq.Title, taskReq.Comment, taskReq.Repeat)
	if err != nil {
		http.Error(w, "Ошибка добавления задачи в базу данных", http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Ошибка получения ID созданной задачи", http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON с ID созданной задачи
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
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

	db, err := sql.Open("sqlite3", "scheduler.db") // Убедитесь, что используете правильный путь
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

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
			http.Error(w, "Ошибка выполнения запроса к базе данных", http.StatusInternalServerError)
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
	limit := 50 // Максимальное количество задач для возврата

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var query string
	var rows *sql.Rows

	// Проверка на поиск по дате
	if _, err := time.Parse("02.01.2006", search); err == nil {
		// Если строка поиска соответствует формату даты
		date, _ := time.Parse("02.01.2006", search)
		query = "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = db.Query(query, date.Format("20060102"), limit)
	} else {
		// Поиск по заголовку и комментарию
		query = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		searchPattern := "%" + search + "%"
		rows, err = db.Query(query, searchPattern, searchPattern, limit)
	}

	if err != nil {
		http.Error(w, "Ошибка выполнения запроса к базе данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Получаем задачи
	tasks := []map[string]string{}
	for rows.Next() {
		var id int
		var date, title, comment, repeat string
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}

		tasks = append(tasks, map[string]string{
			"id":      strconv.Itoa(id),
			"date":    date,
			"title":   title,
			"comment": comment,
			"repeat":  repeat,
		})
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
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Неверный формат даты now", http.StatusBadRequest)
		return
	}

	// Вызываем функцию NextDate с now и dateStr
	nextDateStr, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем следующую дату в формате 20060102
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nextDateStr)
}

// doneTaskHandler обрабатывает запросы на пометку задачи как выполненной
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
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

	db, err := sql.Open("sqlite3", "scheduler.db") // Убедитесь, что используете правильный путь
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Получаем текущую задачу
	var task TaskRequest
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, id)

	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка выполнения запроса к базе данных", http.StatusInternalServerError)
		}
		return
	}

	// Проверяем, является ли задача периодической
	if task.Repeat == "" {
		// Удаляем задачу, если она одноразовая
		deleteQuery := `DELETE FROM scheduler WHERE id = ?`
		_, err = db.Exec(deleteQuery, id)
		if err != nil {
			http.Error(w, "Ошибка удаления задачи", http.StatusInternalServerError)
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
			http.Error(w, "Ошибка обновления задачи", http.StatusInternalServerError)
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
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
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

	db, err := sql.Open("sqlite3", "scheduler.db") // Убедитесь, что используете правильный путь
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Удаляем задачу
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	result, err := db.Exec(deleteQuery, id)
	if err != nil {
		http.Error(w, "Ошибка удаления задачи", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Ошибка получения количества удаленных строк", http.StatusInternalServerError)
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
