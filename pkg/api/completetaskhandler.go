package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/JKasus/go_final_project/pkg/db"
	"github.com/JKasus/go_final_project/pkg/entities"
	"github.com/JKasus/go_final_project/pkg/internal"
)

func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task *entities.Task
	var err error
	var id string
	now := time.Now()

	if idParam := r.URL.Query().Get("id"); idParam != "" {
		id = idParam
		task, err = db.GetTaskById(id)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, err)
			return
		}
	} else {
		err = errors.New("id param is required")
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	if task.Repeat == "" || &task.Repeat == nil {
		err = db.DeleteTask(id)
		writeJSON(w, http.StatusOK, entities.EmptyResponse{})
		if err != nil {
			writeJSON(w, http.StatusBadRequest, err)
			return
		}
	} else {
		task.Date, err = internal.NextDate(now, task.Date, task.Repeat)
		err = db.UpdateTask(task)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, entities.EmptyResponse{})
	}
}
