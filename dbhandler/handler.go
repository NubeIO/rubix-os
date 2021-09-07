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
	db = h.DB
}

func getDb() *database.GormDatabase {
	return db
}
