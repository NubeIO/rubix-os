package main

func (inst *Instance) SyncRubixToFF() (bool, error) {
	inst.rubixpointsyncDebugMsg("Sync Rubix Points with FF Points...")

	/*
		ffNetworksArray, err := inst.GetFFNetworks(inst.config.Job.Networks)
		if err != nil {
			inst.rubixpointsyncErrorMsg(err)
		}

	*/

	err := inst.GetRubixNetworks()
	if err != nil {
		inst.rubixpointsyncErrorMsg(err)
	}

	return true, nil
}
