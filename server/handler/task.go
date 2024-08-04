package handler

import (
	"encoding/json"
	"go_final_project/storage"
	"go_final_project/task"
	"log"
	"net/http"
	"strconv"
	"time"
)

// setJSONContentType позволяет усатновить заголовки во всех обработчиках запросов
func setJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func AddTask(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest POST AddTask")

		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			log.Println("JSON deserialization error", err)
			http.Error(w, `{"error":"JSON deserialization error"}`, http.StatusBadRequest)
			return
		}

		err = task.Check(&t)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		id, err := storage.Add(&t)
		if err != nil {
			log.Printf("can't add task: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := map[string]string{
			"id": strconv.Itoa(id),
		}

		resp, err := json.Marshal(result)
		if err != nil {
			log.Printf("Can`t marshal id: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		setJSONContentType(w)

		_, err = w.Write(resp)
		if err != nil {
			log.Printf("can't write response: %v", err)
		}
		log.Printf("Task %s id:%v added successfully", t.Title, id)
	}
}

func GetTasks(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest GET Tasks")

		var tasks []task.Task
		var err error

		search := r.URL.Query().Get("search")

		if search == "" {
			tasks, err = db.GetList()
			if err != nil {
				log.Printf("can't get tasks: %v", err)
			}
		}

		if search != "" {
			tasks, err = db.SearchTasks(search)
			if err != nil {
				log.Printf("can't find tasks: %v", err)
			}
		}

		if len(tasks) == 0 {
			tasks = []task.Task{}
		}

		result := map[string][]task.Task{
			"tasks": tasks,
		}

		resp, err := json.Marshal(result)
		if err != nil {
			log.Printf("Can`t marshal tasks: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		setJSONContentType(w)

		_, err = w.Write(resp)
		if err != nil {
			log.Printf("can't write response by GetTasks: %v", err)
		} else {
			log.Printf("GetTasks is successful. %d tasks found", len(tasks))
		}
	}
}

func GetTask(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest GET task")

		var err error

		id := r.URL.Query().Get("id")

		if id == "" {
			log.Println("ID is empty")
			json.NewEncoder(w).Encode(map[string]string{"error": "id is empty"})
			return
		}

		_, err = strconv.Atoi(id)
		if err != nil {
			log.Println("incorrect ID, id is not number")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "incorrect task ID"})
			return
		}

		task, err := storage.GetTask(id)
		if err != nil {
			log.Println("can't get task:", err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			http.Error(w, err.Error(), http.StatusNoContent)

			return
		}

		resp, err := json.Marshal(task)
		if err != nil {
			log.Printf("Can`t marshal task: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		setJSONContentType(w)

		_, err = w.Write(resp)
		if err != nil {
			log.Println("can't write response by GetTask:", err)
		} else {
			log.Println("GetTask is successful")
		}
	}
}

func UpdateTask(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest UpdateTask")

		var t task.Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, `{"error":"Can't read request body"}`, http.StatusBadRequest)
			return
		}

		err = task.Check(&t)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		err = storage.Update(t)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		t, err = storage.GetTask(t.ID)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		setJSONContentType(w)

		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func DoneTask(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest DoneTask")

		id := r.URL.Query().Get("id")

		_, err := strconv.Atoi(id)
		if err != nil {
			log.Println("id is not a number:", err)
			json.NewEncoder(w).Encode(map[string]string{"error": "id is not a number"})
			return
		}

		t, err := storage.GetTask(id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "can't get task " + id})
			return
		}

		if t.Repeat == "" {
			log.Println("Repeat is empty, task will delete")
			err = storage.DeleteTask(id)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": "can't delete task " + id})
				return
			}
		}

		if t.Repeat != "" {
			t.Date, err = task.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
				return
			}

			err = storage.Update(t)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
				return
			}
		}

		setJSONContentType(w)
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(map[string]string{})
		if err != nil {
			log.Println("Can't write response by DoneTask:", err)
		} else {
			log.Println("Done task " + id + " is successful")
		}
	}
}

func DeleteTask(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest DelTask")

		id := r.URL.Query().Get("id")

		_, err := strconv.Atoi(id)
		if err != nil {
			log.Println("id is not a number:", err)
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		err = storage.DeleteTask(id)
		if err != nil {
			log.Println("Failed to delete task")
			json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
			return
		}

		setJSONContentType(w)
		if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			log.Println("err encode:", err)
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}

	}
}
