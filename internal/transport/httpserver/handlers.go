package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ASHmanR17/go_final_project/internal/database"

	"github.com/ASHmanR17/go_final_project/internal/service"
	"github.com/go-chi/chi/v5"
)

const (
	DateLayout = "20060102"
	JsonValue  = "application/json; charset=UTF-8"
)

type httpHandler struct {
	taskService service.TaskService
}

func newHTTPHandler(taskService service.TaskService) *httpHandler {
	return &httpHandler{
		taskService: taskService,
	}
}

// NextDate вычисляет следующую дату для задачи в соответствии с указанным правилом
func (h httpHandler) NextDate(w http.ResponseWriter, r *http.Request) {

	// Получаем параметры запроса

	nowString := r.URL.Query().Get("now")
	dateString := r.URL.Query().Get("date")
	repeatString := r.URL.Query().Get("repeat")

	// Проверяем, что все параметры были переданы
	if nowString == "" || dateString == "" || repeatString == "" {
		http.Error(w, "Не переданы все необходимые параметры", http.StatusBadRequest)
		return
	}

	// Преобразуем nowString в формат time.Time
	now, err := time.Parse(DateLayout, nowString)
	if err != nil {
		http.Error(w, "Некорректная дата сейчас", http.StatusBadRequest)
		return
	}
	// Проверим корректность даты
	_, err = time.Parse(DateLayout, dateString)
	if err != nil {
		http.Error(w, "Некорректная дата начала", http.StatusBadRequest)
		return
	}

	// Вычисляем следующую дату
	nextDate, err := service.NextDate(now, dateString, repeatString)
	if err != nil {
		http.Error(w, "Ошибка вычисления даты: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Форматируем ответ
	w.Header().Set("Content-Type", JsonValue)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

// AddTask Обработчик для добавления задачи
func (h httpHandler) AddTask(w http.ResponseWriter, r *http.Request) {

	// создаем объект типа Scheduler
	var task database.Scheduler

	// десериализация JSON
	task, err := h.taskService.TaskFromJson(r.Body)
	if err != nil {
		// Формируем ответ с ошибкой в формате JSON
		response := map[string]string{
			"error": err.Error(),
		}

		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Проверим данные задачи на корректность
	task, err = h.taskService.CheckTask(task)
	if err != nil {
		// Формируем ответ с ошибкой в формате JSON
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// добавляем в базу задачу
	id, err := h.taskService.AddTask(task)
	if err != nil {
		// Формируем ответ с ошибкой в формате JSON
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	// Формируем ответ с id
	response := map[string]string{
		"id": strconv.Itoa(id),
	}
	w.Header().Set("Content-Type", JsonValue)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetTasks Обработчик для получения списка задач
func (h httpHandler) GetTasks(w http.ResponseWriter, r *http.Request) {

	tasks, err := h.taskService.GetTasks()
	if err != nil {
		// Формируем ответ с ошибкой в формате JSON
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(tasks) == 0 {
		tasks = []database.Scheduler{}
	}

	// Преобразование массива задач в JSON
	response := struct {
		Tasks []database.Scheduler `json:"tasks"`
	}{
		Tasks: tasks,
	}
	w.Header().Set("Content-Type", JsonValue)
	json.NewEncoder(w).Encode(response)

}

// GetTask Обработчик для GET-запроса /api/task?id=<идентификатор>
func (h httpHandler) GetTask(w http.ResponseWriter, r *http.Request) {

	// Извлечение Id задачи из URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response := map[string]string{
			"error": fmt.Sprint("неверный id"),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Выполнение SQL-запроса для получения задачи по Id
	task, err := h.taskService.GetTask(idStr)
	if err != nil {
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Преобразование объекта задачи в JSON
	w.Header().Set("Content-Type", JsonValue)
	json.NewEncoder(w).Encode(task)

}

// EditTask Обработчик для редактирования задачи
func (h httpHandler) EditTask(w http.ResponseWriter, r *http.Request) {

	// создаем объект типа Scheduler
	var task database.Scheduler
	// Чтение JSON-данных из запроса (десериализация JSON)
	task, err := h.taskService.TaskFromJson(r.Body)
	if err != nil {
		// Формируем ответ с ошибкой в формате JSON
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	// Проверим данные задачи на корректность
	task, err = h.taskService.CheckTask(task)
	if err != nil {
		// Формируем ответ с ошибкой в формате JSON
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// обновим задачу в базе
	err = h.taskService.UpdateTask(task)
	if err != nil {
		// Формируем ответ с ошибкой в формате JSON
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Формируем ответ с пустым JSON
	response := map[string]string{}
	w.Header().Set("Content-Type", JsonValue)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// DoneTask Обработчик для завершения задачи
func (h httpHandler) DoneTask(w http.ResponseWriter, r *http.Request) {

	// создаем объект типа Scheduler
	var task database.Scheduler
	id := r.URL.Query().Get("id")

	// Получим из базы задачу по Id
	task, err := h.taskService.GetTask(id)
	if err != nil {
		// Формируем ответ с ошибкой в формате JSON
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Если правила повторения нет, удаляем задачу из базы
	if task.Repeat == "" {
		err := h.taskService.DeleteTask(task.Id)
		if err != nil {
			// Формируем ответ с ошибкой в формате JSON
			response := map[string]string{
				"error": err.Error(),
			}
			w.Header().Set("Content-Type", JsonValue)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		// Формируем ответ с пустым JSON
		response := map[string]string{}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	// Получаем текущую дату и время
	currentDate := time.Now()
	// вычисляем следующую дату и заодно проверим правила повторения
	nextDate, err := service.NextDate(currentDate, task.Date, task.Repeat)
	if err != nil {
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// обновим в базе задачу с новой датой
	task.Date = nextDate
	err = h.taskService.UpdateTask(task)
	if err != nil {
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Формируем ответ с пустым JSON
	response := map[string]string{}
	w.Header().Set("Content-Type", JsonValue)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteTask Обработчик для удаления задачи
func (h httpHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil || id <= 0 {
		response := map[string]string{
			"error": fmt.Sprintf("Неверный или отсутствующий Id"),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	// удаляем задачу из базы
	err = h.taskService.DeleteTask(idStr)
	if err != nil {
		response := map[string]string{
			"error": err.Error(),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Формируем ответ с пустым JSON
	response := map[string]string{}
	w.Header().Set("Content-Type", JsonValue)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
