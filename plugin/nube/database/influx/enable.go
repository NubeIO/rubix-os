package main

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	rest "github.com/NubeDev/flow-framework/plugin/nube/protocals/lorawan/restclient"
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
	i.REST = rest.NewChirp(chirpName, chirpPass, ip, port)

	m := new(model.Job)
	fmt.Println(111111)
	m.StartDate = time.Date(2000, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
	m.EndDate = time.Date(2055, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
	m.Frequency = "5s"
	m.PluginConfId = i.pluginUUID
	fmt.Println(111111)
	err = i.jobs.JobAdd(m)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(111111)
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}
