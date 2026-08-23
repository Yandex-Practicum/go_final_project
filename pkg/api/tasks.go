package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/UTKIN-AI/go_final_project/pkg/db"
	"github.com/golang-jwt/jwt/v5"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

type SignInRequest struct {
	Password  string  `json:"password"`
}

// TaskHandler обрабатывает GET запрос /api/tasks
// и возвращает список задач с возможностью поиска.
func tasksHandler(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = writeJSON(w, TasksResp{
		Tasks: tasks,
	})
	if err != nil {
		return 
	}
}

// getTaskHandler обрабатывает GET запрос /api/task
// и возвращает задачу по идентификатору.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, map[string]string {
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "Задача не найдена",
		})
		return
	}

	err = writeJSON(w, task)
	if err != nil {
		return
	}
}

// updateTaskHandler обрабатывает PUT запрос /api/task
// и обновляет существующую задачу.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{
			"error": "Не указано название задачи",
		})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = writeJSON(w, map[string]string{})
	if err != nil {
		return 
	}

}

// deleteTaskHandler обрабатывает DELETE запрос /api/task
// и удаляет задачу по идентификатору.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, map[string]string{
			"error": "не указан идентификатор",
		})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = writeJSON(w, map[string]string{})
	if err != nil {
		return
	}
}

// doneTaskHandler обрабатывает POST запрос /api/task/done
// и отмечает задачу выполненной.
// Для повторяющихся задач устанавливается следующая дата,
// для одноразовой задачи, запись удаляется.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, map[string]string{
			"error": "не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "задача не найдена",
		})
		return 
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error(),
			})
			return
		}
	} else {

		now := time.Now()

		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error(),	
			})
			return
		}
		err = db.UpdateDate(next, id)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error(),
			})
			return
		}
	} 
	err = writeJSON(w, map[string]string{})
	if err != nil {
		return 
	}
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	req := SignInRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	if pass != req.Password {
		writeJSON(w, map[string]string{
			"error": "неверный пароль",
		})
		return
	}

	hash := sha256.Sum256([]byte(pass))
	hashString := hex.EncodeToString(hash[:])

	claims := jwt.MapClaims {
		"hash": hashString,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	secret := getJWTSecret()
	if secret == "" {
		writeJSON(w, map[string]string{
			"error": "JWT_SECRET is not set",
		})
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = writeJSON(w, map[string]string{
		"token": tokenString,
	})
	if err != nil {
		return 
	}
}