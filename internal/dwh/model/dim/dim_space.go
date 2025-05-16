package dim

import "go-cs/internal/dwh/pkg/model"

type DimSpace struct {
	model.DimModel
	model.ChainModel

	SpaceId int64
	// UserId    int64
	SpaceName string
}
