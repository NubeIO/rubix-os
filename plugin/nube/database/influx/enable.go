package main

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"time"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	var arg api.Args
	arg.IpConnection = true
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if err != nil {
		return errors.New("there is no network added please add one")
	}
	i.networkUUID = q.UUID
	if err != nil {
		return errors.New("error on enable lora-plugin")
	}
	i.REST = nil
	m := new(model.Job)
	m.StartDate = time.Date(2000, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
	m.EndDate = time.Date(2055, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
	m.Frequency = "15s"
	m.PluginConfId = i.pluginUUID
	err = i.jobs.JobAdd(m)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}
