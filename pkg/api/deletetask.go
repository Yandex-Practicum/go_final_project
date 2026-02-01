package api

import (
	"errors"
	"net/http"

	"github.com/JKasus/go_final_project/pkg/db"
	"github.com/JKasus/go_final_project/pkg/entities"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if idParam := r.URL.Query().Get("id"); idParam != "" {
		err = db.DeleteTask(idParam)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, err)
			return
		}
	} else {
		err = errors.New("id param is required")
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, entities.EmptyResponse{})
}
