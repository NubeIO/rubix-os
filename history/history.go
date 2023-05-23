package history

import (
	"github.com/NubeIO/rubix-os/database"
	"github.com/go-co-op/gocron"
)

type History struct {
	cron *gocron.Scheduler
	DB   *database.GormDatabase
}
