package cachestore

import (
	"fmt"
	"github.com/patrickmn/go-cache"
)

var store *cache.Cache

type Handler struct {
	Store  *cache.Cache
}

//Init init store
func Init(h *Handler) {
	initStore(h.Store)
}

func initStore(d *cache.Cache) {
	store = d
}

func getStore() *cache.Cache {
	return store
}


func (l *Handler) Get()  {
	store.Set("aa", "aaa", cache.NoExpiration)
	fmt.Println(store.Get("aa"))
}
