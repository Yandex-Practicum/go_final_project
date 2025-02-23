package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func addTask(t *testing.T, task task) string {
	ret, err := postJSON("api/task", map[string]any{
		"date":    task.date,
		"title":   task.title,
		"comment": task.comment,
		"repeat":  task.repeat,
	}, http.MethodPost)
	assert.NoError(t, err)
	assert.NotNil(t, ret["id"])
	id := fmt.Sprint(ret["id"])
	assert.NotEmpty(t, id)
	log.Printf("результат addTask: id= %s", id)
	return id
}

func getTasks(t *testing.T, search string) []map[string]string {
	url := "api/tasks"
	if Search {
		url += "?search=" + search
	}
	body, err := requestJSON(url, nil, http.MethodGet)
	assert.NoError(t, err)

	var m map[string][]map[string]string
	err = json.Unmarshal(body, &m)
	assert.NoError(t, err)
	log.Printf("результат getTasks: %+v", m["tasks"])
	return m["tasks"]
}

func TestTasks(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	now := time.Now()
	_, err := db.Exec("DELETE FROM scheduler")
	assert.NoError(t, err)

	tasks := getTasks(t, "")
	assert.NotNil(t, tasks)
	assert.Empty(t, tasks)

	addTask(t, task{
		date:    now.Format(`20060102`),
		title:   "Просмотр фильма",
		comment: "с попкорном",
		repeat:  "",
	})
	now = now.AddDate(0, 0, 1)
	date := now.Format(`20060102`)
	addTask(t, task{
		date:    date,
		title:   "Сходить в бассейн",
		comment: "",
		repeat:  "",
	})
	addTask(t, task{
		date:    date,
		title:   "Оплатить коммуналку",
		comment: "",
		repeat:  "d 30",
	})
	tasks = getTasks(t, "")
	log.Printf("задач: %d, ожидается 3", len(tasks))
	assert.Equal(t, len(tasks), 3)

	now = now.AddDate(0, 0, 2)
	date = now.Format(`20060102`)
	addTask(t, task{
		date:    date,
		title:   "Поплавать",
		comment: "Бассейн с тренером",
		repeat:  "d 7",
	})
	addTask(t, task{
		date:    date,
		title:   "Позвонить в УК",
		comment: "Разобраться с горячей водой",
		repeat:  "",
	})
	addTask(t, task{
		date:    date,
		title:   "Встретится с Васей",
		comment: "в 18:00",
		repeat:  "",
	})

	tasks = getTasks(t, "")
	log.Printf("задач: %d, ожидается 6", len(tasks))
	assert.Equal(t, len(tasks), 6)

	if !Search {
		return
	}
	tasks = getTasks(t, "УК")
	assert.Equal(t, len(tasks), 1)
	log.Printf("задач: %d, ожидается 1", len(tasks))
	tasks = getTasks(t, now.Format(`02.01.2006`))
	log.Printf("задач: %d, ожидается 3", len(tasks))
	assert.Equal(t, len(tasks), 3)

}
