package main

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.mapLoRa(inst.config.Job.UseExistingNetwork)
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	return nil
}
