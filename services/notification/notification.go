package notification

import (
	"github.com/NubeIO/rubix-os/database"
	"github.com/go-co-op/gocron"
)

type Notification struct {
	cron *gocron.Scheduler
	DB   *database.GormDatabase
}
