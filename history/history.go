package history

import (
	"github.com/NubeIO/flow-framework/database"
	"github.com/go-co-op/gocron"
)

type History struct {
	cron *gocron.Scheduler
	DB   *database.GormDatabase
}
