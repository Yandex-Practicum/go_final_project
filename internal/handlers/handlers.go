package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"todo-server/internal/donetaskrepeat"
	"todo-server/internal/settings"
	"todo-server/internal/store"
	"todo-server/internal/tasks"
)

type ResponseJson struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {

	strnow := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	strRepeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(settings.Template, strnow)
	if err != nil {
		log.Fatal(err)
	}
	nextdate, err := donetaskrepeat.NextDate(now, date, strRepeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(nextdate))
	if err != nil {
		log.Fatal(err)
	}
}

func HandlePostGetPutRequests(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t tasks.Task
		switch r.Method {
		case http.MethodPost:
			err := json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
				return
			}
			id, err := store.CreateTask(t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			resp := ResponseJson{ID: id}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		case http.MethodGet:
			id := r.URL.Query().Get("id")
			task, err := store.GetTask(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(task); err != nil {
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		case http.MethodPut:
			err := json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
				return
			}
			err = store.UpdateTask(t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			err := store.DeleteTask(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, `{"error":"Ошибка метода запроса"}`, http.StatusBadRequest)
			return
		}
	}
}

func HandleTasksGet(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		tasksList, err := store.SearchTask(search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response := map[string][]tasks.Task{
			"tasks": tasksList,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandleTaskDone(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		err := store.DoneTask(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
