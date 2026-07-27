// package db

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"

// 	_ "modernc.org/sqlite"
// )

// const Schema string = `CREATE TABLE scheduler (
//     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     date CHAR(8) NOT NULL DEFAULT "",
//     title VARCHAR(20) NOT NULL DEFAULT "",
//     comment TEXT NOT NULL DEFAULT "",
//     repeat VARCHAR(128) NOT NULL DEFAULT ""
// );

// CREATE INDEX idx_date ON scheduler (date);`

// // CREATE TABLE IF NOT EXISTS - на каждый запуск

// func Init(dbFile string) error {
// 	var install bool

// 	_, err := os.Stat(dbFile)
// 	if err != nil {
// 		install = true
// 	}

// 	dbase, err := sql.Open("sqlite", dbFile)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	defer dbase.Close()

// 	if install {
// 		_, err = dbase.Exec(Schema)
// 		if err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 	}
// 	return err
// }

package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

const Schema string = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(20) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
    
CREATE INDEX idx_date ON scheduler (date);`

func Init(dbFile string) error {

	dbase, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer dbase.Close()

	_, err = dbase.Exec(Schema)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
