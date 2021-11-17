package floweng

import (
	"github.com/NubeIO/flow-framework/database"
	"github.com/NubeIO/flow-framework/floweng/server"
)

var flowEngServer *server.Server

func EngStart(db *database.GormDatabase) {
	flowEngServer = server.NewServer(db)
	flowEngServer.LoadFromDB(db)
	flowEngServer.SetRouter()
}
