package api

import (
	"net/http"
	"strconv"
	"time"

	"go_final_project/db"
)

// Обработчик для отметки о выполнении задачи
func MarkTaskDone(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Неверный идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// Получаем задачу из базы данных
	task, err := db.GetTaskByID(strconv.Itoa(id))
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	if task.Repeat != "" {
		// Рассчитываем следующую дату
		currentDate, err := time.Parse("20060102", task.Date) // Преобразуем строку даты в time.Time
		if err != nil {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}

		nextDate, err := NextDate(currentDate, task.Date, task.Repeat) // Передаем дату задачи и правило повторения
		if err != nil {
			http.Error(w, `{"error":"Ошибка при расчете следующей даты"}`, http.StatusInternalServerError)
			return
		}

		task.Date = nextDate // Обновляем дату задачи в нужном формате
	} else {
		// Удаляем задачу, если она одноразовая
		err = db.DeleteTaskByID(id) // Удаляем задачу по ID
		if err != nil {
			http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}
		w.Write([]byte("{}")) // Возвращаем пустой JSON
		return
	}

	// Обновляем задачу в базе данных
	err = db.UpdateTask(task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{}")) // Возвращаем пустой JSON
}

// Обработчик для удаления задачи
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Неверный идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// Удаляем задачу из базы данных
	err = db.DeleteTaskByID(id)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{}")) // Возвращаем пустой JSON
}
