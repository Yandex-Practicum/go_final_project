package api

import (
	"fmt"
	"go_final_project/pkg/db"
	"net/http"
)

// TaskResponse структура для ответа со списком задач
type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
	Error string     `json:"error,omitempty"`
}

// TasksHendler обработчик для получения списка задач
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, TasksResponse{Error: "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("получен запрос на список задач")

	search := r.URL.Query().Get("search")
	var tasks []*db.Task
	var err error

	if search != "" {
		// Поиск по тексту или дате
		tasks, err = db.SearchTasks(search, 50)
	} else {
		// Получение всех задач
		tasks, err = db.GetTasks(50)
	}

	if err != nil {
		fmt.Printf("Ошибка получения задач: %v\n", err)
		writeJSON(w, TasksResponse{Error: "Ошибка получения задач: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	fmt.Printf("Успешно получено задач: %d\n", len(tasks))

	// Возвращаем успешный ответ
	writeJSON(w, TasksResponse{Tasks: tasks}, http.StatusOK)
}
