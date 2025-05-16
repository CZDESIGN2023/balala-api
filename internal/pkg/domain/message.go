package domain

import (
	"go-cs/internal/utils/oper"
	"time"
)

type MessageType string

type Oper interface {
	GetType() oper.OperatorType
	GetId() int64
}

//修改了， 添加了 ，删除了，移除了 ...

type DomainMessages []DomainMessage

type DomainMessage interface {
	//同一种类型的日志，会被合并记录
	MessageType() MessageType
	SetOper(oper Oper, opTime time.Time)
	GetOper() Oper
	GetOperTime() time.Time
}

type DomainMessageBase struct {
	DomainMessage

	Oper     Oper
	OperTime time.Time
}

func (msg *DomainMessageBase) SetOper(oper Oper, opTime time.Time) {
	msg.Oper = oper
	msg.OperTime = opTime
}

func (msg *DomainMessageBase) GetOper() Oper {
	return msg.Oper
}

func (msg *DomainMessageBase) GetOperTime() time.Time {
	return msg.OperTime
}
