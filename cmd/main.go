package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	utils "github.com/falsefood/go_final_project/internal"
	"github.com/falsefood/go_final_project/internal/handlers"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

var config utils.Config

func main() {
	// Закрепляем путь к БД, фронту и порт
	config = utils.Config{
		WebDir: "./web",
		Port:   ":7540",
		DBPath: "./scheduler.db",
	}

	// Подключаем БД или создаём её
	appPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Не удалось получить путь приложения: %v\n", err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), config.DBPath)
	db, err := utils.OpenDB(dbFile)
	if err != nil {
		log.Fatalf("Ошибка открытия базы данных: %v\n", err)
	}
	defer db.Close()

	// Проверяем существование таблицы и создаем её, если не существует
	exists, err := utils.TableExists(db, "scheduler")
	if err != nil {
		log.Fatalf("Ошибка проверки существования таблицы: %v\n", err)
	}
	if !exists {
		if err := utils.CreateTable(db); err != nil {
			log.Fatalf("Ошибка создания таблицы: %v\n", err)
		}
	}

	// Настройка маршрутов
	router := setupRouter()

	// Запуск сервера
	log.Printf("Запуск сервера на порту %s...\n", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))
}

// setupRouter настраивает маршруты для приложения
func setupRouter() *chi.Mux {
	router := chi.NewRouter()

	// Маршрут для статических файлов (если нужно)
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(config.WebDir, strings.TrimPrefix(r.URL.Path, "/"))
		http.ServeFile(w, r, filePath)
	})

	// API маршруты
	router.Get("/api/nextdate", handlers.NextDateHandler)
	router.Route("/api/task", func(r chi.Router) {
		r.Post("/", utils.WithDB(handlers.CreateTaskHandler))   // POST для создания задачи
		r.Get("/", utils.WithDB(handlers.GetTaskHandler))       // GET для получения задачи (с query-параметром id)
		r.Get("/{id}", utils.WithDB(handlers.GetTaskHandler))   // GET для получения задачи (с параметром в URL)
		r.Delete("/", utils.WithDB(handlers.DeleteTaskHandler)) // DELETE для удаления задачи
		r.Put("/", utils.WithDB(handlers.UpdateTaskHandler))    // PUT для обновления задачи
	})
	router.Get("/api/tasks", utils.WithDB(handlers.GetTasksHandler))          // GET для получения списка задач
	router.Post("/api/task/done", utils.WithDB(handlers.MarkTaskDoneHandler)) // POST для отметки задачи как выполненной

	return router
}
