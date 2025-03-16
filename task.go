package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`             // Обязательный заголовок
	Comment string `json:"comment,omitempty"` // Комментарий
	Repeat  string `json:"repeat,omitempty"`  // Правило повторения
}

// Обработчик добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Разрешаем только POST-запросы
	if r.Method != http.MethodPost {
		log.Printf("Метод не поддерживается") // Логируем ошибку
		JSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	decoder := json.NewDecoder(r.Body)
	// Декодируем JSON-запрос в структуру Task
	if err := decoder.Decode(&task); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	// Получаем текущую дату
	today := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = today
	} else if _, err := time.Parse("20060102", task.Date); err != nil {
		log.Printf("Неверный формат даты: %v", err) // Логируем ошибку
		JSONError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	// Проверяем, что заголовок задачи указан
	if task.Title == "" {
		log.Printf("Не указан заголовок задачи") // Логируем ошибку
		JSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}
	// Проверяем правило повторения на пустоту, если дата в прошлом возвращаем сегодняшний день
	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		}
	}

	// Проверяем правило повторения, если дата в прошлом, пытаемся вычислить новую дату
	if task.Date < today {
		newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Printf("Ошибка в правиле повторения: %v", err) // Логируем ошибку
			JSONError(w, "Ошибка в правиле повторения", http.StatusBadRequest)
			return
		}
		task.Date = newDate
	}

	// SQL-запрос для добавления задачи
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("Ошибка добавления задачи: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка добавления задачи", http.StatusInternalServerError)
		return
	}

	// Получаем ID последней добавленной задачи
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Ошибка получения ID задачи: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка получения ID задачи", http.StatusInternalServerError)
		return
	}

	// Формируем JSON-ответ
	response := map[string]string{"id": fmt.Sprintf("%d", id)}

	// Отправляем успешный JSON-ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Фунуция для полученя списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Разрешаем только GET-запросы
	if r.Method != http.MethodGet {
		log.Printf("Метод не поддерживается") // Логируем ошибку
		JSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем текущую дату в формате YYYYMMDD
	today := time.Now().Format("20060102")

	// SQL-запрос: выбираем ближайшие задачи (до 50 записей)
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= ? ORDER BY date ASC LIMIT 50`
	rows, err := db.Query(query, today)
	if err != nil {
		log.Printf("Ошибка запроса к базе данных: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка запроса к базе данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Список задач
	var tasks []Task

	// Обрабатываем строки результата запроса
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Printf("Ошибка обработки данных: %v", err) // Логируем ошибку
			JSONError(w, "Ошибка обработки данных", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
		if err := rows.Err(); err != nil {
			log.Printf("ООшибка обработки данных: %v", err)
			JSONError(w, "Ошибка обработки данных", http.StatusInternalServerError)
			return
		}
	}

	// Проверяем, пуст ли список задач
	if tasks == nil {
		tasks = []Task{}
	}

	// Формируем JSON-ответ
	response := map[string]interface{}{
		"tasks": tasks,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

// Обработчик получения задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Получаем параметр id из запроса
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		log.Printf("Не указан идентификатор") // Логируем ошибку
		JSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// SQL-запрос для получения задачи по ID
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, taskID)

	// Создаем объект задачи
	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Задача не найдена: %v", err) // Логируем ошибку
			JSONError(w, "Задача не найдена", http.StatusNotFound)
		} else {
			log.Printf("Ошибка получения данных: %v", err) // Логируем ошибку
			JSONError(w, "Ошибка получения данных", http.StatusInternalServerError)
		}
		return
	}

	// Отправляем JSON-ответ с задачей
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var task Task

	// Декодируем JSON-запрос
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	// Проверка: ID должен быть указан
	if task.ID == "" {
		log.Printf("Не указан идентификатор задачи:") // Логируем ошибку
		JSONError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	// Проверка: заголовок не может быть пустым
	if task.Title == "" {
		log.Printf("Не указан заголовок задачи:") // Логируем ошибку
		JSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Проверка: дата должна быть в формате YYYYMMDD
	if _, err := time.Parse("20060102", task.Date); err != nil {
		log.Printf("Неверный формат даты: %v", err) // Логируем ошибку
		JSONError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	// Вызываем NextDate для проверки Repeat и даты , чтобы она не была меньше нынешней
	if task.Repeat != "" {
		today := time.Now().Format("20060102")
		if task.Date < today {
			newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				log.Printf("Ошибка в правиле повторения: %v", err) // Логируем ошибку
				JSONError(w, "Ошибка в правиле повторения", http.StatusBadRequest)
				return
			}
			task.Date = newDate
		} else {
			// проверяем Repeat
			_, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				log.Printf("Ошибка в правиле повторения: %v", err) // Логируем ошибку
				JSONError(w, "Ошибка в правиле повторения", http.StatusBadRequest)
				return
			}
		}

	}

	// Обновляем запись
	query := "UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		log.Printf("Ошибка обновления задачи: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка обновления задачи", http.StatusInternalServerError)
		return
	}

	// Проверяем, была ли обновлена хотя бы одна строка
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Ошибка проверки обновления задачи: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка проверки обновления задачи", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		log.Printf("Задача не найдена: %v", err) // Логируем ошибку
		JSONError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

// функция о выполнении задачи
func taskDoneHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Разрешаем только POST-запросы
	if r.Method != http.MethodPost {
		log.Printf("Метод не поддерживается") // Логируем ошибку
		JSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID задачи из параметров запроса
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		log.Printf("Не указан идентификатор задачи") // Логируем ошибку
		JSONError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	// Получаем данные задачи
	var task Task
	query := "SELECT id, date, repeat FROM scheduler WHERE id = ?"
	err := db.QueryRow(query, taskID).Scan(&task.ID, &task.Date, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Задача не найдена: %v", err) // Логируем ошибку
			JSONError(w, "Задача не найдена", http.StatusNotFound)
		} else {
			log.Printf("Ошибка получения данных: %v", err) // Логируем ошибку
			JSONError(w, "Ошибка получения данных", http.StatusInternalServerError)
		}
		return
	}

	// Если у задачи нет правила повторения — удаляем её
	if task.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", taskID)
		if err != nil {
			log.Printf("Ошибка удаления задачи: %v", err) // Логируем ошибку
			JSONError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
		return
	}

	// Вычисляем новую дату выполнения
	newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		log.Printf("Ошибка в правиле повторения: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка в правиле повторения", http.StatusBadRequest)
		return
	}

	// Обновляем дату выполнения задачи
	_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, taskID)
	if err != nil {
		log.Printf("Ошибка обновления задачи: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка обновления задачи", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}

// Функция для отправки ошибок в JSON-формате
func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// функция для удаления задач
func deleteTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Получаем идентификатор задачи
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		log.Printf("Не указан идентификатор задачи:") // Логируем ошибку
		JSONError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	// Удаляем задачу
	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := db.Exec(query, taskID)
	if err != nil {
		log.Printf("Ошибка удаления задачи: %v", err) // Логируем ошибку
		JSONError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
		return
	}

	// Проверяем, была ли удалена хотя бы одна строка
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Printf("Задача не найдена: %v", err) // Логируем ошибку
		JSONError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
