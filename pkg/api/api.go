package api

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"go_final_project/pkg/db"
	"net/http"
	"os"
	"time"
)

// Регистрация маршрутов
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
	http.HandleFunc("/api/signin", signinHandler)

}
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		writeError(w, "Метод не поддерживается")
		return
	}

	//Читаем id задачи
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор")
		return
	}

	//Получение задачи из БД
	t, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	//Повторения нет - удаляем
	if t.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err.Error())
			return
		}
		w.Write([]byte(`{}`))
		return
	}

	//Парсим дату
	prev, err := time.Parse("20060102", t.Date)
	if err != nil {
		writeError(w, "Неверный формат даты")
		return
	}

	nextDate, err := NextDate(prev, t.Date, t.Repeat)
	if err != nil {
		writeError(w, "Неверное правило повторения")
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		writeError(w, err.Error())
		return
	}

	//Возвращаем пустого Джейсона)
	w.Write([]byte(`{}`))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		task, err := db.GetTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, task)
	case http.MethodDelete:

		id := r.URL.Query().Get("id")

		if err := db.DeleteTask(id); err != nil {

			writeJson(w, map[string]string{"error": err.Error()})
		} else {

			w.Write([]byte(`{}`))
		}
	case http.MethodPut:
		var t db.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			writeError(w, "Неверный формат json")
			return
		}

		//Проверяем ID и заголовок задачи)
		if t.ID == "" {
			writeError(w, "Поле ID не может быть пустым")
			return
		}
		if t.Title == "" {
			writeError(w, "Поле Title не может быть пустым")
			return
		}

		//Проверяем срок годности
		layout := "20060102"
		if _, err := time.Parse(layout, t.Date); err != nil {
			writeError(w, "Неверный формат даты")
			return
		}
		if t.Repeat != "" {
			if _, err := NextDate(time.Now(), t.Date, t.Repeat); err != nil {
				writeError(w, "Неверное правило повторения")
				return
			}
		}

		//Исправляем дату если время пролетело
		parsed, _ := time.Parse(layout, t.Date)
		if t.Repeat != "" {
			nextDate, _ := NextDate(time.Now(), t.Date, t.Repeat)
			if parsed.Before(time.Now()) {
				t.Date = nextDate
			}
		} else if parsed.Before(time.Now()) {
			t.Date = time.Now().Format(layout)
		}

		//сохраняем
		if err := db.UpdateTask(&t); err != nil {
			writeError(w, err.Error())
			return
		}
		w.Write([]byte(`{}`))
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработчик NextDate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат now", http.StatusBadRequest)
			return
		}
	}

	//Без даты - не пускаем
	if dstart == "" {
		http.Error(w, "Поле date не может быть пустым", http.StatusBadRequest)
		return
	}

	//Без правила повторения - тоже
	if repeat == "" {
		http.Error(w, "Поле repeat не может быть пустым", http.StatusBadRequest)
		return
	}

	//Вычисление следующей даты
	result, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(result))
}

func signinHandler(w http.ResponseWriter, r *http.Request) {

	//Читаем секрет
	pass := os.Getenv("TODO_PASSWORD")

	//Если пустой - лови 403
	if pass == "" {
		http.Error(w, "Authentication disabled", http.StatusForbidden)
		return
	}

	var req struct {
		Password string `json:"password"`
	}

	//Декодируем
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//Неправильный пароль - держи 401
	if req.Password != pass {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный пароль"})
		return
	}

	//Правильный пароль - собираем токен Джейсона
	jwtKey := []byte(pass)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pwd_hash": req.Password,
		"exp":      jwt.NewNumericDate(jwt.TimeFunc().Add(8 * 3600e9)),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
