// database/database.go
package database

import (
	"errors"
	"log"
	"os"

	"github.com/Xa4y322/go_final_project/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// type Dbinstance struct {
// 	Db *gorm.DB
// }

var Db *gorm.DB

func getDBPath() string {
	const DBFile = "scheduler.db"
	if val, exists := os.LookupEnv("TODO_DBFILE"); exists {
		return val
	}
	return DBFile
}

func ConnectDB() {
	// создаём подключение к базе данных.
	// В &gorm.Config настраивается логер,
	// который будет сохранять информацию
	// обо всех активностях с базой данных.
	db, err := gorm.Open(sqlite.Open(getDBPath()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database.\n", err)
		os.Exit(1)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migration")
	db.AutoMigrate(&models.Task{})

	// DB = Dbinstance{
	// 	Db: db,
	// }
	Db = db
}

func ReadTasks() ([]models.Task, error) {
	var tasks []models.Task
	result := Db.Order("date").Limit(30).Find(&tasks)
	return tasks, result.Error
}

func ReadTaskByID(id uint) (models.Task, error) {
	var tasks models.Task
	result := Db.Limit(1).Find(&tasks, models.Task{ID: id})
	if result.RowsAffected == 0 {
		return tasks, errors.New("No record found")
	}
	return tasks, result.Error
}
