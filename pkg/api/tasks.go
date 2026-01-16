package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"slices"
	"strconv"
)

// GET
type GetTaskResponse struct {
	Tasks []db.Task `json:"tasks"`
	Error string    `json:"error,omitempty"`
}

// Получение списка задач по запросу GET
func handleGetTasks(w http.ResponseWriter, req *http.Request) {
	// Получаем задачи из Бд
	tasks, err := db.GetTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GetTaskResponse{Error: "Ошибка получения задач из БД"})
		return
	}

	// Если БД пустая, возвращаем пустой слайс
	if tasks == nil {
		tasks = make([]db.Task, 0)
		json.NewEncoder(w).Encode(GetTaskResponse{Tasks: tasks})
		return
	}

	type TaskDateInt struct {
		ID      string
		Date    int
		Title   string
		Comment string
		Repeat  string
	}

	var modifiedTasks []TaskDateInt
	for _, task := range tasks {

		// Конвертация в int (Ближайшие даты все равно будут меньше последних)
		date, err := strconv.Atoi(task.Date)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(GetTaskResponse{Error: "Ошибка: неверный формат времени"})
			return
		}

		// Добавляем в слайс для сортировки элементы
		modifiedTasks = append(modifiedTasks, TaskDateInt{
			ID:      task.ID,
			Date:    date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	// Сортировка по дате
	slices.SortFunc(modifiedTasks, func(task1, task2 TaskDateInt) int {
		if task1.Date < task2.Date {
			return -1
		} else if task1.Date > task2.Date {
			return 1
		}
		return 0
	})

	// Перевод обратно в слайс для ответа
	tasks = nil
	for _, task := range modifiedTasks {
		tasks = append(tasks, db.Task{
			ID:      task.ID,
			Date:    strconv.Itoa(task.Date), // конвертируем обратно в строку
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	// Проверка на длину ответа
	if len(tasks) > 10 {
		cut := len(tasks) - 10
		tasks = tasks[:len(tasks)-cut]
	}

	json.NewEncoder(w).Encode(GetTaskResponse{Tasks: tasks})
}
