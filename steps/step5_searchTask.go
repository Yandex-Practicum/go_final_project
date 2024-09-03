package steps

import (
	"database/sql"
	"net/http"
	"regexp"
	"time"
)

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Err struct {
	Error string `json:"error,omitempty"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "scheduler.db")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, Err{Error: err.Error()})
		return
	}

	err = db.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, Err{Error: err.Error()})
		return
	}

	defer db.Close()

	tasks, err := SearchField(db, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, Err{Error: "ошибка поиска задачи"})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeError(w, Tasks{Tasks: tasks})
}

func SearchField(db *sql.DB, w http.ResponseWriter, r *http.Request) ([]Task, error) {
	search := r.FormValue("search")
	tasks := make([]Task, 0, 20)

	reg, err := regexp.Compile("\\d{2}.\\d{2}.\\d{4}")

	if err != nil {
		return tasks, err
	}

	var rows *sql.Rows
	match := reg.MatchString(search)
	if match {

		parseSearch, err := time.Parse("02.01.2006", search)

		if err != nil {
			return tasks, err
		}

		timeSearch := parseSearch.Format(DateForFormat)

		rows, err = db.Query("SELECT * FROM scheduler WHERE date LIKE :search ORDER BY date LIMIT :limit",
			sql.Named("search", "%"+timeSearch+"%"),
			sql.Named("limit", 20))

		if err != nil {
			return tasks, err
		}

		defer rows.Close()

		for rows.Next() {
			var task Task
			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

			if err != nil {
				return tasks, err
			}

			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			return tasks, err
		}

		return tasks, nil
	}

	rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
		sql.Named("search", "%"+search+"%"),
		sql.Named("limit", 20))

	if err != nil {
		return tasks, err
	}

	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}
