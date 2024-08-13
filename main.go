package main

import (
	"fmt"
	"net/http"

	"github.com/chichimba/go_final_project/tests"
)

func mainHandle(res http.ResponseWriter, req *http.Request) {
	// лог-контроль
	fmt.Println("Получен запрос")
	// запрос в строку
	s := fmt.Sprintf("Host: %s\nPath: %s\nMethod: %s", req.Host, req.URL.Path, req.Method)
	// лог-контроль
	fmt.Println(s)
	// отправка клиенту
	res.Write([]byte(s))
}

func main() {
	// лог-контроль
	fmt.Println("Запускаем сервер")
	mux := http.NewServeMux()
	// вешаем обработчик
	mux.HandleFunc("/main", mainHandle)
	// webDir := "static"
	webDir := ""
	mux.Handle("/*", http.FileServer(http.Dir(webDir)))
	// переменая tests.Port из settings.go
	// !!! заменить на переменную окружения TODO_PORT с помощью os.Getenv(key string) string
	portStr := fmt.Sprintf(":%d", tests.Port)
	// лог-контроль
	fmt.Println(portStr)
	err := http.ListenAndServe(portStr, mux)
	// err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Завершаем работу")
}
