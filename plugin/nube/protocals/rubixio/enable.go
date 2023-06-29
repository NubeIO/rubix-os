package main

import (
	"github.com/NubeIO/rubix-os/args"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, args.Args{})
	if q != nil {
		inst.networkUUID = q.UUID
	}
	if inst.config.EnablePolling {
		if !inst.pollingEnabled {
			var arg polling
			inst.pollingEnabled = true
			arg.enable = true
			go func() error {
				err := inst.polling(arg)
				if err != nil {
					log.Errorf("rubix-io.enable: POLLING ERROR on routine: %v\n", err)
				}
				return nil
			}()
			if err != nil {
				log.Errorf("rubix-io.enable: POLLING ERROR: %v\n", err)
			}
		}
	}
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	if inst.pollingEnabled {
		var arg polling
		inst.pollingEnabled = false
		arg.enable = false
		go func() {
			err := inst.polling(arg)
			if err != nil {
			}
		}()
	}
	return nil
}
