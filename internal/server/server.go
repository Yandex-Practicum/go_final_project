package server

import (
	"database/sql"
	"fmt"
	"go-task-manager-final_project/internal/api/handlers"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	defaultPort      = "7540"  // Порт по умолчанию для запуска сервера
	defaultStaticDir = "./web" // Директория со статическими файлами по умолчанию
	minPort          = 1       // Минимально допустимый номер порта
	maxPort          = 65535   // Максимально допустимый номер порта
)

// GetPort возвращает номер порта из переменной окружения TODO_PORT или значение по умолчанию.
// Проверяет корректность формата и диапазона значения порта.
// Возвращает:
// - int: номер порта (в диапазоне [minPort, maxPort]);
// - error: ошибка, если порт невалидный (не число или вне диапазона).
func GetPort() (int, error) {
	portStr := os.Getenv("TODO_PORT")
	if portStr == "" {
		// Если переменная окружения не задана, используем порт по умолчанию
		portStr = defaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		// Если порт не является числом
		return 0, fmt.Errorf("invalid port format: %s", portStr)
	}
	// Если порт за пределами допустимого диапазона
	if port < minPort || port > maxPort {
		return 0, fmt.Errorf("port out of range [%d, %d]: %d", minPort, maxPort, port)
	}

	// Иначе, возвращаем корректный номер порта, на котором стартуем
	return port, nil
}

// GetStaticDir возвращает путь к директории со статическими файлами.
// Берёт значение из переменной окружения TODO_STATIC_DIR, если она задана.
// Иначе использует значение по умолчанию (defaultStaticDir).
// Возвращает: строку - путь к директории со статическими файлами.
func GetStaticDir() (string, error) {
	dir := os.Getenv("TODO_STATIC_DIR")
	if dir == "" {
		// Если переменная окружения не задана, используем директорию по умолчанию
		dir = defaultStaticDir
	}

	// Проверяем существование директории
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", fmt.Errorf("static directory not found: %s", dir)
	}
	return dir, nil
}

// SetupStaticFileRouting настраивает роутинг для статических файлов в роутере chi.Mux.
// Проверяет существование директории, создаёт файловый сервер и регистрирует обработчик.
// Параметры:
// - r *chi.Mux: роутер chi, в который добавляется обработка статических файлов.
// Возвращает:
// - error: ошибка, если директория не найдена или возникла проблема при настройке.
func SetupStaticFileRouting(r *chi.Mux) error {
	// Получаем путь к директории со статическими файлами
	staticDir, err := GetStaticDir()
	if err != nil {
		return err
	}

	// Создаём файловый сервер для статических файлов
	fs := http.FileServer(http.Dir(staticDir))

	// Настраиваем роутинг: все запросы перенаправляются на статические файлы
	// (префикс "/" удаляется из пути - это позволяет корректно обрабатывать запросы к файлам)
	r.Handle("/*", http.StripPrefix("/", fs))
	log.Printf("Роутинг настроен для статических файлов из %s", staticDir)

	return nil
}

// StartServer запускает HTTP-сервер с заданной конфигурацией.
// Настраивает роутер, подключает обработчики, устанавливает таймауты и запускает сервер.
// Параметры:
// - db *sql.DB: подключение к базе данных, передаваемое обработчикам.
// Возвращает:
// - error: ошибка при конфигурации или запуске сервера (включая проблемы с портом, статикой и тд.).
func StartServer(db *sql.DB) error {
	// Создаём новый роутер chi
	router := chi.NewRouter()

	// Настраиваем обработку статических файлов
	err := SetupStaticFileRouting(router)
	if err != nil {
		return fmt.Errorf("failed to setup static file routing: %w", err)
	}

	// Регистрируем API-обработчики, передавая роутер и подключение к БД
	handlers.Init(router, db)

	// Получаем номер порта для запуска сервера
	port, err := GetPort()
	if err != nil {
		return fmt.Errorf("failed to get port: %w", err)
	}

	// Формируем адрес для прослушивания (например, ":7540")
	address := fmt.Sprintf(":%d", port)

	// Создаём конфигурацию HTTP-сервера
	server := &http.Server{
		Addr:         address,           // Адрес и порт для прослушивания
		Handler:      router,            // Обработчик запросов - наш роутер chi
		ReadTimeout:  5 * time.Second,   // Таймаут на чтение запроса
		WriteTimeout: 10 * time.Second,  // Таймаут на отправку ответа
		IdleTimeout:  120 * time.Second, // Таймаут для неактивных соединений
	}

	// Логируем запуск сервера
	log.Printf("Сервер запущен на http://localhost:%d", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// Логируем ошибку запуска и возвращаем ошибку запуска сервера
		log.Printf("Ошибка при запуске сервера: %v", err)
		return fmt.Errorf("server failed to listen and serve: %w", err)
	}
	return nil
}
