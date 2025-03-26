// package tasks provides tools for working on scheduler tasks
package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"go_final_project/constants"
	"go_final_project/tasks"
)

const TaskLimit = 50

// GetSearchQuery get search query and provide SQL query for search execution
func (d Database) GetTasks(searchStr string) (taskList []tasks.Task, err error) {

	var t tasks.Task
	var searchParam string
	var query string
	var rows *sql.Rows

	log.Printf("searchStr:%v", searchStr)
	searchDate, err := time.Parse("02.01.2006", searchStr)
	switch {
	case searchStr != "" && err == nil:
		// get tasks by date
		log.Println("searchDate")
		searchParam = searchDate.Format(constants.DateFormat)
		log.Printf("searchParam:%v", searchParam)
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = d.db.Query(query, searchParam, TaskLimit)
		if err != nil {
			return taskList, fmt.Errorf("failed to get task by date: %w", err)
		}
	case searchStr != "" && err != nil:
		// get tasks by title/comment
		searchParam = "%" + searchStr + "%"
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		rows, err = d.db.Query(query, searchParam, searchParam, TaskLimit)
		if err != nil {
			return taskList, fmt.Errorf("failed to get task by content: %w", err)
		}
	default:
		// get tasks w/o condition
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
		rows, err = d.db.Query(query, TaskLimit)
		if err != nil {
			return taskList, fmt.Errorf("failed to get tasks: %w", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return taskList, fmt.Errorf("failed to scan tasks: %w", err)
		}
		t.Id = fmt.Sprint(id)
		taskList = append(taskList, t)
		log.Printf("Task: %v, %v", t.Title, t.Date)
		log.Printf("TaskList: %v", taskList)
	}

	if err = rows.Err(); err != nil {
		return taskList, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	log.Printf("TaskList 2: %v", taskList)
	return taskList, nil
}
