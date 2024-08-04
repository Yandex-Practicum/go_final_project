package config

import (
	"log"
	"os"
	"path/filepath"
)

const (
	stdPort    = "7540"
	varEnvPort = "TODO_PORT"

	varEnvDBFile = "TODO_DBFILE"
	dbPath       = "./database"
	dbName       = "scheduler.db"
)

func Port() string {
	port := os.Getenv(varEnvPort) //Getenv извлекает значение переменной окружения, названной ключом
	if port == "" {
		port = stdPort
	}

	log.Printf(`Retrieved port %s from env variable "%s"`, port, varEnvPort)

	return ":" + port
}

func DbPath() string {
	storagePath := os.Getenv(varEnvDBFile)

	if storagePath == "" {
		storagePath = filepath.Join(dbPath, dbName)
		log.Printf(`Database storage address: %s`, storagePath)
	} else {
		log.Printf(`Database storage address %s extract from env variable "%s" `,
			storagePath,
			varEnvDBFile)
	}

	return storagePath
}
