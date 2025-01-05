package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/FunnyFoXD/go_final_project/databases"
	"github.com/FunnyFoXD/go_final_project/helpers"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.Error(w, `{"error":"identifier is empty"}`, http.StatusBadRequest)
		return
	}

	task, err := databases.GetTaskByID(id)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err = databases.DeleteTask(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := helpers.NextDate(time.Now().Truncate(24*time.Hour), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		task.Date = nextDate

		err = databases.UpdateTaskDateByID(id, task.Date)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("{}"))
}
