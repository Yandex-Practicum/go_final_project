package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sclea3/go_final_project/models"
)

// writeJSONError отправляет JSON-ответ с полем "error".
func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// writeTextError отправляет текстовый ответ.
func writeTextError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	io.WriteString(w, msg)
}

// NextDate вычисляет следующую дату для периодической задачи.
// Поддерживаются базовые правила:
//   - "y": добавляем год (повтор каждый год);
//   - "d <число>": добавляем указанное число дней (от 1 до 400).
//
// Параметр now задаёт базис для сравнения (дата должна быть строго больше now).
func NextDate(now time.Time, dateStr, repeat string) (string, error) {
	original, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", errors.New("Неверный формат даты")
	}
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", nil
	}
	if repeat == "y" {
		candidate := original
		// Цикл: пока candidate не больше now, добавляем 1 год.
		for !candidate.After(now) {
			candidate = candidate.AddDate(1, 0, 0)
		}
		return candidate.Format("20060102"), nil
	}
	const prefix = "d "
	if len(repeat) > len(prefix) && repeat[:len(prefix)] == prefix {
		dayInterval, err := strconv.Atoi(repeat[len(prefix):])
		if err != nil {
			return "", errors.New("Неверный интервал для d")
		}
		if dayInterval < 1 || dayInterval > 400 {
			return "", errors.New("Интервал должен быть от 1 до 400")
		}
		candidate := original
		// Для вычисления следующей даты базируемся на исходной дате.
		for !candidate.After(now) {
			candidate = candidate.AddDate(0, 0, dayInterval)
		}
		return candidate.Format("20060102"), nil
	}
	return "", errors.New("Неподдерживаемый формат правила повторения")
}

// NextDateHandler обрабатывает GET-запрос /api/nextdate.
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := strings.TrimSpace(r.URL.Query().Get("repeat"))
	if nowStr == "" || dateStr == "" {
		writeTextError(w, http.StatusBadRequest, "Не указаны все параметры")
		return
	}
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		writeTextError(w, http.StatusBadRequest, "Неверный формат now")
		return
	}
	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		writeTextError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Возвращаем plain text (как требуется тестами).
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, next)
}

// TaskHandler обрабатывает запросы к /api/task.
// GET: возвращает задачу по id.
// POST: добавляет новую задачу.
// PUT: обновляет задачу.
// DELETE: удаляет задачу.
func TaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodGet:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			writeJSONError(w, http.StatusBadRequest, "Не указан идентификатор")
			return
		}
		var task models.Task
		row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", idStr)
		if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			writeJSONError(w, http.StatusNotFound, "Задача не найдена")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(task)
	case http.MethodPost:
		// Добавление задачи
		var newTask models.NewTask
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		if newTask.Title == "" {
			writeJSONError(w, http.StatusBadRequest, "Не указан заголовок задачи")
			return
		}
		if newTask.Date == "" || newTask.Date == "today" {
			newTask.Date = time.Now().Format("20060102")
		} else {
			if _, err := time.Parse("20060102", newTask.Date); err != nil {
				writeJSONError(w, http.StatusBadRequest, "Неверный формат даты")
				return
			}
		}
		// Если дата меньше сегодняшней и задача повторяется, вычисляем новую дату.
		todayStr := time.Now().Format("20060102")
		if newTask.Date < todayStr && newTask.Repeat != "" {
			next, err := NextDate(time.Now(), newTask.Date, newTask.Repeat)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				return
			}
			newTask.Date = next
		} else if newTask.Date < todayStr {
			newTask.Date = todayStr
		}
		// Валидируем правило повторения, если оно указано.
		if newTask.Repeat != "" {
			if _, err := NextDate(time.Now(), newTask.Date, newTask.Repeat); err != nil {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
			newTask.Date, newTask.Title, newTask.Comment, newTask.Repeat)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to get inserted ID")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	case http.MethodPut:
		// Обновление задачи
		var updatedTask struct {
			ID      int64  `json:"id,string"`
			Date    string `json:"date"`
			Title   string `json:"title"`
			Comment string `json:"comment"`
			Repeat  string `json:"repeat"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		if updatedTask.ID == 0 {
			writeJSONError(w, http.StatusBadRequest, "Не указан идентификатор")
			return
		}
		if updatedTask.Title == "" {
			writeJSONError(w, http.StatusBadRequest, "Не указан заголовок задачи")
			return
		}
		if _, err := time.Parse("20060102", updatedTask.Date); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Неверный формат даты")
			return
		}
		// Если правило повторения указано, проверяем его корректность.
		if updatedTask.Repeat != "" {
			if _, err := NextDate(time.Now(), updatedTask.Date, updatedTask.Repeat); err != nil {
				writeJSONError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		result, err := db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
			updatedTask.Date, updatedTask.Title, updatedTask.Comment, updatedTask.Repeat, updatedTask.ID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			writeJSONError(w, http.StatusNotFound, "Задача не найдена")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	case http.MethodDelete:
		// Удаление задачи
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			writeJSONError(w, http.StatusBadRequest, "Не указан идентификатор")
			return
		}
		result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", idStr)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			writeJSONError(w, http.StatusNotFound, "Задача не найдена")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

// ListTasksHandler возвращает список задач в виде {"tasks": [...]}
func ListTasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []models.Task{}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}

// TaskDoneHandler обрабатывает POST-запрос /api/task/done?id=<идентификатор>.
// Если задача одноразовая (repeat пуст), удаляет её.
// Если задача периодическая, вычисляет новую дату следующего выполнения
// и обновляет запись. При этом базис для вычисления берётся из хранимой даты.
func TaskDoneHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}
	var t models.Task
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", idStr)
	if err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		writeJSONError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	// Если задача не повторяется, удаляем её.
	if t.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", idStr)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	// Для периодической задачи используем хранимую дату как базис.
	baseDate, err := time.Parse("20060102", t.Date)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Неверный формат даты в задаче")
		return
	}
	newDate, err := NextDate(baseDate, t.Date, t.Repeat)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, idStr)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
