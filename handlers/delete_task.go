package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/FunnyFoXD/go_final_project/databases"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.Error(w, `{"error":"identifier is empty"}`, http.StatusBadRequest)
		return
	}

	if !isValidID(id) {
		http.Error(w, `{"error":"invalid identifier"}`, http.StatusBadRequest)
		return
	}

	err := databases.DeleteTask(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("{}"))
}

func isValidID(id string) bool {
	matched, _ := regexp.MatchString(`^\d+$`, id)
	return matched
}
