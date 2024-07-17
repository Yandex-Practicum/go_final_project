package handlers

import (
	"github.com/LEbauchoir/go_final_project/database"
)

var dbHelper *database.DbHelper

func SetDBHelper(helper *database.DbHelper) {
	dbHelper = helper
}
