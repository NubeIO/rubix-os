package main

import (
	"context"
)

func (inst *Instance) Enable() error {
	inst.enabled = true
	q, _ := inst.db.GetPluginByPath(pluginPath)
	inst.pluginUUID = q.UUID
	inst.ctx, inst.cancel = context.WithCancel(context.Background())
	go inst.syncPointsLoopTEMPORARY(inst.ctx)
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	inst.cancel()
	return nil
}
