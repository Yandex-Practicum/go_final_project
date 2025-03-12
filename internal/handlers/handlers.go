package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go_final-project/internal/auth"
	"go_final-project/internal/db"
	"go_final-project/internal/logic"
	"go_final-project/internal/task"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetTasksHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, `{"error":"Only GET method is supported"}`, http.StatusMethodNotAllowed)
			return
		}
		search := req.URL.Query().Get("search")
		var dateSearch string

		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			dateSearch = parsedDate.Format("20060102")
		}

		tasks, err := db.GetTasks(dbase, search, dateSearch)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []task.Task{}
		}

		tasksList := make([]map[string]interface{}, len(tasks))
		for i, t := range tasks {
			tasksList[i] = map[string]interface{}{
				"id":      strconv.FormatInt(t.ID, 10),
				"date":    t.Date,
				"title":   t.Title,
				"comment": t.Comment,
				"repeat":  t.Repeat,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasksList})
	}
}

func TaskHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			getTask(dbase, w, req)
		case http.MethodPost:
			addTask(dbase, w, req)
		case http.MethodPut:
			updateTask(dbase, w, req)
		case http.MethodDelete:
			deleteTask(dbase, w, req)
		default:
			sendJSONError(w, "Only GET, POST, PUT, DELETE methods are supported.", http.StatusMethodNotAllowed)
		}
	}
}

func getTask(dbase *sqlx.DB, w http.ResponseWriter, req *http.Request) {
	// Проверяем наличие параметра id
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		sendJSONError(w, "id is required", http.StatusBadRequest)
		return
	}

	// Преобразую id в число
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendJSONError(w, "id error format", http.StatusBadRequest)
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTaskByID(dbase, id)
	if err != nil {
		sendJSONError(w, "issue not found", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"id":      strconv.FormatInt(task.ID, 10),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}

	// Отправляю JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func addTask(dbase *sqlx.DB, w http.ResponseWriter, req *http.Request) {
	var newTask task.Task

	// Декодирую JSON
	err := json.NewDecoder(req.Body).Decode(&newTask)
	if err != nil {
		sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Проверяю есть ли заголовок
	if newTask.Title == "" {
		sendJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}
	now := time.Now().Truncate(24 * time.Hour)
	today := now.Format("20060102")

	// Если дата не передана, ставим сегодняшнюю
	if newTask.Date == "" {
		newTask.Date = today
	} else {
		if _, err := time.Parse("20060102", newTask.Date); err != nil {
			sendJSONError(w, "date error format", http.StatusBadRequest)
			return
		}
	}
	taskDate, _ := time.Parse("20060102", newTask.Date)

	// Обрабатываю повторяющиеся задачи
	if newTask.Repeat != "" {
		if taskDate.Before(now) {
			nextDate, err := logic.NextDate(now, newTask.Date, newTask.Repeat)
			if err != nil || nextDate == "" {
				sendJSONError(w, "Invalid repeat format or no valid next date found", http.StatusBadRequest)
				return
			}
			newTask.Date = nextDate
		}
	} else if taskDate.Before(now) {
		newTask.Date = today
	}

	// Добавляю задачу в БД
	id, err := db.AddTask(dbase, &newTask)
	if err != nil {
		sendJSONError(w, "Failed to save task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": strconv.FormatInt(id, 10)})
}

func updateTask(dbase *sqlx.DB, w http.ResponseWriter, req *http.Request) {
	// Декодирую JSON в промежуточную структуру
	var rawData map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&rawData)
	if err != nil {
		sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Проверяю и конвертирую ID
	idValue, ok := rawData["id"]
	if !ok {
		sendJSONError(w, "id is required", http.StatusBadRequest)
		return
	}

	var id int64
	switch v := idValue.(type) {
	case string:
		id, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			sendJSONError(w, "id error format", http.StatusBadRequest)
			return
		}
	case float64:
		id = int64(v)
	default:
		sendJSONError(w, "id error format", http.StatusBadRequest)
		return
	}

	// Удаляю строковый ID из исходного JSON и подставляю int64 ID
	rawData["id"] = id

	// Преобразую обратно в JSON
	updatedJSON, err := json.Marshal(rawData)
	if err != nil {
		sendJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Декодирую JSON в Task
	var updatedTask task.Task
	err = json.Unmarshal(updatedJSON, &updatedTask)
	if err != nil {
		sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Загружаю текущую версию задачи из базы
	existingTask, err := db.GetTaskByID(dbase, id)
	if err != nil {
		sendJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	// Проверяю, что title передан и не пустой
	if updatedTask.Title == "" {
		sendJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Обновляю только переданные поля, сохраняя старые значения
	if updatedTask.Date == "" {
		updatedTask.Date = existingTask.Date
	}
	if updatedTask.Comment == "" {
		updatedTask.Comment = existingTask.Comment
	}
	if updatedTask.Repeat == "" {
		updatedTask.Repeat = existingTask.Repeat
	}

	// Проверяю корректность даты
	if _, err := time.Parse("20060102", updatedTask.Date); err != nil {
		sendJSONError(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Проверяю корректность repeat
	if updatedTask.Repeat != "" {
		_, err := logic.NextDate(time.Now(), updatedTask.Date, updatedTask.Repeat)
		if err != nil {
			sendJSONError(w, "Invalid repeat format", http.StatusBadRequest)
			return
		}
	}

	// Обновляю задачу в БД
	err = db.UpdateTask(dbase, &updatedTask)
	if err != nil {
		sendJSONError(w, "issue not found", http.StatusBadRequest)
		return
	}

	// Отправляю пустой JSON (успешное выполнение)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{})
}

func MarkTaskDoneHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			sendJSONError(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		idStr := req.URL.Query().Get("id")
		if idStr == "" {
			sendJSONError(w, "id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			sendJSONError(w, "id error format", http.StatusBadRequest)
			return
		}

		// Получаю задачу из БД
		task, err := db.GetTaskByID(dbase, id)
		if err != nil {
			sendJSONError(w, "task not found", http.StatusNotFound)
			return
		}

		// Если у задачи нет повторения - удаляем
		if task.Repeat == "" {
			err = db.DeleteTask(dbase, id)
			if err != nil {
				sendJSONError(w, "Failed to delete task", http.StatusInternalServerError)
				return
			}
		} else {
			// Рассчитываю новую дату выполнения
			nextDate, err := logic.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				sendJSONError(w, "Failed to calculate next date", http.StatusBadRequest)
				return
			}
			task.Date = nextDate

			// Обновляю задачу
			err = db.UpdateTask(dbase, task)
			if err != nil {
				sendJSONError(w, "Failed to update task", http.StatusInternalServerError)
				return
			}
		}

		// Отправляю успешный JSON-ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func deleteTask(dbase *sqlx.DB, w http.ResponseWriter, req *http.Request) {
	// Проверяю, передан ли ID
	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		sendJSONError(w, "id is required", http.StatusBadRequest)
		return
	}

	// Преобразую ID в число
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		sendJSONError(w, "id error format", http.StatusBadRequest)
		return
	}

	// Удаляю задачу
	err = db.DeleteTask(dbase, id)
	if err != nil {
		sendJSONError(w, "Task not found", http.StatusNotFound)
		return
	}

	// Отправляю пустой JSON (успешное выполнение)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{})
}

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// SignInHandler обработчик аутентификации
func SignInHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid method"}`, http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if request.Password != os.Getenv("TODO_PASSWORD") {
		sendJSONError(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		sendJSONError(w, "Token error", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
