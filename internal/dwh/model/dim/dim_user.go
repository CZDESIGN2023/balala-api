package dim

import "go-cs/internal/dwh/pkg/model"

type DimUser struct {
	model.DimModel
	model.ChainModel

	UserId       int64
	UserName     string
	UserNickName string
	UserPinyin   string
}
