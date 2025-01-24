package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	utils "github.com/falsefood/go_final_project/internal"
)

func MarkTaskDoneHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "Идентификатор задачи не указан")
		return
	}

	taskID, err := strconv.Atoi(id)
	if err != nil || taskID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "ID задачи должен быть положительным числом")
		return
	}

	var task utils.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = db.QueryRow(query, taskID).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, "Задача не найдена")
		} else {
			utils.WriteError(w, http.StatusInternalServerError, "Ошибка при получении задачи")
		}
		return
	}

	if task.Repeat == "" {
		if _, err := db.Exec("DELETE FROM scheduler WHERE id = ?", taskID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Ошибка при удалении задачи")
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{})
		return
	}

	nextDateStr, err := nextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Ошибка при расчёте следующей даты: %v", err))
		return
	}

	if _, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDateStr, taskID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка при обновлении задачи")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
