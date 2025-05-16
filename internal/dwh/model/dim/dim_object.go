package dim

import "go-cs/internal/dwh/pkg/model"

type DimObject struct {
	model.DimModel
	model.ChainModel

	SpaceId    int64
	ObjectId   int64
	ObjectName string
}
