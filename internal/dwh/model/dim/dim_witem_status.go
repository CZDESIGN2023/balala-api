package dim

import "go-cs/internal/dwh/pkg/model"

type DimWitemStatus struct {
	model.DimModel
	model.ChainModel

	SpaceId    int64
	StatusId   int64
	StatusName string
	StatusKey  string
	StatusVal  string
	StatusType int32
	FlowScope  string
}
