package handler

import (
	"context"
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/patrickmn/go-cache"
)

var db *database.GormDatabase
var bus eventbus.BusService
var busCTX context.Context
var store *cache.Cache

type Handler struct {
	BUS    eventbus.BusService
	BusCTX context.Context
	Store  *cache.Cache
	DB     *database.GormDatabase
}

func InitHandler(h *Handler) {
	initDb(h.DB)
	initBus(h.BUS)
	initCTX(h.BusCTX)
	initStore(h.Store)
}

func initDb(d *database.GormDatabase) {
	db = d
}

func initBus(d eventbus.BusService) {
	bus = d
}

func initCTX(d context.Context) {
	busCTX = d
}

func initStore(d *cache.Cache) {
	store = d
}

func getDb() *database.GormDatabase {
	return db
}

func getBus() eventbus.BusService {
	return bus
}

func getCTX() context.Context {
	return busCTX
}
