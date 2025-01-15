package httpserver

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ASHmanR17/go_final_project/internal/models"
	"github.com/ASHmanR17/go_final_project/internal/services"
	"github.com/go-chi/chi/v5"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// handleNextDate Обработчик GET-запросов для /api/nextdate
func handleNextDate(w http.ResponseWriter, r *http.Request) {
	// Проверка, что метод запроса - GET
	if r.Method != "GET" {
		http.Error(w, "Метод запроса не поддерживается", http.StatusMethodNotAllowed)
		return
	}
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
	now, err := time.Parse("20060102", nowString)
	if err != nil {
		http.Error(w, "Некорректная дата сейчас", http.StatusBadRequest)
		return
	}
	// Проверим корректность даты
	_, err = time.Parse("20060102", dateString)
	if err != nil {
		http.Error(w, "Некорректная дата начала", http.StatusBadRequest)
		return
	}

	// Вычисляем следующую дату
	nextDate, err := services.NextDate(now, dateString, repeatString)
	if err != nil {
		http.Error(w, "Ошибка вычисления даты: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Форматируем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
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

// handleAddTask Обработчик для добавления задачи
func handleAddTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// создаем объект типа Scheduler
		var task models.Scheduler
		var buf bytes.Buffer

		// читаем тело запроса
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			// Формируем ответ с ошибкой в формате JSON
			response := map[string]string{
				"error": err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		// десериализуем JSON в структуру task
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			// Формируем ответ с ошибкой в формате JSON
			response := map[string]string{
				"error": fmt.Sprintf("ошибка десериализации JSON: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		// Проверим наличие заголовка Title
		if task.Title == "" {
			// Формируем ответ с ошибкой в формате JSON
			response := map[string]string{
				"error": fmt.Sprintf("не указан заголовок задачи: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		// Проверим корректность даты
		// Регулярное выражение для проверки формата даты
		validFormat := regexp.MustCompile(`^\d{8}$`)
		// Получаем текущую дату и время
		currentDate := time.Now()

		if task.Date == "" { // Проверка пустой даты
			// Форматируем текушую дату в формате "YYYYMMDD" и присваиваем её полю Date
			task.Date = currentDate.Format("20060102")
		} else if !validFormat.MatchString(task.Date) { // Проверка формата даты
			// Формируем ответ с ошибкой в формате JSON
			response := map[string]string{
				"error": fmt.Sprintf("Дата не соответствует формату YYYYMMDD: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		} else if _, err := time.Parse("20060102", task.Date); err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("Неккоректная дата: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}

		// вычисляем следующую дату и заодно проверим правила повторения
		nextDate, err := services.NextDate(currentDate, task.Date, task.Repeat)
		if err != nil {
			response := map[string]string{
				"error": err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		//если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число
		if nextDate == "" {
			task.Date = currentDate.Format("20060102")
		}
		taskTime, _ := time.Parse("20060102", task.Date)

		//если дата задачи меньше сегодняшней, подставляем сегодняшнюю дату
		if taskTime.Before(currentDate) {
			task.Date = currentDate.Format("20060102")
		}

		// добавляем в базу задачу
		res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
			sql.Named("date", task.Date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat),
		)
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("ошибка добавления задачи в базу: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		// верните идентификатор последней добавленной записи
		id, err := res.LastInsertId()
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("не получил идентификатор последней добавленной записи: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		// Формируем ответ с id
		response := map[string]string{
			"id": strconv.FormatInt(id, 10),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	}
}

// handleGetTasks Обработчик для получения списка задач
func handleGetTasks(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO Добавьте возможность выбрать задачи через строку поиска
		// Предел отображаемых задач
		var limit string
		limit = "20"
		// SQL-запрос для получения ближайших задач
		query := `SELECT * FROM scheduler ORDER BY date LIMIT :limit`
		rows, err := db.Query(query, sql.Named("limit", limit))
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("ошибка в запросе к базе данных: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				response := map[string]string{
					"error": fmt.Sprintf("ошибка закрытия Rows: %s", err),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(response)
				return
			}
		}(rows)

		// Массив для хранения задач
		var tasks []models.Scheduler

		// Чтение данных из базы данных и создание массива задач
		for rows.Next() {
			var task models.Scheduler
			err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				response := map[string]string{
					"error": fmt.Sprintf("ошибка при чтении (Rows.Scan) базы данных: %s", err),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(response)
				return
			}
			tasks = append(tasks, task)
		}

		// Закрытие соединения с базой данных
		err = rows.Close()
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("ошибка закрытия базы данных: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Преобразование массива задач в JSON
		if len(tasks) == 0 {
			tasks = []models.Scheduler{}
		}

		response := struct {
			Tasks []models.Scheduler `json:"tasks"`
		}{
			Tasks: tasks,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// handleGetTask Обработчик для GET-запроса /api/task?id=<идентификатор>
func handleGetTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлечение ID задачи из URL
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			response := map[string]string{
				"error": fmt.Sprintf("Задача не найдена: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Выполнение SQL-запроса для получения задачи по ID
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`
		row := db.QueryRow(query, sql.Named("id", idStr))

		// Чтение данных из базы данных и создание объекта задачи
		var task models.Scheduler
		err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("ошибка при чтении (Row.Scan) базы данных: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Преобразование объекта задачи в JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

// handleEditTask Обработчик для редактирования задачи
func handleEditTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// создаем объект типа Scheduler
		var task models.Scheduler
		// Чтение JSON-данных из запроса
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("Ошибка декодирования JSON: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(response)
			return
		}
		// Извлечение ID задачи
		var idStr string
		idStr = task.ID
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			response := map[string]string{
				"error": fmt.Sprintf("Неверный или отсутствующий ID: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Проверка существования задачи по ID
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
		row := db.QueryRow(query, id)
		var existingTask models.Scheduler
		err = row.Scan(&existingTask.ID, &existingTask.Date, &existingTask.Title, &existingTask.Comment, &existingTask.Repeat)
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("Ошибка при чтении базы данных: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Сравнение текущих данных с данными из запроса
		if existingTask.ID != idStr {
			response := map[string]string{
				"error": "ID в запросе не совпадает с ID существующей задачи",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Проверим наличие заголовка Title
		if task.Title == "" {
			// Формируем ответ с ошибкой в формате JSON
			response := map[string]string{
				"error": fmt.Sprintf("не указан заголовок задачи: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		// Проверим корректность даты
		// Регулярное выражение для проверки формата даты
		validFormat := regexp.MustCompile(`^\d{8}$`)
		// Получаем текущую дату и время
		currentDate := time.Now()

		if task.Date == "" { // Проверка пустой даты
			// Форматируем текушую дату в формате "YYYYMMDD" и присваиваем её полю Date
			task.Date = currentDate.Format("20060102")
		} else if !validFormat.MatchString(task.Date) { // Проверка формата даты
			// Формируем ответ с ошибкой в формате JSON
			response := map[string]string{
				"error": fmt.Sprintf("Дата не соответствует формату YYYYMMDD: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		} else if _, err := time.Parse("20060102", task.Date); err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("Неккоректная дата: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}

		// вычисляем следующую дату и заодно проверим правила повторения
		nextDate, err := services.NextDate(currentDate, task.Date, task.Repeat)
		if err != nil {
			response := map[string]string{
				"error": err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		//если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число
		if nextDate == "" {
			task.Date = currentDate.Format("20060102")
		}
		taskTime, _ := time.Parse("20060102", task.Date)

		//если дата задачи меньше сегодняшней, подставляем сегодняшнюю дату
		if taskTime.Before(currentDate) {
			task.Date = currentDate.Format("20060102")
		}

		// добавляем в базу задачу
		query = "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
		_, err = db.Exec(query,
			sql.Named("date", task.Date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat),
			sql.Named("id", task.ID),
		)
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("ошибка изменения задачи: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Формируем ответ с пустым JSON
		response := map[string]string{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	}
}

// handleDoneTask Обработчик для завершения задачи
func handleDoneTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// создаем объект типа Scheduler
		var task models.Scheduler
		id := r.URL.Query().Get("id")

		// Проверка существования задачи по ID
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
		row := db.QueryRow(query, id)

		err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("Ошибка при чтении базы данных: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		// Если правила повторения нет, удаляем задачу из базы
		if task.Repeat == "" {
			query = "DELETE FROM scheduler WHERE id = :id"
			_, err = db.Exec(query, sql.Named("id", id))
			if err != nil {
				response := map[string]string{
					"error": fmt.Sprintf("Ошибка удаления из базы: %s", err),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(response)
				return
			}
			// Формируем ответ с пустым JSON
			response := map[string]string{}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(response)
			return
		}
		// Получаем текущую дату и время
		currentDate := time.Now()
		// вычисляем следующую дату и заодно проверим правила повторения
		nextDate, err := services.NextDate(currentDate, task.Date, task.Repeat)
		if err != nil {
			response := map[string]string{
				"error": err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}

		// добавляем в базу дату следующего выполнения
		query = "UPDATE scheduler SET date = :date WHERE id = :id"
		_, err = db.Exec(query,
			sql.Named("date", nextDate),
			sql.Named("id", task.ID),
		)
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("ошибка изменения даты в задаче: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Формируем ответ с пустым JSON
		response := map[string]string{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	}
}

// handleDeleteTask Обработчик для удаления задачи
func handleDeleteTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// создаем объект типа Scheduler
		//var task models.Scheduler
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			response := map[string]string{
				"error": "Отсутствует обязательный параметр 'id'.",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(response)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			response := map[string]string{
				"error": fmt.Sprintf("Неверный или отсутствующий ID: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(response)
			return
		}

		//// Проверка существования задачи по ID
		//query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
		//row := db.QueryRow(query, id)
		//
		//err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		//if err != nil {
		//	response := map[string]string{
		//		"error": fmt.Sprintf("Ошибка при чтении базы данных: %s", err),
		//	}
		//	w.Header().Set("Content-Type", "application/json")
		//	w.WriteHeader(500)
		//	json.NewEncoder(w).Encode(response)
		//	return
		//}
		// удаляем задачу из базы
		query := "DELETE FROM scheduler WHERE id = :id"
		_, err = db.Exec(query, sql.Named("id", id))
		if err != nil {
			response := map[string]string{
				"error": fmt.Sprintf("Ошибка удаления из базы: %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response)
			return
		}
		// Формируем ответ с пустым JSON
		response := map[string]string{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	}
}
