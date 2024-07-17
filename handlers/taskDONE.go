package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func TaskDonePOST(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	if len(taskID) == 0 {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(taskID); err != nil {
		http.Error(w, `{"error":"Неверный формат Id"}`, http.StatusBadRequest)
		log.Printf("Error: Неверный формат Id: %v", taskID)
		return
	}

	newID, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, `{"error":"не парсится ID"}`, http.StatusBadRequest)
		log.Println("Error: не парсится ID")
		return
	}

	task, err := dbHelper.GetTask(newID)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err := dbHelper.DeleteTask(newID)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при расчете следующей даты"}`, http.StatusBadRequest)
			return
		}

		task.Date = nextDate
		err = dbHelper.UpdateTask(task)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	resp := []byte(`{}`)
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Ошибка при ответе: %v", err)
	}
}
