package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ErrResponse struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"error_message"`
}

type ReturnResponse struct {
	ID int64 `json:"id"`
}

func sendJson(rw http.ResponseWriter, status int, r any) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	rw.Header().Set("Content-Type", "application/json") // set content-type so our clients know how to read our response
	rw.WriteHeader(status)                              // write our status header with the proper http code
	// write the marshalled json into the response.
	// as per documentation, this is a final call in handling requests and will finish the handling process.
	_, err = rw.Write(b)
	return err
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createTaskHandler(w, r)

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		editTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

func createTask(t *Task) (int64, error) {

	db, err := sql.Open("sqlite3", DBFile)
	defer db.Close()
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?) RETURNING id"
	result, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	return lastId, nil
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// read data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var task *Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		panic(err)
	}

	// validate

	if valid, message := task.IsValid(); !valid {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: message,
		}); err != nil {
			panic(err) // couldn't even send json, panic
		}
	}
	// execute
	valid, message, newTask := task.MakeValid()

	if !valid {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: message,
		}); err != nil {
			panic(err) // couldn't even send json, panic
		}
	}

	newTaskID, err := createTask(newTask)
	if err != nil {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: err.Error(),
		}); err != nil {
			panic(err) // couldn't even send json, panic
		}
	}

	// send json
	if err := sendJson(w, 200, ReturnResponse{
		ID: newTaskID,
	}); err != nil {
		panic(err)
	}
}

func getTask(id int) (*Task, error) {
	db, err := sql.Open("sqlite3", DBFile)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM scheduler WHERE id=?"
	row := db.QueryRow(query, id)

	var task *Task
	switch err := row.Scan(task); err {
	case sql.ErrNoRows:
		return nil, nil // no task, but no error too
	case nil:
		return task, nil // no error, return task
	default:
		return nil, err // any other error, return it as is
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// get input
	idStr := r.URL.Query().Get("id") // r.FormValue extracts from the request body too, if available, we don't need that

	// validate
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: "invalid id",
		}); err != nil {
			panic(err)
		}
	}

	// execute
	task, err := getTask(id)

	// handle error
	if err != nil {
		if err := sendJson(w, 500, ErrResponse{
			Error:        true,
			ErrorMessage: "internal server error",
		}); err != nil {
			panic(err)
		}
	}

	// handle 404
	if task == nil {
		if err := sendJson(w, 404, ErrResponse{
			Error:        true,
			ErrorMessage: "no task with this id",
		}); err != nil {
			panic(err)
		}
	}

	// return
	if err := sendJson(w, 200, task); err != nil {
		panic(err)
	}
}

func getTasks(limit int) (*Tasks, error) {
	db, err := sql.Open("sqlite3", DBFile)
	defer db.Close()
	if err != nil {
		fmt.Println(1, err)
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", limit)
	if err != nil {
		fmt.Println(2, err)
		return nil, err
	}

	tasks := &Tasks{}

	defer rows.Close()

	for rows.Next() {
		task := &Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(3, err)
			return nil, err
		}
		tasks.Tasks = append(tasks.Tasks, *task)
	}
	if tasks.Tasks == nil {
		tasks.Tasks = []Task{}
	}
	return tasks, nil
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := getTasks(tasksLimit)
	if err != nil {
		if err := sendJson(w, 500, ErrResponse{
			Error:        true,
			ErrorMessage: "internal server error",
		}); err != nil {
			panic(err)
		}
	}
	if err := sendJson(w, 200, tasks); err != nil {
		panic(err)
	}
}

func editTask(id int, t *Task) error {
	db, err := sql.Open("sqlite3", DBFile)
	defer db.Close()
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	_, err = db.Exec(query,
		t.Date,
		t.Title,
		t.Comment,
		t.Repeat,
		id)

	if err != nil {
		return err
	}
	return nil
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	// read
	idStr := r.URL.Query().Get("id")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var task *Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		panic(err)
	}

	// validate
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: "invalid id",
		}); err != nil {
			panic(err)
		}
	}
	valid, message := task.IsValid()
	if !valid {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: message,
		}); err != nil {
			panic(err) // couldn't even send json, panic
		}
	}

	valid, message, newTask := task.MakeValid()
	if !valid {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: message,
		}); err != nil {
			panic(err) // couldn't even send json, panic
		}
	}
	//fmt.Println(2, valid, message)

	// execute
	err = editTask(id, newTask)
	if err != nil {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: err.Error(),
		}); err != nil {
			panic(err) // couldn't even send json, panic
		}
	}
	// return
	if err := sendJson(w, 200, task); err != nil {
		panic(err)
	}
}

func deleteTask(id int) error {
	db, err := sql.Open("sqlite3", DBFile)
	defer db.Close()
	if err != nil {
		return err
	}

	query := "DELETE FROM scheduler WHERE id=?"
	_, err = db.Exec(query, id)

	if err != nil {
		return err

	}
	return nil
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// read
	idStr := r.URL.Query().Get("id")

	// validate
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: "invalid id",
		}); err != nil {
			panic(err)
		}
	}
	// execute
	err = deleteTask(id)
	if err != nil {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: "invalid id",
		}); err != nil {
			panic(err)
		}
	}
	// return
	if err := sendJson(w, 200, struct{}{}); err != nil {
		panic(err)
	}
}

func doneTask(id int, t *Task) error {
	db, err := sql.Open("sqlite3", DBFile)
	defer db.Close()
	if err != nil {
		return err
	}

	newDate, err := NextDate(time.Now(), t.Date, t.Repeat)

	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	_, err = db.Exec(query, newDate, id)

	if err != nil {
		return err
	}
	return nil
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var task *Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		panic(err)
	}

	// validate
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: "invalid id",
		}); err != nil {
			panic(err)
		}
	}
	// execute
	err = doneTask(id, task)
	if err != nil {
		if err := sendJson(w, 400, ErrResponse{
			Error:        true,
			ErrorMessage: "invalid id",
		}); err != nil {
			panic(err)
		}
	}
	// return
	if err := sendJson(w, 200, struct{}{}); err != nil {
		panic(err)
	}
}
