package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

// POST
type AddTaskResponse struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// GET by ID
type GetByIdResponse struct {
	ID      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PUT
type UpdateTaskResponce struct {
	Error string `json:"error,omitempty"`
}

// Добавление задачи по запросу POST
func handlePostTask(w http.ResponseWriter, req *http.Request) {
	// Парсинг json
	var r db.Task
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "Ошибка десериализации JSON"})
		return
	}

	// Проверка на пустой title
	if r.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "Нет заголовка задачи"})
		return
	}

	now := time.Now()
	today := now.Format(format)
	if r.Date == "" {
		// Если дата не указана, то ставим сегодня
		r.Date = today
	} else {
		// Если дата указана
		date, err := time.Parse(format, r.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(AddTaskResponse{Error: "Дата не в формате 20060102"})
			return
		}
		if afterNow(now, date) {
			// Если правила повторения нет, то ставим сегодня
			if r.Repeat == "" {
				r.Date = today
			} else {
				// Если правило повторения есть, то сразу ищем следующую дату
				nextDate, err := NextDate(now, r.Date, r.Repeat)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(AddTaskResponse{Error: err.Error()})
					return
				}
				r.Date = nextDate
			}
		}
	}

	// Добавляем структуру в БД
	id, err := db.AddTask(&r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "Ошибка добавления задачи в БД"})
		return
	}

	// Ответ
	json.NewEncoder(w).Encode(AddTaskResponse{ID: id})
}

// Сохранение изменений
func handlePutTask(w http.ResponseWriter, req *http.Request) {
	// Парсинг json
	var r db.Task
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Ошибка десериализации JSON"})
		return
	}

	// Проверка на пустой title
	if r.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Нет заголовка задачи"})
		return
	}

	now := time.Now()
	today := now.Format(format)
	if r.Date == "" {
		// Если дата не указана, то ставим сегодня
		r.Date = today
	} else {
		// Если дата указана
		date, err := time.Parse(format, r.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Дата не в формате 20060102"})
			return
		}
		if afterNow(now, date) {
			// Если правила повторения нет, то ставим сегодня
			if r.Repeat == "" {
				r.Date = today
			} else {
				// Если правило повторения есть, то сразу ищем следующую дату
				nextDate, err := NextDate(now, r.Date, r.Repeat)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(UpdateTaskResponce{Error: err.Error()})
					return
				}
				r.Date = nextDate
			}
		}
	}

	// Обновляем структуру в БД
	err = db.UpdateTask(&r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Ошибка добавления задачи в БД"})
		return
	}

	// Ответ
	json.NewEncoder(w).Encode(UpdateTaskResponce{})
}

// Возвращает все параментры задачи по id
func handleGetSingleTask(w http.ResponseWriter, req *http.Request) {
	// Получаем id из URL
	id := req.URL.Query().Get("id")

	// Получаем задачу по id
	task, err := db.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GetByIdResponse{Error: "Ошибка получения задачи из БД"})
		return
	}
	if task == nil {
		//w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetByIdResponse{Error: "Задача не найдена"})
		return
	}

	json.NewEncoder(w).Encode(GetByIdResponse{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	})
}

// Удаление задачи
func handleDelete(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем id из запроса
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "id задачи не указан"})
		return
	}

	// Проверка на существование задачи
	task, err := db.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Ошибка проверки задачи"})
		return
	}
	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Задача не найдена"})
		return
	}

	// Непосредственно удаление
	err = db.DeleteTask(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(UpdateTaskResponce{Error: "Не удалось удалить задачу"})
		return
	}
	json.NewEncoder(w).Encode(UpdateTaskResponce{})
}
