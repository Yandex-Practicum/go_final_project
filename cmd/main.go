package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"github.com/AsyaBiryukova/go_final_project/api"
	db "github.com/AsyaBiryukova/go_final_project/database"
)

func main() {
	//Подгружаем переменные среды
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	// Создаём бд, если её нет
	if !db.DbExists() {
		db.InstallDB()
	}

	// Запускаем бд
	db.StartDB()
	defer db.DB.Close()

	//Адрес для запуска сервера
	ip := ""
	port := os.Getenv("TODO_PORT")
	addr := fmt.Sprintf("%s:%s", ip, port)

	//Роутер
	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", api.GetNextDateHandler)
	r.Get("/api/tasks", auth(api.GetTasksHandler))
	r.Post("/api/task/done", auth(api.PostTaskDoneHandler))
	r.Post("/api/signin", api.PostSigninHandler)
	r.Handle("/api/task", auth(api.TaskHandler))

	//Запускаем сервер

	err = http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// смотрим наличие пароля

		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			err := api.GetAndVerifyToken(r)
			if err != nil {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})

}
