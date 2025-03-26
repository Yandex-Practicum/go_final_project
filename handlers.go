package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)


func TaskHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    switch r.Method {
    case http.MethodGet:
        GetTasksHandler(w, r)
    case http.MethodPost:
        AddTaskHandler(w, r)
    case http.MethodPut:
        UpdateTaskHandler(w, r)
    case http.MethodDelete:
        handleDeleteTask(w, r)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
    }
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req struct {
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
		return
	}

	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task title is required"})
		return
	}

	now := time.Now()
	currentDate := now.Format("20060102")
	useDate := req.Date

	// Обработка специального значения "today"
	if useDate == "today" {
		useDate = currentDate
	}

	// Если дата не указана - используем сегодняшнюю
	if useDate == "" {
		useDate = currentDate
	} else {
		// Проверяем формат даты (если не "today")
		if _, err := time.Parse("20060102", useDate); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format (expected YYYYMMDD)"})
			return
		}
	}

	// Если дата в прошлом
	if useDate < currentDate {
		if req.Repeat != "" {
			// Для повторяющихся задач вычисляем следующую валидную дату
			next, err := NextDate(now, useDate, req.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			useDate = next
		} else {
			// Для разовых задач используем сегодняшнюю дату
			useDate = currentDate
		}
	}

	// Проверка правила повторения
	if req.Repeat != "" {
		if _, err := NextDate(now, useDate, req.Repeat); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
	}

	// Сохраняем в БД
	res, err := DB.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		useDate, req.Title, req.Comment, req.Repeat,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}


type Task struct {
    ID      string  `json:"id"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    response := struct {
        Tasks []Task `json:"tasks"`
    }{
        Tasks: []Task{},
    }

    rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    defer rows.Close()

    for rows.Next() {
        var id int64
        var task Task
        if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }
        task.ID = strconv.FormatInt(id, 10) // Конвертируем int64 в string
        response.Tasks = append(response.Tasks, task)
    }

    json.NewEncoder(w).Encode(response)
}

// Добавляем в handlers.go

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPut {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
        return
    }

    var task Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
        return
    }

    if task.ID == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор"})
        return
    }

    if task.Title == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Task title is required"})
        return
    }

    // Проверка формата даты
    if _, err := time.Parse("20060102", task.Date); err != nil && task.Date != "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format (expected YYYYMMDD)"})
        return
    }

    // Проверка правила повторения (используем существующую функцию NextDate)
    if task.Repeat != "" {
        now := time.Now()
        if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }
    }

    // Конвертируем строковый ID в int64
    id, err := strconv.ParseInt(task.ID, 10, 64)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
        return
    }

    // Обновляем задачу в БД
    res, err := DB.Exec(
        "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
        task.Date, task.Title, task.Comment, task.Repeat, id,
    )
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
        return
    }

    // Проверяем, что задача была действительно обновлена
    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Задача не найдена"})
        return
    }

    // Возвращаем пустой JSON в случае успеха
    json.NewEncoder(w).Encode(map[string]interface{}{})
}


// Обработчик для POST /api/task/done
func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id := r.URL.Query().Get("id")
    if id == "" {
        json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор"})
        return
    }

    // Получаем данные задачи
    var task struct {
        Date   string
        Repeat string
    }
    err := DB.QueryRow("SELECT date, repeat FROM scheduler WHERE id = ?", id).
        Scan(&task.Date, &task.Repeat)
    
    if err != nil {
        if err == sql.ErrNoRows {
            json.NewEncoder(w).Encode(map[string]string{"error": "Задача не найдена"})
        } else {
            json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка базы данных"})
        }
        return
    }

    if task.Repeat == "" {
        // Удаляем одноразовую задачу
        _, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
        if err != nil {
            json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при удалении задачи"})
            return
        }
    } else {
        // Для повторяющейся задачи вычисляем следующую дату
        nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
        if err != nil {
            json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при расчете следующей даты"})
            return
        }

        // Обновляем дату выполнения
        _, err = DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
        if err != nil {
            json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при обновлении задачи"})
            return
        }
    }

    // Успешный ответ
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(struct{}{})
}

// Обработчик для DELETE /api/task
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id := r.URL.Query().Get("id")
    if id == "" {
        json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор"})
        return
    }

    // Проверяем существование задачи
    var exists bool
    err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM scheduler WHERE id = ?)", id).Scan(&exists)
    if err != nil || !exists {
        json.NewEncoder(w).Encode(map[string]string{"error": "Задача не найдена"})
        return
    }

    // Удаляем задачу
    _, err = DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
    if err != nil {
        json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при удалении задачи"})
        return
    }

    // Успешный ответ
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(struct{}{})
}