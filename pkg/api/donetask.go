package api

import (
	"fmt"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

// DoneTaskHandler обработчик для отметки задачи выполненной
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Поддерживаем и GET и POST для совместимости с тестами
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	fmt.Printf("Отмечаем задачу выполненной ID: %s (Method: %s)\n", id, r.Method)

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Если задача без повторения - удаляем её
	if task.Repeat == "" {
		fmt.Printf("Задача без повторения - удаляем ID: %s\n", id)
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
		fmt.Printf("Задача ID %s удалена\n", id)
	} else {
		// Если задача с повторением - вычисляем следующую дату
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": fmt.Sprintf("Ошибка вычисления следующей даты: %v", err)}, http.StatusBadRequest)
			return
		}

		fmt.Printf("Задача с повторением - обновляем дату с %s на %s\n", task.Date, nextDate)

		// Обновляем дату задачи
		if err := db.UpdateDate(id, nextDate); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		fmt.Printf("Дата задачи ID %s обновлена на %s\n", id, nextDate)
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

// DeleteTaskHandler обработчик для удаления задачи
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	fmt.Printf("Удаляем задачу ID: %s\n", id)

	// Удаляем задачу из БД
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	fmt.Printf("Задача ID %s удалена\n", id)

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
