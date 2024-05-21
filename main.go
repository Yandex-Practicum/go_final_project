package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

// Ошибки
var (
	ErrWrongFormat   = errors.New("wrong format")
	ErrExceededLimit = errors.New("limit exceeded")
	ErrInvalidValue  = errors.New("invalid value")
	ErrNoTitle       = errors.New("no title")
	ErrNoID          = errors.New("no such id")
	ErrTaskNotFound  = errors.New("task not found")
	ErrWrongPassword = errors.New("wrong password")
)

const format = "20060102"

func init() {
	// поддержка .env файла
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		serverPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(filepath.Dir(serverPath), "scheduler.db")
	}
	log.Println("Using datatbase :" + dbFile)
	dbExists := true
	if _, err := os.Stat(dbFile); err != nil {
		dbExists = false
	}
	//если базы данных нет создаем новую с таблицей и индексом
	if !dbExists {
		log.Println("Creating table... ")
		createSchedulerTable()
		time.Sleep(2 * time.Second)
	}
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(static.Serve("/", static.LocalFile("./web", false)))
	router.GET("/api/nextdate", NextDateHandler) // обработчик для тестов NextDate()
	router.POST("/api/signin", authHandler)
	//обработчики, для которых  требуется аутентификаия
	reqAuth := router.Group("/api")
	reqAuth.Use(auth())
	reqAuth.GET("/tasks", getRecentTasksHandler)
	reqAuth.GET("/task", getTaskByIDHandler)
	reqAuth.PUT("/task", updateTaskHandler)
	reqAuth.POST("/task", postTaskHandler)
	reqAuth.POST("/task/done", doneTaskHandeler)
	reqAuth.DELETE("/task", deleteTaskHandler)
	//проверка наличия Env var TODO_PORT
	port, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		//если такой переменной нет присваиваем порт по умолчанию (7540)
		port = "7540"
	}
	port = ":" + port
	if err := router.Run(port); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
