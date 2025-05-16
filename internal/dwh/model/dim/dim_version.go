package dim

import "go-cs/internal/dwh/pkg/model"

type DimVersion struct {
	model.DimModel
	model.ChainModel

	SpaceId     int64
	VersionId   int64
	VersionName string
}
