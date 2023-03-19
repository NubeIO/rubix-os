package dbhandler

import (
	"encoding/json"
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

func (h *Handler) CreatePoint(body *model.Point, updatePoint bool) (
	*model.Point, error) {
	pnt, err := getDb().CreatePoint(body)
	if err != nil {
		return nil, err
	}
	if updatePoint {
		pnt, err = h.UpdatePoint(pnt.UUID, pnt) // MARC: UpdatePoint is called here so that the PresentValue and Priority are updated to use the fallback value. Otherwise, they are left as Null and the Edge28 Outputs are left floating.
		if err != nil {
			return nil, err
		}
	}
	return pnt, nil
}

func createPointDeepCopy(point model.Point) model.Point {
	var outputPoint model.Point
	out, _ := json.Marshal(point)
	_ = json.Unmarshal(out, &outputPoint)
	return outputPoint
}

func (h *Handler) UpdatePoint(uuid string, body *model.Point) (*model.Point, error) {
	return getDb().UpdatePoint(uuid, body, true)
}

// TODO: This was only added to allow for the EnableWriteable property to be updated.  It can be removed (along with the code at the bottom of SYSTEM plugin enable().
func (h *Handler) UpdatePointPlugin(uuid string, body *model.Point) (*model.Point, error) {
	return getDb().UpdatePointPlugin(uuid, body)
}

func (h *Handler) PointWrite(uuid string, pointWriter *model.PointWriter) (
	returnPoint *model.Point, isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
	return getDb().PointWrite(uuid, pointWriter, nil, false)
}

// UpdatePointErrors will only update the error properties of the point, all other properties will not be updated.
func (h *Handler) UpdatePointErrors(uuid string, body *model.Point) error {
	return getDb().UpdatePointErrors(uuid, body)
}

// UpdatePointSuccess will only update the error properties of the point, all other properties will not be updated.
func (h *Handler) UpdatePointSuccess(uuid string, body *model.Point) error {
	return getDb().UpdatePointSuccess(uuid, body)
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
