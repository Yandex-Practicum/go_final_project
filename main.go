package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var DBFile = "./scheduler.db"

// daysInMonth - функция, которая принимает год и месяц в качестве аргументов и возвращает количество дней в этом месяце.
func daysInMonth(year, month int) int {
	// Используем switch-case, чтобы определить количество дней в каждом месяце.
	switch month {
	// Месяцы с 31 дня.
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	// Месяцы с 30 дня.
	case 4, 6, 9, 11:
		return 30
	// Февраль.
	case 2:
		// Если год является высокосным, то в феврале 29 дней.
		if isLeapYear(year) {
			return 29
		}
		// Иначе в феврале 28 дней.
		return 28
	// Неверный номер месяца.
	default:
		return 0
	}
}

// isLeapYear - функция, которая принимает год в качестве аргумента и возвращает true, если год является высокосным, и false, если нет.
func isLeapYear(year int) bool {
	// Если год делится на 400 без остатка, то это высокосный год.
	if year%400 == 0 {
		return true
	}
	// Если год делится на 100 без остатка, то это не высокосный год.
	if year%100 == 0 {
		return false
	}
	// Если год делится на 4 без остатка, то это высокосный год.
	if year%4 == 0 {
		return true
	}
	// Иначе это не высокосный год.
	return false
}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		err := errors.New("invalid repeat rule") //дефотлтная ошибка для вывода при некорректном вводе
		return "", err                           //пустой repeat
	}

	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err //неверный формат даты
	}

	commands := strings.Split(repeat, " ") //вычленяем команды для повторов
	ruleType := commands[0]                //берем тип команды

	newDate := parsedDate
	dateFound := false
	switch ruleType {

	case "d":
		if len(commands) < 2 {
			return "", err
		}
		repeatPeriod, err := strconv.Atoi(commands[1]) //находим дни для переноса задачи

		if err != nil {
			return "", err
		}
		if repeatPeriod > 400 {
			err := errors.New("invalid repeat rule") //дефотлтная ошибка для вывода при некорректном вводе
			return "", err                           //если агрумент для >400
		}
		for dateFound == false {
			if newDate.After(now) && newDate.After(parsedDate) { //если ближайшая дата перевалила за текущий день
				dateFound = true //ну вот тут вопрос как правильнее, так как при return мы все равно выходим из функции, стоит ли поднимать флаг и стоит ли делать break
				return newDate.Format("20060102"), nil
			}
			newDate = newDate.AddDate(0, 0, repeatPeriod)
		}
		return "", err

	case "y":
		for dateFound == false {
			if newDate.After(now) && newDate.After(parsedDate) { //если ближайшая дата перевалила за текущий день
				dateFound = true

				return newDate.Format("20060102"), nil
			}
			newDate = newDate.AddDate(1, 0, 0) //находим ближайшую новую дату
		}
		return "", err

	case "w":

		if len(commands) < 2 {
			err := errors.New("invalid repeat rule") //дефолтная ошибка для вывода при некорректном вводе
			return "", err
		}

		daysToRepeatStr := commands[1]                           //находим агрумент для повтора при команде повтора на неделю
		daysToRepeatArray := strings.Split(daysToRepeatStr, ",") //записываем в массив номера дней недели
		var daysToRepeatInts []int                               //записываем в массив номера дней недели

		for _, day := range daysToRepeatArray {

			dayInt, err := strconv.Atoi(day) //проверяем что день недели корректный
			if err != nil {
				return "", err
			}
			if dayInt > 7 { // все еще проверяем
				err := errors.New("invalid repeat rule") //дефолтная ошибка для вывода при некорректном вводе
				return "", err
			}
			daysToRepeatInts = append(daysToRepeatInts, dayInt)
		}

		sort.Ints(daysToRepeatInts)
		weekAdd := 0
		for dateFound == false {
			sunday := time.Date(newDate.Year(), time.Month(newDate.Month()), newDate.Day()-int(newDate.Weekday())+(weekAdd*7), 0, 0, 0, 0, time.UTC)
			for _, day := range daysToRepeatInts {

				// Вычисляем конкретный день недели

				newDate := sunday.AddDate(0, 0, day)

				if newDate.After(now) && newDate.After(parsedDate) { //если ближайшая дата перевалила за текущий день
					dateFound = true
					return newDate.Format("20060102"), nil
				}
			}
			weekAdd += 1
		}
		return "", err

	case "m":
		daysToRepeatStr := commands[1]                           //получаем аргументы для повторяемых дней
		daysToRepeatArray := strings.Split(daysToRepeatStr, ",") //переносим дни в массив
		var daysToRepeatInts []int
		for _, day := range daysToRepeatArray {
			dayInt, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}
			if dayInt == 0 || dayInt > 31 || dayInt < -2 { // все еще проверяем
				err := errors.New("invalid repeat rule")
				return "", err
			}
			daysToRepeatInts = append(daysToRepeatInts, dayInt)
		}

		everyMonthSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

		monthSlice := make([]int, 1, 1)

		if len(commands) > 2 { //если есть аргумент для месяца
			monthToRepeatStr := commands[2]                            //получаем агрументы для повторяемых месяцев
			monthToRepeatArray := strings.Split(monthToRepeatStr, ",") // переносим месяцы в массив
			for _, month := range monthToRepeatArray {
				monthInt, err := strconv.Atoi(month) //получаем номер месяца из строки
				if err != nil {
					return "", err
				}
				if monthInt > 12 || monthInt <= 0 {
					err := errors.New("invalid repeat rule") //дефолтная ошибка для вывода при некорректном вводе
					return "", err
				}
				monthSlice = append(monthSlice, monthInt)
			}
		} else {
			monthSlice = everyMonthSlice
		}
		sort.Ints(monthSlice)
		sort.Ints(daysToRepeatInts)

		year := newDate.Year()
		for dateFound == false {
			for _, month := range monthSlice {
				for _, day := range daysToRepeatInts {
					newDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

					if day < 0 && day > -3 {
						newDate = time.Date(year, time.Month(month), day+1, 0, 0, 0, 0, time.UTC)
					}
					if day == 0 || day < -2 {
						err := errors.New("invalid repeat rule") //дефолтная ошибка для вывода при некорректном вводе
						return "", err
					}

					if newDate.After(now) && newDate.After(parsedDate) { //если ближайшая дата перевалила за текущий день
						dateFound = true
						return newDate.Format("20060102"), err
					}
				}
			}
			year++
		}

		return "", err

	default:
		err := errors.New("invalid repeat rule") //дефотлтная ошибка для вывода при некорректном вводе
		return "", err
	}
}

func nextDateHandler(res http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		return
	}
	response, err := NextDate(nowTime, date, repeat)
	io.WriteString(res, response)
}

func CreateDB() {
	db, err := sql.Open("sqlite3", DBFile)

	if err != nil {
		return
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			"id" INTEGER PRIMARY KEY,
			"date" CHAR(8) NOT NULL,
			"title" VARCHAR(255) NOT NULL,
			"comment" TEXT,
			"repeat" VARCHAR(128) NOT NULL
			)
		`)
	if err != nil {
		return
	}
	_, err = db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_date ON sсheduler (date)
		`)
	if err != nil {
		return
	}
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

func addTask(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var buf bytes.Buffer
		task := &Task{}
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &task)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			responseError := struct {
				Error string `json:"error"`
			}{Error: err.Error()}
			json.NewEncoder(res).Encode(responseError)
			return
		}
		if task.Title == "" {
			responseError := struct {
				Error string `json:"error"`
			}{Error: "Не указан заголовок задачи"}
			json.NewEncoder(res).Encode(responseError)
			return
		}

		date := time.Now()

		if task.Date != "" {
			date, err = time.Parse("20060102", task.Date)
			if err != nil {
				responseError := struct {
					Error string `json:"error"`
				}{Error: err.Error()}
				json.NewEncoder(res).Encode(responseError)
				return
			}
		} else {
			task.Date = date.Format("20060102")
		}

		if date.Before(time.Now()) {
			if task.Repeat == "" {
				task.Date = time.Now().Format("20060102")
			} else {
				task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					responseError := struct {
						Error string `json:"error"`
					}{Error: err.Error()}
					json.NewEncoder(res).Encode(responseError)
					return
				}
			}
		}

		db, err := sql.Open("sqlite3", DBFile)
		defer db.Close()
		query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
		result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)

		if err != nil {
			return
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return
		}
		response := struct {
			ID int64 `json:"id"`
		}{ID: lastId}
		json.NewEncoder(res).Encode(response)

	}
}

func getTasks(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		db, err := sql.Open("sqlite3", DBFile)
		defer db.Close()

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			responseError := struct {
				Error string `json:"error"`
			}{Error: err.Error()}
			json.NewEncoder(res).Encode(responseError)
			return
		}

		rows, err := db.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", 10)
		if err != nil {
			responseError := struct {
				Error string `json:"error"`
			}{Error: err.Error()}
			json.NewEncoder(res).Encode(responseError)
			return
		}
		tasks := &Tasks{}
		for rows.Next() {
			task := &Task{}

			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				responseError := struct {
					Error string `json:"error"`
				}{Error: err.Error()}
				json.NewEncoder(res).Encode(responseError)
				return
			}
			tasks.Tasks = append(tasks.Tasks, *task)
		}
		if tasks.Tasks == nil {
			tasks.Tasks = []Task{}
		}
		data, err := json.Marshal(&tasks)
		if err != nil {
			responseError := struct {
				Error string `json:"error"`
			}{Error: err.Error()}
			json.NewEncoder(res).Encode(responseError)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write(data)
	}
}
func main() {
	debug := false
	if debug {
		line := []string{"20240125", "ooops", "20240129"}

		now, err := time.Parse("20060102", "20240126")
		if err != nil {
			return
		}
		date := line[0]
		repeat := line[1]
		fmt.Println(NextDate(now, date, repeat))
		return
	}
	var webDir = "./web"

	var port = ":7540"
	if len(os.Getenv("TODO_PORT")) > 0 {
		port = os.Getenv("TODO_PORT")
	}

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Getenv("TODO_DBFILE")) > 0 {
		DBFile = filepath.Join(filepath.Dir(appPath), "TODO_DBFILE")
	}

	_, err = os.Stat(DBFile)

	var install bool
	if err != nil {
		install = true
	}
	if install {
		CreateDB()
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/task", addTask)
	http.HandleFunc("/api/tasks", getTasks)
	http.HandleFunc("/api/nextdate", nextDateHandler)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")

}
