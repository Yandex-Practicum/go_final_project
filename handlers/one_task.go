package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go_final_project/sqlite"
	"go_final_project/tasks"
)

// GetOneTaskHandler - возвращает одну задачу по выданному признаку (по id)
func GetOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var errRes sqlite.StringError
	var err error
	var result []byte

	id := r.FormValue("id")

	// Проверяем не пустой ли id
	switch {
	case id == "":
		errRes.StrEr = "Не указан идентификатор"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	default:
		// Проверяем id
		data, err := tasks.IDValidation(id)
		if err != nil {
			fmt.Println("Ошибка конвертации входящего значения api/tasks. ", err)
		}

		// Идём искать одну задачу по входным данным
		result, err = sqlite.TodoStorage.GetTask(data)
		if err != nil {
			fmt.Println("Ошибка чтения из БД ", err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(result)
	if err != nil {
		fmt.Println("Ошибка формирования ответа. ", err)
	}
}

// PostOneTaskHandler -  записывает в базу одну задачу
func PostOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	var errRes sqlite.StringError
	var result []byte

	// Читаем сообщение
	httpData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Не прочитано тело запроса api/tasks. ", err)
	}

	// Проверяем полученную информацию о задаче
	data, err := tasks.TaskDataValidation(httpData)

	switch {
	// Если не смогли валедировать входящие данные
	case err != nil:
		fmt.Println("Ошибка конвертации входящего значения api/tasks. ", err)
		errRes.StrEr = fmt.Sprint(err)
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	// Если заголовок пуст возвращаем ошибку
	case data.Title == "":
		errRes.StrEr = "Не указан заголовок задачи"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	// Идём записывать задачу в базу
	default:
		result, err = sqlite.TodoStorage.AddTask(data)
		if err != nil {
			fmt.Println("Ошибка записи в БД ", err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(result)
	if err != nil {
		fmt.Println("Ошибка формирования ответа. ", err)
	}
}

// PutOneTaskHandler - изменяет в базе одну задачу
func PutOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var errRes sqlite.StringError
	var result []byte

	// Читаем сообщение
	httpData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Не прочитано тело запроса api/tasks. ", err)
	}

	// Проверяем полученную информацию о задаче
	data, err := tasks.TaskDataValidation(httpData)

	switch {
	case err != nil:
		fmt.Println("Ошибка конвертации входящего значения api/tasks. ", err)
		errRes.StrEr = "Задача не найдена"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
		// Идём записывать задачу в базу
	default:
		result, err = sqlite.TodoStorage.UpdateTask(data)
		if err != nil {
			fmt.Println("Ошибка записи в БД ", err)
			errRes.StrEr = "Задача не найдена"
			result, err = json.Marshal(errRes)
			if err != nil {
				fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(result)
	if err != nil {
		fmt.Println("Ошибка формирования ответа. ", err)
	}
}

// DoneOneTaskHandler - завершает в базе одну задачу
func DoneOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var errRes sqlite.StringError
	var err error
	var result []byte

	id := r.FormValue("id")

	// Проверяем не пустой ли id
	switch {
	case id == "":
		errRes.StrEr = "Не указан идентификатор"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	default:
		// Проверяем id
		data, err := tasks.IDValidation(id)
		if err != nil {
			fmt.Println("Ошибка конвертации входящего значения api/tasks. ", err)
		}

		// Идём закрывать одну задачу по входным данным
		result, err = sqlite.TodoStorage.DoneTasks(data)
		if err != nil {
			fmt.Println("Ошибка чтения из БД ", err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(result)
	if err != nil {
		fmt.Println("Ошибка формирования ответа. ", err)
	}
}

// DeleteOneTaskHandler - удаляет из базы одну задачу по выданному признаку (по id)
func DeleteOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var errRes sqlite.StringError
	var err error
	var result []byte

	id := r.FormValue("id")

	// Проверяем не пустой ли id
	switch {
	case id == "":
		errRes.StrEr = "Не указан идентификатор"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	default:
		// Проверяем id
		data, err := tasks.IDValidation(id)

		if err != nil {
			fmt.Println("Ошибка конвертации входящего значения api/tasks. ", err)
			errRes.StrEr = "Не верно указан идентификатор"
			result, err = json.Marshal(errRes)
		} else {
			// Идём искать одну задачу по входным данным
			result, err = sqlite.TodoStorage.DeleteTask(data)
			if err != nil {
				fmt.Println("Ошибка чтения из БД ", err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(result)
	if err != nil {
		fmt.Println("Ошибка формирования ответа. ", err)
	}
}
