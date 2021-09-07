package dbhandler

import (
	"github.com/NubeDev/flow-framework/database"
)

var db *database.GormDatabase


type Handler struct {
	DB     *database.GormDatabase
}


//Init give db access
func Init(h *Handler) {
	initDb(h.DB)
}

func initDb(d *database.GormDatabase) {
	db = d
}

func getDb() *database.GormDatabase {
	return db
}
