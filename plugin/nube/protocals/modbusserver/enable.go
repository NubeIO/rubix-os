package main

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.log(false, "modbus-server:", "Enable():", "enable modbus system")
	go inst.serverInit()
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	return nil
}
