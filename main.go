package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AlexeyVilkov/go_final_project/database"
	"github.com/AlexeyVilkov/go_final_project/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// подключаемся к БД
	_, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// отключаемся по завершении функции
	defer database.Db.Close()

	// устанавливаем значение переменной окружения TODO_PORT для проверки
	//os.Setenv("TODO_PORT", "8080")

	port := ":7540" //значение по умолчанию, из tests: strconv.Itoa(tests.Port)

	if value, exists := os.LookupEnv("TODO_PORT"); exists {
		port = ":" + value
	}

	r := http.FileServer(http.Dir("web"))
	http.Handle("/", r)
	http.HandleFunc("/api/nextdate", handlers.GetNextDate)
	http.HandleFunc("/api/task", handlers.ActionTask)
	http.HandleFunc("/api/tasks", handlers.GetTasks)
	http.HandleFunc("/api/task/done", handlers.DoneTask)

	fmt.Println("Server is listening...", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
