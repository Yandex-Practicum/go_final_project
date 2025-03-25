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

// SendErrorResponse отправляет ошибку в формате JSON и статус сервера(?)
func SendErrorResponse(res http.ResponseWriter, errorMessage string) {
	response := Response{Error: errorMessage}
	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(res).Encode(response); err != nil {
		http.Error(res, "не удалось обработать ошибку", http.StatusInternalServerError)
	}
}

// SendJSONResponse отправляет ответ в формате JSON
func SendJSONResponse(res http.ResponseWriter, response interface{}) {
	resp, err := json.Marshal(response)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при сериализации JSON: %v", err))
		return
	}

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(resp)
}

// NextDateHandler принимает запрос в формате
// "/api/nextdate?now=<20060102>&date=<20060102>&repeat=<правило>"
// и возвращает следующую дату задачи или ошибку
func NextDateHandler(res http.ResponseWriter, req *http.Request) {
	nowStr := req.FormValue("now")
	dateStr := req.FormValue("date")
	repeatStr := req.FormValue("repeat")

	t, err := time.Parse(timeFormat, nowStr)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при парсинге времени now: %v", err))
		return
	}

	taskDay, err := DateParse(t, dateStr, repeatStr)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при парсинге даты: %v", err))
		return
	}

	_, _ = res.Write([]byte(taskDay))
}

// PostTask POST-обработчик /api/task, который добавляет задачу в базу данных.
// Запрос передается в формате JSON. Возвращает JSON с полем id, в случае успешного
// добавления задачи, или error с текстом ошибки
func PostTask(res http.ResponseWriter, req *http.Request) {
	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка десериализации JSON: %v", err))
		return
	}

	dateInt, err := addTaskRules(&task)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при добавлении задачи: %v", err))
		return
	}

	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		dateInt, task.Title, task.Comment, task.Repeat)
	if err != nil {
		SendErrorResponse(res, "не удалось вставить задачу в таблицу")
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		SendErrorResponse(res, "ошибка получения последнего ID")
		return
	}

	response := Response{ID: fmt.Sprintf("%d", id)}

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
		SendErrorResponse(res, "не удалось получить задачу")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			SendErrorResponse(res, fmt.Sprintf("ошибка при извлечении данных: %v", err))
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при итерации: %v", err))
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

// GetTaskId GET-обработчки запроса /api/task?id=<идентификатор>.
// Возращает JSON-объект со всеми плолями задачи с указанным идентификатором.
func GetTaskId(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	var task Task

	rows, err := db.Query("SELECT *FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		SendErrorResponse(res, "задача не найдена")
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			SendErrorResponse(res, fmt.Sprintf("ошибка при извлечении данных: %v", err))
			return
		}
	}

	if task == (Task{}) {
		SendErrorResponse(res, "не указан идентификатор")
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
		SendErrorResponse(res, fmt.Sprintf("ошибка десериализации JSON: %v", err))
		return
	}

	err = db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", taskUpdate.ID)).
		Scan(&taskCurrent.ID, &taskCurrent.Date, &taskCurrent.Title, &taskCurrent.Comment, &taskCurrent.Repeat)
	if err != nil {
		SendErrorResponse(res, "ошибка при извлечении данных")
		return
	}

	dateInt, err := addTaskRules(&taskUpdate)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при добавлении задачи: %v", err))
		return
	}

	if taskUpdate.Title == "" {
		SendErrorResponse(res, "не указан заголовок задачи")
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
		SendErrorResponse(res, "ошибка при обновлении данных")
		return
	}

	if taskCurrent == (Task{}) {
		SendErrorResponse(res, "не указан идентификатор")
		return
	}

	SendJSONResponse(res, taskCurrent)
}

// TaskDone обработчик для POST-запроса /api/task/done, который делает
// задачу выполненой. Одноразовая задача с пустым полем repeat удаляется.
// Возвращает пустой JSON или ошибку
func TaskDone(res http.ResponseWriter, req *http.Request) {
	// Получаем id из url-пути
	id := req.URL.Query().Get("id")

	if id == "" {
		SendErrorResponse(res, "не указан идентификатор")
		return
	}

	// Создаем экзмплер структуры типа Task и заполняем его поля
	// значениями из таблицы с указнным id
	var taskCurrent Task

	err := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id)).
		Scan(&taskCurrent.ID, &taskCurrent.Date, &taskCurrent.Title, &taskCurrent.Comment, &taskCurrent.Repeat)
	if err != nil {
		SendErrorResponse(res, "задача не найдена")
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
		SendErrorResponse(res, "ошибка при получении даты")
		return
	}

	// Приводим время к типу int и обновляем строку в таблице scheduler
	dateInt, err := strconv.Atoi(date)
	if err != nil {
		SendErrorResponse(res, fmt.Sprintf("ошибка при преобразовании: %v", err))
	}

	_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
		sql.Named("date", dateInt),
		sql.Named("id", id))

	if err != nil {
		SendErrorResponse(res, "ошибка при обновлении данных")
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
		SendErrorResponse(res, "ошибка при удалении задачи")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		SendErrorResponse(res, "ошибка при получении количества удалённых строк")
		return
	}

	if rowsAffected == 0 {
		SendErrorResponse(res, "задача с таким id не найдена")
		return
	}

	SendJSONResponse(res, struct{}{})
}

// SignIn обработчик POST-запроса /api/signin. Получает JSON с полем password.
// Если пароль совпадает, возвращает формирует JWT и передает его в поле JSON-объекта.
// Если пароль невернный или произошла ошибка, возвращает JSON с текстом ошибки
func SignIn(res http.ResponseWriter, req *http.Request) {
	var p Password

	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) > 0 { // Получаем пароль из JSON
		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			SendErrorResponse(res, fmt.Sprintf("ошибка десериализации JSON: %v", err))
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
			response := Response{Token: signedToken}
			SendJSONResponse(res, response)
		} else {
			SendErrorResponse(res, "невереный пароль")
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
				SendErrorResponse(res, fmt.Sprint("токен не валиден"))
				return
			}

			// Получаем хэш из полезной нагрузки токена
			result, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				SendErrorResponse(res, fmt.Sprint("не удалось выполнить приведение типа к  jwt.MapClaims"))
				return
			}

			hashRow := result["hash"]

			hash, ok := hashRow.(string)
			if !ok {
				SendErrorResponse(res, fmt.Sprint("не удалось выполнить приведение типа к string"))
				return
			}

			// Получаем хэш, как в функции SignIn
			hashedPass := sha256.Sum256([]byte(pass))
			expectedHash := hex.EncodeToString(hashedPass[:])

			if hash != expectedHash {
				SendErrorResponse(res, fmt.Sprint("ошибка валидации хэша"))
				return
			}
		}
		next(res, req)
	})
}
