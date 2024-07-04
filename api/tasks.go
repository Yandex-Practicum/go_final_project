package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Masteker/go_final_project/database/db"
)

// GetTasksHandler обрабатывает запросы к /api/tasks с методом GET.
// Если пользователь авторизован, возвращает JSON {"tasks": Task} содержащий последние добавленные задачи, или
// последние добавленные задачи соответствующие поисковому запросу search. В случае ошибки возвращает JSON {"error": error}.
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []db.Task
	var err error
	var date time.Time

	// write отправляет клиенту ответ либо ошибку, в формате json
	write := func() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var resp []byte
		if err != nil {
			writeErr(err, w)
			return
		} else {
			if len(tasks) == 0 {
				tasksResp := map[string][]db.Task{
					"tasks": {},
				}
				resp, err = json.Marshal(tasksResp)
			} else {
				tasksResp := map[string][]db.Task{
					"tasks": tasks,
				}
				resp, err = json.Marshal(tasksResp)

			}

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

	// Проверяем есть ли поисковой зарпос
	q := r.URL.Query()
	search := q.Get("search")
	// Проверяем может ли поисковой запрос содержать поиск по дате
	isDate, _ := regexp.Match("[0-9]{2}.[0-9]{2}.[0-9]{4}", []byte(search))

	switch {
	case len(search) == 0:
		tasks, err = dbh.GetTasksList()

	case isDate:
		date, err = time.Parse("02.01.2006", search)
		if err == nil {
			search = date.Format(dateFormat)
			tasks, err = dbh.GetTasksList(search)
			break
		}
		fallthrough

	default:
		search = fmt.Sprint("%" + search + "%")
		tasks, err = dbh.GetTasksList(search)

	}

	if err != nil {
		log.Println(err)
	}

	write()

}
