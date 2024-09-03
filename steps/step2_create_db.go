package steps

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func CreateDB() {
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	file := os.Getenv("DBFILE")
	var dbfile string

	if len(file) > 0 {
		dbfile = file
	} else {
		dbfile = os.Getenv("DBFILE")
	}

	dbFile := filepath.Join(filepath.Dir(path), dbfile)
	_, err = os.Stat(dbFile)

	var success bool

	if err != nil {
		success = true
	}

	if success {
		db, err := sql.Open("sqlite", dbFile)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer db.Close()

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler
		(id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8) NOT NULL DEFAULT '',
		  title VARCHAR(128) NOT NULL DEFAULT '', 
		  comment VARCHAR(256) NOT NULL DEFAULT '',
		  repeat VARCHAR(128) NOT NULL DEFAULT '')`, `CREATE INDEX date_index ON scheduler(date)`)

		if err != nil {
			log.Fatal(err)
		}
	}

}
