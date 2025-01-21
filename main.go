package main

import (
	//"fmt"
	//"os"
	//"path/filepath"

	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

var (
	port   = "7540"
	webDir = "./web"
)

const dateConst = "20060102"

type Task struct {
	Date    string `json:"date"`    // дата задачи в формате
	Title   string `json:"title"`   // заголовок задачи
	Comment string `json:"comment"` // комментарий к задаче
	Repeat  string `json:"repeat"`  // правило повторения
}

type TaskWithID struct {
	Id      string `json:"id"`      // уникальный номер задачи
	Date    string `json:"date"`    // дата задачи в формате
	Title   string `json:"title"`   // заголовок задачи
	Comment string `json:"comment"` // комментарий к задаче
	Repeat  string `json:"repeat"`  // правило повторения

}

// расчет следующей даты задачи
func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("не задано правило повтора")
	}

	// парсим время
	dateParse, err := time.Parse(dateConst, date)
	if err != nil {
		return "", err
	}

	// ежедневное выполнение задачи
	if repeat == "d 1" && !dateParse.After(now) {
		return now.Format(dateConst), nil
	}

	for {
		// добавление одного года
		if repeat == "y" {
			dateParse = dateParse.AddDate(1, 0, 0)

			// вычисление количества дней и их добавление
		} else if strings.HasPrefix(repeat, "d ") {
			daysString := strings.Split(repeat, " ")
			daysInt, err := strconv.Atoi(daysString[1])
			if daysInt < 1 || daysInt > 400 || err != nil {
				return "", err
			}
			dateParse = dateParse.AddDate(0, 0, daysInt)
		} else {
			return "", errors.New("incorrect repeat rule")
		}

		//проверка новой даты
		if dateParse.After(now) {
			return dateParse.Format(dateConst), nil
		}
	}
}

// проектирование и создание базы данных
func initDB() (*sql.DB, error) {

	// определяем путь к базе данных
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	// проверяем, существует ли файл базы данных
	_, err = os.Stat(dbFile)
	install := os.IsNotExist(err)

	// открываем базу данных
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// создаем таблицу и индекс, если файл базы данных отсутствует
	if install {
		createTableQuery := `
            CREATE TABLE IF NOT EXISTS scheduler (
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
				date INTEGER NOT NULL DEFAULT 20060102,
    			title VARCHAR(64) NOT NULL DEFAULT "",
    			comment VARCHAR(256) NOT NULL DEFAULT "",
   	 			repeat VARCHAR(128) NOT NULL DEFAULT ""
			);
			CREATE INDEX IF NOT EXISTS date_id ON scheduler(date);
		`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			return nil, err
		}
		log.Println("Таблица scheduler создана успешно.")
	}

	return db, nil
}

// GET даты из запроса и расчет новой
func nextDateHandler(w http.ResponseWriter, r *http.Request) {

	// получение параметров из строки запроса
	nowParam := r.URL.Query().Get("now")
	repeatParam := r.URL.Query().Get("repeat")
	dateParam := r.URL.Query().Get("date")

	// проверка и парсинг параметра now
	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		http.Error(w, "invalid 'now' date format", http.StatusBadRequest)
		return
	}

	// вызов функции NextDate для получения следующей даты задачи
	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ответ с результатом
	w.Write([]byte(nextDate))
}

// добавление новой задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	now := time.Now()
	nowFormat := now.Format(dateConst)

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}
	// проверка поля Title
	if task.Title == "" {
		http.Error(w, `{"error":"Поле 'title' пустое"}`, http.StatusBadRequest)
		return
	}

	// проверка пустого поля Date
	if task.Date == "" {
		task.Date = nowFormat
	} else {
		// проверка парсинга даты в формате 20060102
		dateParse, err := time.Parse(dateConst, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}

		dateParse = dateParse.Truncate(24 * time.Hour)
		now = now.Truncate(24 * time.Hour)

		// если дата меньше сегодняшнего числа
		if dateParse.Before(now) {
			// если правило Repeat "пустое" ставим сегодняшнюю дату
			if task.Repeat == "" {
				task.Date = nowFormat
			} else {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Неверный формат правила"}`, http.StatusInternalServerError)
					return
				}
				// проверяем правило повторения
				if task.Repeat != "" {
					_, err := NextDate(now, task.Date, task.Repeat)
					if err != nil {
						http.Error(w, `{"error":"Неверный формат правила"}`, http.StatusInternalServerError)
						return
					}
				}
				// если задача ежедневная
				if task.Repeat == "d 1" && nextDate == nowFormat {
					task.Date = nowFormat
				} else {
					task.Date = nextDate
				}
			}
		}
	}

	// добавление задачи в БД
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при добавлении задачи в базу данных"}`, http.StatusInternalServerError)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error":"Ошибка получения id задачи"}`, http.StatusInternalServerError)
	}

	// формируем ответ
	response := map[string]string{"id": fmt.Sprintf("%d", id)}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error":"Ошибка конвертирования данных в JSON"}`, http.StatusInternalServerError)
		return
	}

}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	var count int
	// подключение к БД
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//подсчитываем количество записей
	err = db.QueryRow("SELECT COUNT(*) FROM scheduler").Scan(&count)
	if err != nil {
		http.Error(w, `{"error":"Ошибка подсчета количества задач"}`, http.StatusInternalServerError)
		return
	}

	//инициализация пустого массива
	tasks := []TaskWithID{}

	if count > 0 {
		// запрашиваем записи в порядке возрастания по дате
		rows, err := db.Query(`SELECT id, date, title, comment, repeat
			FROM scheduler 
			ORDER BY date ASC LIMIT 50
		`)
		if err != nil {
			http.Error(w, `{"error":"Ошибка подсчета количества задач"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var task TaskWithID
			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Ошибка считывания данных из дб"}`, http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}
	}
	// формируем ответ
	response := map[string]interface{}{
		"tasks": tasks,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error":"Ошибка сериализации данных в JSON"}`, http.StatusInternalServerError)
		return
	}
}

// получение задачи по Id
func getTaskIdHandler(w http.ResponseWriter, r *http.Request) {

	// подключение к БД
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// получение Id из запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Ошибка запроса Id"}`, http.StatusBadRequest)
		return
	}

	// получение задачи по Id из БД
	var task TaskWithID
	err = db.QueryRow(`SELECT id, date, title, comment, repeat	
		FROM scheduler 
		WHERE id = ?	`, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusInternalServerError)
		} else {
			http.Error(w, `{"error":"Ошибка запроса к базе данных"}`, http.StatusInternalServerError)

		}
		return
	}

	// формируем ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка сериализации данных в JSON"}`, http.StatusInternalServerError)
		return
	}
}

// изменение задачи
func editTaskHandler(w http.ResponseWriter, r *http.Request) {

	// подключение к БД
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var task TaskWithID

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка десериализации данных JSON"}`, http.StatusBadRequest)
		return
	}
	// проверяем поля Id и Title
	if task.Id == "" {
		http.Error(w, `{"error":"Id не задан"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Task не задано"}`, http.StatusBadRequest)
		return
	}

	// обработка поля Date
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateConst)
	} else {
		parsedDate, err := time.Parse(dateConst, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(now) && task.Repeat != "" {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Неверный формат правила"}`, http.StatusInternalServerError)
				return
			}
			task.Date = nextDate
		}
	}

	// редактируем задачу в БД
	res, err := db.Exec(`UPDATE scheduler 
		SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?
		`, task.Date, task.Title, task.Comment, task.Repeat, task.Id)

	if err != nil {
		http.Error(w, `{"error":"Ошибка редактирования задачи"}`, http.StatusInternalServerError)
		return
	}

	// проверяем редактирование задачи
	affected, err := res.RowsAffected()
	if err != nil || affected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusInternalServerError)
		return
	}

	// возвращаем пустой JSON-ответ, если все хорошо
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// получение Id из запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Ошибка запроса Id"}`, http.StatusBadRequest)
		return
	}

	// подключение к БД
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// получение задачи по Id из БД
	var task TaskWithID
	err = db.QueryRow(`SELECT id, date, title, comment, repeat	
	FROM scheduler 
	WHERE id = ?	`, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusInternalServerError)
		} else {
			http.Error(w, `{"error":"Ошибка запроса к базе данных"}`, http.StatusInternalServerError)
		}
		return
	}

	// удаление одноразовой задачи
	if task.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			log.Printf("Error deleting task: %v", err)
			http.Error(w, `{"error":"Задача не удаляется"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// вычисление даты для периодической задачи
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Error calculating next date: %v", err)
			http.Error(w, `{"error":"Ошибка расчета даты для повторяющейся задачи"}`, http.StatusInternalServerError)
			return
		}

		// обновление даты выполнения задачи в базе данных
		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		if err != nil {
			log.Printf("Error updating task date: %v", err)
			http.Error(w, `{"error":"Невозможно обновить дату для повторяющейся задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	// возвращаем пустой JSON-ответ, если все хорошо
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

// удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	// получение Id из запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Ошибка запроса Id"}`, http.StatusBadRequest)
		return
	}

	// подключение к БД
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// удаление задачи из БД
	res, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, `{"error":"Задача не удаляется"}`, http.StatusInternalServerError)
		return
	}

	// проверка изменений в БД
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"Ошибка подсчета удаленных строк"}`, http.StatusInternalServerError)
		return
	}
	// если ни одна задача не удалена, возвращаем ошибку
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusInternalServerError)
		return
	}

	// возвращаем пустой JSON-ответ, если все хорошо
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

func main() {

	// инициализируем БД
	db, err := initDB()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", nextDateHandler)
	r.Get("/api/task", getTaskIdHandler)
	r.Post("/api/task", addTaskHandler)
	r.Put("/api/task", editTaskHandler)
	r.Get("/api/tasks", getTaskHandler)
	r.Post("/api/task/done", doneTaskHandler)
	r.Delete("/api/task", deleteTaskHandler)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}

}
