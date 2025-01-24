package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	utils "github.com/falsefood/go_final_project/internal"
)

// Получение задачи по ID
func GetTaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("--- Функция getTaskHandler вызвана ---")

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		log.Println("Ошибка: Не указан идентификатор задачи")
		utils.WriteError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	id, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		log.Printf("Ошибка преобразования ID: %v\n", err)
		utils.WriteError(w, http.StatusBadRequest, "Неверный формат идентификатора")
		return
	}

	var task utils.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Ошибка: Задача не найдена")
			utils.WriteError(w, http.StatusNotFound, "Задача не найдена")
		} else {
			log.Printf("Ошибка выполнения SQL-запроса: %v\n", err)
			utils.WriteError(w, http.StatusInternalServerError, "Ошибка при получении задачи")
		}
		return
	}

	log.Printf("Получена задача: %+v\n", task)

	utils.WriteJSON(w, http.StatusOK, task)
}

// Получение всех задач
func GetTasksHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	const limit = 10

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := db.Query(query, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка чтения задач")
		return
	}
	defer rows.Close()

	tasks := make([]utils.Task, 0)
	for rows.Next() {
		var task utils.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Ошибка чтения задач")
			return
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка чтения задач")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string][]utils.Task{"tasks": tasks})
}
