package module

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	parentArgs "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/database"
	"github.com/NubeIO/rubix-os/module/common"
	log "github.com/sirupsen/logrus"
)

type dbHelper struct{}

func (*dbHelper) GetWithoutParam(path, args string) ([]byte, error) {
	var out interface{}
	var err error
	apiArgs, err := parentArgs.DeserializeArgs(args)
	if err != nil {
		return nil, err
	}
	if path == "networks" {
		out, err = database.GlobalGormDatabase.GetNetworks(*apiArgs)
	} else if path == "devices" {
		out, err = database.GlobalGormDatabase.GetDevices(*apiArgs)
	} else if path == "points" {
		out, err = database.GlobalGormDatabase.GetPoints(*apiArgs)
	} else if path == "one_device_by_args" {
		out, err = database.GlobalGormDatabase.GetOneDeviceByArgs(*apiArgs)
	} else if path == "one_point_by_args" {
		out, err = database.GlobalGormDatabase.GetOneDeviceByArgs(*apiArgs)
	} else if path == "schedules" {
		out, err = database.GlobalGormDatabase.GetSchedules()
	} else {
		return nil, errors.New("not found")
	}
	return marshal(err, out)
}

func (*dbHelper) Get(path, uuid, args string) ([]byte, error) {
	var out interface{}
	var err error
	apiArgs, err := parentArgs.DeserializeArgs(args)
	if err != nil {
		return nil, err
	}
	if path == "networks" {
		out, err = database.GlobalGormDatabase.GetNetwork(uuid, *apiArgs)
	} else if path == "devices" {
		out, err = database.GlobalGormDatabase.GetDevice(uuid, *apiArgs)
	} else if path == "points" {
		out, err = database.GlobalGormDatabase.GetPoint(uuid, *apiArgs)
	} else if path == "networks_by_plugin_name" {
		out, err = database.GlobalGormDatabase.GetNetworksByPluginName(uuid, *apiArgs)
	} else if path == "network_by_name" {
		name := uuid
		out, err = database.GlobalGormDatabase.GetNetworkByName(name, *apiArgs)
	} else if path == "plugin_by_path" {
		out, err = database.GlobalGormDatabase.GetPluginByPath(uuid)
	} else if path == "plugin_by_id" {
		out, err = database.GlobalGormDatabase.GetPlugin(uuid)
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
	} else if path == "network_errors" {
		network := model.Network{}
		err = json.Unmarshal(body, &network)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = database.GlobalGormDatabase.UpdateNetworkErrors(uuid, &network)
	} else if path == "device_errors" {
		device := model.Device{}
		err = json.Unmarshal(body, &device)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = database.GlobalGormDatabase.UpdateDeviceErrors(uuid, &device)
	} else if path == "point_errors" {
		point := model.Point{}
		err = json.Unmarshal(body, &point)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = database.GlobalGormDatabase.UpdatePointErrors(uuid, &point)
	} else if path == "point_success" {
		point := model.Point{}
		err = json.Unmarshal(body, &point)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = database.GlobalGormDatabase.UpdatePointSuccess(uuid, &point)
	} else if path == "point_write" {
		pw := model.PointWriter{}
		err = json.Unmarshal(body, &pw)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		point, isPresentValueChange, isWriteValueChange, isPriorityChanged, err := database.GlobalGormDatabase.PointWrite(uuid, &pw)
		if err != nil {
			return nil, err
		}
		out = common.PointWriteResponse{
			Point:                *point,
			IsPresentValueChange: isPresentValueChange,
			IsWriteValueChange:   isWriteValueChange,
			IsPriorityChanged:    isPriorityChanged,
		}
	} else if path == "schedules" {
		sch := model.Schedule{}
		err = json.Unmarshal(body, &sch)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		out, err = database.GlobalGormDatabase.UpdateScheduleAllProps(uuid, &sch)
		if err != nil {
			return nil, err
		}
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

func (*dbHelper) SetErrorsForAll(path, uuid, message, messageLevel, messageCode string, doPoints bool) error {
	var err error
	if path == "devices_on_network" {
		err = database.GlobalGormDatabase.SetErrorsForAllDevicesOnNetwork(
			uuid,
			message,
			messageLevel,
			messageCode,
			doPoints,
		)
	} else if path == "points_on_device" {
		err = database.GlobalGormDatabase.SetErrorsForAllPointsOnDevice(
			uuid,
			message,
			messageLevel,
			messageCode,
		)
	}
	if err != nil {
		return err
	}
	return nil
}

func (*dbHelper) ClearErrorsForAll(path, uuid string, doPoints bool) error {
	var err error
	if path == "devices_on_network" {
		err = database.GlobalGormDatabase.ClearErrorsForAllDevicesOnNetwork(uuid, doPoints)
	} else if path == "points_on_device" {
		err = database.GlobalGormDatabase.ClearErrorsForAllPointsOnDevice(uuid)
	}
	if err != nil {
		return err
	}
	return nil
}

func (*dbHelper) WizardNewNetworkDevicePoint(plugin string, net, dev, pnt []byte) (bool, error) {
	var err error
	network := model.Network{}
	err = json.Unmarshal(net, &network)
	if err != nil {
		log.Error(err)
		return false, err
	}
	device := model.Device{}
	err = json.Unmarshal(dev, &device)
	if err != nil {
		log.Error(err)
		return false, err
	}
	point := model.Point{}
	err = json.Unmarshal(pnt, &point)
	if err != nil {
		log.Error(err)
		return false, err
	}
	_, err = database.GlobalGormDatabase.WizardNewNetworkDevicePoint(plugin, &network, &device, &point)
	if err != nil {
		return false, err
	}
	return true, nil
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
