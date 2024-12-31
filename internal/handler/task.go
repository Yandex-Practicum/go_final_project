package handler

import (
	"encoding/json"
	"go_final_project/internal/models"
	"go_final_project/internal/repository"
	"go_final_project/internal/scheduler"
	"log"
	"net/http"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeErrorJSON(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if err := scheduler.ValidateAndProcessTask(&task, now); err != nil {
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := repository.AddTaskToDB(&task)
	if err != nil {
		writeErrorJSON(w, "ошибка добавления задачи в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	if err != nil {
		log.Printf("JSON Encoding Error: %v", err)
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeErrorJSON(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if err := scheduler.ValidateAndProcessTask(&task, now); err != nil {
		writeErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := repository.GetTaskById(task.Id)
	if err != nil {
		writeErrorJSON(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}

	err = repository.UpdateTaskToDB(&task)
	if err != nil {
		writeErrorJSON(w, "ошибка редактирования задачи в базе данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(map[string]string{})
	if err != nil {
		log.Printf("JSON Encoding Error: %v", err)
	}
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок для JSON-ответа
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Получаем параметры поиска из строки запроса
	id := r.URL.Query().Get("id")

	// Получаем задачи из репозитория
	task, err := repository.GetTaskById(id)
	if err != nil {
		writeErrorJSON(w, "ошибка получения задачи из базы данных", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func GetTasksListHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок для JSON-ответа
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Получаем параметры поиска из строки запроса
	search := r.URL.Query().Get("search")

	// Получаем задачи из репозитория
	tasks, err := repository.GetTasks(search)
	if err != nil {
		writeErrorJSON(w, "ошибка получения задач из базы данных", http.StatusInternalServerError)
		return
	}

	// Возвращаем список задач
	if tasks == nil {
		tasks = []models.Task{}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}

func DoneDeleteTask(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок для JSON-ответа
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Получаем параметры поиска из строки запроса
	id := r.URL.Query().Get("id")

	// Получаем задачи из репозитория
	task, err := repository.GetTaskById(id)
	if err != nil {
		writeErrorJSON(w, "ошибка получения задачи из базы данных", http.StatusInternalServerError)
		return
	}
	if r.Method == "DELETE" || task.Repeat == "" {
		err = repository.DeleteTask(task.Id)
		if err != nil {
			writeErrorJSON(w, "ошибка удаления задачи из базы данных", http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		newDate, err := scheduler.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeErrorJSON(w, "не удалось получить новую дату", http.StatusBadRequest)
			return
		}
		err = repository.UpdateDateToDB(newDate, task)
		if err != nil {
			writeErrorJSON(w, "не удалось обновить дату", http.StatusBadRequest)
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]string{})
}

func writeErrorJSON(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка encoding JSON error response: %v", err)
	}
}
