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

func AutherENV() string {
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		return ""
	}
	return pass
}
