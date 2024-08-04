package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id,omitempty" db:"id"`
	Date    string `json:"date,omitempty" db:"date"`
	Title   string `json:"title,omitempty" db:"title"`
	Comment string `json:"comment,omitempty" db:"comment"`
	Repeat  string `json:"repeat,omitempty" db:"repeat"`
}

type CreateTaskResponse struct {
	Id    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
	Error string `json:"error,omitempty"`
}

const DATE_FORMAT = "20060102"

func startServer() {
	const webDir = "web"
	const PORT = "7540"
	portString := os.Getenv("TODO_PORT")
	if len(portString) == 0 {
		portString = PORT
	}
	portString = fmt.Sprint(":", portString)
	log.Println("listen", portString)

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)

	log.Fatal(http.ListenAndServe(portString, nil))
}

func NextDate(now time.Time, date string, repeat string) (time.Time, error) {
	startDate, err := time.Parse(DATE_FORMAT, date)
	if err != nil {
		return time.Time{}, err
	}
	if repeat == "" {
		return time.Time{}, errors.New("empty repeat value")
	}
	nextDate := startDate
	if repeat == "y" {
		for {
			nextDate = nextDate.AddDate(1, 0, 0)
			if nextDate.After(now) {
				return nextDate, nil
			}
		}
	}
	if strings.HasPrefix(repeat, "d ") {
		daysStr := repeat[2:]
		s, err := strconv.ParseInt(daysStr, 10, 32)
		if err != nil {
			return time.Time{}, err
		}
		days := int(s)
		if days > 400 {
			return time.Time{}, errors.New("days can't be greater than 400")
		}
		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if nextDate.After(now) {
				return nextDate, nil
			}
		}
	}
	return time.Time{}, errors.New("unexpected repeat value")
}

func getDate(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}

	date, err := time.Parse(DATE_FORMAT, s)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil

}

func setErrorResponse(w http.ResponseWriter, s string, err error) {
	response := CreateTaskResponse{}
	response.Error = s
	if err != nil {
		response.Error += ": " + err.Error()
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func setOkResponse(w http.ResponseWriter) {
	response := CreateTaskResponse{}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	s := q["now"][0]
	now, err := time.Parse(DATE_FORMAT, s)
	if err != nil {
		http.Error(w, "cant parse now: "+s, http.StatusBadRequest)
	}
	date := q["date"][0]
	repeat := q["repeat"][0]

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "cant get next date: "+err.Error(), http.StatusBadRequest)
	}
	fmt.Fprint(w, nextDate.Format(DATE_FORMAT))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		addTask(w, r)
	} else if r.Method == "GET" {
		getTaskById(w, r)
	} else if r.Method == "PUT" {
		putTask(w, r)
	} else if r.Method == "DELETE" {
		deleteTask(w, r)
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	s, ok := q["id"]
	if !ok {
		setErrorResponse(w, "missing id url parameter", nil)
		return
	}
	id, err := strconv.ParseInt(s[0], 10, 32)
	if err != nil {
		setErrorResponse(w, "failed to convert id to int", err)
		return
	}
	err = deleteTaskById(int(id))
	if err != nil {
		setErrorResponse(w, "failed to delete task with id = "+strconv.Itoa(int(id)), err)
		return
	}

	response := CreateTaskResponse{}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		setErrorResponse(w, "failed to marshal response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func getTaskById(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	s, ok := q["id"]
	if !ok {
		setErrorResponse(w, "missing id url parameter", nil)
		return
	}
	id, err := strconv.ParseInt(s[0], 10, 32)
	if err != nil {
		setErrorResponse(w, "failed to convert id to int", err)
		return
	}
	task, err := loadTaskById(id)
	if err != nil {
		setErrorResponse(w, "task with id = "+strconv.Itoa(int(id))+" not found", err)
		return
	}

	jsonResponse, err := json.Marshal(task)
	if err != nil {
		setErrorResponse(w, "failed to marshal response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func addTask(w http.ResponseWriter, r *http.Request) {

	response := CreateTaskResponse{}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var task Task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		setErrorResponse(w, "json deserialization error", err)
		return
	}

	if task.Title == "" {
		setErrorResponse(w, "missing title", nil)
		return

	}

	date, err := getDate(task.Date)
	if err != nil {
		setErrorResponse(w, "bad date format", err)
		return
	}
	strDate := date.Format(DATE_FORMAT)
	if strDate < time.Now().Format(DATE_FORMAT) {
		if task.Repeat == "" {
			date = time.Now()
		} else {
			date, err = NextDate(time.Now(), strDate, task.Repeat)
			if err != nil {
				setErrorResponse(w, "failed to get next date", err)
				return
			}
		}
		strDate = date.Format(DATE_FORMAT)
	}

	id, err := insertTask(strDate, task.Title, task.Comment, task.Repeat)
	if err != nil {
		setErrorResponse(w, "failed to insert task into db", err)
		return
	}
	response.Id = id

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func putTask(w http.ResponseWriter, r *http.Request) {

	response := CreateTaskResponse{}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var task Task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		setErrorResponse(w, "json deserialization error", err)
		return
	}

	if task.Title == "" {
		setErrorResponse(w, "missing title", nil)
		return

	}

	date, err := getDate(task.Date)
	if err != nil {
		setErrorResponse(w, "bad date format", err)
		return
	}
	strDate := date.Format(DATE_FORMAT)
	if strDate < time.Now().Format(DATE_FORMAT) {
		if task.Repeat == "" {
			date = time.Now()
		} else {
			date, err = NextDate(time.Now(), strDate, task.Repeat)
			if err != nil {
				setErrorResponse(w, "failed to get next date", err)
				return
			}
		}
		strDate = date.Format(DATE_FORMAT)
	}
	id, err := strconv.ParseInt(task.ID, 10, 32)
	if err != nil {
		setErrorResponse(w, "failed to update task into db", err)
		return
	}

	err = updateTask(int(id), strDate, task.Title, task.Comment, task.Repeat)
	if err != nil {
		setErrorResponse(w, "failed to update task into db", err)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	s, ok := q["id"]
	if !ok {
		setErrorResponse(w, "missing id url parameter", nil)
		return
	}
	id, err := strconv.ParseInt(s[0], 10, 32)
	if err != nil {
		setErrorResponse(w, "failed to convert id to int", err)
		return
	}

	task, err := loadTaskById(id)
	if err != nil {
		setErrorResponse(w, "failed to load task by id", err)
		return
	}
	if task.Repeat == "" {
		deleteTaskById(int(id))
		setOkResponse(w)
		return
	}

	date, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		setErrorResponse(w, "failed to get next date", err)
		return
	}
	strDate := date.Format(DATE_FORMAT)

	err = updateTask(int(id), strDate, task.Title, task.Comment, task.Repeat)
	if err != nil {
		setErrorResponse(w, "failed to update task into db", err)
		return
	}

	setOkResponse(w)

}

func tasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := getAllTasks()
	if err != nil {
		setErrorResponse(w, "failed to load tasks from db", err)
		return
	}

	response := TasksResponse{}
	response.Tasks = tasks

	jsonResponse, err := json.Marshal(response)
	if err != nil {

		setErrorResponse(w, "failed to marshal response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
