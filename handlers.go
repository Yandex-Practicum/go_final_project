package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
)

// Task описывает поля таблицы scheduler
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// JSONObject используется при формиравании JSON-объекта
type JSONObject struct {
	ID    string `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// TaskResponse используется для передачи списка задач в формате JSON
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// Password используется для передачи пароля в формате JSON
type Password struct {
	Pass string `json:"password"`
}

// SendErrorResponse отправляет ошибку в формате JSON и статус сервера
func SendErrorResponse(res http.ResponseWriter, errorMessage string, statusCode int) {
	response := JSONObject{Error: errorMessage}

	resp, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(res, "ошибка при сериализации JSON", http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(statusCode)
	_, err = res.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

// SendJSONResponse отправляет ответ в формате JSON
func SendJSONResponse(res http.ResponseWriter, response interface{}) {
	resp, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при сериализации JSON", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

// AddTaskRules проверяет необходимые условия при добавлении задачи, а именно:
// корректность формата даты и наличие поля title. Возвращает дату в формате int и ошибку
func AddTaskRules(t *Task) (int, error) {
	// Если время не указано, устанавливаем текущее время
	if t.Date == "" {
		t.Date = time.Now().Format(timeFormat)
	}

	// Проверяем формат времени
	_, err := time.Parse(timeFormat, t.Date)
	if err != nil {
		return 0, fmt.Errorf("дата представлена в формате, отличном от 20060102: %w", err)
	}

	if t.Date < time.Now().Format(timeFormat) {
		if t.Repeat == "" {
			t.Date = time.Now().Format(timeFormat)
		} else {
			taskDay, err := DateParse(now, t.Date, t.Repeat)
			if err != nil {
				return 0, fmt.Errorf("ошибка при парсинге даты: %w", err)
			}
			t.Date = taskDay
		}
	}

	// Переводим дату в int
	dateInt, err := strconv.Atoi(t.Date)
	if err != nil {
		return 0, fmt.Errorf("не удалось конвертировать дату в число: %w", err)
	}

	if t.Title == "" {
		return 0, fmt.Errorf("не указан заголовок задачи")
	}

	return dateInt, nil
}

// NextDateHandler принимает запрос в формате
// "/api/nextdate?now=<20060102>&date=<20060102>&repeat=<правило>"
// и возвращает следующую дату задачи или ошибку
func NextDateHandler(res http.ResponseWriter, req *http.Request) {
	// Получаем параметры now, date, repeat из запроса
	nowStr := req.FormValue("now")
	dateStr := req.FormValue("date")
	repeatStr := req.FormValue("repeat")

	// Проверяем передано ли время в нужном формате
	t, err := time.Parse(timeFormat, nowStr)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при парсинге времени now", http.StatusBadRequest)
		return
	}

	// Получаем ближайшую дату задачи
	taskDay, err := DateParse(t, dateStr, repeatStr)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при парсинге даты", http.StatusBadRequest)
		return
	}

	_, _ = res.Write([]byte(taskDay))
}

// PostTask обработчик для POST-запроса /api/task, который добавляет задачу в базу данных.
// Запрос передается в формате JSON. Возвращает JSON с полем id, в случае успешного
// добавления задачи, или error с текстом ошибки
func PostTask(res http.ResponseWriter, req *http.Request) {
	// Создаем экземпляр структуры Task и заполняем его поля значениями,
	// полученными в формате JSON
	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	dateInt, err := AddTaskRules(&task)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при добавлении задачи", http.StatusBadRequest)
		return
	}

	// Добавляем задачу в таблицу scheduler
	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", dateInt),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "не удалось вставить задачу в таблицу", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка получения последнего ID", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ
	response := JSONObject{ID: fmt.Sprintf("%d", id)}

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(res).Encode(response)
}

// GetTasks обработчик для GET-запроса /api/tasks, который возвращает список ближайших
// задач в формате JSON.
// Дополнительно обрабатывает параметр search в строке поиска
// Поиск происходит по слову и по дате
func GetTasks(res http.ResponseWriter, req *http.Request) {
	var (
		tasks []Task
		query string
		args  []interface{}
	)

	limit := 50
	searchStr := req.URL.Query().Get("search")

	if searchStr == "" {
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit"
		args = append(args, limit)
	} else {
		t, err := time.Parse("02.01.2006", searchStr)
		if err == nil {
			query = "SELECT * FROM scheduler WHERE date = :date LIMIT :limit"
			args = append(args, t.Format(timeFormat), limit)
		} else {
			query = "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
			args = append(args, "%"+searchStr+"%", limit)
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "не удалось получить задачу", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			SendErrorResponse(res, "ошибка при извлечении данных", http.StatusInternalServerError)
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при итерации", http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		response := TasksResponse{Tasks: []Task{}}
		SendJSONResponse(res, response)
		return
	}

	response := TasksResponse{Tasks: tasks}
	SendJSONResponse(res, response)
}

// GetTaskId обработчик для GET-запроса /api/task?id=<идентификатор>.
// Возращает JSON-объект со всеми полями задачи с указанным идентификатором.
func GetTaskId(res http.ResponseWriter, req *http.Request) {
	// Получаем id из запроса
	id := req.URL.Query().Get("id")

	// Создаем экзмепляр структуры Task и заполняем его поля значениеми из таблицы scheduler
	var task Task

	rows, err := db.Query("SELECT *FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			SendErrorResponse(res, "задача не найдена", http.StatusNotFound)
			return
		}
		SendErrorResponse(res, "ошибка на стороне сервера", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			SendErrorResponse(res, "ошибка при извлечении данных", http.StatusInternalServerError)
			return
		}
	}

	if task == (Task{}) {
		SendErrorResponse(res, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	SendJSONResponse(res, task)
}

// PutTask реализует редактирование задач.
// PUT-обработчик /api/task, который отправляет значение
// всех полей в виде JSON-объекта. Возвращает JSON cо структурой Task или ошибку
func PutTask(res http.ResponseWriter, req *http.Request) {
	var (
		taskCurrent Task
		taskUpdate  Task
	)

	err := json.NewDecoder(req.Body).Decode(&taskUpdate)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	err = db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", taskUpdate.ID)).
		Scan(&taskCurrent.ID, &taskCurrent.Date, &taskCurrent.Title, &taskCurrent.Comment, &taskCurrent.Repeat)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при извлечении данных", http.StatusInternalServerError)
		return
	}

	dateInt, err := AddTaskRules(&taskUpdate)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при добавлении задачи", http.StatusBadRequest)
		return
	}

	// Если переданные в JSON поля пустые, оставляем поля преждней структуры
	if taskUpdate.Title == "" {
		SendErrorResponse(res, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if taskUpdate.Date != "" {
		taskCurrent.Date = taskUpdate.Date
	}

	if taskUpdate.Title != "" {
		taskCurrent.Title = taskUpdate.Title
	}

	if taskUpdate.Comment != "" {
		taskCurrent.Comment = taskUpdate.Comment
	}

	if taskUpdate.Repeat != "" {
		taskCurrent.Repeat = taskUpdate.Repeat
	}

	// Обновляем поля в таблице
	_, err = db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment , repeat = :repeat WHERE id = :id",
		sql.Named("date", dateInt),
		sql.Named("title", taskCurrent.Title),
		sql.Named("comment", taskCurrent.Comment),
		sql.Named("repeat", taskCurrent.Repeat),
		sql.Named("id", taskCurrent.ID))

	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при обновлении данных", http.StatusInternalServerError)
		return
	}

	if taskCurrent == (Task{}) {
		SendErrorResponse(res, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Отправляем JSON со структурой taskCurrent
	SendJSONResponse(res, taskCurrent)
}

// TaskDone обработчик для POST-запроса /api/task/done, который делает
// задачу выполненой. Одноразовая задача с пустым полем repeat удаляется.
// Возвращает пустой JSON или ошибку
func TaskDone(res http.ResponseWriter, req *http.Request) {
	// Получаем id из запроса
	id := req.URL.Query().Get("id")

	if id == "" {
		SendErrorResponse(res, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Создаем экземпляр структуры типа Task и заполняем его поля
	// значениями из таблицы с указанным id
	var taskCurrent Task

	err := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id)).
		Scan(&taskCurrent.ID, &taskCurrent.Date, &taskCurrent.Title, &taskCurrent.Comment, &taskCurrent.Repeat)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			SendErrorResponse(res, "задача не найдена", http.StatusNotFound)
			return
		}
		SendErrorResponse(res, "ошибка на стороне сервера", http.StatusInternalServerError)
		return
	}

	// Удаляем задачу с пустым полем repeat
	if taskCurrent.Repeat == "" {
		DeleteTask(res, req)
		return
	}

	// Получаем ближайшее время задачи
	date, err := DateParse(now, taskCurrent.Date, taskCurrent.Repeat)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при получении даты", http.StatusBadRequest)
		return
	}

	// Приводим время к типу int и обновляем строку в таблице scheduler
	dateInt, err := strconv.Atoi(date)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при преобразовании", http.StatusInternalServerError)
	}

	_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
		sql.Named("date", dateInt),
		sql.Named("id", id))

	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при обновлении данных", http.StatusInternalServerError)
		return
	}

	taskCurrent.Date = date

	// В случае успешного обновления поля таблицы, отправляем пустой JSON
	SendJSONResponse(res, struct{}{})
}

// DeleteTask обработчик DELETE-запроса /api/task/done?id=<идентификатор>.
// Удаляет задачу из таблицы scheduler. Возвращает пустой JSON или ошибку.
func DeleteTask(res http.ResponseWriter, req *http.Request) {
	// Получаем id из запроса
	id := req.URL.Query().Get("id")

	// Удаляем задачу по полученному id
	result, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при удалении задачи", http.StatusInternalServerError)
		return
	}

	// Проверяем полявились ли изменения в таблице
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка при получении количества удалённых строк", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		SendErrorResponse(res, "задача с таким id не найдена", http.StatusNotFound)
		return
	}

	// В случае успешного удаления, отправляем пустой JSON
	SendJSONResponse(res, struct{}{})
}

// SignIn обработчик POST-запроса /api/signin. Получает JSON с полем password.
// Если пароль совпадает, формирует JWT и передает его в поле JSON-объекта.
// Если пароль невернный или произошла ошибка, возвращает JSON с текстом ошибки
func SignIn(res http.ResponseWriter, req *http.Request) {
	var p Password

	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) == 0 {
		return
	}

	// Получаем пароль из JSON
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		log.Println(err)
		SendErrorResponse(res, "ошибка десериализации JSON", http.StatusUnauthorized)
		return
	}

	if p.Pass != pass {
		SendErrorResponse(res, "неверный пароль", http.StatusUnauthorized)
		return
	}

	// Создаем токен и в качетсве полезной нагрузки передаем хэш
	secret := []byte("secret_key")

	hashedPass := sha256.Sum256([]byte(p.Pass))
	hashString := hex.EncodeToString(hashedPass[:])

	claims := jwt.MapClaims{
		"hash": hashString,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		log.Println(err)
		http.Error(res, "ошибка получения подписанного токена", http.StatusUnauthorized)
		return
	}
	response := JSONObject{Token: signedToken}
	SendJSONResponse(res, response)
}

// Auth проводит проверку аутентификации для API-запросов
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) == 0 {
			return
		}

		var jwtCookie string // JWT-токен из куки
		// Получаем куку
		cookie, err := req.Cookie("token")
		if err == nil {
			jwtCookie = cookie.Value
		}

		//валидация и проверка JWT-токена
		secret := []byte("secret_key")
		jwtToken, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			log.Println(err)
			http.Error(res, "Ошибка при парсинге токена", http.StatusUnauthorized)
			return
		}

		if !jwtToken.Valid {
			SendErrorResponse(res, "токен не валиден", http.StatusUnauthorized)
			return
		}

		// Получаем хэш из полезной нагрузки токена
		result, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			SendErrorResponse(res, "не удалось выполнить приведение типа к  jwt.MapClaims", http.StatusUnauthorized)
			return
		}

		hashRow := result["hash"]

		hash, ok := hashRow.(string)
		if !ok {
			SendErrorResponse(res, "не удалось выполнить приведение типа к string", http.StatusUnauthorized)
			return
		}

		// Получаем хэш, как в функции SignIn
		hashedPass := sha256.Sum256([]byte(pass))
		expectedHash := hex.EncodeToString(hashedPass[:])

		if hash != expectedHash {
			SendErrorResponse(res, "ошибка валидации хэша", http.StatusUnauthorized)
			return
		}

		next(res, req)
	})
}
