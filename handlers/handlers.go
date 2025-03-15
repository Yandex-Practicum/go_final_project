package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/qwsnxnjene/go_final_project/db"
	_ "modernc.org/sqlite"
)

// Лимит задач, которые будут возвращаться при поиске
const TaskLimit = 50

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// NextDateHandler() обрабатывает GET-запросы по адресу /api/nextdate
func NextDateHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	date, repeat := r.FormValue("date"), r.FormValue("repeat")
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}
	newDate, err := NextDate(now, date, repeat)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	} else {
		rw.Write([]byte(newDate))
	}
}

// TaskHandler() обрабатывает запросы по адресу /api/task
func TaskHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(rw, r)
	case http.MethodGet:
		taskByIdHandler(rw, r)
	case http.MethodPut:
		updateTaskHandler(rw, r)
	case http.MethodDelete:
		deleteTaskHandler(rw, r)
	}
}

// TasksHandler() обрабатывает GET-запросы по адресу /api/tasks
func TasksHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

	toSearch := r.FormValue("search")

	tasks := []Task{}
	dB := db.Database
	var query string
	var rows *sql.Rows
	if len(toSearch) == 0 {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`
		rows, err := dB.Query(query, sql.Named("limit", TaskLimit))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
			return
		}
		defer rows.Close()
	} else {
		searchTime, err := time.Parse("02.01.2006", toSearch)
		if err != nil {
			//значит запрос не временной
			query = `SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`
			rows, err = dB.Query(query, sql.Named("limit", TaskLimit), sql.Named("search", "%"+toSearch+"%"))
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
				return
			}
			defer rows.Close()
		} else {
			timeToFind := searchTime.Format("20060102")
			query = `SELECT * FROM scheduler WHERE date = :date LIMIT :limit`
			rows, err = dB.Query(query, sql.Named("limit", TaskLimit), sql.Named("date", timeToFind))
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
				return
			}
			defer rows.Close()
		}

	}

	for rows.Next() {
		var t Task

		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
			return
		}
		tasks = append(tasks, t)
	}
	err := rows.Err()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
		return
	}

	var jsonTask struct {
		Tasks []Task `json:"tasks"`
	}

	jsonTask.Tasks = tasks
	data, err := json.Marshal(jsonTask)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка сериализации %v"}`, err)))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

// TaskDoneHandler() обрабатывает POST-запросы по адресу /api/task/done
func TaskDoneHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.FormValue("id")
	if len(id) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte((`{"error":"не указан идентификатор"}`)))
		return
	}

	dB := db.Database

	idInt, err := strconv.Atoi(id)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с id %v"}`, err)))
		return
	}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`
	row := dB.QueryRow(query, sql.Named("id", idInt))

	var task Task
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"error":"запись не найдена"}`))
			return
		}
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
		return
	}

	if len(task.Repeat) == 0 {
		query := `DELETE FROM scheduler WHERE id = :id`
		res, err := dB.Exec(query, sql.Named("id", id))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
			return
		}
		if rows, err := res.RowsAffected(); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
			return
		} else if rows == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"error":"задача не найдена"}`))
			return
		}
	} else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с датой %v"}`, err)))
			return
		}

		query := `UPDATE scheduler SET date = :date WHERE id = :id`
		res, err := dB.Exec(query,
			sql.Named("id", id),
			sql.Named("date", nextDate),
		)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
			return
		}
		if rows, err := res.RowsAffected(); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка работы с БД %v"}`, err)))
			return
		} else if rows == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"error":"задача не найдена"}`))
			return
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{}`))
}

// SignInHandler() обрабатывает POST-запросы по адресу /api/signin
func SignInHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var p struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf(`{"error":"ошибка десериализации %v"}`, err)))
		return
	}
	defer r.Body.Close()

	var ans string
	if p.Password == os.Getenv("TODO_PASSWORD") {
		secret := []byte(os.Getenv("TODO_PASSWORD"))

		hash := sha256.Sum256([]byte(os.Getenv("TODO_PASSWORD")))

		claims := jwt.MapClaims{
			"hash": hex.EncodeToString(hash[:]),
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := jwtToken.SignedString(secret)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Errorf("failed to sign jwt: %s", err).Error()))
			return
		}

		ans = fmt.Sprintf(`{"token":"%v"}`, signedToken)
		rw.WriteHeader(http.StatusOK)
		fmt.Println(signedToken)
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		ans = `{"error":"Неверный пароль"}`
	}
	rw.Write([]byte(ans))
}
