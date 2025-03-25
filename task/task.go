package task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

var DateFormat = "20060102"

//var db *sql.DB

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Структура для ответа
type Response struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func ParseDate(dateStr string) (time.Time, error) {
	// Попробуем распарсить дату, если она передана
	if dateStr != "" {
		return time.Parse(DateFormat, dateStr)
	}
	// Если дата не передана, используем сегодняшнюю
	return time.Now(), nil
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		err := errors.New("не указано правило повторения")
		return "", err
	}

	dateForm, err := time.Parse(DateFormat, date)
	if err != nil {
		err := errors.New("указан неверный формат времени")
		return "", err
	}
	rules := strings.Split(repeat, " ")

	switch rules[0] {
	case "d":
		if len(rules) != 2 {
			err = errors.New("неподдерживаемый формат правила повторения")
			return "", err
		}
		days, err := strconv.Atoi(rules[1])
		if err != nil {
			err = errors.New("неподдерживаемый формат правила повторения")
			return "", err
		}

		if days > 400 || days < 1 {
			err = errors.New("недопустимое количество дней")
			return "", err
		}
		//Вычисляем новую дату
		if days == 1 {
			return now.Format(DateFormat), nil
		}
		for {
			dateForm = dateForm.AddDate(0, 0, days)
			if dateForm.After(now) {
				break
			}
		}
		return dateForm.Format(DateFormat), nil

	case "y":
		for {
			dateForm = dateForm.AddDate(1, 0, 0)
			if dateForm.After(now) {
				break
			}
		}
		return dateForm.Format(DateFormat), nil

	//case "w":
	//	dateForm = dateForm.AddDate(0, 0, 0)
	/*case "m":
	var returnDates string
	if len(rules) == 2{
		days := rules[1]
		day := strings.Split(days, ",")
		for _,d := range day {
			dInt, err := strconv.Atoi(d)
			if err != nil {
				err := errors.New("указан неверный формат времени")
				return "", err
			}
			dateRepeat := int(dateForm.Year())*1000 + int(dateForm.Month())*100 + dInt
			dateRepeatForm, err := time.Parse(DateFormat, string(dateRepeat))
			if err != nil {
				err := errors.New("указан неверный формат времени")
				return "", err
			}
			if dateForm.Before(dateRepeatForm){
				returnDates += dateForm.Format(DateFormat)
			}
		}
		return returnDates, nil
	}
	if len(rules) == 3{
		days := rules[1]
		day := strings.Split(days, ",")
		months := rules[2]
		month := strings.Split(months, ",")
		for _, m := range month{
			mInt, err := strconv.Atoi(m)
			if err != nil {
				err := errors.New("указан неверный формат времени")
				return "", err
			}
			for _,d := range day {
				dInt, err := strconv.Atoi(d)
				if err != nil {
					err := errors.New("указан неверный формат времени")
					return "", err
				}
				dateRepeat := dateForm.Year()*1000 + mInt*100 + dInt
				dateRepeatForm, err := time.Parse(DateFormat, string(dateRepeat))
				if err != nil {
					err := errors.New("указан неверный формат времени")
					return "", err
				}
				if dateForm.Before(dateRepeatForm){
					returnDates += dateForm.Format(DateFormat)
				}
			}
		}
		return returnDates, nil
	}
	*/

	default:
		err = errors.New("недопустимое правило повторения")
		return "", err
	}
}

func AddTask(w http.ResponseWriter, req *http.Request) {
	var task Task
	// Декодируем JSON из запроса
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}
	// Проверяем обязательное поле title
	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверяем и парсим дату
	date, err := ParseDate(task.Date)
	if err != nil {
		http.Error(w, `{"error":"Неправильный формат времени"}`, http.StatusBadRequest)
		return
	}

	var dateStr string
	// Если дата задачи меньше сегодняшней, вычисляем новую дату с учетом повторения
	if date.Before(time.Now()) {
		if task.Repeat == "" {
			// Если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число
			date = time.Now()
			dateStr = date.Format(DateFormat)
		} else {
			// Если повторение указано, вычисляем следующую дату с учетом правила
			dateStr, err = NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Неправильное вычисление даты"}`, http.StatusBadRequest)
				return
			}
		}
	} else {
		dateStr = date.Format(DateFormat)
	}
	task.Date = dateStr
	//:date, :title, :comment, :repeat
	//Добавляем задачу в базу данных
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return

	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"insert failed"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error":"scan id failed"}`, http.StatusInternalServerError)
		return
	}

	// Формируем ответ

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]interface{}{"id": id}
	json.NewEncoder(w).Encode(response)
	// Возвращаем ответ в формате JSON
	/*w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]interface{}{"id": id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Ошибка при формировании ответа"}`, http.StatusInternalServerError)
		return
	}*/
}
