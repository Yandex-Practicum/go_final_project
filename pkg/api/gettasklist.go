package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/JKasus/go_final_project/pkg/constants"
	"github.com/JKasus/go_final_project/pkg/db"
	"github.com/JKasus/go_final_project/pkg/entities"
)

func getTaskListHandler(w http.ResponseWriter, r *http.Request) {

	var filter entities.Filter
	var searchValue string
	filter.Limit = constants.DefaultLimit

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		offset, err := strconv.ParseInt(offsetParam, 10, 64)
		if err != nil {
			err = errors.New("failed to parse 'offset' parameter")
			writeJSON(w, http.StatusBadRequest, err)
			return
		}
		filter.Offset = offset
	}

	if searchParam := r.URL.Query().Get("search"); searchParam != "" {
		searchValue = searchParam
	}

	taskList, err := db.GetTaskList(&filter, searchValue)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	if taskList == nil {
		writeJSON(w, http.StatusOK, []entities.Task{})
	} else {
		writeJSON(w, http.StatusOK, taskList)
	}
}
