package cachestore

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var store *cache.Cache

//StoreKeys used to try and keep track of the keys
var StoreKeys = struct {
	PluginUUID string
	PluginName string //would be plugin_name_lora
}{
	PluginUUID: "plugin_uuid",
	PluginName: "plugin_name",
}


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

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (l *Handler) Get(key string) (interface{}, bool) {
	value, found := store.Get(key)
	return value, found
}


// Set an item to the cache, replacing any existing item. If the duration is 0
// (DefaultExpiration), the cache's default expiration time is used. If it is -1
// (NoExpiration), the item never expires.
func (l *Handler) Set(key string, value interface{}, d time.Duration) {
	store.Set(key, value, d)
}
