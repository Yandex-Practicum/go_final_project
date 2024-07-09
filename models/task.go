package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VesuvyX/go_final_project/handlers"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type ResponseTaskId struct {
	ID string `json:"id"`
}

var db *sql.DB

//var dbGorm *gorm.DB

func SetDB(database *sql.DB) {
	db = database
}

func TaskAddPOST(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error reading request body: %v", err) // Лог
		return
	}

	// десериализуем JSON в task
	// 1. проверка десериализации
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error unmarshaling JSON: %v", err) // Лог
		return
	}

	// 2. проверка заголовка
	if len(task.Title) == 0 {
		http.Error(w, `{"error":"Заголовок пуст"}`, http.StatusBadRequest)
		log.Println("Error: Заголовок пуст") // Лог
		return
	}

	// 3. проверка формата даты (не 20060102)
	// 4. проверка правила повторения
	if len(task.Date) == 0 {
		task.Date = time.Now().Format("20060102")
	} else {
		if _, err := time.Parse("20060102", task.Date); err != nil {
			http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
			log.Printf("Error: Дата указана в неверном формате: %v", task.Date) // Лог
			return
		}

		if len(task.Repeat) > 0 {
			if !strings.HasPrefix(task.Repeat, "d ") && task.Repeat != "y" {
				http.Error(w, `{"error":"Неверное значение для repeat"}`, http.StatusBadRequest)
				log.Printf("Error: Неверное значение для repeat: %v", task.Repeat) // Лог
				return
			}

			now := time.Now()
			nextDate, err := handlers.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Printf("Error in NextDate: %v", err) // Лог
				return
			} else if task.Date < time.Now().Format("20060102") {
				task.Date = nextDate
			}
		}

		if task.Date < time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		}

	}

	lastId, err := AddTask(task)
	if err != nil {
		http.Error(w, `{"error":"ошибка добавления задачи в БД"}`, http.StatusInternalServerError)
		log.Printf("Ошибка добавления задачи в DB: %v", err) // Лог
		return
	}

	taskId, err := json.Marshal(ResponseTaskId{ID: lastId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshaling response: %v", err) // Лог
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(taskId)

	log.Printf("Task added successfully: %+v", task) // Лог
}

func AddTask(t Task) (string, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

func TasksShowGET(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	var err error

	if tasks, err = TasksShow(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Ошибка при выводе Tasks из БД: %v", err) // Лог
		return
	}

	if tasks == nil {
		tasks = []Task{}
	}

	tasksData, err := json.Marshal(Tasks{Tasks: tasks})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshaling tasksData: %v", err) // Лог
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(tasksData)

	log.Printf("Список задач: %v", tasks) // Вывод задач в лог
}

func TasksShow() ([]Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 20"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса к БД: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки строк: %v", err)
	}
	return tasks, nil
}

func ReadTaskByIdGET(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if len(id) == 0 {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		log.Println("Error: Не указан идентификатор задачи")
		return
	}

	// проверка на Макс ID
	var maxID int
	maxRow := db.QueryRow(`SELECT MAX(id) FROM scheduler`)
	maxRow.Scan(&maxID)
	if err := maxRow.Err(); err != nil {
		http.Error(w, `{"error":"Неверный формат Id"}`, http.StatusBadRequest)
		log.Println("Error: Неверный формат Id") // Лог
		return
	}
	newID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error":"не парсится ID"}`, http.StatusBadRequest)
		log.Println("Error: не парсится ID") // Лог
		return
	}
	if newID > maxID {
		http.Error(w, `{"error":"новый ID больше, чем строк в БД"}`, http.StatusBadRequest)
		log.Println("Error: новый ID больше, чем строк в БД") // Лог
		return
	}

	taskData, err := ReadTaskById(newID)
	log.Printf("Взята задача: %v", taskData) // Лог
	if err != nil {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusBadRequest)
		log.Printf("Задача не найдена: %v", err) // Лог
		return
	}

	responseData, err := json.Marshal(taskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshaling task data: %v", err) // Лог
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)

	log.Printf("Взята задача: %v", taskData) // Лог
}

func ReadTaskById(id int) (Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?"
	row := db.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("задача с id %v не найдена", id)
		}
		log.Printf("ошибка выполнения запроса к БД: %v", err) // Лог
		return Task{}, fmt.Errorf("ошибка выполнения запроса к БД: %v", err)
	}

	return task, nil
}

func TaskUpdatePUT(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Ошибка при чтении body: %v", err) // Лог
		return
	}

	// десериализуем JSON в task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error unmarshaling JSON: %v", err) // Лог
		return
	}

	// проверка ID
	if len(task.ID) == 0 {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		log.Println("Error: Не указан идентификатор задачи") // Лог
		return
	}

	// проверка ID
	if _, err := strconv.Atoi(task.ID); err != nil {
		http.Error(w, `{"error":"Неверный формат Id"}`, http.StatusBadRequest)
		log.Println("Error: Неверный формат Id") // Лог
		return
	}

	// проверка на Макс ID
	var maxID int
	maxRow := db.QueryRow(`SELECT MAX(id) FROM scheduler`)
	maxRow.Scan(&maxID)
	if err = maxRow.Err(); err != nil {
		http.Error(w, `{"error":"Неверный формат Id"}`, http.StatusBadRequest)
		log.Println("Error: Неверный формат Id") // Лог
		return
	}
	newID, err := strconv.Atoi(task.ID)
	if err != nil {
		http.Error(w, `{"error":"не парсится ID"}`, http.StatusBadRequest)
		log.Println("Error: не парсится ID") // Лог
		return
	}
	if newID > maxID {
		http.Error(w, `{"error":"новый ID больше, чем строк в БД"}`, http.StatusBadRequest)
		log.Println("Error: новый ID больше, чем строк в БД") // Лог
		return
	}

	// 2. проверка заголовка
	if len(task.Title) == 0 {
		http.Error(w, `{"error":"Заголовок пуст"}`, http.StatusBadRequest)
		log.Println("Error: Заголовок пуст") // Лог
		return
	}

	// 3. проверка формата даты (не 20060102)
	// 4. проверка правила повторения
	if len(task.Date) == 0 {
		task.Date = time.Now().Format("20060102")
	} else {
		if _, err := time.Parse("20060102", task.Date); err != nil {
			http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
			log.Printf("Error: Дата указана в неверном формате: %v", task.Date) // Лог
			return
		}

		if len(task.Repeat) > 0 {
			if !strings.HasPrefix(task.Repeat, "d ") && task.Repeat != "y" {
				http.Error(w, `{"error":"Неверное значение для repeat"}`, http.StatusBadRequest)
				log.Printf("Error: Неверное значение для repeat: %v", task.Repeat) // Лог
				return
			}

			now := time.Now()
			nextDate, err := handlers.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Printf("Ошибка в NextDate: %v", err) // Лог
				return
			} else if task.Date < time.Now().Format("20060102") {
				task.Date = nextDate
			}
		}

		if task.Date < time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		}

	}

	if err := UpdateTask(task); err != nil {
		http.Error(w, `{"error":"Ошибка обновления задачи в БД"}`, http.StatusInternalServerError)
		log.Printf("Ошибка обновления задачи в БД: %v", err) // Лог
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))

	log.Printf("Задание успешно обновлено: %v", task) // Лог
}

func UpdateTask(task Task) error {
	query := "UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?"
	_, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса к БД: %v", err)

	}
	return nil
}

func TaskDonePOST(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	if len(taskID) == 0 {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// проверка ID
	if _, err := strconv.Atoi(taskID); err != nil {
		http.Error(w, `{"error":"Неверный формат Id"}`, http.StatusBadRequest)
		log.Printf("Error: Неверный формат Id: %v", taskID) // Лог
		return
	}

	// проверка на Макс ID
	var maxID int
	maxRow := db.QueryRow(`SELECT MAX(id) FROM scheduler`)
	maxRow.Scan(&maxID)
	if err := maxRow.Err(); err != nil {
		http.Error(w, `{"error":"Неверный формат Id"}`, http.StatusBadRequest)
		log.Println("Error: Неверный формат Id") // Лог
		return
	}
	newID, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, `{"error":"не парсится ID"}`, http.StatusBadRequest)
		log.Println("Error: не парсится ID") // Лог
		return
	}
	if newID > maxID {
		http.Error(w, `{"error":"новый ID больше, чем строк в БД"}`, http.StatusBadRequest)
		log.Println("Error: новый ID больше, чем строк в БД") // Лог
		return
	}

	task, err := ReadTaskById(newID)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err := DeleteTask(newID)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := handlers.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при расчете следующей даты"}`, http.StatusBadRequest)
			return
		}

		task.Date = nextDate
		err = UpdateTask(task)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func TaskDELETE(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	if len(taskID) == 0 {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// проверка ID
	if _, err := strconv.Atoi(taskID); err != nil {
		http.Error(w, `{"error":"Неверный формат Id"}`, http.StatusBadRequest)
		log.Printf("Error: Неверный формат Id: %v", taskID) // Лог
		return
	}

	// проверка на Макс ID
	var maxID int
	maxRow := db.QueryRow(`SELECT MAX(id) FROM scheduler`)
	maxRow.Scan(&maxID)
	if err := maxRow.Err(); err != nil {
		http.Error(w, `{"error":"Неверный формат Id"}`, http.StatusBadRequest)
		log.Println("Error: Неверный формат Id") // Лог
		return
	}
	newID, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, `{"error":"не парсится ID"}`, http.StatusBadRequest)
		log.Println("Error: не парсится ID") // Лог
		return
	}
	if newID > maxID {
		http.Error(w, `{"error":"новый ID больше, чем строк в БД"}`, http.StatusBadRequest)
		log.Println("Error: новый ID больше, чем строк в БД") // Лог
		return
	}

	err = DeleteTask(newID)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
		log.Println("Error: Ошибка при удалении задачи") // Лог
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(`{}`))
}

func DeleteTask(taskID int) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	log.Printf("Удалена задача с ID: %v", taskID) // Лог
	_, err := db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса удаления к БД: %v", err)
	}

	return err
}
