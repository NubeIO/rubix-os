package main

// Enable implements plugin.Plugin
func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	//port := fmt.Sprintf("%d", p)
	//_, _, foundPort := linixpingport.PingPort("0.0.0.0", port, 1, false)
	//if !foundPort {
	//	go inst.serverInit()
	//}
	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	return nil
}
