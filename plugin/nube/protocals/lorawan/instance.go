package main

func (inst *Instance) setUUID() {
	q, err := inst.db.GetPluginByPath(pluginPath)
	if err != nil {
		return
	}
	inst.pluginUUID = q.UUID
}
