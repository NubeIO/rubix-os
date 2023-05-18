package module

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/database"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

type dbHelper struct{}

func (*dbHelper) GetWithoutParam(path, args string) ([]byte, error) {
	apiArgs := parseArgs(args)
	var out interface{}
	var err error
	if path == "networks" {
		out, err = database.GlobalGormDatabase.GetNetworks(apiArgs)
	} else if path == "devices" {
		out, err = database.GlobalGormDatabase.GetDevices(apiArgs)
	} else if path == "points" {
		out, err = database.GlobalGormDatabase.GetPoints(apiArgs)
	} else if path == "flow_networks" {
		out, err = database.GlobalGormDatabase.GetFlowNetworks(apiArgs)
	} else {
		return nil, errors.New("not found")
	}
	return marshal(err, out)
}

func (*dbHelper) Get(path, uuid, args string) ([]byte, error) {
	apiArgs := parseArgs(args)
	var out interface{}
	var err error
	if path == "networks" {
		out, err = database.GlobalGormDatabase.GetNetwork(uuid, apiArgs)
	} else if path == "devices" {
		out, err = database.GlobalGormDatabase.GetNetwork(uuid, apiArgs)
	} else if path == "points" {
		out, err = database.GlobalGormDatabase.GetPoint(uuid, apiArgs)
	} else if path == "networks_by_plugin_name" {
		out, err = database.GlobalGormDatabase.GetNetworkByPluginName(uuid, apiArgs)
	} else {
		return nil, errors.New("not found")
	}
	return marshal(err, out)
}

func (*dbHelper) Post(path string, body []byte) ([]byte, error) {
	var out interface{}
	var err error
	if path == "networks" {
		network := model.Network{}
		err = json.Unmarshal(body, &network)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		out, err = database.GlobalGormDatabase.CreateNetwork(&network)
	} else if path == "devices" {
		device := model.Device{}
		err = json.Unmarshal(body, &device)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		out, err = database.GlobalGormDatabase.CreateDevice(&device)
	} else if path == "points" {
		point := model.Point{}
		err = json.Unmarshal(body, &point)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		out, err = database.GlobalGormDatabase.CreatePoint(&point)
	} else {
		return nil, errors.New("not found")
	}
	return marshal(err, out)
}

func (*dbHelper) Put(path, uuid string, body []byte) ([]byte, error) {
	return nil, nil
}

func (*dbHelper) Patch(path, uuid string, body []byte) ([]byte, error) {
	var out interface{}
	var err error
	if path == "networks" {
		network := model.Network{}
		err = json.Unmarshal(body, &network)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		out, err = database.GlobalGormDatabase.UpdateNetwork(uuid, &network)
	} else if path == "devices" {
		device := model.Device{}
		err = json.Unmarshal(body, &device)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		out, err = database.GlobalGormDatabase.UpdateDevice(uuid, &device)
	} else if path == "points" {
		point := model.Point{}
		err = json.Unmarshal(body, &point)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		out, err = database.GlobalGormDatabase.UpdatePoint(uuid, &point)
	} else {
		return nil, errors.New("not found")
	}
	return marshal(err, out)
}

func (*dbHelper) Delete(path, uuid string) ([]byte, error) {
	var out interface{}
	var err error
	if path == "networks" {
		out, err = database.GlobalGormDatabase.DeleteNetwork(uuid)
	} else if path == "devices" {
		out, err = database.GlobalGormDatabase.DeleteDevice(uuid)
	} else if path == "points" {
		out, err = database.GlobalGormDatabase.DeletePoint(uuid)
	} else {
		return nil, errors.New("not found")
	}
	return marshal(err, out)
}

func marshal(err error, out interface{}) ([]byte, error) {
	if err != nil {
		log.Error(err)
		return nil, err
	}
	o, e := json.Marshal(out)
	if e != nil {
		return nil, e
	}
	return o, nil
}
