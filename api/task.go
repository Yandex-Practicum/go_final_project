package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	db "github.com/AsyaBiryukova/go_final_project/database"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodGet:
		getTask(w, r)
	case http.MethodPost:
		postTask(w, r)
	case http.MethodPut:
		putTask(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	}
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	var err error
	var id int64

	write := func() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			writeErr(err, w)
			return
		} else {
			idResp := map[string]string{
				"id": strconv.Itoa(int(id)),
			}
			resp, err := json.Marshal(idResp)
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(resp)
			if err != nil {
				log.Println(err)
			}
			return
		}

	}

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		write()
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		write()
		return
	}

	task, err = formatTask(task)
	if err != nil {
		write()
		return
	}

	id, err = db.AddTask(task)
	write()
}

// PutTaskHandler обрабатывает запрос с методом PUT.
// Если пользователь авторизован и задача существует, и отправлена в корректном формате, обновляет поля задачи в базе данных.
// Возвращает пустой JSON {} или JSON {"error": error} в случае ошибки.
func putTask(w http.ResponseWriter, r *http.Request) {
	var updatedTask db.Task
	var buf bytes.Buffer
	var err error

	write := func() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			writeErr(err, w)
			return
		} else {
			writeEmptyJson(w)
			return
		}

	}

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		write()
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &updatedTask); err != nil {
		write()
		return
	}

	updatedTask, err = formatTask(updatedTask)
	if err != nil {
		write()
		return
	}

	err = db.PutTask(updatedTask)
	write()

}

// GetTaskHandler обрабатывает запрос с методом GET.
// Если пользователь авторизован, возвращает задачу с указанным ID.
// Возвращает JSON {"task":Task}, или JSON {"error": error} при ошибке.
func getTask(w http.ResponseWriter, r *http.Request) {
	var err error
	var task db.Task

	write := func() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var resp []byte
		if err != nil {
			writeErr(err, w)
			return
		} else {
			resp, err = json.Marshal(task)
		}

		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		if err != nil {
			log.Println(err)
		}

	}

	q := r.URL.Query()
	id := q.Get("id")

	task, err = db.GetTaskByID(id)
	if err != nil {
		log.Println(err)
	}
	write()

}

// DeleteTaskHandler обрабатывает запрос к api/task с методом DELETE.
// Если пользователь авторизован и id существует, удаляет задачу.
// При успешном выполнение возвращает пустой JSON {}. Иначе возвращает JSON {"error":error}.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	var err error

	q := r.URL.Query()
	id := q.Get("id")
	isID := isID(id)
	if !isID {
		writeErr(fmt.Errorf("некорректный формат id"), w)
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		writeErr(err, w)
		return
	}
	writeEmptyJson(w)

}
