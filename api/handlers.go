package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sclea3/go_final_project/models"
	"github.com/Sclea3/go_final_project/scheduler"
)

const DateLayout = scheduler.DateLayout

// DB – глобальное подключение к базе данных, которое устанавливается из main.
var DB *sql.DB

// respondWithJSON отправляет JSON-ответ.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError отправляет ошибку в формате JSON.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// NextDateHandler обрабатывает GET-запросы по маршруту /api/nextdate.
// Ожидает параметры now, date и repeat.
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(DateLayout, nowStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Неверный формат параметра now")
		return
	}

	next, err := scheduler.NextDate(now, dateStr, repeat)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"date": next})
}

// TaskHandler обрабатывает запросы по маршруту /api/task.
// Поддерживает GET (получение задачи по id), POST (добавление задачи) и PUT (редактирование задачи).
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTask(w, r)
	case http.MethodPost:
		addTask(w, r)
	case http.MethodPut:
		editTask(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

// getTask возвращает параметры задачи по id.
func getTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	var task models.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Задача не найдена")
		return
	}
	respondWithJSON(w, http.StatusOK, task)
}

// addTask добавляет новую задачу.
func addTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Ошибка десериализации JSON")
		return
	}
	if strings.TrimSpace(task.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Не указан заголовок задачи")
		return
	}

	// Если дата не указана, берём сегодняшнюю дату.
	today := time.Now().Format(DateLayout)
	if strings.TrimSpace(task.Date) == "" {
		task.Date = today
	} else {
		parsedDate, err := time.Parse(DateLayout, task.Date)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат даты")
			return
		}
		todayTime, _ := time.Parse(DateLayout, today)
		if parsedDate.Before(todayTime) {
			if strings.TrimSpace(task.Repeat) != "" {
				next, err := scheduler.NextDate(todayTime, task.Date, task.Repeat)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				task.Date = next
			} else {
				task.Date = today
			}
		}
	}

	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(lastID, 10)})
}

// editTask обновляет существующую задачу.
func editTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Ошибка десериализации JSON")
		return
	}
	if strings.TrimSpace(task.ID) == "" {
		respondWithError(w, http.StatusBadRequest, "Не указан идентификатор задачи")
		return
	}
	if strings.TrimSpace(task.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Не указан заголовок задачи")
		return
	}

	today := time.Now().Format(DateLayout)
	todayTime, _ := time.Parse(DateLayout, today)
	parsedDate, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Неверный формат даты")
		return
	}
	if parsedDate.Before(todayTime) {
		if strings.TrimSpace(task.Repeat) != "" {
			next, err := scheduler.NextDate(todayTime, task.Date, task.Repeat)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			task.Date = next
		} else {
			task.Date = today
		}
	}

	res, err := DB.Exec("UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		respondWithError(w, http.StatusNotFound, "Задача не найдена")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{})
}

// TasksHandler возвращает список задач. Поддерживается фильтрация по параметру search.
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var rows *sql.Rows
	var err error
	limit := 50

	// Если параметр search соответствует формату даты "02.01.2006", преобразуем его в формат DateLayout.
	if t, errParse := time.Parse("02.01.2006", search); errParse == nil {
		search = t.Format(DateLayout)
		rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?", search, limit)
	} else if search != "" {
		likeParam := "%" + search + "%"
		rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?", likeParam, likeParam, limit)
	} else {
		rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		tasks = append(tasks, t)
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"tasks": tasks})
}

// DoneHandler обрабатывает запросы по маршруту /api/task/done для отметки задачи как выполненной (POST)
// или удаления задачи (DELETE).
func DoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}
	switch r.Method {
	case http.MethodPost:
		// Отметить задачу как выполненную.
		row := DB.QueryRow("SELECT id, date, repeat FROM scheduler WHERE id = ?", id)
		var task models.Task
		err := row.Scan(&task.ID, &task.Date, &task.Repeat)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Задача не найдена")
			return
		}
		todayTime, _ := time.Parse(DateLayout, time.Now().Format(DateLayout))
		if strings.TrimSpace(task.Repeat) == "" {
			// Одноразовая задача – удаляем.
			_, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			// Повторяющаяся задача – вычисляем новую дату.
			newDate, err := scheduler.NextDate(todayTime, task.Date, task.Repeat)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			_, err = DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, id)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		respondWithJSON(w, http.StatusOK, map[string]string{})
	case http.MethodDelete:
		_, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{})
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}
