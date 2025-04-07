package api

import (
	"encoding/json"
	"net/http"

	"github.com/Stanislav-id05/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // запросим максимум 50 задач
	if err != nil {
		writeJsonError(w, err) // функция, которая возвращает ошибку в JSON
		return
	}

	// Если нет задач, создаем пустой срез
	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func writeJsonError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
