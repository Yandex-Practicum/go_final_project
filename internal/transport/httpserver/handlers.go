package httpserver

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ASHmanR17/go_final_project/internal/database"
	"github.com/ASHmanR17/go_final_project/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

const (
	DateLayout = "20060102"
	JsonValue  = "application/json; charset=UTF-8"
)

// создаём секретный ключ для подписи.
var secret = "my_secret_key"
var req struct {
	Password string `json:"password"`
}

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
	// Получаем параметр search из строки запроса
	search := r.URL.Query().Get("search")

	tasks, err := h.taskService.GetTasks(search)
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

	id := r.URL.Query().Get("id")

	err := h.taskService.DoneTask(id)
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

func (h httpHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	// Получаем пароль из тела запроса

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response := map[string]string{
			"error": fmt.Sprintf("Неверный формат запроса"),
		}
		w.Header().Set("Content-Type", JsonValue)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	// Сравниваем пароль с хранимым в переменной окружения
	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}
	if req.Password == pass {
		fmt.Println("Пароль введен верно")
	}
	// проверяем cookie
	cookie, err := r.Cookie("token")
	if err == nil {
		tokenString := cookie.Value
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{
			"pwd": hashPassword(req.Password),
		}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err == nil && token.Valid { // Возвращаем токен в формате JSON
			response := map[string]string{
				"token": tokenString,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Создаем JWT-токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pwd": hashPassword(req.Password),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	// Возвращаем токен в формате JSON
	response := map[string]string{
		"token": tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func hashPassword(password string) string {
	// Создаем хэш пароля
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtCook string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			fmt.Println(cookie)
			if err == nil {
				jwtCook = cookie.Value
			}
			// здесь код для валидации и проверки JWT-токена
			// Проверяем токен на валидность
			token, err := jwt.ParseWithClaims(jwtCook, jwt.MapClaims{
				"pwd": hashPassword(req.Password),
			}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("неверный метод подписи: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				fmt.Println("Authentification required", err)
				return
			}

		}
		next(w, r)
	})
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
