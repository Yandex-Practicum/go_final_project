package db

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"os"
	"path/filepath"

	"github.com/askoren1/go_final_project/internal/repository"
)

//  Инициализациия базы данных
func Migration(repo *repository.Repository) { //функция для создания таблицы в базе данных, если она еще не существует
	appPath, err := os.Executable() //Получаем путь к исполняемому файлу приложения
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db") // Конструируем полный путь к файлу БД scheduler.db

	_, err = os.Stat(dbFile) //Проверяем, существует ли файл базы данных по указанному пути
	var install bool
	if err != nil {
		install = true
	}

	if install {
		if err := repo.CreateScheduler(); err != nil { //Вызываем метод CreateScheduler() у репозитория для создания таблицы
			log.Fatal(err)
		}
	} else {
		fmt.Println("База данных уже существует.")
	}
}
