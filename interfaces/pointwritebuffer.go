package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type PointWriteBuffer struct {
	UUID                  string             `json:"uuid"`
	Body                  *model.PointWriter `json:"body"`
	AfterRealDeviceUpdate bool               `json:"after_real_device_update"`
	CurrentWriterUUID     *string            `json:"current_writer_uuid"`
	ForceWrite            bool               `json:"force_write"`
}
