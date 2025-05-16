package biz

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	v1 "go-cs/api/space_work_item/v1"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/biz/command"
	biz_utils "go-cs/internal/biz/pkg"
	"go-cs/internal/consts"
	perm "go-cs/internal/domain/perm"
	comment "go-cs/internal/domain/space_work_item_comment"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/locker"
	"go-cs/pkg/stream"
	"math"
	"time"
)

func (s *SpaceWorkItemUsecase) AddComment(ctx context.Context, uid int64, req *v1.AddWorkItemCommentRequest) (*v1.AddWorkItemCommentReplyData, error) {

	req.ReferUserIds = stream.Unique(req.ReferUserIds) //去重

	workItemId := req.WorkItemId

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, nil, nil)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	_, err = s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, uid)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	r, err := s.addWorkItemCommentCommand.Execute(ctx, &command.AddCommentVo{
		OperUid:        uid,
		SpaceId:        workItem.SpaceId,
		WorkItemId:     workItemId,
		Content:        req.Content,
		ReferUserIds:   req.ReferUserIds,
		ReplyCommentId: req.ReplyCommentId,
	})

	if err != nil {
		return nil, err
	}

	return &v1.AddWorkItemCommentReplyData{
		Id: r.CommentId,
	}, nil
}

func (s *SpaceWorkItemUsecase) UpdateComment(ctx context.Context, uid int64, req *v1.UpdateWorkItemCommentRequest) (*v1.UpdateWorkItemCommentReplyData, error) {
	itemComment, err := s.commentRepo.GetComment(ctx, req.Id)
	if err != nil && !errs.IsDbRecordNotFoundErr(err) {
		return nil, errs.Internal(ctx, err)
	}
	if errs.IsDbRecordNotFoundErr(err) {
		return nil, errs.Business(ctx, "此评论已删除")
	}
	if itemComment.UserId != uid {
		return nil, errs.NoPerm(ctx)
	}

	itemComment.UpdateContent(req.Content, req.ReferUserIds, shared.UserOper(uid))

	err = s.commentRepo.SaveComment(ctx, itemComment)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	_ = s.commentRepo.SetUserReadTime(ctx, uid, itemComment.WorkItemId, time.Now())

	s.domainMessageProducer.Send(ctx, itemComment.GetMessages())

	return &v1.UpdateWorkItemCommentReplyData{
		Id:        itemComment.Id,
		Content:   itemComment.Content,
		CreatedAt: itemComment.CreatedAt,
		UpdatedAt: itemComment.UpdatedAt,
	}, nil
}

func (s *SpaceWorkItemUsecase) CommentList(ctx context.Context, uid int64, req *v1.WorkItemCommentListRequest) (*v1.WorkItemCommentListReplyData, error) {

	workItemId := req.WorkItemId
	workItem, err := s.repo.GetWorkItem(ctx, workItemId, nil, nil)
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, uid)
	if member == nil || err != nil {
		return nil, errs.NoPerm(ctx)
	}

	if req.Pos == 0 {
		err := s.commentRepo.SetUserReadTime(ctx, uid, workItemId, time.Now())
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}
	}

	if req.Pos == 0 && req.Order == "DESC" {
		req.Pos = math.MaxInt32
	}

	if req.Size == 0 {
		req.Size = 20
	}

	var pos = int(req.Pos)
	var size = int(req.Size)

	list, err := s.commentRepo.QCommentPagination(ctx, req.WorkItemId, pos, size+1, req.Order)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var hasNext bool
	var nextPos int64
	if len(list) > size {
		hasNext = true
		nextPos = list[size-1].Id
		list = list[:size]
	}

	var replyCommentIds []int64
	for _, v := range list {
		if v.ReplyCommentId != 0 {
			replyCommentIds = append(replyCommentIds, v.ReplyCommentId)
		}
	}
	replyCommentMap, err := s.commentRepo.CommentMap(ctx, replyCommentIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var replyUserIds []int64
	for _, v := range replyCommentMap {
		if v.UserId != 0 {
			replyUserIds = append(replyUserIds, v.UserId)
		}
	}

	userIds := stream.Map(list, func(v *comment.SpaceWorkItemComment) int64 {
		return v.UserId
	})

	emojiUserIds := stream.Flat(stream.Map(list, func(v *comment.SpaceWorkItemComment) []int64 {
		return stream.Flat(stream.Map(v.Emojis, func(v comment.EmojiInfo) []int64 {
			return utils.ToInt64Array(v.UserIds)
		}))
	}))

	allUserIds := stream.Of(userIds).Concat(emojiUserIds...).Concat(replyUserIds...).Unique().List()
	userMap, err := s.userRepo.UserMap(ctx, allUserIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var items []*v1.WorkItemCommentListReplyData_Item
	for _, v := range list {
		replyComment := replyCommentMap[v.ReplyCommentId]
		var apiReplyComment *v1.WorkItemCommentListReplyData_Item

		if replyComment != nil {
			apiReplyComment = &v1.WorkItemCommentListReplyData_Item{
				Id:        replyComment.Id,
				Content:   replyComment.Content,
				CreatedAt: replyComment.CreatedAt,
				UpdatedAt: replyComment.UpdatedAt,
				User:      biz_utils.ToSimpleUser(userMap[replyComment.UserId]),
			}
		}

		emojis := stream.Map(v.Emojis, func(v comment.EmojiInfo) *v1.WorkItemCommentListReplyData_Emoji {
			users := stream.Map(v.UserIds, func(v string) *rsp.SimpleUserInfo {
				return biz_utils.ToSimpleUser(userMap[cast.ToInt64(v)])
			})
			return &v1.WorkItemCommentListReplyData_Emoji{
				Id:    v.Id,
				Users: users,
			}
		})

		items = append(items, &v1.WorkItemCommentListReplyData_Item{
			Id:             v.Id,
			Content:        v.Content,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			User:           biz_utils.ToSimpleUser(userMap[v.UserId]),
			ReplyComment:   apiReplyComment,
			ReplyCommentId: v.ReplyCommentId,
			Emojis:         emojis,
		})
	}

	return &v1.WorkItemCommentListReplyData{
		Items:   items,
		HasNext: hasNext,
		NextPos: nextPos,
		Order:   req.Order,
		Total:   int64(workItem.CommentNum),
	}, nil
}

func (s *SpaceWorkItemUsecase) DeleteComment(ctx context.Context, uid int64, id int64) error {
	itemComment, err := s.commentRepo.GetComment(ctx, id)
	if err != nil && !errs.IsDbRecordNotFoundErr(err) {
		return errs.Internal(ctx, err)
	}
	if errs.IsDbRecordNotFoundErr(err) {
		return errs.Business(ctx, "此评论已删除")
	}

	workItem, err := s.repo.GetWorkItem(ctx, itemComment.WorkItemId, nil, nil)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, uid)
	if member == nil || err != nil {
		return errs.NoPerm(ctx)
	}

	if !perm.Instance().Check(member.RoleId, consts.PERM_DeleteComment) && uid != itemComment.UserId {
		return errs.NoPerm(ctx)
	}

	config, err := s.spaceRepo.GetSpaceConfig(ctx, workItem.SpaceId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if config.CommentDeletable == 0 {
		return errs.Business(ctx, "评论删除功能已关闭，请联系项目管理员开启")
	}

	status, err := s.witemStatusService.GetWorkItemStatusItem(ctx, workItem.WorkItemStatus.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if status.StatusType == consts.WorkItemStatusType_Archived && config.CommentDeletableWhenArchived == 0 {
		return errs.Business(ctx, "任务已归档，无法删除评论")
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) error {
		_, err = s.commentRepo.DelComment(ctx, id)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		_, err := s.repo.IncrCommentNum(ctx, itemComment.WorkItemId, -1)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = s.repo.SaveWorkItem(ctx, workItem)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	itemComment.OnDelete(shared.UserOper(uid))

	s.domainMessageProducer.Send(ctx, itemComment.GetMessages())

	return nil
}

func (s *SpaceWorkItemUsecase) CommentEmojiAdd(ctx context.Context, uid int64, commentId int64, emoji string) error {
	lock := locker.Lock(fmt.Sprintf("commentEmoji:%v", commentId))
	lock.Lock()
	defer lock.Unlock()

	itemComment, err := s.commentRepo.GetComment(ctx, commentId)
	if err != nil && !errs.IsDbRecordNotFoundErr(err) {
		return errs.Internal(ctx, err)
	}

	if errs.IsDbRecordNotFoundErr(err) {
		return errs.Business(ctx, "此评论已删除")
	}

	itemComment.AddEmoji(emoji, uid, shared.UserOper(uid))

	err = s.commentRepo.SaveComment(ctx, itemComment)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, itemComment.GetMessages())

	return nil
}

func (s *SpaceWorkItemUsecase) CommentEmojiRemove(ctx context.Context, uid int64, commentId int64, emojiId string) error {
	lock := locker.Lock(fmt.Sprintf("commentEmoji:%v", commentId))
	lock.Lock()
	defer lock.Unlock()

	itemComment, err := s.commentRepo.GetComment(ctx, commentId)
	if err != nil && !errs.IsDbRecordNotFoundErr(err) {
		return errs.Internal(ctx, err)
	}

	if errs.IsDbRecordNotFoundErr(err) {
		return errs.Business(ctx, "此评论已删除")
	}

	itemComment.RemoveEmoji(emojiId, uid, shared.UserOper(uid))

	err = s.commentRepo.SaveComment(ctx, itemComment)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, itemComment.GetMessages())

	return nil
}
