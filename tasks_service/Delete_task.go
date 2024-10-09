package tasks_service

import (
	"database/sql"
	"fmt"
	"net/http"
)

// УДАЛЕНИЕ ЗАДАЧ
func deleteTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан id"}`, http.StatusBadRequest)
		return
	}

	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM scheduler WHERE id = ?)`, id).Scan(&exists)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка чтения БД: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, `{"error":"Задачи с таким ID не существует"}`, http.StatusNotFound)
		return
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	_, err = db.Exec(query, id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка удаления: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{}`)
}
