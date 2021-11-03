package main

import "fmt"

// Enable implements plugin.Plugin.Enable
func (i *Instance) Enable() error {
	fmt.Print("Enable")
	i.enabled = true
	i.setUUID()
	i.BusServ()
	job, err := i.createJob()
	if err != nil {
		return err
	}
	err = i.jobs.JobAdd(job)
	if err != nil {
		return err
	}
	return nil
}

// Disable implements plugin.Plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}
