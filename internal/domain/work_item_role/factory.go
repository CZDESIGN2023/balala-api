package work_item_role

import (
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"

	"github.com/google/uuid"
)

func NewWorkItemRole(id int64, spaceId int64, workItemTypeId int64, name string, key string, ranking int64, isSys int32, uid int64, flowScope consts.FlowScope, oper shared.Oper) *WorkItemRole {

	ins := &WorkItemRole{
		Id:             id,
		UserId:         uid,
		Uuid:           uuid.NewString(),
		SpaceId:        spaceId,
		WorkItemTypeId: workItemTypeId,
		Key:            key,
		Name:           name,
		Status:         1,
		Ranking:        ranking,
		IsSys:          isSys,
		FlowScope:      flowScope,
	}

	ins.AddMessage(oper, &domain_message.CreateWorkItemRole{
		SpaceId:          spaceId,
		WorkItemRoleName: name,
		WorkItemRoleId:   id,
		FlowScope:        flowScope,
	})

	return ins
}
