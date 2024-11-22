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
		task.Date = time.Now().Format("20060102") // указзываем сегодняшнюю дату при пустом поле
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

func (c *TaskController) GetTaskId(w http.ResponseWriter, r *http.Request) {
	// Получаем id задачи из URL запроса
	taskId := r.URL.Query().Get("id")
	if taskId == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(map[string]string{"error": "такой задачи нет"})
		w.Write(response)
		return
	}

	var task controller.Task
	err := c.db.QueryRow("SELECT id, date, title, comment, repeat FROM tasks WHERE id = ?", taskId).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			response, _ := json.Marshal(map[string]string{"error": "Задача не найдена"})
			w.Write(response)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)

	/*w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}*/
}

func (c *TaskController) UpdateTaskId(w http.ResponseWriter, r *http.Request) {
	var task controller.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(map[string]string{"error": "Неверный формат данных"})
		w.Write(jsonResponse)
		return
	}

	// проверяем id
	if task.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(map[string]string{"error": "id задачи не указан"})
		w.Write(jsonResponse)
		return
	}

	// проверяем title
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(map[string]string{"error": "заголовок задачи пуст"})
		w.Write(jsonResponse)
		return
	}

	// проверяем дату
	var date time.Time

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		date, err = time.Parse("20060102", task.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(map[string]string{"error": "неверный формат даты"})
			w.Write(jsonResponse)
			return
		}
		if date.Before(now) {
			if date.Before(now) {
				if task.Repeat == "" {
					task.Date = now.Format("20060102") // устанавливаем сегодняшнюю дату
				}
			}
		}
	}

	// Проверяем правило
	if task.Repeat != "" {
		_, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(map[string]string{"error": "некорректный формат правила"})
			w.Write(jsonResponse)
			return
		}
	}

	// Выполняем запрос обновления записи к db
	res, err := c.db.Exec("UPDATE tasks SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(map[string]string{"error": "Ошибка обновления"})
		w.Write(jsonResponse)
		return
	}

	changes, err := res.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse, _ := json.Marshal(map[string]string{"error": "ошибка проверки обновления"})
		w.Write(jsonResponse)
		return
	}

	if changes == 0 {
		w.WriteHeader(http.StatusNotFound)
		jsonResponse, _ := json.Marshal(map[string]string{"error": "такой задачи нет"})
		w.Write(jsonResponse)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	jsonResp, _ := json.Marshal(map[string]string{"message": "задача обновлена"})
	w.Write(jsonResp)
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
