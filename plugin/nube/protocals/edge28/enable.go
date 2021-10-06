package main

import (
	"os"
	"os/signal"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.setUUID()
	if !i.brokerEnabled {

	}
	i.brokerEnabled = true
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
