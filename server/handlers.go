package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/tasks"
	"net/http"
	"strconv"
	"time"
)

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка: Неверный формат даты 'now': %s", err), http.StatusBadRequest)
		return
	}

	nextDate, err := tasks.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка: %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", nextDate)
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task tasks.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, "Ошибка десериализации JSON: ", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format(tasks.Format)
	} else {
		_, err = time.Parse(tasks.Format, task.Date)
		if err != nil {
			writeError(w, "Неверный формат даты, ожидается формат "+tasks.Format, http.StatusBadRequest)
			return
		}
	}

	if task.Repeat == "" {
		if task.Date < time.Now().Format(tasks.Format) {
			task.Date = time.Now().Format(tasks.Format)
		}
	} else {
		nextDate, err := tasks.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Неверный формат правила повторения: "+err.Error(), http.StatusBadRequest)
			return
		}
		if task.Date < time.Now().Format(tasks.Format) {
			task.Date = nextDate
		}
	}

	id, err := tasks.AddTask(task)
	if err != nil {
		writeError(w, "Ошибка сохранения задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeResponse(w, map[string]int64{"id": id}, http.StatusCreated)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}
	task, err := tasks.GetTaskByID(strconv.Itoa(id))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeResponse(w, task, http.StatusOK)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var updatedTask tasks.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask, err = tasks.FormatTask(updatedTask)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tasks.PutTask(updatedTask)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, map[string]string{}, http.StatusCreated)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasksList, err := tasks.GetTasks("./scheduler.db", search)
	if err != nil {
		writeError(w, "Ошибка поиска задач: ", http.StatusInternalServerError)
		return
	}

	writeResponse(w, map[string][]tasks.Task{"tasks": tasksList}, http.StatusOK)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		writeError(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := tasks.DeleteTask(strconv.Itoa(id)); err != nil {
		writeError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
		return
	}

	writeResponse(w, struct{}{}, http.StatusOK)
}

func MarkTasksAsDoneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	task, err := tasks.GetTaskByID(strconv.Itoa(id))
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "Задача не найдена", http.StatusNotFound)
		} else {
			writeError(w, "Ошибка получения задачи", http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat != "" {
		nextDate, err := tasks.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, "Ошибка расчета следующей даты", http.StatusInternalServerError)
			return
		}
		task.Date = nextDate
		err = tasks.PutTask(task)
		if err != nil {
			writeError(w, "Ошибка обновления задачи", http.StatusInternalServerError)
			return
		}
	} else {
		err = tasks.DeleteTask(strconv.Itoa(id))
		if err != nil {
			writeError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
			return
		}
	}

	writeResponse(w, struct{}{}, http.StatusOK)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
