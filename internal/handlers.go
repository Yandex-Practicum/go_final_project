package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	database "github.com/sandrinasava/go_final_project/internal/db"
	"github.com/sandrinasava/go_final_project/internal/models"
	"github.com/sandrinasava/go_final_project/internal/services"
	_ "modernc.org/sqlite"
)

const limit = "15"

func sendErrorResponse(res http.ResponseWriter, message string, statusCode int) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(map[string]string{"error": message})
}

// обработчик для NextDate
func NextDateHandle(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Ожидается Get запрос", http.StatusMethodNotAllowed)
		return
	}

	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	// вызов функции NextDate
	nextDate, err := services.NextDate(now, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// Успешный ответ
	fmt.Fprintf(res, nextDate)
}

// oбработчик  для api/tasks
func TasksHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			sendErrorResponse(res, "неподходящий метод запроса", http.StatusMethodNotAllowed)
			return
		}

		// Инициализирую слайс как пустой слайс
		tasksSlice := []models.Task{}

		search := req.FormValue("search")
		if search != "" {
			log.Printf("search = %s", search)
			D, err := time.Parse("02.01.2006", search)
			if err != nil {
				log.Printf("парсинг даты не удался")
				//если это не дата, ищу соответствие в столбцах title и comment
				selectTask := `
	              SELECT * FROM scheduler
                  WHERE title LIKE '%' || ? || '%'
                  OR comment LIKE '%' || ? || '%'
                  ORDER BY date ASC LIMIT ?;`

				tasksSlice, err = database.FindTasks(db, selectTask, search, search, limit)
				if err != nil {
					sendErrorResponse(res, err.Error(), http.StatusBadRequest)
					return
				}
				log.Printf("tasksSlice = %+v", tasksSlice)
			} else {
				// ищу по дате
				Dstr := D.Format(services.Format)

				log.Printf("searchDate = %s", Dstr)
				selectTask := `
	            SELECT * FROM scheduler WHERE date LIKE ? ORDER BY date ASC LIMIT ?`
				tasksSlice, err = database.FindTasks(db, selectTask, Dstr, limit)
				if err != nil {
					sendErrorResponse(res, err.Error(), http.StatusBadRequest)
					return
				}
			}

			// если параметра search нет, ищу все ближайшие задачи
		} else {
			selectTask := `
	         SELECT * FROM scheduler ORDER BY date ASC LIMIT ?`
			var err error
			tasksSlice, err = database.FindTasks(db, selectTask, limit)
			if err != nil {
				sendErrorResponse(res, err.Error(), http.StatusBadRequest)
				return
			}
		}
		// Структура для ответа
		type TasksResponse struct {
			Tasks []models.Task `json:"tasks"`
		}

		response := TasksResponse{Tasks: tasksSlice}

		err := json.NewEncoder(res).Encode(response)
		if err != nil {
			log.Printf("Ошибка кодирования в JSON")
			sendErrorResponse(res, "Ошибка кодирования в JSON", http.StatusBadRequest)
			return
		}
		return
	}
}

// oбработчик  для api/task
func TaskHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Получен запрос: %s %s", req.Method, req.URL.Path)

		switch req.Method {

		case http.MethodDelete:
			id := req.FormValue("id")

			// проверяю сущ-е id
			var exists int
			err := db.QueryRow("SELECT 1 FROM scheduler WHERE id = ?", id).Scan(&exists)
			if err != nil {
				sendErrorResponse(res, "записи с указанным id нет", http.StatusBadRequest)
				return
			}
			//удаляю задачу
			_, err = db.Exec("DELETE FROM scheduler WHERE id = ?;", id)
			if err != nil {
				sendErrorResponse(res, "неудачный DELETE запрос", http.StatusBadRequest)
				return
			}
			// если все успешно, отправляю поустой json
			res.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(res).Encode(map[string]string{})
			return

		case http.MethodPost:

			if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
				sendErrorResponse(res, "запрос не содержит json", http.StatusUnsupportedMediaType)
				return
			}

			var task models.Task

			// декод-ю тело запроса
			err := json.NewDecoder(req.Body).Decode(&task)
			if err != nil {
				sendErrorResponse(res, "неудачное декодир-е json", http.StatusBadRequest)
				return
			}

			// ищу новую дату
			date, err := CheckTaskAndFindDate(task)
			if err != nil {
				sendErrorResponse(res, err.Error(), http.StatusBadRequest)
				return
			}
			//добавляю задачу в бд
			ID, err := database.InsertAndReturnID(db, date, task.Title, task.Comment, task.Repeat)
			if err != nil {
				log.Printf("ошибка при добавлении задачи")
				sendErrorResponse(res, "ошибка при добавлении задачи", http.StatusBadRequest)
				return
			}
			res.Header().Set("Content-Type", "application/json; charset=UTF-8")
			response := map[string]interface{}{"id": ID}
			err = json.NewEncoder(res).Encode(response)
			if err != nil {
				sendErrorResponse(res, "Ошибка кодирования в JSON", http.StatusBadRequest)
				return
			}

			return

		case http.MethodGet:
			id := req.FormValue("id")
			if id != "" {
				selectTask := `
		        SELECT * FROM scheduler WHERE id = ?`
				t := models.Task{}
				err := db.QueryRow(selectTask, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
				if err != nil {
					sendErrorResponse(res, "неуспешный select запрос", http.StatusBadRequest)
					return
				}
				err = json.NewEncoder(res).Encode(t)
				if err != nil {
					log.Printf("Ошибка кодирования в JSON")
					sendErrorResponse(res, "Ошибка кодирования в JSON", http.StatusBadRequest)
					return

				}
			} else {
				sendErrorResponse(res, "недостаточно параметров", http.StatusMethodNotAllowed)
				return
			}

		case http.MethodPut:
			var task models.Task

			contentType := req.Header.Get("Content-Type")
			log.Println("Content-Type:", contentType)
			if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
				log.Printf("запрос не содержит json")
				sendErrorResponse(res, "запрос не содержит json", http.StatusUnsupportedMediaType)
				return
			}

			err := json.NewDecoder(req.Body).Decode(&task)
			if err != nil {
				log.Printf("неудачное декодир-е json")
				sendErrorResponse(res, "неудачное декодир-е json", http.StatusBadRequest)
				return
			}

			// ищу новую дату
			date, err := CheckTaskAndFindDate(task)
			if err != nil {
				sendErrorResponse(res, err.Error(), http.StatusBadRequest)
				return
			}

			// проверяю сущ-е id
			var exists int
			err = db.QueryRow("SELECT 1 FROM scheduler WHERE id = ?;", task.ID).Scan(&exists)
			if err != nil {
				sendErrorResponse(res, "записи с указанным id нет", http.StatusBadRequest)
				return
			}
			// изменяю данные в дб
			selectTask := `
	              UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 where id = $5;
                  `
			_, err = db.Exec(selectTask, date, task.Title, task.Comment, task.Repeat, task.ID)
			if err != nil {
				sendErrorResponse(res, "неудачный update запрос", http.StatusBadRequest)
				return
			}
			// если все успешно, отправляю поустой json
			res.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(res).Encode(map[string]string{})
			return

		default:
			sendErrorResponse(res, "неподходящий метод запроса", http.StatusMethodNotAllowed)
			return

		}

	}
}

// oбработчик  для api/task/done
func TaskDoneHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		id := req.FormValue("id")

		// проверяю сущ-е id
		var exists int
		err := db.QueryRow("SELECT 1 FROM scheduler WHERE id = ?;", id).Scan(&exists)
		if err != nil {
			sendErrorResponse(res, "записи с указанным id нет", http.StatusBadRequest)
			return
		}

		var t models.Task
		//нахожу задачу
		err = db.QueryRow("SELECT date, repeat FROM scheduler WHERE id = ?;", id).Scan(&t.Date, &t.Repeat)
		if err != nil {
			sendErrorResponse(res, "неуспешный select запрос", http.StatusBadRequest)
			return
		}
		// если repeat пустой - удаляю задачу
		if t.Repeat == "" {
			_, err := db.Exec("DELETE FROM scheduler WHERE id = ?;", id)
			if err != nil {
				sendErrorResponse(res, "неудачный DELETE запрос", http.StatusBadRequest)
				return
			}
			// если все успешно, отправляю поустой json
			res.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(res).Encode(map[string]string{})
			return
		}
		// если repeat есть, ищу новую дату(так как задача уже была в бд, проверок на валидность не делаю)
		nowTime := time.Now()
		now := nowTime.Format(services.Format)
		date, err := services.NextDate(now, t.Date, t.Repeat)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		// выполняю UPDATE запрос
		updateTask := `
			        UPDATE scheduler SET date = ?  where id = ?;
			       `
		_, err = db.Exec(updateTask, date, id)
		if err != nil {
			sendErrorResponse(res, "неудачный UPDATE запрос", http.StatusBadRequest)
			return
		}
		// если все успешно, отправляю поустой json
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(res).Encode(map[string]string{})
		return

	}
}
