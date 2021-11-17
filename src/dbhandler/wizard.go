package dbhandler

import "github.com/NubeIO/flow-framework/model"

func (h *Handler) WizardNewNetworkDevicePoint(plugin string, net *model.Network, dev *model.Device, pnt *model.Point) (bool, error) {
	_, err := getDb().WizardNewNetworkDevicePoint(plugin, net, dev, pnt)
	if err != nil {
		return false, err
	}
	return true, nil
}
