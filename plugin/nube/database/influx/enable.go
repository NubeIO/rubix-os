package main

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
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

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}
