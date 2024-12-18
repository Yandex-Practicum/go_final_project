package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"final_project/database"
	"final_project/models"
)

func main() {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addTaskHandler(w, r)
		case http.MethodGet:
			getTaskHandler(w, r)
		case http.MethodPut:
			updateTaskHandler(w, r)
		case http.MethodDelete:
			deleteTaskHandler(w, r)
		}
	})
	http.HandleFunc("/api/task/done", markTaskAsDoneHandler)
	http.HandleFunc("/api/nextdate", handler)
	http.HandleFunc("/api/tasks", getTasksHandler)

	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", errors.New("некорректная дата")
	}

	var nextDate time.Time

	switch {
	case repeat == "y":
		nextDate = date.AddDate(1, 0, 0)

	case strings.HasPrefix(repeat, "d "):
		daysStr := strings.TrimPrefix(repeat, "d ")
		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("некорректное правило повторения: " + repeat)
		}

		nextDate = date
		for {
			if nextDate.After(now) {
				break
			}
			nextDate = nextDate.AddDate(0, 0, days)
		}

	default:
		return "", errors.New("некорректное правило повторения: " + repeat)
	}

	return nextDate.Format("20060102"), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Некорректная текущая дата", http.StatusBadRequest)
		return
	}

	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("handle:", nextDate)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	var taskDate time.Time
	var err error
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}
	taskDate, err = time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, "Дата представлена в неправильном формате", http.StatusBadRequest)
		return
	}

	today := time.Now()
	if taskDate.Before(today) {
		if task.Repeat == "" {
			task.Date = today.Format("20060102")
			taskDate = today
		} else {
			nextDateStr, err := NextDate(today, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			taskDate, err = time.Parse("20060102", nextDateStr)
			if err != nil {
				http.Error(w, "Ошибка при вычислении следующей даты", http.StatusInternalServerError)
				return
			}
		}
	}

	id, err := database.AddTask(taskDate.Format("20060102"), task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var taskResponses []models.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, models.TaskResponse{
			ID:      strconv.FormatInt(task.ID, 10),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	response := models.TasksResponse{Tasks: taskResponses}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := database.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	response := models.TaskResponse{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.ID == 0 {
		http.Error(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}
	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, "Дата представлена в неправильном формате", http.StatusBadRequest)
		return
	}

	err = database.UpdateTask(task)
	if err != nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func markTaskAsDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := database.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err = database.DeleteTask(id)
		if err != nil {
			http.Error(w, "Ошибка при удалении задачи", http.StatusInternalServerError)
			return
		}
	} else {
		nextDateStr, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task.Date = nextDateStr

		err = database.UpdateTask(task)
		if err != nil {
			http.Error(w, "Ошибка при обновлении задачи", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	err := database.DeleteTask(id)
	if err != nil {
		http.Error(w, "Ошибка при удалении задачи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}
