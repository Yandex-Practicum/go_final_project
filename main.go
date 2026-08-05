package main

import (
	"proj/pkg/api/service"
	"proj/pkg/db"
	"proj/pkg/server"
)

func main() {
	db.Init(service.PathDbManual())
	server.ServerRun()
}
