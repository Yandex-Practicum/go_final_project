package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/internal/db"
	"go_final_project/internal/repository"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

var format = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		var buf bytes.Buffer
		now := time.Now()
		formatNow := now.Format(format)

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if buf.Len() == 0 {
			errPost := map[string]string{"error": "Пустое тело запроса"}

			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			fmt.Print("поле title пустое")
			errPost := map[string]string{"error": "поле title пустое"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if task.Date == "" {
			task.Date = formatNow
		}

		if len(task.Date) > 8 || len(task.Date) < 8 {
			log.Printf("Некорректное количество символов")
			errPost := map[string]string{"error": "Некорректное количество символов"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		inputMonth, _ := strconv.Atoi(task.Date[4:6])
		inputDay, _ := strconv.Atoi(task.Date[6:8])

		if inputMonth < 1 || inputMonth > 12 {
			log.Printf("Некорректный месяц '%d' в дате '%s'", inputMonth, task.Date)
			errPost := map[string]string{"error": "некорректный месяц"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if inputDay < 1 || inputDay > 31 {
			log.Printf("Некорректный день '%d' в дате '%s'", inputDay, task.Date)
			errPost := map[string]string{"error": "некорректный день"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		parsedDate, err := time.Parse(format, task.Date)
		if err != nil {
			log.Printf("ошибка парсинга даты '%s': %v", task.Date, err)
			errPost := map[string]string{"error": "ошибка парсинга даты"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if parsedDate.Format(format) != task.Date {
			log.Printf("Дата '%s' не соответствует формату '%s'", task.Date, format)
			errPost := map[string]string{"error": "некорректная дата"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if task.Repeat == "" {
			task.Date = formatNow
		}

		if task.Repeat != "" {
			if task.Date != formatNow {
				task.Date, err = NextDate(now, task.Date, task.Repeat)
				if err != nil {
					fmt.Print("не удалось вычислить следующую дату")
					errPost := map[string]string{"error": "не удалось вычислить следующую дату"}
					if err := json.NewEncoder(w).Encode(errPost); err != nil {
						log.Printf("Ошибка кодирования JSON: %v", err)
						http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
						return
					}
				}
			} else {
				task.Date = formatNow
			}
		}

		dataInt, err := strconv.Atoi(task.Date)
		if err != nil {
			log.Printf("ошибка конвертации date '%s' в число: %v", task.Date, err)
			return
		}

		formatNowInt, err := strconv.Atoi(formatNow)
		if err != nil {
			log.Printf("ошибка конвертации текущей даты '%s' в число: %v", formatNow, err)
			return
		}

		if dataInt < formatNowInt {
			if task.Repeat == "" {
				task.Date = formatNow
			} else {
				task.Date, err = NextDate(now, task.Date, task.Repeat)
				if err != nil {
					fmt.Print("не удалось вычислить следующую дату")
					errPost := map[string]string{"error": "не удалось вычислить следующую дату"}
					if err := json.NewEncoder(w).Encode(errPost); err != nil {
						log.Printf("Ошибка кодирования JSON: %v", err)
						http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
						return
					}
					return
				}
			}
		}

		query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
		res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)

		if err != nil {
			log.Printf("Ошибка при выполнении запроса: %v", err)
			http.Error(w, "Ошибка записи в базу данных", http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, fmt.Sprintf("не удалось получить ID последней вставленной записи: %v", err), http.StatusInternalServerError)
			return
		}
		idStr := strconv.FormatInt(id, 10)
		task.ID = idStr

		response := map[string]string{"id": task.ID}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Ошибка кодирования JSON: %v", err)
			http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
			return
		}
	}
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	var nextDate string

	nowString := now.Format(format)

	dataInt, err := strconv.Atoi(date)
	if err != nil {
		return "", fmt.Errorf("ошибка конвертации date '%s' в число: %w", date, err)
	}

	nowInt, err := strconv.Atoi(nowString)
	if err != nil {
		return "", fmt.Errorf("ошибка конвертации now '%s' в число: %w", nowString, err)
	}

	parsedTime, err := time.Parse(format, date)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга даты '%s': %w", date, err)
	}

	if repeat == "" {
		return "", fmt.Errorf("пустая строка в колонке repeat")
	}

	if repeat == "y" {
		if dataInt <= nowInt {
			for dataInt <= nowInt {
				parsedTime = parsedTime.AddDate(1, 0, 0)
				dataInt, err = strconv.Atoi(parsedTime.Format("20060102"))
				if err != nil {
					return "", fmt.Errorf("ошибка преобразования даты в число: %w", err)
				}
				nextDate = parsedTime.Format("20060102")
			}
			return nextDate, nil
		}
		if dataInt > nowInt {
			parsedTime = parsedTime.AddDate(1, 0, 0)
			nextDate = parsedTime.Format("20060102")
			return nextDate, nil
		}
	}

	repeatSlices := strings.Split(repeat, " ")
	if len(repeatSlices) != 2 {
		return "", fmt.Errorf("неверный формат repeat: '%s'", repeat)
	}

	repeatDay, err := strconv.Atoi(repeatSlices[1])
	if err != nil {
		return "", fmt.Errorf("ошибка конвертации repeat '%s' в число: %w", repeatSlices[1], err)
	}

	if repeatSlices[0] != "d" && repeatSlices[0] != "y" {
		return "", fmt.Errorf("недопустимый символ в repeat: '%s'", repeatSlices[0])
	}

	if repeatDay > 400 {
		return "", fmt.Errorf("превышен максимально допустимый интервал: %d", repeatDay)
	}

	if repeatSlices[0] == "d" {
		if dataInt <= nowInt {
			for dataInt <= nowInt {
				parsedTime = parsedTime.AddDate(0, 0, repeatDay)
				dataInt, err = strconv.Atoi(parsedTime.Format("20060102"))
				if err != nil {
					return "", fmt.Errorf("ошибка преобразования даты в число: %w", err)
				}
				nextDate = parsedTime.Format(format)
			}
			return nextDate, nil
		}
		if dataInt > nowInt {
			parsedTime = parsedTime.AddDate(0, 0, repeatDay)
			nextDate = parsedTime.Format("20060102")
			return nextDate, nil
		}
	}

	return "", fmt.Errorf(
		"не удалось вычислить следующую дату: now=%s, date=%s, repeat=%s",
		now.Format(format),
		date,
		repeat,
	)
}

func repeatTask(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nowParam := query.Get("now")
	date := query.Get("date")
	repeat := query.Get("repeat")

	paresedNow, err := time.Parse(format, nowParam)
	if err != nil {
		log.Printf("Ошибка парсинга параметра now='%s': %v", nowParam, err)
		return
	}

	repeatDecoded, err := url.QueryUnescape(repeat)
	if err != nil {
		log.Printf("Ошибка декодирования repeat: %v", err)
		return
	}

	nextDate, err := NextDate(paresedNow, date, repeatDecoded)
	if err != nil {
		log.Printf("Ошибка в NextDate: %v (now=%s, date=%s, repeat=%s)", err, nowParam, date, repeat)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

func GetTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Tasks struct {
			Tasks []Task `json:"tasks"`
		}

		query := "SELECT * FROM scheduler ORDER BY date, repeat LIMIT 10"
		rows, err := db.Query(query)

		if err != nil {
			errPost := map[string]string{"error": "Ошибка записи в базу данных"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				errPost := map[string]string{"error": "Ошибка сканирования строк"}
				if err := json.NewEncoder(w).Encode(errPost); err != nil {
					log.Printf("Ошибка кодирования JSON: %v", err)
					http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
					return
				}
				return
			}
			tasks = append(tasks, task)
		}

		response := Tasks{
			Tasks: tasks,
		}

		if len(tasks) == 0 {
			response.Tasks = []Task{}
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		result := r.URL.Query()
		id := result.Get("id")

		if id == "" {
			errPost := map[string]string{"error": "Не указан идентификатор"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		query := "SELECT * FROM scheduler WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			errPost := map[string]string{"error": "Задача не найдена"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if err := json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PutTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		var buf bytes.Buffer
		now := time.Now()
		formatNow := now.Format(format)

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if buf.Len() == 0 {
			errPost := map[string]string{"error": "Пустое тело запроса"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			errPost := map[string]string{"error": "Ошибка парсинга"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if task.ID == "" {
			errPost := map[string]string{"error": "Не указан идентификатор"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		_, err = strconv.Atoi(task.ID)
		if err != nil {
			errPost := map[string]string{"error": "Неверный формат ID"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		queryID := "SELECT COUNT(*) FROM scheduler WHERE id = ?"
		var count int
		err = db.QueryRow(queryID, task.ID).Scan(&count)
		if err != nil {
			errPost := map[string]string{"error": "Ошибка выполнения запроса"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if count == 0 {
			log.Printf("Задача с id=%s не найдена", task.ID)
			errPost := map[string]string{"error": "Задача с данным ID не найдена"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if task.Title == "" {
			fmt.Print("поле title пустое")
			errPost := map[string]string{"error": "поле title пустое"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if task.Date == "" {
			task.Date = formatNow
		}

		if len(task.Date) > 8 || len(task.Date) < 8 {
			log.Printf("Некорректное количество символов")
			errPost := map[string]string{"error": "Некорректное количество символов"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		inputMonth, _ := strconv.Atoi(task.Date[4:6])
		inputDay, _ := strconv.Atoi(task.Date[6:8])

		if inputMonth < 1 || inputMonth > 12 {
			log.Printf("Некорректный месяц '%d' в дате '%s'", inputMonth, task.Date)
			errPost := map[string]string{"error": "некорректный месяц"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if inputDay < 1 || inputDay > 31 {
			log.Printf("Некорректный день '%d' в дате '%s'", inputDay, task.Date)
			errPost := map[string]string{"error": "некорректный день"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		parsedDate, err := time.Parse(format, task.Date)
		if err != nil {
			log.Printf("ошибка парсинга даты '%s': %v", task.Date, err)
			errPost := map[string]string{"error": "ошибка парсинга даты"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if parsedDate.Format(format) != task.Date {
			log.Printf("Дата '%s' не соответствует формату '%s'", task.Date, format)
			errPost := map[string]string{"error": "некорректная дата"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if task.Repeat == "" {
			task.Date = formatNow
		}

		if task.Repeat != "" {
			if task.Date != formatNow {
				task.Date, err = NextDate(now, task.Date, task.Repeat)
				if err != nil {
					fmt.Print("не удалось вычислить следующую дату")
					errPost := map[string]string{"error": "не удалось вычислить следующую дату"}
					if err := json.NewEncoder(w).Encode(errPost); err != nil {
						log.Printf("Ошибка кодирования JSON: %v", err)
						http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
						return
					}
				}
			} else {
				task.Date = formatNow
			}
		}

		dataInt, err := strconv.Atoi(task.Date)
		if err != nil {
			log.Printf("ошибка конвертации date '%s' в число: %v", task.Date, err)
			return
		}

		formatNowInt, err := strconv.Atoi(formatNow)
		if err != nil {
			log.Printf("ошибка конвертации текущей даты '%s' в число: %v", formatNow, err)
			return
		}

		if dataInt < formatNowInt {
			if task.Repeat == "" {
				task.Date = formatNow
			} else {
				task.Date, err = NextDate(now, task.Date, task.Repeat)
				if err != nil {
					fmt.Print("не удалось вычислить следующую дату")
					errPost := map[string]string{"error": "не удалось вычислить следующую дату"}
					if err := json.NewEncoder(w).Encode(errPost); err != nil {
						log.Printf("Ошибка кодирования JSON: %v", err)
						http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
						return
					}
					return
				}
			}
		}

		query := "UPDATE scheduler SET id = ?, date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
		_, err = db.Exec(query, task.ID, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			log.Printf("Ошибка выполнения запроса: %v", err)
			http.Error(w, "Ошибка обновления записи", http.StatusInternalServerError)
			return
		}

		if err != nil {
			fmt.Print("Задача не найдена")
			errPost := map[string]string{"error": "Задача не найдена"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		} else {
			errPost := map[string]string{}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
		}
	}
}

func DoneTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		now := time.Now()

		result := r.URL.Query()
		id := result.Get("id")

		if id == "" {
			errPost := map[string]string{"error": "Не указан идентификатор"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		_, err := strconv.Atoi(id)
		if err != nil {
			errPost := map[string]string{"error": "Неверный формат ID"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		queryID := "SELECT COUNT(*) FROM scheduler WHERE id = ?"
		var count int
		err = db.QueryRow(queryID, id).Scan(&count)
		if err != nil {
			errPost := map[string]string{"error": "Ошибка выполнения запроса"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if count == 0 {
			log.Printf("Задача с id=%s не найдена", id)
			errPost := map[string]string{"error": "Задача с данным ID не найдена"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		queryDelete := "DELETE FROM scheduler WHERE id = ?"
		queryInsert := "UPDATE scheduler SET date = ? WHERE id = ?"

		err = db.QueryRow("SELECT id, title, date, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Title, &task.Date, &task.Repeat)
		if err != nil {
			errPost := map[string]string{}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		fmt.Print(task)

		log.Printf("Обработка задачи ID: %s, Текущая дата: %s, Заданная дата: %s, Повтор: %s",
			id, now.Format("20060102"), task.Date, task.Repeat)

		if task.Repeat != "" {
			task.Date, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				log.Print("Невыполнено обновление даты")
				errPost := map[string]string{"error": "Невыполнено обновление даты"}
				if err := json.NewEncoder(w).Encode(errPost); err != nil {
					log.Printf("Ошибка кодирования JSON: %v", err)
					http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
					return
				}
				return
			}
			log.Printf("Обновление даты для задачи ID: %s, Новая дата: %s", id, task.Date)

			_, err = db.Exec(queryInsert, task.Date, id)
			if err != nil {
				log.Print("Ошибка изменения даты")
				errPost := map[string]string{"error": "Ошибка изменения даты"}
				if err := json.NewEncoder(w).Encode(errPost); err != nil {
					log.Printf("Ошибка кодирования JSON: %v", err)
					http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
					return
				}
				return
			} else {
				errPost := map[string]string{}
				if err := json.NewEncoder(w).Encode(errPost); err != nil {
					log.Printf("Ошибка кодирования JSON: %v", err)
					http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
					return
				}
			}
		}

		if task.Repeat == "" {
			_, err := db.Exec(queryDelete, id)
			if err != nil {
				log.Print("Ошибка удаления задачи")
				errPost := map[string]string{"error": "Ошибка удаления задачи"}
				if err := json.NewEncoder(w).Encode(errPost); err != nil {
					log.Printf("Ошибка кодирования JSON: %v", err)
					http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
					return
				}
			} else {
				errPost := map[string]string{}
				if err := json.NewEncoder(w).Encode(errPost); err != nil {
					log.Printf("Ошибка кодирования JSON: %v", err)
					http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

func DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := r.URL.Query()
		id := result.Get("id")

		if id == "" {
			errPost := map[string]string{"error": "Не указан идентификатор"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		_, err := strconv.Atoi(id)
		if err != nil {
			errPost := map[string]string{"error": "Неверный формат ID"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		queryID := "SELECT COUNT(*) FROM scheduler WHERE id = ?"
		var count int
		err = db.QueryRow(queryID, id).Scan(&count)
		if err != nil {
			errPost := map[string]string{"error": "Ошибка выполнения запроса"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		}

		if count == 0 {
			log.Printf("Задача с id=%s не найдена", id)
			errPost := map[string]string{"error": "Задача с данным ID не найдена"}
			if err := json.NewEncoder(w).Encode(errPost); err != nil {
				log.Printf("Ошибка кодирования JSON: %v", err)
				http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
				return
			}
			return
		} else {
			query := "DELETE FROM scheduler WHERE id = ?"
			_, err := db.Exec(query, id)
			if err != nil {
				log.Print("Ошибка удаления задачи")
				errPost := map[string]string{"error": "Ошибка удаления задачи"}
				if err := json.NewEncoder(w).Encode(errPost); err != nil {
					log.Printf("Ошибка кодирования JSON: %v", err)
					http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
					return
				}
			} else {
				errPost := map[string]string{}
				if err := json.NewEncoder(w).Encode(errPost); err != nil {
					log.Printf("Ошибка кодирования JSON: %v", err)
					http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

func main() {
	db := db.New()

	rep := repository.New(db)
	migration(rep)

	r := chi.NewRouter()

	r.Get("/api/nextdate", repeatTask)
	r.Post("/api/task", AddTask(db))
	r.Get("/api/tasks", GetTasks(db))
	r.Get("/api/task", GetTask(db))
	r.Put("/api/task", PutTask(db))
	r.Post("/api/task/done", DoneTask(db))
	r.Delete("/api/task", DeleteTask(db))

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":7540", r)
	if err != nil {
		fmt.Println(err)
	}
}

func migration(rep *repository.Repository) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install {
		if err := rep.CreateScheduler(); err != nil {
			log.Fatal(err)
		}
	}
}
