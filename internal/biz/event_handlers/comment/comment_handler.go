package comment

import (
	"context"
	"fmt"
	"go-cs/internal/biz/command"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	witem_status_repo "go-cs/internal/domain/work_item_status/repo"
	"go-cs/internal/pkg/domain"
	"go-cs/internal/utils"

	"github.com/go-kratos/kratos/v2/log"
)

type CommentHandler struct {
	log                   *log.Helper
	workItemStatusRepo    witem_status_repo.WorkItemStatusRepo
	addWorkItemCommentCmd *command.AddWorkItemCommentCmd
	domainMessageConsumer *domain_message.DomainMessageConsumer
}

func NewCommentHandler(
	logger log.Logger,
	workItemStatusRepo witem_status_repo.WorkItemStatusRepo,
	addWorkItemCommentCmd *command.AddWorkItemCommentCmd,
	domainMessageConsumer *domain_message.DomainMessageConsumer,
) *CommentHandler {

	moduleName := "CommentHandler"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &CommentHandler{
		log:                   hlog,
		workItemStatusRepo:    workItemStatusRepo,
		addWorkItemCommentCmd: addWorkItemCommentCmd,
		domainMessageConsumer: domainMessageConsumer,
	}
}

func (s *CommentHandler) Init() {
	s.domainMessageConsumer.SetMessageListener("comment_handler", s.domainMessagePublishEventHandler)
}

func (s *CommentHandler) domainMessagePublishEventHandler(evt domain_message.DomainMessagePublishEvent) {

	for _, v := range evt.Messages {
		switch v.MessageType() {
		case domain_message.Message_Type_WorkItem_Status_Change:
			s.changeWorkItemStatusByDomainMessage(v)
		case domain_message.Message_Type_WorkItem_FlowNode_Rollback:
			s.rollbackWorkItemFlowNodeByDomainMessage(v)
		}
	}
}

func (s *CommentHandler) rollbackWorkItemFlowNodeByDomainMessage(msg domain.DomainMessage) {

	s.log.Infof("comment.rollbackWorkItemFlowNodeByDomainMessage: %+v", msg)
	e := msg.(*domain_message.RollbackWorkItemFlowNode)

	ctx := context.Background()
	s.addWorkItemCommentCmd.Execute(ctx, &command.AddCommentVo{
		OperUid:    e.Oper.GetId(),
		SpaceId:    e.SpaceId,
		WorkItemId: e.WorkItemId,
		Content:    "已回滚：" + e.Reason,
	})

}

func (s *CommentHandler) changeWorkItemStatusByDomainMessage(msg domain.DomainMessage) {
	s.log.Infof("comment.changeWorkItemStatusByDomainMessage: %+v", msg)

	e := msg.(*domain_message.ChangeWorkItemStatus)

	ctx := context.Background()

	//没原因的，就都不评论了
	if e.Reason == "" && e.Remark == "" {
		return
	}

	oldStatus, err := s.workItemStatusRepo.GetWorkItemStatusItem(ctx, e.OldWorkItemStatusId)
	if err != nil {
		s.log.Error(err)
		return
	}

	newStatus, err := s.workItemStatusRepo.GetWorkItemStatusItem(ctx, e.NewWorkItemStatusId)
	if err != nil {
		s.log.Error(err)
		return
	}

	if newStatus.IsClose() {

		s.addWorkItemCommentCmd.Execute(ctx, &command.AddCommentVo{
			OperUid:    e.Oper.GetId(),
			SpaceId:    e.SpaceId,
			WorkItemId: e.WorkItemId,
			Content:    "已关闭：" + e.Reason,
		})
	} else if newStatus.IsTerminated() {

		s.addWorkItemCommentCmd.Execute(ctx, &command.AddCommentVo{
			OperUid:    e.Oper.GetId(),
			SpaceId:    e.SpaceId,
			WorkItemId: e.WorkItemId,
			Content:    "已终止：" + e.Reason,
		})

	} else if oldStatus.IsClose() && !newStatus.IsArchivedTypeState() {

		s.addWorkItemCommentCmd.Execute(ctx, &command.AddCommentVo{
			OperUid:    e.Oper.GetId(),
			SpaceId:    e.SpaceId,
			WorkItemId: e.WorkItemId,
			Content:    "已重启：" + e.Reason,
		})

	} else if oldStatus.IsTerminated() && !newStatus.IsArchivedTypeState() {

		s.addWorkItemCommentCmd.Execute(ctx, &command.AddCommentVo{
			OperUid:    e.Oper.GetId(),
			SpaceId:    e.SpaceId,
			WorkItemId: e.WorkItemId,
			Content:    "已恢复：" + e.Reason,
		})

	} else if oldStatus.IsCompleted() && !newStatus.IsArchivedTypeState() {

		s.addWorkItemCommentCmd.Execute(ctx, &command.AddCommentVo{
			OperUid:    e.Oper.GetId(),
			SpaceId:    e.SpaceId,
			WorkItemId: e.WorkItemId,
			Content:    "已重启：" + e.Reason,
		})
	} else if e.WorkItemTypeKey == consts.WorkItemTypeKey_StateTask {
		content := fmt.Sprintf("从 [%s] 流转至：[%s]", oldStatus.Name, newStatus.Name)
		if e.Reason != "" {
			content += "\n原因：" + e.Reason
		}
		if e.Remark != "" {
			content += "\n备注：" + e.Remark
		}
		s.addWorkItemCommentCmd.Execute(ctx, &command.AddCommentVo{
			OperUid:    e.Oper.GetId(),
			SpaceId:    e.SpaceId,
			WorkItemId: e.WorkItemId,
			Content:    fmt.Sprintf("<pre>%s</pre>", content),
		})
	}

}
