package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetPoints(args api.Args) ([]*model.Point, error) {
	q, err := getDb().GetPoints(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPointsByDeviceUUID(deviceUUID string, args api.Args) ([]*model.Point, error) {
	args.DeviceUUID = &deviceUUID
	q, err := getDb().GetPoints(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPoint(uuid string, args api.Args) (*model.Point, error) {
	q, err := getDb().GetPoint(uuid, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) CreatePoint(body *model.Point, fromPlugin, updatePoint bool) (
	*model.Point, error) {
	pnt, err := getDb().CreatePoint(body, fromPlugin)
	if err != nil {
		return nil, err
	}
	if updatePoint {
		pnt, err = getDb().UpdatePoint(pnt.UUID, pnt, fromPlugin, false) // MARC: UpdatePoint is called here so that the PresentValue and Priority are updated to use the fallback value.  Otherwise they are left as Null and the Edge28 Outputs are left floating.
		if err != nil {
			return nil, err
		}
	}
	return pnt, nil
}

func (h *Handler) UpdatePoint(uuid string, body *model.Point, afterRealDeviceUpdate bool) (*model.Point, error) {
	return getDb().UpdatePoint(uuid, body, true, afterRealDeviceUpdate)
}

func (h *Handler) PointWrite(uuid string, pointWriter *model.PointWriter, afterRealDeviceUpdate bool) (
	returnPoint *model.Point, isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
	return getDb().PointWrite(uuid, pointWriter, true, afterRealDeviceUpdate, nil, false)
}

// UpdatePointErrors will only update the error properties of the point, all other properties will not be updated.
func (h *Handler) UpdatePointErrors(uuid string, body *model.Point) error {
	return getDb().UpdatePointErrors(uuid, body)
}

func (h *Handler) GetOnePointByArgs(args api.Args) (*model.Point, error) {
	return getDb().GetOnePointByArgs(args)
}

func (h *Handler) DeletePoint(uuid string) (bool, error) {
	_, err := getDb().DeletePoint(uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}
