package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ikamenev/database"
	"github.com/Ikamenev/model"
	"github.com/Ikamenev/service"

	"log"
	"net/http"
	"strconv"
	"time"
)

func NextDate(w http.ResponseWriter, r *http.Request) {

	now, err := time.Parse(model.DatePattern, r.FormValue("now"))
	if err != nil {
		http.Error(w, fmt.Sprintf(""), http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nextDate, err := service.NextDate(now, date, repeat)

	if err != nil {
		http.Error(w, fmt.Sprintf(""), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))

	if err != nil {
		http.Error(w, fmt.Errorf("writing tasks data error: %w", err).Error(), http.StatusBadRequest)
	}
}

func responseWithError(w http.ResponseWriter, errorText string, err error) {
	errorResponse := model.ErrorResponse{
		Error: fmt.Errorf("%s: %w", errorText, err).Error()}
	errorData, _ := json.Marshal(errorResponse)
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(errorData)

	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
	}
}

func TaskAddPOST(w http.ResponseWriter, r *http.Request) {
	var taskData model.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "body getting error", err)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &taskData); err != nil {
		responseWithError(w, "JSON encoding error", err)
		return
	}

	if len(taskData.Date) == 0 {
		taskData.Date = time.Now().Format(model.DatePattern)
	} else {
		date, err := time.Parse(model.DatePattern, taskData.Date)
		if err != nil {
			responseWithError(w, "bad data format", err)
			return
		}

		if date.Before(time.Now()) {
			taskData.Date = time.Now().Format(model.DatePattern)
		}
	}

	if len(taskData.Title) == 0 {
		responseWithError(w, "invalid title", errors.New("title is empty"))
		return
	}

	if len(taskData.Repeat) > 0 {
		if _, err := service.NextDate(time.Now(), taskData.Date, taskData.Repeat); err != nil {
			responseWithError(w, "invalid repeat format", errors.New("no such format"))
			return
		}
	}

	taskId, err := database.InsertTask(taskData)
	if err != nil {
		responseWithError(w, "failed to create task", err)
		return
	}

	taskIdData, err := json.Marshal(model.TaskIdResponse{Id: taskId})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(taskIdData)
	log.Println(fmt.Sprintf("Added task with id=%d", taskId))

	if err != nil {
		responseWithError(w, "writing task id error", err)
	}
}

func TasksReadGET(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var tasks []model.Task

	if len(search) > 0 {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			tasks, err = database.SearchTasks(search)
		} else {
			tasks, err = database.SearchTasksByDate(date.Format(model.DatePattern))
		}
	} else {
		err := errors.New("")
		tasks, err = database.ReadTasks()
		if err != nil {
			responseWithError(w, "failed to get tasks", err)
			return
		}
	}

	tasksData, err := json.Marshal(model.Tasks{Tasks: tasks})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
	log.Println(fmt.Sprintf("Read %d tasks", len(tasks)))

	if err != nil {
		responseWithError(w, "writing tasks error", err)
	}
}

func TaskReadGET(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := database.ReadTask(id)
	if err != nil {
		responseWithError(w, "failed to get task", err)
		return
	}

	tasksData, err := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
	log.Println(fmt.Sprintf("Read task with id=%s", id))

	if err != nil {
		responseWithError(w, "writing task error", err)
	}
}

func TaskUpdatePUT(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "body getting error", err)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &task); err != nil {
		responseWithError(w, "JSON encoding error", err)
		return
	}

	if len(task.Id) == 0 {
		responseWithError(w, "invalid id", errors.New("id is empty"))
		return
	}

	if _, err := strconv.Atoi(task.Id); err != nil {
		responseWithError(w, "invalid id", err)
		return
	}

	if _, err := time.Parse(model.DatePattern, task.Date); err != nil {
		responseWithError(w, "invalid date", err)
		return
	}

	if len(task.Title) == 0 {
		responseWithError(w, "invalid title", errors.New("title is empty"))
		return
	}

	if len(task.Repeat) > 0 {
		if _, err := service.NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			responseWithError(w, "invalid repeat format", errors.New("no such format"))
			return
		}
	}

	_, err := database.UpdateTask(task)
	if err != nil {
		responseWithError(w, "invalid title", errors.New("failed to update task"))
		return
	}

	taskIdData, err := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskIdData)
	log.Println(fmt.Sprintf("Updated task with id=%s", task.Id))

	if err != nil {
		responseWithError(w, "updating task error", err)
		return
	}
}

func TaskDonePOST(w http.ResponseWriter, r *http.Request) {
	task, err := database.ReadTask(r.URL.Query().Get("id"))
	if err != nil {
		responseWithError(w, "failed to get task", err)
		return
	}

	if len(task.Repeat) == 0 {
		err = database.DeleteTaskDb(task.Id)
		if err != nil {
			responseWithError(w, "failed to delete task", err)
			return
		}
	} else {
		task.Date, err = service.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			responseWithError(w, "failed to get next date", err)
			return
		}

		_, err = database.UpdateTask(task)
		if err != nil {
			responseWithError(w, "failed to update task", err)
			return
		}
	}

	tasksData, err := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
	log.Println(fmt.Sprintf("Done task with id=%s", task.Id))

	if err != nil {
		responseWithError(w, "writing task error", err)
	}
}

func TaskDELETE(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := database.DeleteTaskDb(id)
	if err != nil {
		responseWithError(w, "failed to delete task", err)
		return
	}

	tasksData, err := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
	log.Println(fmt.Sprintf("Deleted task with id=%s", id))

	if err != nil {
		responseWithError(w, "writing task error", err)
		return
	}
}
