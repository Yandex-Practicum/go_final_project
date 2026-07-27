package service

import (
	"os"
	"path/filepath"
)

func PortGlobal() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		return "7540"
	}
	return port
}

func PathDbManual() string {
	envPath := os.Getenv("TODO_DBFILE")
	if envPath != "" {
		return envPath
	}

	dir, err := os.Getwd()
	if err != nil {
		return filepath.Join("pkg", "db", "scheduler.db")
	}

	if filepath.Base(dir) == "tests" {
		return filepath.Join(dir, "..", "pkg", "db", "scheduler.db")
	}

	return filepath.Join(dir, "pkg", "db", "scheduler.db")
}
