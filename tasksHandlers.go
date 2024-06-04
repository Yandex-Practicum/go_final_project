package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func checkInputJSON(res http.ResponseWriter, task *Task) error {
	var err error
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if task.Title == "" {
		err := errors.New("Title error")
		response.Error = err.Error()
		json.NewEncoder(res).Encode(response)
		return err
	}

	date := time.Now()

	if task.Date != "" {
		date, err = time.Parse(dataFormat, task.Date)
		if err != nil {
			err := errors.New("Date error")
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return err
		}
	} else {
		task.Date = date.Format(dataFormat)
	}

	if date.Before(time.Now()) {
		if task.Repeat == "" {
			task.Date = time.Now().Format(dataFormat)
		} else {
			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				err := errors.New("Date error")
				response.Error = err.Error()
				json.NewEncoder(res).Encode(response)
				return err
			}
		}
	}
	return nil
}

func taskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var buf bytes.Buffer
		task := &Task{}
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			err := errors.New("JSON read error")
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &task)
		if err != nil {
			err := errors.New("JSON Unmarshalling error")
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}

		err = checkInputJSON(res, task)
		if err != nil {
			err := errors.New("Title error")
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
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

	case http.MethodGet:
		id := req.FormValue("id")
		if id == "" {
			err := errors.New("Не укаазан ID задачи")
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}

		db, err := sql.Open("sqlite3", DBFile)
		defer db.Close()
		if err != nil {
			return
		}

		query := "SELECT * FROM scheduler WHERE id=?"
		row, err := db.Query(query, id)

		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return

		}
		task := &Task{}
		for row.Next() {
			err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(res).Encode(response)
				return
			}
		}
		row.Close()
		if err = row.Err(); err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}

		data, err := json.Marshal(&task)
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}
		res.WriteHeader(http.StatusOK)
		_, err = res.Write(data)
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}

	case http.MethodPut:
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

		err = checkInputJSON(res, task)
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}

		db, err := sql.Open("sqlite3", DBFile)
		defer db.Close()
		query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
		_, err = db.Exec(query,
			task.Date,
			task.Title,
			task.Comment,
			task.Repeat,
			task.ID)

		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}
		json.NewEncoder(res).Encode(response)
		return

	case http.MethodDelete:
		id := req.FormValue("id")
		if id == "" {
			err := errors.New("Не укаазан ID задачи")
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}

		db, err := sql.Open("sqlite3", DBFile)
		defer db.Close()
		if err != nil {
			return
		}

		query := "DELETE FROM scheduler WHERE id=?"
		_, err = db.Exec(query, id)

		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return

		}
		json.NewEncoder(res).Encode(response)
	}
}

func doneTask(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		id := req.FormValue("id")
		if id == "" {
			err := errors.New("Не укаазан ID задачи")
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}

		db, err := sql.Open("sqlite3", DBFile)
		defer db.Close()
		if err != nil {
			return
		}

		query := "SELECT * FROM scheduler WHERE id=?"
		row, err := db.Query(query, id)

		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return

		}
		task := &Task{}
		for row.Next() {

			err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(res).Encode(response)
				return
			}
		}
		row.Close()
		if task.Repeat == "" {
			query := "DELETE FROM scheduler WHERE id=?"
			_, err := db.Exec(query, id)

			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(res).Encode(response)
				return
			}
		} else {
			newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(res).Encode(response)
				return
			}
			query := "UPDATE scheduler SET date=? WHERE id=?"
			_, err = db.Exec(query, newDate, id)

			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(res).Encode(response)
				return
			}
		}
		_, err = json.Marshal(response)
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}
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
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}

		rows, err := db.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", tasksLimit)
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}
		tasks := &Tasks{}
		for rows.Next() {
			task := &Task{}

			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				response.Error = err.Error()
				json.NewEncoder(res).Encode(response)
				return
			}
			tasks.Tasks = append(tasks.Tasks, *task)
		}
		rows.Close()
		if tasks.Tasks == nil {
			tasks.Tasks = []Task{}
		}
		data, err := json.Marshal(&tasks)
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}
		res.WriteHeader(http.StatusOK)
		_, err = res.Write(data)
		if err != nil {
			response.Error = err.Error()
			json.NewEncoder(res).Encode(response)
			return
		}
	}
}
