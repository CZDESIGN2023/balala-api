package repo

import (
	"context"
	db "go-cs/internal/bean/biz"
)

type OperLogRepo interface {
	AddOperLog(ctx context.Context, in *db.OperLog) error
}
