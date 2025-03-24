// Шаг 1 ТЗ - в main.go обработчики + инициализация
// Шаг 2 ТЗ - в db.go все про БД + func GetTaskByID, func UpdateTask, func DeleteTaskByID
// Шаг 3-4 ТЗ - в nextdate.go func NextDate с правилами повторения + func AddTaskHandler добавление задач
// Шаг 5 ТЗ - в task.go func GetTasksHandler получение задач
// Шаг 6 ТЗ - в getPut.go func TaskHandler(получения и обновления задач) + func getTask(по идентификатору) + func updateTask
// Шаг 7 ТЗ - в doneDel func MarkTaskDone отметка о выполнении + func DeleteTask удаление
// все тесты проходят

package main

import (
	"log"
	"net/http"

	"go_final_project/api"
	"go_final_project/db"
)

func main() {
	const port = ":7540"
	const webDir = "./web"

	// Инициализируем базу данных
	database, err := db.InitDB(false)
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}
	defer database.Close()

	log.Println("База данных успешно инициализирована.")

	// Настраиваем файловый сервер для обслуживания статических файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Регистрация обработчиков для API
	http.HandleFunc("/api/nextdate", api.NextDateHandler) // Используем обработчик из nextdate.go
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		api.GetTasksHandler(w, r, database) // Обработчик для получения задач
	})
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.AddTaskHandler(w, r, database) // Обработчик для добавления задач
		} else if r.Method == http.MethodDelete {
			api.DeleteTask(w, r) // Обработчик для удаления задачи
		} else {
			api.TaskHandler(w, r) // Обработчик для получения и обновления задач
		}
	})
	http.HandleFunc("/api/task/", api.TaskHandler)      // Обработчик для получения и обновления задач
	http.HandleFunc("/api/task/done", api.MarkTaskDone) // Регистрация обработчика для отметки о выполнении задачи

	log.Printf("Сервер запущен на http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
