package oper

import (
	"context"

	"github.com/google/uuid"
)

type OperLogCtxKey struct{}

type OperData interface {
	GetOperInfo() *OperResultInfo
	SetOperLogger(*OperLogger)
}

func NewOperLogger() *OperLogger {
	return &OperLogger{
		UUID:        uuid.NewString(),
		RequestInfo: &OperRequestInfo{},
		operLogs:    make([]OperData, 0),
	}
}

func NewOperLoggerWithCtx(ctx context.Context) (context.Context, *OperLogger) {
	opLogger := NewOperLogger()
	newCtx := context.WithValue(ctx, OperLogCtxKey{}, opLogger)
	return newCtx, opLogger
}

func GetOperLoggerFormCtx(ctx context.Context) *OperLogger {
	operLoggerV := ctx.Value(OperLogCtxKey{})
	if operLogger, isOk := operLoggerV.(*OperLogger); isOk {
		return operLogger
	}
	return nil
}

type OperLogger struct {
	UUID        string
	Operator    *OperUser
	RequestInfo *OperRequestInfo
	operLogs    []OperData
}

type OperRequestInfo struct {
	RequestMethod string
	OperParam     string
	OperUrl       string
	OperIp        string
	OperLocation  string
}

// 操作人
type OperUser struct {
	OperType         int
	OperUid          int64
	OperUname        string
	OperUserNickName string
}

type OperResultInfo struct {
	SpaceId   int64
	SpaceName string

	//谁
	//如何操作了 1 新增 2 修改  3 删除
	BusinessType BusinessType
	ShowType     ShowType
	OperatorType OperatorType
	//操作项 标题, 类型，关联id
	ModuleTitle string
	ModuleType  ModuleType
	ModuleFlags []ModuleFlag
	ModuleId    int
	//操作结果描述
	OperMsg string
}

func (p *OperLogger) Add(in OperData) {
	p.operLogs = append(p.operLogs, in)
}

func (p *OperLogger) GetLogs() []OperData {
	return p.operLogs
}
