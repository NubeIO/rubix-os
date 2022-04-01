package main

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	log "github.com/sirupsen/logrus"
)

// Enable implements plugin.Plugin
func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.BusServ()
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if q != nil {
		inst.networkUUID = q.UUID
	} else {
		inst.networkUUID = "NA"
	}
	if inst.config.EnablePolling {
		if !inst.pollingEnabled {
			var arg polling
			inst.pollingEnabled = true
			inst.PointWriteModeTest()
			arg.enable = true
			go func() error {
				err := inst.PollingTCP(arg)
				if err != nil {
					log.Errorf("modbus:  PLUGIN Enable POLLING ERROR: %v\n", err)
				}
				return nil
			}()
			if err != nil {
				log.Errorf("modbus: POLLING ERROR: %v\n", err)
			}
		}
	}
	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	if inst.pollingEnabled {
		var arg polling
		inst.pollingEnabled = false
		arg.enable = false
		go func() {
			err := inst.PollingTCP(arg)
			if err != nil {
				log.Errorf("modbus:  PLUGIN Disable POLLING ERROR: %v\n", err)
			}
		}()
		if err != nil {
			return errors.New("error on starting polling")
		}
	}
	return nil
}
