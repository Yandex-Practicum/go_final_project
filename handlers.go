package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Task представляет задачу
type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskWithID struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// nextDateHandler обрабатывает запрос для расчета следующей даты задачи
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметров из строки запроса
	nowParam := r.URL.Query().Get("now")
	repeatParam := r.URL.Query().Get("repeat")
	dateParam := r.URL.Query().Get("date")

	// Проверка и парсинг параметра now
	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		http.Error(w, "Invalid 'now' date format, must be YYYYMMDD", http.StatusBadRequest)
		return
	}

	// Вызов функции NextDate для получения следующей даты задачи
	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ответ с результатом
	w.Write([]byte(nextDate))
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "JSON Deserialization Error", http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля title
	if task.Title == "" {
		sendErrorResponse(w, "Task title is not specified")
		return
	}

	// Определение текущей даты
	now := time.Now()
	today := now.Format(TimeFormat)

	// Проверка и обработка поля date
	if task.Date == "" {
		task.Date = today
	} else {
		parsedDate, err := time.Parse(TimeFormat, task.Date)
		if err != nil {
			sendErrorResponse(w, "Date is in an incorrect format")
			return
		}

		parsedDate = parsedDate.Truncate(24 * time.Hour)
		now = now.Truncate(24 * time.Hour)

		// Если дата меньше сегодняшней
		if parsedDate.Before(now) {
			if task.Repeat == "" {
				// Если нет правила повторения, подставляем сегодняшнюю дату
				task.Date = today
			} else {
				// Применяем правило повторения, чтобы получить следующую дату
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					sendErrorResponse(w, "Error in repetition rule")
					return
				}

				// Если правило повторения - "d 1" и сегодняшняя дата допустима, устанавливаем её
				if task.Repeat == "d 1" && nextDate == today {
					task.Date = today
				} else {
					task.Date = nextDate
				}
			}
		}
	}

	// Проверка правила повторения с использованием NextDate()
	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			sendErrorResponse(w, "Invalid repetition rule")
			return
		}
	}

	// Добавление задачи в базу данных
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		sendErrorResponse(w, "Error adding task to the database")
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		sendErrorResponse(w, "Error retrieving task ID")
		return
	}

	// Формируем успешный ответ с ID задачи
	response := map[string]string{"id": fmt.Sprintf("%d", id)}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse отправляет JSON-ответ с полем error
func sendErrorResponse(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Выполняем запрос для подсчета задач в таблице
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM scheduler").Scan(&count)
	if err != nil {
		log.Printf("Error counting tasks in database: %v", err)
		sendErrorResponse(w, "Error retrieving tasks count from database")
		return
	}

	// Инициализация массива задач, даже если их нет
	tasks := []TaskWithID{}

	// Если в таблице есть задачи, выполняем запрос для их получения
	if count > 0 {
		rows, err := db.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date ASC
			LIMIT 50
		`)
		if err != nil {
			log.Printf("Database query error: %v", err)
			sendErrorResponse(w, "Error retrieving tasks from database")
			return
		}
		defer rows.Close()

		// Заполняем массив задач, если записи найдены
		for rows.Next() {
			var task TaskWithID
			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				log.Printf("Error scanning row: %v", err)
				sendErrorResponse(w, "Error reading tasks from database")
				return
			}
			tasks = append(tasks, task)
		}

		// Проверка на наличие ошибок после завершения итерации
		if err = rows.Err(); err != nil {
			log.Printf("Row iteration error: %v", err)
			sendErrorResponse(w, "Error iterating tasks")
			return
		}
	}

	// Формирование JSON-ответа
	response := map[string]interface{}{
		"tasks": tasks,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("JSON encoding error: %v", err)
		sendErrorResponse(w, "Error encoding tasks to JSON")
		return
	}
}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(w, "Task ID is not specified")
		return
	}

	var task TaskWithID
	err := db.QueryRow(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Task not found")
		} else {
			log.Printf("Database error: %v", err)
			sendErrorResponse(w, "Error retrieving task")
		}
		return
	}

	// Формирование JSON-ответа с задачей
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("JSON encoding error: %v", err)
		sendErrorResponse(w, "Error encoding task to JSON")
	}
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var task TaskWithID
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "JSON Deserialization Error", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		sendErrorResponse(w, "Task ID is not specified")
		return
	}
	if task.Title == "" {
		sendErrorResponse(w, "Task title is not specified")
		return
	}

	// Проверка и обработка поля date
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(TimeFormat)
	} else {
		parsedDate, err := time.Parse(TimeFormat, task.Date)
		if err != nil {
			sendErrorResponse(w, "Date is in an incorrect format")
			return
		}
		if parsedDate.Before(now) && task.Repeat != "" {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				sendErrorResponse(w, "Error in repetition rule")
				return
			}
			task.Date = nextDate
		}
	}

	// Обновление записи в базе данных
	res, err := db.Exec(`
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`, task.Date, task.Title, task.Comment, task.Repeat, task.ID)

	if err != nil {
		log.Printf("Database error: %v", err)
		sendErrorResponse(w, "Error updating task")
		return
	}

	affected, err := res.RowsAffected()
	if err != nil || affected == 0 {
		sendErrorResponse(w, "Task not found")
		return
	}

	// Возвращаем пустой JSON-ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

// markTaskDoneHandler обрабатывает POST-запрос для отметки задачи выполненной
func markTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	// Получение идентификатора задачи
	id := r.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(w, "Task ID is not specified")
		return
	}

	// Извлечение задачи из базы данных
	var task TaskWithID
	err := db.QueryRow(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Task not found")
		} else {
			log.Printf("Database error: %v", err)
			sendErrorResponse(w, "Error retrieving task")
		}
		return
	}

	// Если задача одноразовая (repeat пуст), удаляем её
	if task.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			log.Printf("Error deleting task: %v", err)
			sendErrorResponse(w, "Error deleting task")
			return
		}
	} else {
		// Для периодической задачи вычисляем следующую дату
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Error calculating next date: %v", err)
			sendErrorResponse(w, "Error calculating next date for recurring task")
			return
		}

		// Обновляем дату выполнения задачи в базе данных
		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		if err != nil {
			log.Printf("Error updating task date: %v", err)
			sendErrorResponse(w, "Error updating task date")
			return
		}
	}

	// Возвращаем пустой JSON-ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

// deleteTaskHandler обрабатывает DELETE-запрос для удаления задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получение идентификатора задачи
	id := r.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(w, "Task ID is not specified")
		return
	}

	// Выполнение запроса на удаление и получение количества затронутых строк
	res, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		sendErrorResponse(w, "Error deleting task")
		return
	}

	// Проверяем, была ли удалена хотя бы одна строка
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving affected rows: %v", err)
		sendErrorResponse(w, "Error retrieving affected rows")
		return
	}

	if rowsAffected == 0 {
		// Если строк с таким id не найдено, возвращаем ошибку
		sendErrorResponse(w, "Task not found")
		return
	}

	// Возвращаем пустой JSON-ответ в случае успешного удаления
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}
