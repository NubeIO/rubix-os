package module

import (
	"encoding/json"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/database"
)

type dbHelper struct{}

func (*dbHelper) Sum(a, b int64) (int64, error) {
	return a + b, nil
}

func (*dbHelper) CallAPI(path, args string) ([]byte, error) {
	networks, err := database.GlobalGormDatabase.GetNetworks(api.Args{})
	nets, err := json.Marshal(networks)
	return nets, err
}
