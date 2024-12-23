package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

func SetupRouter() http.Handler {
	// Настройка роутера
	r := mux.NewRouter()
	//r.HandleFunc("/auth/register", handlers.RegisterHandler).Methods("POST")
	//r.HandleFunc("/auth/login", handlers.LoginHandler).Methods("POST")

	// Обработка статических файлов
	staticDir := "./web/"
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	return r
}
