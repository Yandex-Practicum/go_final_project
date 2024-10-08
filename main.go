package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "go_final_project/nextdate"

	_ "modernc.org/sqlite"
)

func main() {

	// создаем БД и таблицу
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("База данных будет создана")
			install = true
		} else {
			log.Println("не получилось проверить файл")
			log.Fatal(err)
		}
	}
	log.Println("База данных была создана ранее")

	if install {
		Table := `CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT, 
			repeat VARCHAR(128) NOT NULL
			);`
		_, err = db.Exec(Table)
		if err != nil {
			log.Fatal(err)
		}

		Index := `CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);`
		_, err = db.Exec(Index)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Определение порта при запуске сервера
	port := "7540"
	env := os.Getenv("TODO_PORT")
	if len(env) != 0 {
		port = env
	}
	port = ":" + port

	// Запускаем веб-сервер
	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", NextDateHandler)

	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}
