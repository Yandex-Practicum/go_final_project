package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signInHandler)
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodPut:
		updateTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Метод не поддерживается",
		})
	}
}
