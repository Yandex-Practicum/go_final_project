package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Xa4y322/go_final_project/database"
	"github.com/Xa4y322/go_final_project/models"
)

const DateFormat string = "20060102"

func nextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeate is empty")

	}
	if strings.Contains(repeat, "d ") {
		daysToAdd, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil {
			return "", err
		}
		if daysToAdd > 400 {
			return "", errors.New("repeat period in days max 400")
		}

		beginDate, err := time.Parse(DateFormat, date)
		if err != nil {
			return "", err
		}

		newDate := beginDate.AddDate(0, 0, daysToAdd)
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, daysToAdd)
		}

		return newDate.Format(DateFormat), nil

	}
	if repeat == "y" {
		beginDate, err := time.Parse(DateFormat, date)
		if err != nil {
			return "", err
		}

		newDate := beginDate.AddDate(1, 0, 0)

		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}

		return newDate.Format(DateFormat), nil
	} else {
		return "", errors.New("repeat string has wrong format")
	}
}

func NextDateGET(w http.ResponseWriter, r *http.Request) {
	pNow, err := time.Parse(DateFormat, r.FormValue("now"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pDate := r.FormValue("date")
	pRepeat := r.FormValue("repeat")
	newDate, err := nextDate(pNow, pDate, pRepeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(newDate))

	if err != nil {
		http.Error(w, fmt.Sprintf("nextDateGET error: %v", err), http.StatusBadRequest)
	}
}

func responseWithError(w http.ResponseWriter, err error) {
	log.Printf("Error: %w", err)

	error, _ := json.Marshal(models.ResponseError{Error: err.Error()})

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(error)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func TaskPost(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		responseWithError(w, err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		responseWithError(w, err)
		return
	}

	//проверки корректного заполнения полей
	if len(task.Title) == 0 {
		responseWithError(w, errors.New("title is empty"))
		return
	}

	if len(task.Date) == 0 {
		task.Date = time.Now().Format(DateFormat)
	} else {
		_, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			responseWithError(w, errors.New("bad date"))
			return
		}

		if len(task.Repeat) > 0 {
			if nextDate, err := nextDate(time.Now(), task.Date, task.Repeat); err != nil {
				responseWithError(w, err)
				return
			} else if task.Date < time.Now().Format(DateFormat) {
				task.Date = nextDate
			}
		}

		if task.Date < time.Now().Format(DateFormat) {
			task.Date = time.Now().Format(DateFormat)
		}
	}

	if result := database.Db.Create(&task); result.Error != nil {
		log.Fatalf("Err: %s", result.Error)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	taskID, _ := json.Marshal(models.ResponseTaskId{Id: task.ID})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(taskID); err != nil {
		responseWithError(w, err)
		return
	}

	log.Printf("Add task id=%d", task.ID)
}

func TasksRead(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	var err error

	if tasks, err = database.ReadTasks(); err != nil {
		responseWithError(w, err)
		return
	}

	tasksData, err := json.Marshal(models.Tasks{Tasks: tasks})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
	log.Println(fmt.Sprintf("Read %d tasks", len(tasks)))

	if err != nil {
		responseWithError(w, err)
	}
}

func TaskReadByID(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var err error
	var id uint64

	if len(r.FormValue("id")) == 0 {
		responseWithError(w, errors.New("id not specified"))
		return
	}

	if id, err = strconv.ParseUint(r.FormValue("id"), 10, 64); err != nil {
		responseWithError(w, err)
		return
	}

	if task, err = database.ReadTaskByID(uint(id)); err != nil {
		responseWithError(w, err)
		return
	}

	tasksData, err := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)

	if err != nil {
		responseWithError(w, err)
		return
	}

	log.Println(fmt.Sprintf("Read task. Id = %d", task.ID))
}

func TaskUpdate(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		responseWithError(w, err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		responseWithError(w, err)
		return
	}

	if _, err := database.ReadTaskByID(task.ID); err != nil {
		log.Fatalf("Err: %s", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	//проверки корректного заполнения полей
	if len(task.Title) == 0 {
		responseWithError(w, errors.New("title is empty"))
		return
	}

	if len(task.Date) == 0 {
		task.Date = time.Now().Format(DateFormat)
	} else {
		_, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			responseWithError(w, errors.New("bad date"))
			return
		}

		if len(task.Repeat) > 0 {
			if nextDate, err := nextDate(time.Now(), task.Date, task.Repeat); err != nil {
				responseWithError(w, err)
				return
			} else if task.Date < time.Now().Format(DateFormat) {
				task.Date = nextDate
			}
		}

		if task.Date < time.Now().Format(DateFormat) {
			task.Date = time.Now().Format(DateFormat)
		}
	}

	if result := database.Db.Updates(&task); result.Error != nil {
		responseWithError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("{}")); err != nil {
		responseWithError(w, err)
		return
	}

	log.Printf("Update task id=%d", task.ID)
}

func TaskDone(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var err error
	var id uint64

	if len(r.FormValue("id")) == 0 {
		responseWithError(w, errors.New("id not specified"))
		return
	}

	if id, err = strconv.ParseUint(r.FormValue("id"), 10, 64); err != nil {
		responseWithError(w, err)
		return
	}

	if task, err = database.ReadTaskByID(uint(id)); err != nil {
		responseWithError(w, err)
		return
	}

	if len(task.Repeat) > 0 {
		if nextDate, err := nextDate(time.Now(), task.Date, task.Repeat); err != nil {
			responseWithError(w, err)
			return
		} else {
			task.Date = nextDate
		}
		if result := database.Db.Save(&task); result.Error != nil {
			responseWithError(w, result.Error)
			return
		}
	} else {
		if result := database.Db.Delete(&task); result.Error != nil {
			responseWithError(w, result.Error)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("{}")); err != nil {
		responseWithError(w, err)
		return
	}

	log.Printf("Task marked as done id=%d", task.ID)
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var err error
	var id uint64

	if len(r.FormValue("id")) == 0 {
		responseWithError(w, errors.New("id not specified"))
		return
	}

	if id, err = strconv.ParseUint(r.FormValue("id"), 10, 64); err != nil {
		responseWithError(w, err)
		return
	}

	if task, err = database.ReadTaskByID(uint(id)); err != nil {
		responseWithError(w, err)
		return
	}

	if result := database.Db.Delete(&task); result.Error != nil {
		responseWithError(w, result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("{}"))
	log.Println(fmt.Sprintf("Delete task. Id = %d", task.ID))

	if err != nil {
		responseWithError(w, err)
	}
}
