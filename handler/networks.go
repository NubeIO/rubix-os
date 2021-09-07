package handler

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
)

func (l *Handler) GetNetworks() []*model.Network {
	users, err := getDb().GetNetworks(false, false)
	if err != nil {
		return nil
	}
	if err != nil {
		fmt.Println(err, "ERR")
	}
	//store.Set("aa", "aaa", cache.NoExpiration)
	//fmt.Println(store.Get("aa"))
	return users
}
