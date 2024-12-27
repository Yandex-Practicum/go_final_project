package main

import (
	"net/http"
)

func main() {

	appPath, err := os.Getwd() 
	if err != nil {
		log.Fatal(err)
	}
	println(appPath)

	dbFile := filepath.Join(appPath, "appscheduler.db")
	_, err := os.Stat(dbFile)
	println(dbFile)

	ver install bool
	if err != nil {
		install = true
	}

	if install == true { //если файл БД не существует, то создаем его

		_, err := os.Create(dbFile) //создаем файл БД
		if err != nil {
			log.Fatal(err) 
		}

		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY, date TEXT, title TEXT, comment TEXT, repeat TEXT(128), status TEXT)") // создаем таблицу с данными
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_date ON scheduler (date)")
		if err != nil {
			log.Fatal(err)
		}
	
	}

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}

}
