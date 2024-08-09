package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

type TaskPutDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (h *Handler) handlePutTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskPutDTO
	err := json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		utils.RespondWithError(w, "Ошибка десериализации JSON")
		return
	}
	taskId, err := strconv.ParseInt(taskDTO.ID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, "Неверный ID задачи")
		return
	}

	taskRequest := models.Task{
		ID:      taskId,
		Date:    taskDTO.Date,
		Title:   taskDTO.Title,
		Comment: taskDTO.Comment,
		Repeat:  taskDTO.Repeat,
	}

	task, err := validateTask(&taskRequest)
	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	query := `UPDATE scheduler
    		  SET
    			date = ?,
    			title = ?,
    			comment = ?,
    			repeat = ?
			  WHERE id = ?`
	updateResult, err := h.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		utils.RespondWithError(w, "Задача не найдена")
		return
	}

	rowsAffected, err := updateResult.RowsAffected()
	if err != nil || rowsAffected == 0 {
		utils.RespondWithError(w, "Задача не найдена")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}
