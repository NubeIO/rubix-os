package dbhandler

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

func (h *Handler) GetPointMetaTags() ([]*model.PointMetaTag, error) {
	return getDb().GetPointMetaTags()
}
