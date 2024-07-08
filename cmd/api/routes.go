package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Scheduler struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	Id    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
	r.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))

	r.Get("/api/nextdate", handleNextDate)
	r.Post("/api/task", handlePostTask)
	r.Get("/api/tasks", handleGetTask)
	r.Get("/api/task", handleGetTaskByID)
	r.Put("/api/task", handleUpdateTask)
	r.Post("/api/task/done", handleDoneTask)
	r.Delete("/api/task", handleDeleteTask)

	return r
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	var response Response
	var scheduler Scheduler
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&scheduler.ID, &scheduler.Date, &scheduler.Title, &scheduler.Comment, &scheduler.Repeat)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", scheduler.ID))
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&response)
	w.WriteHeader(http.StatusOK)
}

func handleDoneTask(w http.ResponseWriter, r *http.Request) {
	var response Response
	var scheduler Scheduler
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&scheduler.ID, &scheduler.Date, &scheduler.Title, &scheduler.Comment, &scheduler.Repeat)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dateNow := time.Now().Format(timeTemplate)
	dateNow_, err := time.Parse(timeTemplate, dateNow)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if scheduler.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", scheduler.ID))
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(w).Encode(&response)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		scheduler.Date, err = NextDate(dateNow_, scheduler.Date, scheduler.Repeat)
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(w).Encode(&response)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
			sql.Named("date", scheduler.Date),
			sql.Named("title", scheduler.Title),
			sql.Named("comment", scheduler.Comment),
			sql.Named("repeat", scheduler.Repeat),
			sql.Named("id", scheduler.ID))
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(w).Encode(&response)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	_, err = json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&response)
	w.WriteHeader(http.StatusOK)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var response Response
	var scheduler Scheduler
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &scheduler); err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if scheduler.ID == "" {
		response.Error = "Задача не найдена"
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if scheduler.Title == "" {
		response.Error = "Не указан заголовок задачи"
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dateNow := time.Now().Format(timeTemplate)

	if scheduler.Date == "" {
		scheduler.Date = dateNow
	}

	var date_ time.Time

	date_, err = time.Parse(timeTemplate, scheduler.Date)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dateNow_, err := time.Parse(timeTemplate, dateNow)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if date_.Compare(dateNow_) == -1 {
		if scheduler.Repeat == "" {
			scheduler.Date = dateNow_.Format(timeTemplate)
		} else {
			scheduler.Date, err = NextDate(dateNow_, scheduler.Date, scheduler.Repeat)
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(w).Encode(&response)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", scheduler.ID))
	err = row.Scan(&scheduler.ID)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", scheduler.Date),
		sql.Named("title", scheduler.Title),
		sql.Named("comment", scheduler.Comment),
		sql.Named("repeat", scheduler.Repeat),
		sql.Named("id", scheduler.ID))
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&response)
	w.WriteHeader(http.StatusOK)
}

func handleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var response Response
	var scheduler Scheduler
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&scheduler.ID, &scheduler.Date, &scheduler.Title, &scheduler.Comment, &scheduler.Repeat)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&scheduler)
	w.WriteHeader(http.StatusFound)
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	var response Response
	var scheduler Scheduler

	type Tasks struct {
		Tasks []interface{} `json:"tasks"`
	}
	var tasks Tasks

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	search := r.URL.Query().Get("search")
	limit := 50

	var rows *sql.Rows

	if search == "" {
		rows, err = db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", limit))
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(w).Encode(&response)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			search = `%` + search + `%`
			rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
				sql.Named("search", search),
				sql.Named("limit", limit))
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(w).Encode(&response)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			rows, err = db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit",
				sql.Named("date", date.Format(timeTemplate)),
				sql.Named("limit", limit))
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(w).Encode(&response)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&scheduler.ID, &scheduler.Date, &scheduler.Title, &scheduler.Comment, &scheduler.Repeat)
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(w).Encode(&response)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tasks.Tasks = append(tasks.Tasks, scheduler)
	}

	if tasks.Tasks == nil {
		tasks.Tasks = []interface{}{}
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func handlePostTask(w http.ResponseWriter, r *http.Request) {
	var scheduler Scheduler
	var buf bytes.Buffer
	var response Response

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &scheduler); err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if scheduler.Title == "" {
		response.Error = "Не указан заголовок задачи"
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dateNow := time.Now().Format(timeTemplate)

	if scheduler.Date == "" {
		scheduler.Date = dateNow
	}

	var date_ time.Time

	date_, err = time.Parse(timeTemplate, scheduler.Date)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dateNow_, err := time.Parse(timeTemplate, dateNow)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if date_.Compare(dateNow_) == -1 { // || date_.Compare(dateNow_) == 0 {
		if scheduler.Repeat == "" {
			scheduler.Date = dateNow_.Format(timeTemplate)
		} else {
			scheduler.Date, err = NextDate(dateNow_, scheduler.Date, scheduler.Repeat)
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(w).Encode(&response)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", scheduler.Date),
		sql.Named("title", scheduler.Title),
		sql.Named("comment", scheduler.Comment),
		sql.Named("repeat", scheduler.Repeat))
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Id = strconv.FormatInt(id, 10)
	_, err = json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(&response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&response)
	w.WriteHeader(http.StatusCreated)
}

func handleNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nowDate, err := time.Parse(timeTemplate, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	nextDate, err := NextDate(nowDate, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("не указано правило повторения")
	}

	date_, err := time.Parse(timeTemplate, date)
	if err != nil {
		return "", err
	}

	s := strings.Split(repeat, " ")

	daysLater := date_

	switch repeat[0] {
	case 'd':
		if len(s) != 2 {
			return "", errors.New("не указан интервал в днях")
		}

		//проверяем корректность введенных дней
		num, err := strconv.Atoi(s[1])
		if err != nil {
			return "", err
		}

		if !(num > 0 && num <= 400) {
			return "", errors.New("неверный диапазон дней")
		}

		//вычисляем новую дату
		for {
			daysLater = daysLater.AddDate(0, 0, num)
			res := daysLater.Compare(now)
			if res == 0 || res == 1 {
				break
			}
		}
	case 'y':
		for {
			daysLater = daysLater.AddDate(1, 0, 0)
			res := daysLater.Compare(now)
			if res == 0 || res == 1 {
				break
			}
		}
	case 'w':
		if len(s) != 2 {
			return "", errors.New("не указаны дни недели")
		}

		weekDays := strings.Split(s[1], ",")

		var weekDaysNums []int

		//проверяем корректность введенных дней недели
		for _, day := range weekDays {
			num, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}

			if !(num >= 1 && num <= 7) {
				return "", errors.New("неверный диапазон дней недели")
			}
			weekDaysNums = append(weekDaysNums, num)
		}

		var dayOfWeek = map[string]int{
			"Monday":    1,
			"Tuesday":   2,
			"Wednesday": 3,
			"Thursday":  4,
			"Friday":    5,
			"Saturday":  6,
			"Sunday":    7,
		}

		//вычисляем новую дату
		for {
			daysLater = daysLater.AddDate(0, 0, 1)
			res := daysLater.Compare(now)
			weekDayNum := dayOfWeek[daysLater.Weekday().String()]
			//containsDay := slices.Contains(weekDaysNums, weekDayNum)
			if res == 1 && slices.Contains(weekDaysNums, weekDayNum) {
				break
			}
		}
	case 'm':
		if len(s) < 2 || len(s) > 3 {
			return "", errors.New("некорректные параметры повторения")
		}

		days := strings.Split(s[1], ",")

		var daysNums []int
		//проверяем корректность введенных дней
		for _, day := range days {
			num, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}

			if !(num >= 1 && num <= 31) && num != -1 && num != -2 {
				return "", errors.New("неверный диапазон дней")
			}
			daysNums = append(daysNums, num)
		}
		var months []string
		if len(s) > 2 {
			months = strings.Split(s[2], ",")
		}

		var monthsNums []int
		//проверяем корректность введенных месяцев
		for _, month := range months {
			num, err := strconv.Atoi(month)
			if err != nil {
				return "", err
			}

			if !(num >= 1 && num <= 12) {
				return "", errors.New("неверный диапазон месяцев")
			}
			monthsNums = append(monthsNums, num)
		}

		if len(monthsNums) == 0 {
			monthsNums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		}

		var dateNxt time.Time
		var dates []time.Time
		year, _, _ := daysLater.Date()

		for _, day := range daysNums {
			for _, month := range monthsNums {
				for y := year; y < year+5; y++ {
					if day == -1 || day == -2 {
						firstOfMonth := time.Date(y, time.Month(month), 1, 0, 0, 0, 0, daysLater.Location())
						dateNxt = time.Date(y, time.Month(month), firstOfMonth.AddDate(0, 1, day).Day(), daysLater.Hour(),
							daysLater.Minute(), daysLater.Second(), daysLater.Nanosecond(), daysLater.Location())
					} else {
						dateNxt = time.Date(y, time.Month(month), day, daysLater.Hour(),
							daysLater.Minute(), daysLater.Second(), daysLater.Nanosecond(), daysLater.Location())
						if dateNxt.Month() != time.Month(month) || dateNxt.Day() != day {
							continue
						}
					}
					cmp := dateNxt.Compare(daysLater)
					if cmp == 1 {
						dates = append(dates, dateNxt)
					}
				}
			}
		}

		sort.Sort(ByDate(dates))
		for _, d := range dates {
			cmp := d.Compare(now)
			if cmp == 1 {
				daysLater = d
				break
			}
		}
	default:
		return "", errors.New("неподдерживаемый формат повторения")
	}

	return daysLater.Format(timeTemplate), nil
}

type ByDate []time.Time

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Before(a[j]) }
