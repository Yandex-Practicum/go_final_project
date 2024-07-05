package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Masteker/go_final_project/api"
	"github.com/Masteker/go_final_project/database/auth"
	"github.com/Masteker/go_final_project/database/db"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Загружаем переменные среды
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	// .env сам подгружается если мы используем docker compose для запуска, но для тестов удобнее запускать код напрямую, поэтому оставил godotenv

	dbFile := os.Getenv("TODO_DBFILE")

	// Если бд не существует, создаём
	if !db.DbExists(dbFile) {
		err = db.InstallDB()
		if err != nil {
			log.Println(err)
		}
	}

	// Запуск бд
	dbStorage, err := db.StartDB()
	defer func() {
		err := dbStorage.CloseDB()
		if err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		log.Fatal(err)
	}
	api.ApiInit(dbStorage)

	// Адрес для запуска сервера
	ip := ""
	port := os.Getenv("TODO_PORT")
	addr := fmt.Sprintf("%s:%s", ip, port)

	// Router
	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	r.Get("/api/nextdate", api.GetNextDateHandler)
	r.Get("/api/tasks", auth.Auth(api.GetTasksHandler))
	r.Post("/api/task/done", auth.Auth(api.PostTaskDoneHandler))
	r.Post("/api/signin", auth.Auth(api.PostSigninHandler))
	r.Handle("/api/task", auth.Auth(api.TaskHandler))

	// Запуск сервера
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Server running on %s\n", port)

}
