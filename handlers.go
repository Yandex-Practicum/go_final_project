package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	Tasks []Task `json:"tasks"`
	Pass  string `json:"password"`
}

// SendErrorResponse отправляет ошибку в формате JSON и статус сервера
func SendErrorResponse(res http.ResponseWriter, errorMessage string, statusCode int) {
	response := JSONObject{Error: errorMessage}
	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(res).Encode(response); err != nil {
		http.Error(res, "не удалось обработать ошибку", statusCode)
	}
}

// SendJSONResponse отправляет ответ в формате JSON
func SendJSONResponse(res http.ResponseWriter, response interface{}) {
	resp, err := json.Marshal(response)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при сериализации JSON: %v", err), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(resp)
}

// AddTaskRules проверяет необходимые условия при добавлении задачи, а именно:
// корректность формата даты и наличие поля title. Возвращает дату в формате int и ошибку
func AddTaskRules(t *Task) (int, error) {
	if t.Date == "" {
		t.Date = time.Now().Format(timeFormat)
	}
	_, err := time.Parse(timeFormat, t.Date)
	if err != nil {
		return 0, fmt.Errorf("дата представлена в формате, отличном от 20060102")
	}

	if t.Date < time.Now().Format(timeFormat) {
		if t.Repeat == "" {
			t.Date = time.Now().Format(timeFormat)
		} else {
			taskDay, err := DateParse(now, t.Date, t.Repeat)
			if err != nil {
				return 0, fmt.Errorf("ошибка при парсинге даты: %v", err)
			}
			t.Date = taskDay
		}
	}

	dateInt, err := strconv.Atoi(t.Date)
	if err != nil {
		return 0, fmt.Errorf("не удалось конвертировать дату в число")
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
	// Получаем параметры now, daye, repeat из запроса
	nowStr := req.FormValue("now")
	dateStr := req.FormValue("date")
	repeatStr := req.FormValue("repeat")

	// Проверяем передано ли время в нужном формате
	t, err := time.Parse(timeFormat, nowStr)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при парсинге времени now: %v", err), http.StatusInternalServerError)
		return
	}

	// Получаем ближайшую дату задачи
	taskDay, err := DateParse(t, dateStr, repeatStr)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при парсинге даты: %v", err), http.StatusInternalServerError)
		return
	}

	_, _ = res.Write([]byte(taskDay))
}

// PostTask POST-обработчик /api/task, который добавляет задачу в базу данных.
// Запрос передается в формате JSON. Возвращает JSON с полем id, в случае успешного
// добавления задачи, или error с текстом ошибки
func PostTask(res http.ResponseWriter, req *http.Request) {
	// Создаем экзмепляр структуры Task и заполняем его поля значениеми,
	// полученными в формате JSON
	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка десериализации JSON: %v", err), http.StatusInternalServerError)
		return
	}

	dateInt, err := AddTaskRules(&task)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при добавлении задачи: %v", err), http.StatusInternalServerError)
		return
	}

	// Добавляем задачу в таблицу scheduler
	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		dateInt, task.Title, task.Comment, task.Repeat)
	if err != nil {
		SendErrorResponse(res, "не удалось вставить задачу в таблицу", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		SendErrorResponse(res, "ошибка получения последнего ID", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ
	response := JSONObject{ID: fmt.Sprintf("%d", id)}
	//SendJSONResponse(res, response)

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(res).Encode(response)
}

// GetTasks GET-обработчик /api/tasks, который возвращает список ближайших
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
		SendErrorResponse(res, "не удалось получить задачу", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			SendErrorResponse(res, fmt.Sprintf("ошибка при извлечении данных: %v", err), http.StatusInternalServerError)
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при итерации: %v", err), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		response := JSONObject{Tasks: []Task{}}
		SendJSONResponse(res, response)
		return
	}

	response := JSONObject{Tasks: tasks}
	SendJSONResponse(res, response)
}

// GetTaskId GET-обработчки запроса /api/task?id=<идентификатор>.
// Возращает JSON-объект со всеми плолями задачи с указанным идентификатором.
func GetTaskId(res http.ResponseWriter, req *http.Request) {
	// Получаем id из url
	id := req.URL.Query().Get("id")

	// Создаем экзмепляр структуры Task и заполняем его поля значениеми из таблицы scheduler
	var task Task

	rows, err := db.Query("SELECT *FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		SendErrorResponse(res, "задача не найдена", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			SendErrorResponse(res, fmt.Sprintf("ошибка при извлечении данных: %v", err), http.StatusInternalServerError)
			return
		}
	}

	if task == (Task{}) {
		SendErrorResponse(res, "не указан идентификатор", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(res, task)
}

// PutTask реализует редактирование задач.
// PUT-обработчик /api/task, который отправляет значение
// всех полей в виде JSON-объекта. Возвращает пустой JSON или ошибку
func PutTask(res http.ResponseWriter, req *http.Request) {
	var (
		taskCurrent Task
		taskUpdate  Task
	)

	err := json.NewDecoder(req.Body).Decode(&taskUpdate)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка десериализации JSON: %v", err), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", taskUpdate.ID)).
		Scan(&taskCurrent.ID, &taskCurrent.Date, &taskCurrent.Title, &taskCurrent.Comment, &taskCurrent.Repeat)
	if err != nil {
		SendErrorResponse(res, "ошибка при извлечении данных", http.StatusInternalServerError)
		return
	}

	dateInt, err := AddTaskRules(&taskUpdate)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при добавлении задачи: %v", err), http.StatusInternalServerError)
		return
	}

	if taskUpdate.Title == "" {
		SendErrorResponse(res, "не указан заголовок задачи", http.StatusInternalServerError)
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

	_, err = db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment , repeat = :repeat WHERE id = :id",
		sql.Named("date", dateInt),
		sql.Named("title", taskCurrent.Title),
		sql.Named("comment", taskCurrent.Comment),
		sql.Named("repeat", taskCurrent.Repeat),
		sql.Named("id", taskCurrent.ID))

	if err != nil {
		SendErrorResponse(res, "ошибка при обновлении данных", http.StatusInternalServerError)
		return
	}

	if taskCurrent == (Task{}) {
		SendErrorResponse(res, "не указан идентификатор", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(res, taskCurrent)
}

// TaskDone обработчик для POST-запроса /api/task/done, который делает
// задачу выполненой. Одноразовая задача с пустым полем repeat удаляется.
// Возвращает пустой JSON или ошибку
func TaskDone(res http.ResponseWriter, req *http.Request) {
	// Получаем id из запроса
	id := req.URL.Query().Get("id")

	if id == "" {
		SendErrorResponse(res, "не указан идентификатор", http.StatusInternalServerError)
		return
	}

	// Создаем экзмплер структуры типа Task и заполняем его поля
	// значениями из таблицы с указнным id
	var taskCurrent Task

	err := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id)).
		Scan(&taskCurrent.ID, &taskCurrent.Date, &taskCurrent.Title, &taskCurrent.Comment, &taskCurrent.Repeat)
	if err != nil {
		SendErrorResponse(res, "задача не найдена", http.StatusInternalServerError)
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
		SendErrorResponse(res, "ошибка при получении даты", http.StatusInternalServerError)
		return
	}

	// Приводим время к типу int и обновляем строку в таблице scheduler
	dateInt, err := strconv.Atoi(date)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при преобразовании: %v", err), http.StatusInternalServerError)
	}

	_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
		sql.Named("date", dateInt),
		sql.Named("id", id))

	if err != nil {
		SendErrorResponse(res, "ошибка при обновлении данных", http.StatusInternalServerError)
		return
	}

	taskCurrent.Date = date
	// В случае успешного обновления столбца (?) таблицы, отправляем пустой JSON
	SendJSONResponse(res, struct{}{})
}

// DeleteTask обработчик DELETE-запроса /api/task/done?id=<идентификатор>.
// Удаляет задачу из таблицы scheduler. Возвращает пустой JSON или ошибку.
func DeleteTask(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	result, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		SendErrorResponse(res, "ошибка при удалении задачи", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		SendErrorResponse(res, "ошибка при получении количества удалённых строк", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		SendErrorResponse(res, "задача с таким id не найдена", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(res, struct{}{})
}

// SignIn обработчик POST-запроса /api/signin. Получает JSON с полем password.
// Если пароль совпадает, возвращает формирует JWT и передает его в поле JSON-объекта.
// Если пароль невернный или произошла ошибка, возвращает JSON с текстом ошибки
func SignIn(res http.ResponseWriter, req *http.Request) {
	var p JSONObject

	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) > 0 { // Получаем пароль из JSON
		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			SendErrorResponse(res, fmt.Sprintf("ошибка десериализации JSON: %v", err), http.StatusUnauthorized)
			return
		}

		if p.Pass == pass { // Создаем токен и в качетсве полезной нагрузки передаем хэш
			secret := []byte("secret_key")

			hashedPass := sha256.Sum256([]byte(p.Pass))
			hashString := hex.EncodeToString(hashedPass[:])

			claims := jwt.MapClaims{
				"hash": hashString,
			}

			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			signedToken, err := jwtToken.SignedString(secret)
			if err != nil {
				http.Error(res, "ошибка получения подписанного токена", http.StatusBadRequest)
			}
			response := JSONObject{Token: signedToken}
			SendJSONResponse(res, response)
		} else {
			SendErrorResponse(res, "невереный пароль", http.StatusUnauthorized)
		}
	}
}

// Auth проводит проверку аутентификации для API-запросов
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
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
				http.Error(res, fmt.Sprintf("Ошибка при парсинге токена: %v", err), http.StatusUnauthorized)
				return
			}

			if !jwtToken.Valid {
				SendErrorResponse(res, fmt.Sprint("токен не валиден"), http.StatusUnauthorized)
				return
			}

			// Получаем хэш из полезной нагрузки токена
			result, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				SendErrorResponse(res, fmt.Sprint("не удалось выполнить приведение типа к  jwt.MapClaims"), http.StatusUnauthorized)
				return
			}

			hashRow := result["hash"]

			hash, ok := hashRow.(string)
			if !ok {
				SendErrorResponse(res, fmt.Sprint("не удалось выполнить приведение типа к string"), http.StatusUnauthorized)
				return
			}

			// Получаем хэш, как в функции SignIn
			hashedPass := sha256.Sum256([]byte(pass))
			expectedHash := hex.EncodeToString(hashedPass[:])

			if hash != expectedHash {
				SendErrorResponse(res, fmt.Sprint("ошибка валидации хэша"), http.StatusUnauthorized)
				return
			}
		}
		next(res, req)
	})
}
