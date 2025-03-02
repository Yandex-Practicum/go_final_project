package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AlexeyVilkov/go_final_project/database"
	"github.com/AlexeyVilkov/go_final_project/date"
	"github.com/AlexeyVilkov/go_final_project/model"

	_ "github.com/mattn/go-sqlite3"
)

// выносим формат даты в константу
const dateFormat = "20060102"

func JSONError(writer http.ResponseWriter, message string, code int) {
	http.Error(writer, fmt.Sprintf(`{"error":"%s"}`, message), code)
}

func GetNextDate(w http.ResponseWriter, r *http.Request) {
	sNow := r.FormValue("now")
	if sNow == "" {
		http.Error(w, "Параметр now не найден", http.StatusBadRequest)
		return
	}

	sDate := r.FormValue("date")
	if sDate == "" {
		http.Error(w, "Параметр date не найден", http.StatusBadRequest)
		return
	}

	sRepeat := r.FormValue("repeat")
	if sRepeat == "" {
		http.Error(w, "Параметр repeat не найден", http.StatusBadRequest)
		return
	}

	tNow, err := time.Parse(dateFormat, sNow)
	if err != nil {
		http.Error(w, "Ошибка преобразования параметра now в дату", http.StatusUnprocessableEntity)
		return
	}

	// получение даты следующей задачи
	nextDate, err := date.NextDate(tNow, sDate, sRepeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

func ActionTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var taskID model.TaskID
	var buf bytes.Buffer
	var resp []byte
	var err error

	now := time.Now()

	// в зависимости от метода запроса определяем обработчик
	switch r.Method {
	case http.MethodGet:
		// получаем информацию по задаче
		id := r.FormValue("id")

		// Проверяем не пустой ли id
		switch {
		case id == "":
			fmt.Println("Не указан идентификатор")
			JSONError(w, "Не указан идентификатор", http.StatusNotFound)
			return
		default:
			task, err = database.GetTaskByID(id)
			if err != nil {
				fmt.Println("Ошибка получения задачи по id: ", err.Error())
				JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			resp, err = json.Marshal(task)
			if err != nil {
				fmt.Println("Не удалось упаковать ошибку в JSON: ", err.Error())
				JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	case http.MethodPost:
		// создаём новую задачу

		// проверяем, что передан json
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			fmt.Println("Необходимо передать json в запросе")
			JSONError(w, "Необходимо передать json в запросе", http.StatusUnsupportedMediaType)
			return
		}

		// читаем тело запроса
		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			fmt.Println("Ошибка получения тела запроса: ", err.Error())
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// проверяем успешную десериализацию
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			fmt.Println("Ошибка десериализации тела запроса: ", err.Error())
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// если дата не передана, то сегодняшнее число
		if len(task.Date) == 0 {
			task.Date = now.Format(dateFormat)
		}

		// проверяем, что дата корректно передана
		if _, err = time.Parse(dateFormat, task.Date); err != nil {
			fmt.Println("Ошибка преобразования параметра date в дату")
			JSONError(w, "Ошибка преобразования параметра date в дату", http.StatusUnprocessableEntity)
			return
		}

		// если дата меньше сегодняшней, то применяем NextDate
		if task.Date < now.Format(dateFormat) {
			if task.Repeat == "" {
				// если правила повтора нет, ставим сегодняшнюю дату
				task.Date = now.Format(dateFormat)
			} else {
				// получаем следующую дату
				task.Date, err = date.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					fmt.Println("Ошибка при получении следующей даты: ", err.Error())
					JSONError(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}

		// проверяем, что заголовок задачи обязательно передан
		if task.Title == "" {
			fmt.Println("Заголовок задачи title обязателен, но не передан")
			JSONError(w, "Заголовок задачи title обязателен, но не передан", http.StatusUnprocessableEntity)
			return
		}

		id, err := database.PostTask(task)
		if err != nil {
			fmt.Println("Ошибка сохранения задачи: ", err.Error())
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if id > 0 {
			taskID.Id = strconv.Itoa(id)
			resp, err = json.Marshal(taskID)
			if err != nil {
				fmt.Println("Не удалось упаковать id в JSON: ", err.Error())
				JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)

	case http.MethodPut:
		// редактируем задачу

		// проверяем, что передан json
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			fmt.Println("Необходимо передать json в запросе")
			JSONError(w, "Необходимо передать json в запросе", http.StatusUnsupportedMediaType)
			return
		}

		// читаем тело запроса
		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			fmt.Println("Ошибка получения тела запроса: ", err.Error())
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// проверяем успешную десериализацию
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			fmt.Println("Ошибка десериализации тела запроса: ", err.Error())
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// проверяем, что id передан
		if task.Id == "" {
			fmt.Println("Идентификатор задачи id обязателен, но не передан")
			JSONError(w, "Идентификатор задачи id обязателен, но не передан", http.StatusUnprocessableEntity)
			return
		}

		// если дата не передана, то сегодняшнее число
		if len(task.Date) == 0 {
			task.Date = now.Format(dateFormat)
		}

		// проверяем, что дата корректно передана
		if _, err = time.Parse(dateFormat, task.Date); err != nil {
			fmt.Println("Ошибка преобразования параметра date в дату")
			JSONError(w, "Ошибка преобразования параметра date в дату", http.StatusUnprocessableEntity)
			return
		}

		// если дата меньше сегодняшней, то применяем NextDate
		if task.Date < now.Format(dateFormat) {
			if task.Repeat == "" {
				// если правила повтора нет, ставим сегодняшнюю дату
				task.Date = now.Format(dateFormat)
			} else {
				// получаем следующую дату
				task.Date, err = date.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					fmt.Println("Ошибка при получении следующей даты: ", err.Error())
					JSONError(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}

		// проверяем, что заголовок задачи обязательно передан
		if task.Title == "" {
			fmt.Println("Заголовок задачи title обязателен, но не передан")
			JSONError(w, "Заголовок задачи title обязателен, но не передан", http.StatusUnprocessableEntity)
			return
		}

		err = database.UpdateTask(task)
		if err != nil {
			fmt.Println("Ошибка редактирования задачи: ", err.Error())
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp = []byte(`{}`)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	case http.MethodDelete:
		// удаляем задачу

		// получаем id задачи
		id := r.FormValue("id")

		// проверяем не пустой ли id
		if id == "" {
			fmt.Println("Не указан идентификатор")
			JSONError(w, "Не указан идентификатор", http.StatusNotFound)
			return
		}

		err = database.DeleteTask(id)
		if err != nil {
			fmt.Println("Ошибка удаления задачи: ", err.Error())
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp = []byte(`{}`)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	default:
		http.Error(w, "Метод запроса не определен", http.StatusMethodNotAllowed)
		return
	}

}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	var resp []byte

	result, err := database.ListTasks(100)
	if err != nil {
		fmt.Println("Ошибка при получении списка задач: ", err.Error())
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err = json.Marshal(result)
	if err != nil {
		fmt.Println("Не удалось упаковать результат в JSON: ", err.Error())
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	// выполнение задачи
	var resp []byte
	var err error

	// получаем id задачи
	id := r.FormValue("id")

	// проверяем не пустой ли id
	if id == "" {
		fmt.Println("Не указан идентификатор")
		JSONError(w, "Не указан идентификатор", http.StatusNotFound)
		return
	}

	err = database.DoneTask(id)
	if err != nil {
		fmt.Println("Ошибка выполнения задачи: ", err.Error())
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp = []byte(`{}`)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
