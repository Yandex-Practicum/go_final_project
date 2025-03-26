package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/qavaleria/go_final_project/models"
	"github.com/qavaleria/go_final_project/tasks"
	"log"
	"strings"
	"time"
)

func GetTasks(db *sqlx.DB, search string, limit int) ([]models.Task, error) {
	var rows *sqlx.Rows
	var err error

	// Убираем лишние пробелы в поисковой строке
	search = strings.TrimSpace(search)
	if search != "" {
		// Проверяем, является ли поисковая строка датой в формате "02.01.2006"
		if searchDate, err := time.Parse("02.01.2006", search); err == nil {
			searchDateStr := searchDate.Format(tasks.FormatDate)
			rows, err = db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :searchDate ORDER BY date LIMIT :limit`, map[string]interface{}{
				"searchDate": searchDateStr,
				"limit":      limit,
			})
			if err != nil {
				log.Printf("Ошибка выполнения запроса с датой: %v", err)
				return nil, err
			}
		} else {
			//search = "%" + strings.ToLower(search) + "%"
			rows, err = db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler WHERE LOWER(title) LIKE :search OR LOWER(comment) LIKE :search ORDER BY date LIMIT :limit`, map[string]interface{}{
				"search": "%" + search + "%",
				"limit":  limit,
			})
			if err != nil {
				log.Printf("Ошибка выполнения запроса с поиском: %v", err)
				return nil, err
			}
		}
	} else {
		// Если поисковая строка пустая, просто выбираем все задачи
		rows, err = db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`, map[string]interface{}{
			"limit": limit,
		})
		if err != nil {
			log.Printf("Ошибка выполнения запроса без поиска: %v", err)
			return nil, err
		}
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Ошибка закрытия rows: %v", err)
		}
	}()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.StructScan(&task)
		if err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Ошибка чтения строк: %v", err)
		return nil, err
	}

	// Возвращаем пустой массив задач, если ни одной задачи не найдено
	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}
