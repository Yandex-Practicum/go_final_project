package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Memonagi/go_final_project/constants"
	"github.com/Memonagi/go_final_project/date"
	"github.com/Memonagi/go_final_project/task"
)

// GetNextDate GET-обработчик для получения следующей даты
func GetNextDate(w http.ResponseWriter, r *http.Request) {
	nowReq := r.FormValue("now")
	dateReq := r.FormValue("date")
	repeatReq := r.FormValue("repeat")

	now, err := time.Parse(constants.DateFormat, nowReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nextDate, err := date.NextDate(now, dateReq, repeatReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		http.Error(w, fmt.Errorf("writing tasks data error: %w", err).Error(), http.StatusBadRequest)
	}
}

// PostAddTask POST-обработчик для добавления новой задачи
func PostAddTask(w http.ResponseWriter, r *http.Request) {
	var taskStruct constants.Task
	var err error
	var dateOfTask string

	err = json.NewDecoder(r.Body).Decode(&taskStruct)
	if err != nil {
		http.Error(w, "ошибка десериализации JSON", http.StatusInternalServerError)
		return
	}

	titleOfTask, err := task.CheckTitle(taskStruct.Title)
	if err != nil {
		response := constants.ErrorResponse{Error: err.Error()}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	taskStruct.Title = titleOfTask

	now := time.Now()
	if taskStruct.Repeat == "" {
		dateOfTask, err = task.CheckDate(taskStruct.Date)
		if err != nil {
			response := constants.ErrorResponse{Error: err.Error()}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		taskStruct.Date = dateOfTask
	} else {
		err = task.CheckRepeat(taskStruct.Repeat)
		if err != nil {
			response := constants.ErrorResponse{Error: err.Error()}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		dateOfTask, err = task.CheckDate(taskStruct.Date)
		if err != nil {
			response := constants.ErrorResponse{Error: err.Error()}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if dateOfTask == now.Format(constants.DateFormat) {
			taskStruct.Date = dateOfTask
		} else {
			nextDate, err := date.NextDate(now, dateOfTask, taskStruct.Repeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			taskStruct.Date = nextDate
		}
	}

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		http.Error(w, "не удается подключиться к базе данных", http.StatusBadRequest)
		return
	}
	defer db.Close()

	taskId, err := task.AddTask(db, taskStruct)
	if err != nil {
		http.Error(w, "не удалось добавить задачу", http.StatusBadRequest)
		return
	}

	response := constants.TaskIdResponse{Id: taskId}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "ошибка сериализации JSON", http.StatusInternalServerError)
	}
}

// GetAddTasks GET-обработчик для получения списка ближайших задач
func GetAddTasks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		response := constants.ErrorResponse{Error: "не удается подключиться к базе данных"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", constants.Limit)
	if err != nil {
		response := constants.ErrorResponse{Error: "не удается получить информацию из базы данных"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	tasks := []constants.Task{}

	for rows.Next() {

		var taskStruct constants.Task

		if err = rows.Scan(&taskStruct.ID, &taskStruct.Date, &taskStruct.Title, &taskStruct.Comment, &taskStruct.Repeat); err != nil {
			response := constants.ErrorResponse{Error: "не удается обработать информацию из базы данных"}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		tasks = append(tasks, taskStruct)
	}
	if err := rows.Err(); err != nil {
		response := constants.ErrorResponse{Error: "не удается обработать информацию из базы данных"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if tasks == nil {
		tasks = []constants.Task{}
	}

	response := constants.TaskResponse{Tasks: tasks}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		response := constants.ErrorResponse{Error: "ошибка сериализации JSON"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// GetTaskId GET-обработчик для получения задачи по ее id
func GetTaskId(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		response := constants.ErrorResponse{Error: "не удается подключиться к базе данных"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	id := r.URL.Query().Get("id")

	if id == "" {
		response := constants.ErrorResponse{Error: "задача не найдена"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var taskStruct constants.Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	if err = db.QueryRow(query, id).Scan(&taskStruct.ID, &taskStruct.Date, &taskStruct.Title, &taskStruct.Comment, &taskStruct.Repeat); err != nil {
		response := constants.ErrorResponse{Error: "не удается получить информацию из базы данных"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(taskStruct)
	if err != nil {
		response := constants.ErrorResponse{Error: "ошибка сериализации JSON"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// UpdateTaskId PUT-обработчик для редактирования задачи
func UpdateTaskId(w http.ResponseWriter, r *http.Request) {
	var taskStruct constants.Task

	if err := json.NewDecoder(r.Body).Decode(&taskStruct); err != nil {
		response := constants.ErrorResponse{Error: "ошибка десериализации JSON"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		response := constants.ErrorResponse{Error: "не удается подключиться к базе данных"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	if taskStruct.Title == "" {
		response := constants.ErrorResponse{Error: "заголовок не может быть пустым"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	now := time.Now()
	if taskStruct.Date == "" {
		taskStruct.Date = now.Format(constants.DateFormat)
	} else {
		dateOfTask, err := time.Parse(constants.DateFormat, taskStruct.Date)
		if err != nil {
			response := constants.ErrorResponse{Error: "неправильный формат даты"}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if dateOfTask.Before(now) {
			if taskStruct.Repeat == "" {
				taskStruct.Date = now.Format(constants.DateFormat)
			} else {
				nextDate, err := date.NextDate(now, taskStruct.Date, taskStruct.Repeat)
				if err != nil {
					response := constants.ErrorResponse{Error: "неправильный формат даты"}
					if err := json.NewEncoder(w).Encode(response); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				}
				taskStruct.Date = nextDate
			}
		}
	}

	if err := task.CheckRepeat(taskStruct.Repeat); err != nil {
		response := constants.ErrorResponse{Error: "неправильный формат правила повторения"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	row, err := db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", taskStruct.Date, taskStruct.Title, taskStruct.Comment, taskStruct.Repeat, taskStruct.ID)
	if err != nil {
		response := constants.ErrorResponse{Error: "не удалось обновить информацию"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	checkRow, err := row.RowsAffected()
	if err != nil || checkRow == 0 {
		response := constants.ErrorResponse{Error: "не удалось обновить информацию"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(taskStruct); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TaskDone POST-обработчик для выполнения задачи
func TaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		response := constants.ErrorResponse{Error: "не указан ID"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		response := constants.ErrorResponse{Error: "не удается подключиться к базе данных"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	var taskStruct constants.Task

	if err := db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(&taskStruct.ID, &taskStruct.Date, &taskStruct.Title, &taskStruct.Comment, &taskStruct.Repeat); err != nil {
		response := constants.ErrorResponse{Error: "задача не найдена"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	switch taskStruct.Repeat {
	case "":
		_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			response := constants.ErrorResponse{Error: "не удалось удалить задачу"}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	default:
		now := time.Now()
		nextDate, err := date.NextDate(now, taskStruct.Date, taskStruct.Repeat)
		if err != nil {
			response := constants.ErrorResponse{Error: "неправильный формат даты"}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		if err != nil {
			response := constants.ErrorResponse{Error: "не удалось перенести задачу на другую дату"}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	response := struct{}{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteTask DELETE-обработчик для удаления задачи
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		response := constants.ErrorResponse{Error: "не указан ID"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		response := constants.ErrorResponse{Error: "не удается подключиться к базе данных"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	row, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		response := constants.ErrorResponse{Error: "не удалось удалить задачу"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	checkRow, err := row.RowsAffected()
	if err != nil || checkRow == 0 {
		response := constants.ErrorResponse{Error: "не удалось удалить задачу"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		response := constants.ErrorResponse{Error: "ошибка сериализации JSON"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
