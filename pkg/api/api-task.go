package api

import "net/http"

// taskHandler распределяет запросы эндпойнта "/api/task" по типу
func taskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	}

}
