package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	dbase "go1f/pkg/db" 
)

type TasksResp struct {
	Tasks []*dbase.Task `json:"tasks"`
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	type errorjson struct {
		Error string `json:"error"`
	}

	tasks, err := dbase.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "dbase.Tasks(50)"})
		return
	}
	if tasks == nil {
		tasks = []*dbase.Task{} // Преобразуем nil в пустой срез
	}

	// Создаем структуру ответа с ключом "tasks"
	response := TasksResp{Tasks: tasks}

	resp, err := json.MarshalIndent(response, "", "    ")
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// ******************************
func getoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	type errorjson struct {
		Error string `json:"error"`
	}
	id := r.URL.Query().Get("id") // Теперь получаем id из query-параметра

	task, err := dbase.GetTask(id) // в параметре максимальное количество записей
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "getoneTaskHandler"})
		return
	}
	resp, err := json.MarshalIndent(task, "", "    ")
	// в заголовок записываем тип контента,  данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// ********************************************************
func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	type errorjson struct {
		Error string `json:"error"`
	}
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "read body error"})
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateAndAdjustTask(&task, w) {
		return // Ошибка уже обработана в функции
	}

	err = dbase.UpdateTask(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "UpdateTask"})
		return
	}
	// в заголовок записываем тип контента,  данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write([]byte("{}"))
}

// *******************************************
func delTaskHandler(w http.ResponseWriter, r *http.Request) {
	type errorjson struct {
		Error string `json:"error"`
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "ID query error"})
		return
	}
	err := dbase.DeleteTask(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "DeleteTask error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write([]byte("{}"))
}

// *******************************************
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	type errorjson struct {
		Error string `json:"error"`
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "ID query error"})
		return
	}

	task, err := dbase.GetTask(idStr) // в параметре максимальное количество записей
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "getoneTaskHandler"})
		return
	}
	//Одноразовая задача с пустым полем repeat удаляется.
	if len(task.Repeat) == 0 {
		err := dbase.DeleteTask(idStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorjson{Error: "DeleteTask error"})
			return
		}
	} else {
		var nextd string
		now := time.Now()
		nextd, err = NextDate(now, task.Date, task.Repeat)
		err = dbase.UpdateDate(nextd, idStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorjson{Error: "DeleteTask error"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write([]byte("{}"))
}
