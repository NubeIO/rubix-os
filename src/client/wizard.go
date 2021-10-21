package client

import (
	"fmt"
)

func (a *FlowClient) WizardMasterSlavePointMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error) {
	resp, err := a.client.R().
		SetPathParams(map[string]string{"global_uuid": globalUUID}).
		Post("/api/database/wizard/mapping/master_slave/points/consumer/{global_uuid}")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return false, fmt.Errorf("WizardMasterSlaveMappingOnConsumerSideByProducerSide: %s", err)
		} else {
			return false, fmt.Errorf("WizardMasterSlaveMappingOnConsumerSideByProducerSide: %s", resp)
		}
	}
	if resp.IsError() {
		return false, fmt.Errorf("WizardMasterSlaveMappingOnConsumerSideByProducerSide: %s", resp)
	}
	return true, nil
}

func (a *FlowClient) WizardRemotePointMappingOnConsumerSideByProducerSide(globalUUID string) (bool, error) {
	resp, err := a.client.R().
		SetPathParams(map[string]string{"global_uuid": globalUUID}).
		Post("/api/database/wizard/mapping/remote/points/consumer/{global_uuid}")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return false, fmt.Errorf("WizardRemotePointMappingOnConsumerSideByProducerSide: %s", err)
		} else {
			return false, fmt.Errorf("WizardRemotePointMappingOnConsumerSideByProducerSide: %s", resp)
		}
	}
	if resp.IsError() {
		return false, fmt.Errorf("WizardRemotePointMappingOnConsumerSideByProducerSide: %s", resp)
	}
	return true, nil
}
