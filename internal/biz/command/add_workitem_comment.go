package command

import (
	"context"
	domain_message "go-cs/internal/domain/pkg/message"
	member_repo "go-cs/internal/domain/space_member/repo"
	comment_repo "go-cs/internal/domain/space_work_item_comment/repo"
	comment_service "go-cs/internal/domain/space_work_item_comment/service"
	witem_repo "go-cs/internal/domain/work_item/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils/errs"
	"time"
)

type AddWorkItemCommentCmd struct {
	tm              trans.Transaction
	commentService  *comment_service.SpaceWorkItemCommentService
	witemRepo       witem_repo.WorkItemRepo
	spaceMemberRepo member_repo.SpaceMemberRepo
	commentRepo     comment_repo.SpaceWorkItemCommentRepo

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewAddWorkItemCommentCommand(
	tm trans.Transaction,
	commentService *comment_service.SpaceWorkItemCommentService,
	witemRepo witem_repo.WorkItemRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,
	commentRepo comment_repo.SpaceWorkItemCommentRepo,
) *AddWorkItemCommentCmd {
	return &AddWorkItemCommentCmd{
		commentService:  commentService,
		witemRepo:       witemRepo,
		spaceMemberRepo: spaceMemberRepo,
		commentRepo:     commentRepo,
		tm:              tm,
	}
}

type AddCommentVo struct {
	OperUid        int64
	SpaceId        int64
	WorkItemId     int64
	Content        string
	ReferUserIds   []int64
	ReplyCommentId int64
}

type AddCommentResultVo struct {
	CommentId int64
}

func (cmd *AddWorkItemCommentCmd) Execute(ctx context.Context, in *AddCommentVo) (*AddCommentResultVo, error) {

	workItem, err := cmd.witemRepo.GetWorkItem(ctx, in.WorkItemId, nil, nil)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	newComment, err := cmd.commentService.CreateComment(ctx, &comment_service.CreateCommentRequest{
		UserId:         in.OperUid,
		WorkItemId:     in.WorkItemId,
		Content:        in.Content,
		ReferUserIds:   in.ReferUserIds,
		ReplyCommentId: in.ReplyCommentId,
	})

	if err != nil {
		return nil, err
	}

	err = cmd.tm.InTx(ctx, func(ctx context.Context) error {

		err = cmd.witemRepo.SaveWorkItem(ctx, workItem)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		_, err := cmd.witemRepo.IncrCommentNum(ctx, workItem.Id, 1)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = cmd.commentRepo.CreateComment(ctx, newComment)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 增加未读数
	allMemberIds, _ := cmd.spaceMemberRepo.GetSpaceAllMemberIds(ctx, in.SpaceId)
	_ = cmd.commentRepo.IncrUnreadNumForUser(ctx, in.WorkItemId, allMemberIds)
	_ = cmd.commentRepo.SetUserReadTime(ctx, in.OperUid, in.WorkItemId, time.Now())

	newComment.AddMessage(shared.UserOper(newComment.UserId), &domain_message.CreateComment{
		WorkItemId:     in.WorkItemId,
		CommentId:      newComment.Id,
		UserId:         newComment.UserId,
		Content:        newComment.Content,
		ReferUserIds:   newComment.ReferUserIds,
		ReplyCommentId: newComment.ReplyCommentId,
	})

	cmd.domainMessageProducer.Send(ctx, newComment.GetMessages())

	return &AddCommentResultVo{
		CommentId: newComment.Id,
	}, nil
}
