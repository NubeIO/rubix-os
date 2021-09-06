package database

import log "github.com/sirupsen/logrus"

type DBHandler struct{}

func (db *DBHandler) Test() {
	networks, _ := gormDatabase.GetFlowNetworks(true)
	log.Info("Flow Networks are: ", networks)
}
