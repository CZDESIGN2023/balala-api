package domain

import (
	"go-cs/pkg/stream"
	"time"
)

type DomainEntity struct {
	_diffs    PropDiffs
	_messages DomainMessages
}

func (ar *DomainEntity) GetDiffs() PropDiffs {
	diffs := stream.Unique(ar._diffs)
	ar._diffs = make(PropDiffs, 0)
	return diffs
}

func (ar *DomainEntity) AddDiff(diff ...PropDiff) {
	ar._diffs = append(ar._diffs, diff...)
}

func (ar *DomainEntity) HasDiffs() bool {
	return len(ar._diffs) > 0
}

func (ar *DomainEntity) AddMessage(oper Oper, log DomainMessage) {
	if oper == nil || log == nil {
		return
	}

	log.SetOper(oper, time.Now())
	ar._messages = append(ar._messages, log)
}

func (ar *DomainEntity) GetMessages() DomainMessages {
	return ar._messages
}
