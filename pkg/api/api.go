package api

import (
	"encoding/json"
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneHandler)
}

func taskHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodPost:
		// Добавление задачи
		handlePostTask(w, req)
	case http.MethodGet:
		if req.URL.Query().Get("id") != "" {
			// Достать одну задачу по id
			handleGetSingleTask(w, req)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GetByIdResponse{Error: "Не указан id задачи"})
		}
	case http.MethodPut:
		// Сохранение задачи после изменений
		handlePutTask(w, req)
	case http.MethodDelete:
		// Удаление задачи
		handleDelete(w, req)
	default:
		// Ошибка если метод не реализован
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "Метод не поддерживается"})
	}
}

// Отдельный хендлер для tasks чтобы запросы GET без id друг друга не перекрывали
func tasksHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		handleGetTasks(w, req)
	} else {
		// Ошибка если метод не реализован
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "Метод не поддерживается"})
	}
}

// Отдельный хендлер для done чтобы запросы POST друг друга не перекрывали
func doneHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		handleTaskDone(w, req)
	} else {
		// Ошибка если метод не реализован
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "Метод не поддерживается"})
	}
}
