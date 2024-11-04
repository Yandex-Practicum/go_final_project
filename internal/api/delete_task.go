package api

import (
	"net/http"

	"go_final_project/internal/db"
)

// DeleteTaskHandler обрабатывает DELETE-запросы к /api/task для удаления задачи.
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	// Проверка, указан ли идентификатор
	if taskID == "" {
		http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// Подключение к базе данных
	dbConn, err := db.SetupDB()
	if err != nil {
		http.Error(w, `{"error": "Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Удаляем задачу из базы данных
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	result, err := dbConn.Exec(deleteQuery, taskID)
	if err != nil {
		http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
		return
	}

	// Проверяем, была ли удалена хотя бы одна строка
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON в случае успешного удаления
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}
