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

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("не задано правило повтора")
	}

	// Парсим время
	dateParse, err := time.Parse(dateConst, date)
	if err != nil {
		return "", err
	}

	// Ежедневное выполнение задачи
	if repeat == "d 1" && !dateParse.After(now) {
		return now.Format(dateConst), nil
	}

	for {
		// Дополнение одного года
		if repeat == "y" {
			dateParse = dateParse.AddDate(1, 0, 0)

			// Вычисление количества дней и их добавление
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

		//Проверка новой даты
		if dateParse.After(now) {
			return dateParse.Format(dateConst), nil
		}
	}
}

func initDB() (*sql.DB, error) {
	// Определяем путь к базе данных
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	// Проверяем, существует ли файл базы данных
	_, err = os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем базу данных
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу и индекс, если файл базы данных отсутствует
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

func nextDateHandler(w http.ResponseWriter, r *http.Request) {

	// Получение параметров из строки запроса
	nowParam := r.URL.Query().Get("now")
	repeatParam := r.URL.Query().Get("repeat")
	dateParam := r.URL.Query().Get("date")

	// Проверка и парсинг параметра now
	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		http.Error(w, "invalid 'now' date format", http.StatusBadRequest)
		return
	}

	// Вызов функции NextDate для получения следующей даты задачи
	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ответ с результатом
	w.Write([]byte(nextDate))
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task

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

	timeNow := time.Now()

	// проверка пустого поля Date
	if task.Date == "" {
		task.Date = timeNow.Format(dateConst)
	} else {
		// проверка парсинга даты в формате 20060102
		dateParse, err := time.Parse(dateConst, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}

		// если дата меньше сегодняшнего числа
		if dateParse.Before(timeNow) {
			// если правило Repeat "пустое" ставим сегодняшнюю дату
			if task.Repeat == "" {
				task.Date = timeNow.Format(dateConst)
			} else {
				nextDate, err := NextDate(timeNow, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Неверный формат правила"}`, http.StatusInternalServerError)
					return
				}
				// проверяем правило повторения
				if task.Repeat != "" {
					_, err := NextDate(timeNow, task.Date, task.Repeat)
					if err != nil {
						http.Error(w, `{"error":"Неверный формат правила"}`, http.StatusInternalServerError)
						return
					}
				}
				// если задача ежедневная
				if task.Repeat == "d 1" && nextDate == timeNow.Format(dateConst) {
					task.Date = timeNow.Format(dateConst)
				} else {
					task.Date = nextDate
				}
			}
		}
	}

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при добавлении задачи в базу данных"}`, http.StatusInternalServerError)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error":"Ошибка получения id задачи"}`, http.StatusInternalServerError)
	}

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
	// подключение к бд
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

func getTaskIdHandler(w http.ResponseWriter, r *http.Request) {
	// подключение к бд
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// получение Id из запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Ошибка запроса Id к бд"}`, http.StatusBadRequest)
		return
	}

	// получение задачи по Id из бд
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

func editTaskHandler(w http.ResponseWriter, r *http.Request) {

	// подключение к бд
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

func main() {

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

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}

}
