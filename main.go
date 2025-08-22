package main

import (
	"fmt"
	dbase "go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	var err error
	dbFile := "scheduler.db"

	err = dbase.Init(dbFile)
	if err != nil {
		fmt.Println(err)
	}

	defer dbase.DB.Close()

	err = server.Run()
	if err != nil {
		panic(err)
	}
}
