package client

import "github.com/NubeIO/flow-framework/nresty"

func (inst *FlowClient) WizardMasterSlavePointMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error) {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetPathParams(map[string]string{"global_uuid": globalUUID}).
		Post("/api/database/wizard/mapping/master_slave/points/consumer/{global_uuid}"))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (inst *FlowClient) WizardP2PMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error) {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetPathParams(map[string]string{"global_uuid": globalUUID}).
		Post("/api/database/wizard/mapping/p2p/points/consumer/{global_uuid}"))
	if err != nil {
		return false, err
	}
	return true, nil
}
