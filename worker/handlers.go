package worker

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"pwd/pkg/db"
	"strconv"
	"time"

	"pwd/internal/controller"
	"pwd/internal/nextdate"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

}

type TaskController struct {
	db *sql.DB
}

func NewTaskController(db *sql.DB) *TaskController {
	return &TaskController{db: db}
}

func (c *TaskController) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasks(c.db)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type response struct {
		Tasks []controller.Task `json:"tasks"`
	}

	resp := response{Tasks: tasks}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (c *TaskController) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task controller.Task // экзмпляр структуры со значениями
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		controller.ResponseError(w, "некорректный запрос")
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		controller.ResponseError(w, "ошибка десериализации json")
	}

	if task.Title == "" {
		controller.ResponseError(w, "поле title обязательно")
		return
	}

	// проверяем дату
	var date time.Time

	if task.Date == "" {
		task.Date = time.Now().Format("20060102") // указзываем сегодняшнюю дату
	}
	date, err = time.Parse("20060102", task.Date)
	if err != nil {
		controller.ResponseError(w, "некорректный формат времени")
		return
	}

	now := time.Now()
	// если дата меньше сегодняшней

	if date.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102") // устанавливаем сегодняшнюю дату
		} else {
			// вычисляем следующую дату
			nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				controller.ResponseError(w, "ошибка вычисления следующей даты")
				return
			}
			task.Date = nextDate
		}
	}

	id, err := db.AddTask(c.db, task)
	if err != nil {
		controller.ResponseError(w, "ошибка при добавлении задачи")
		log.Println(err)
		return
	}
	taskId := strconv.Itoa(id)

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": taskId})

}

func (c *TaskController) NextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "недостаточно значений", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "недопустимый формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
