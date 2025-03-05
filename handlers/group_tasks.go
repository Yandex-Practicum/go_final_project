package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go_final_project/sqlite"
)

// GetTasksHandler - возвращает группу задач
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	var errRes sqlite.StringError
	var result []byte
	var err error

	search := r.FormValue("search")

	result, err = sqlite.TodoStorage.FindTasks(search) // FindTasks(search)
	if err != nil {
		fmt.Println("Ошибка чтения из БД ", err)
		errRes.StrEr = fmt.Sprint(err)
		result, err = json.Marshal(errRes)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(result)
	if err != nil {
		fmt.Println("Ошибка формирования ответа. ", err)
	}
}
