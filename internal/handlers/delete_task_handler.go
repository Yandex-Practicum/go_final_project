package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	utils "github.com/falsefood/go_final_project/internal"
)

func DeleteTaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("--- Функция deleteTaskHandler вызвана ---")

	taskID := r.URL.Query().Get("id")
	log.Printf("Полученный ID задачи: %s", taskID)

	if taskID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Идентификатор задачи не указан")
		return
	}

	id, err := strconv.Atoi(taskID)
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "ID задачи должен быть положительным числом")
		return
	}

	log.Printf("Попытка удаления задачи с ID: %d", id)

	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка при удалении задачи")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка при проверке удаления задачи")
		return
	}
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	log.Printf("Задача с ID %d успешно удалена", id)
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
