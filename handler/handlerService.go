package handler

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/patrickmn/go-cache"
	"time"
)


type Handler struct {
	SourceField           string
	SkipErrRecordNotFound bool
	NotificationService   eventbus.NotificationService
	ctx                   context.Context
	store                 *cache.Cache
	//DB 					  database.GormDatabase //TODO will cause cycle import
}

func NewHandler() *Handler {
	notificationService := eventbus.NewNotificationService(eventbus.BUS)
	c := eventbus.BusContext
	store := cache.New(5*time.Minute, 10*time.Minute)
	return &Handler{
		SkipErrRecordNotFound: true,
		NotificationService:   notificationService,
		ctx:                   c,
		store:                 store,
	}
}

func (l *Handler) Get() string {
	l.store.Set("sssss", "sss", cache.NoExpiration)
	fmt.Println(l.store.Get("sssss"))

	//l.DB.GetIntegrationsList()
	//list, err := l.db.GetIntegrationsList()
	//if err != nil {
	//	fmt.Println(err)
	//	return ""
	//}
	//fmt.Println(list)
	return "l"
}
