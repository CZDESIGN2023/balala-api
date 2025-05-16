package domain

import "go-cs/internal/utils/oper"

type UserOper int64

func (o UserOper) GetType() oper.OperatorType {
	return oper.OperatorTypeUser
}

func (o UserOper) GetId() int64 {
	return int64(o)
}

type SysOper int64

func (o SysOper) GetType() oper.OperatorType {
	return oper.OperatorTypeSys
}

func (o SysOper) GetId() int64 {
	return int64(o)
}
