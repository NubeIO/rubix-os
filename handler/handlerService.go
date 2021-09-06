package handler

import (
	"context"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"time"
)

type Handler struct {
	SourceField           string
	SkipErrRecordNotFound bool
	NotificationService   eventbus.NotificationService
	ctx                   context.Context
	store                 *cache.Cache
	db                    *DB
}

type DBHandler interface {
	Test()
}

type DB struct {
	DBHandler DBHandler
}

func CustomHandler(db *DB) *Handler {
	notificationService := eventbus.NewNotificationService(eventbus.BUS)
	c := eventbus.BusContext
	store := cache.New(5*time.Minute, 10*time.Minute)
	return &Handler{
		SkipErrRecordNotFound: true,
		NotificationService:   notificationService,
		ctx:                   c,
		store:                 store,
		db:                    db,
	}
}

func (l *Handler) Get() string {
	l.store.Set("sssss", "sss", cache.NoExpiration)
	v, isExist := l.store.Get("sssss")
	log.Info("Get store value: ", v, ", which does exist: ", isExist)
	l.db.DBHandler.Test()
	return "l"
}
