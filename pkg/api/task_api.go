package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	//Читаем id задачи
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID not specified", http.StatusBadRequest)
		return
	}

	//Получение задачи из БД
	t, err := db.GetTask(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}

	//Повторения нет - удаляем
	if t.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			http.Error(w, "Task not found", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(`{}`))
		return
	}

	//Парсим дату
	prev, err := time.Parse(Layout, t.Date)
	if err != nil {
		http.Error(w, "Incorrect date format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(prev, t.Date, t.Repeat)
	if err != nil {
		http.Error(w, "Incorrect rule of repeat", http.StatusBadRequest)
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		http.Error(w, "Failed to edit date", http.StatusBadRequest)
		return
	}

	//Возвращаем пустого Джейсона
	_, _ = w.Write([]byte(`{}`))
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		task, err := db.GetTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, task)
	case http.MethodDelete:

		id := r.URL.Query().Get("id")

		if err := db.DeleteTask(id); err != nil {

			writeJson(w, map[string]string{"error": err.Error()})
			return
		}

		_, _ = w.Write([]byte(`{}`))

	case http.MethodPut:
		var t db.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			writeJson(w, map[string]string{"error": "Invalid JSON format"})
			return
		}

		//Проверяем ID и заголовок задачи
		if t.ID == "" {
			writeJson(w, map[string]string{"error": "The ID field cannot be empty."})
			return
		}
		if t.Title == "" {
			writeJson(w, map[string]string{"error": "The ID field cannot be empty."})
			return
		}

		//Проверяем срок годности
		if _, err := time.Parse(Layout, t.Date); err != nil {
			writeJson(w, map[string]string{"error": "incorrect date format"})

			return
		}
		if t.Repeat != "" {
			if _, err := NextDate(time.Now(), t.Date, t.Repeat); err != nil {
				writeJson(w, map[string]string{"error": "Incorrect rule of repeat"})

				return
			}
		}

		//Исправляем дату если время пролетело
		parsed, _ := time.Parse(Layout, t.Date)
		if t.Repeat != "" {
			nextDate, _ := NextDate(time.Now(), t.Date, t.Repeat)
			if parsed.Before(time.Now()) {
				t.Date = nextDate
			}
			parsed.Before(time.Now())
			t.Date = time.Now().Format(Layout)
		}

		//сохраняем
		if err := db.UpdateTask(&t); err != nil {
			writeJson(w, map[string]string{"error": "Failed to edit date"})
			return
		}
		_, _ = w.Write([]byte(`{}`))
	default:
		writeJson(w, map[string]string{"error": "Method not allowed"})

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// AddTaskHandler Отвечает за добавление задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	//Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": "JSON parsing error"})
		return
	}

	//Десериализация
	if err := json.Unmarshal(body, &task); err != nil {
		writeJson(w, map[string]string{"error": "JSON parsing error"})
		return
	}

	//Валидация
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "The title cannot be empty"})
		return
	}

	//Установим дату "Сегодня"
	now := time.Now()
	nowDate := now.Format(Layout)

	if task.Date == "" {
		task.Date = nowDate
	}

	//Проверка даты
	t, err := time.Parse(Layout, task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "Invalid date format"})
		return
	}

	today := time.Now().Format(Layout)
	if task.Date < today {
		task.Date = today
	}

	taskDate := t.Format(Layout)

	//Обработка повторения задач
	if task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			writeJson(w, map[string]string{"error": "Incorrect rule of repeat"})

			return
		}

		if taskDate < nowDate {
			nextDate, _ := NextDate(now, task.Date, task.Repeat)
			task.Date = nextDate
		}

		if taskDate < nowDate {
			task.Date = nowDate
		}
	}

	//Сохраняем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Error adding task"})
		return
	}

	// Возвращаем ID новой задачи
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	payload := map[string]any{"id": id}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Encode response: %v", err)
	}
}
