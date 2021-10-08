package main

import (
	edgerest "github.com/NubeDev/flow-framework/plugin/nube/protocals/edge28/restclient"
	"os"
	"os/signal"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.setUUID()

	i.rest = edgerest.NewNoAuth("192.168.15.101", "5000")
	i.enabled = true
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
