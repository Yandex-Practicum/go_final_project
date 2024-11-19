package migration

import (
	"final_project/internal/db"
	"final_project/internal/repository"
)

func Migration() *repository.Repository {
	database := db.CreateDB()
	rep := repository.New(database)
	if db.Install {
		rep.CreateScheduler()
		rep.CreateSchedulerIndex()
	}
	return rep
}
