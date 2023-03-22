package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type SyncDevice struct {
	DeviceName     string
	DeviceTags     []*model.Tag
	DeviceMetaTags []*model.DeviceMetaTag
}
