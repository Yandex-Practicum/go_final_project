package service

import "os"

func PortGlobal() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		return "7540"
	}
	return port
}

func PathDbManual() string {
	path := os.Getenv("TODO_DBFILE")
	if path == "" {
		return "scheduler.db"
	}
	return path
}
