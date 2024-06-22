package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Загружаем переменные среды
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	// Если бд не существует, создаём
	if !dbExists() {
		installDB()
	}

	// Адрес для запуска сервера
	ip := ""
	port := os.Getenv("TODO_PORT")
	addr := fmt.Sprintf("%s:%s", ip, port)

	r := chi.NewRouter()

	r.Get("/api/nextdate", getNextDate)

	// Запуска сервера
	fmt.Println("Запускаем сервер")
	err = http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
	NextDate(time.Now(), "", "")

	fmt.Println("Завершаем работу")
}
