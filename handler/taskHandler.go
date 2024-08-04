package handler

import (
	"database/sql"
	"encoding/json"
	"go_final_project/storage"
	"net/http"
)

var ResponseStatus int

// TaskHandler возвращает обработчик для создания и обновления задач
func TaskHandler(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query().Get("id")

	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		http.Error(w, "error opening database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var response []byte

	switch req.Method {
	case http.MethodGet:
		if param == "" {
			http.Error(w, `{"error":"inncorect id"}`, http.StatusBadRequest)
			return
		}
		response, ResponseStatus, err = storage.TaskID(db, param)
		defer db.Close()
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}

	case http.MethodPost:
		response, ResponseStatus, err = storage.AddTask(db, req)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}

	case http.MethodPut:
		response, ResponseStatus, err = storage.UptadeTaskID(db, req)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
	case http.MethodDelete:
		ResponseStatus, err = storage.DeleteTask(db, param)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
		// если прошло всё успешно, то возвращаем пустой json
		str := map[string]interface{}{}
		response, err = json.Marshal(str)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
