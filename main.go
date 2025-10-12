package main

import (
	"log"
	"net/http"

	"go_final_project/pkg/api"
	"go_final_project/pkg/auth"
	"go_final_project/pkg/db"
)

func main() {
	// Инициализируем конфигурацию
	api.InitConfig()
	config := api.GetConfig()

	// используем конфиг для порта и БД
	port := config.Port
	dbFile := config.DBfile

	// Инициализируем базу данных
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	defer db.Close()

	// Инициализируем API обработчики
	api.Init()

	// Настраиваем файловый сервер с проверкой аутентификации для главной страницы
	webDir := "./web"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && auth.IsAuthEnabled() {
			// Проверяем есть ли токен в куках
			cookie, err := r.Cookie("token")
			if err != nil {
				// Нет куки - перенаправление на логин
				http.Redirect(w, r, "/login.html", http.StatusFound)
				return
			}

			// Проверяем валидность токена
			valid, err := auth.ValidateToken(cookie.Value)
			if err != nil || !valid {
				// Невалидный логин - перенаправляем на логин
				http.Redirect(w, r, "/ligin.html", http.StatusFound)
				return
			}
		}
		// Во всех остальных случаях отдаем статические файлы
		http.FileServer(http.Dir(webDir)).ServeHTTP(w, r)
	})

	// Запускаем сервер
	log.Printf("Сервер запущен на порту %s", port)
	log.Printf("База данных: %s", dbFile)
	if auth.IsAuthEnabled() {
		log.Printf("Аутентификация ВКЛЮЧЕНА")
	} else {
		log.Printf("Аутентификация ОТКЛЮЧЕНА")
	}
	log.Printf("Откройте http://localgost:%s в браузере", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

}
