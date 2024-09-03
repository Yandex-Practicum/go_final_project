package main

import (
	"go_final_project/steps"
	"os"

	_ "modernc.org/sqlite"
)

const Port = 7540

func main() {
	//mux := http.NewServeMux()
	//mux.Handle("/", http.FileServer(http.Dir("./web")))

	//сделать взятие значения порта из переменной окружения
	//err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", tests.Port), nil)
	//if err != nil {
	//log.Fatal(err)
	//}
	os.Setenv("TODO_PASSWORD", "0330")
	steps.StartServer()
	steps.CreateDB()

	//db, err := sql.Open("sqlite", "scheduler.db")
	//if err != nil {
	//log.Fatalf("Error while opening database: %v", err)
	//}
	//defer db.Close()
	//router := chi.NewRouter()

	//fileServer := http.FileServer(http.Dir("./web"))
	//router.Handle("/*", http.StripPrefix("/", fileServer))
	//router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//fileServer.ServeHTTP(w, r)
	//})

	//router.Route("/api", func(r chi.Router) {
	//r.Get("/nextdate", nextDateWMHandler)

	//r.Route("/task", func(rt chi.Router) {
	//rt.Get("/", getTaskHandler(db))
	//rt.Post("/", addTaskHandler(db))
	//rt.Post("/done", doneTaskHandler(db))
	//rt.Delete("/", deleteTaskHandler(db))
	//rt.Put("/", updateTaskHandler(db))

	//})

	//r.Route("/tasks", func(rts chi.Router) {
	//rts.Get("/", getTasksHandler(db))
	//})

	//})

	//log.Println("Run on port:", Port)

	//err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", Port), router)
	//if err != nil {
	//log.Fatal(err)
	//}

}
