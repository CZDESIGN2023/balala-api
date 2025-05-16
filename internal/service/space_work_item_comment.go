package service

import (
	"context"
	"go-cs/api/comm"
	v1 "go-cs/api/space_work_item/v1"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
)

func (s *SpaceWorkItemService) AddWorkItemComment(ctx context.Context, req *v1.AddWorkItemCommentRequest) (*v1.AddWorkItemCommentReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.AddWorkItemCommentReply, error) {
		return &v1.AddWorkItemCommentReply{Result: &v1.AddWorkItemCommentReply_Error{Error: err}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId
	if uid == 0 {
		return reply(errs.NotLogin(ctx))
	}

	if req.WorkItemId <= 0 {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if req.Content == "" {
		return reply(errs.Param(ctx, "Content"))
	}

	if req.Content == "" {
		return reply(errs.Param(ctx, "Content"))
	}

	data, err := s.uc.AddComment(ctx, uid, req)
	if err != nil {
		if len(req.Content) >= 1000 {
			return reply(errs.Business(ctx, "输入已超出最大字符数量"))
		}
		return reply(errs.Cast(err))
	}

	return &v1.AddWorkItemCommentReply{Result: &v1.AddWorkItemCommentReply_Data{Data: data}}, nil
}

func (s *SpaceWorkItemService) UpdateWorkItemComment(ctx context.Context, req *v1.UpdateWorkItemCommentRequest) (*v1.UpdateWorkItemCommentReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.UpdateWorkItemCommentReply, error) {
		return &v1.UpdateWorkItemCommentReply{Result: &v1.UpdateWorkItemCommentReply_Error{Error: err}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId
	if uid == 0 {
		return nil, errs.NotLogin(ctx)
	}

	if req.Content == "" {
		return reply(errs.Param(ctx, "Content"))
	}

	data, err := s.uc.UpdateComment(ctx, uid, req)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.UpdateWorkItemCommentReply{Result: &v1.UpdateWorkItemCommentReply_Data{Data: data}}, nil
}

func (s *SpaceWorkItemService) WorkItemCommentList(ctx context.Context, req *v1.WorkItemCommentListRequest) (*v1.WorkItemCommentListReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.WorkItemCommentListReply, error) {
		return &v1.WorkItemCommentListReply{Result: &v1.WorkItemCommentListReply_Error{Error: err}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId
	if uid == 0 {
		return nil, errs.NotLogin(ctx)
	}

	if req.WorkItemId <= 0 {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if req.Order == "" {
		req.Order = "DESC"
	}

	if req.Order != "DESC" && req.Order != "ASC" {
		return reply(errs.Param(ctx, "Order"))
	}

	data, err := s.uc.CommentList(ctx, uid, req)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.WorkItemCommentListReply{Result: &v1.WorkItemCommentListReply_Data{Data: data}}, nil
}

func (s *SpaceWorkItemService) DeleteWorkItemComment(ctx context.Context, req *v1.DeleteWorkItemCommentRequest) (*v1.DeleteWorkItemCommentReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.DeleteWorkItemCommentReply, error) {
		return &v1.DeleteWorkItemCommentReply{Result: &v1.DeleteWorkItemCommentReply_Error{Error: err}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	if req.Id <= 0 {
		return reply(errs.Param(ctx, "Id"))
	}

	err := s.uc.DeleteComment(ctx, uid, req.Id)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.DeleteWorkItemCommentReply{Result: nil}, nil
}

func (s *SpaceWorkItemService) CommentEmojiAdd(ctx context.Context, req *v1.CommentEmojiAddRequest) (*v1.CommentEmojiAddReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.CommentEmojiAddReply, error) {
		return &v1.CommentEmojiAddReply{Result: &v1.CommentEmojiAddReply_Error{Error: err}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	if req.CommentId <= 0 {
		return reply(errs.Param(ctx, "Id"))
	}

	if req.EmojiId == "" {
		return reply(errs.Param(ctx, "EmojiId"))
	}

	err := s.uc.CommentEmojiAdd(ctx, uid, req.CommentId, req.EmojiId)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.CommentEmojiAddReply{Result: nil}, nil
}

func (s *SpaceWorkItemService) CommentEmojiRemove(ctx context.Context, req *v1.CommentEmojiRemoveRequest) (*v1.CommentEmojiRemoveReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.CommentEmojiRemoveReply, error) {
		return &v1.CommentEmojiRemoveReply{Result: &v1.CommentEmojiRemoveReply_Error{Error: err}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	if req.CommentId <= 0 {
		return reply(errs.Param(ctx, "Id"))
	}

	if req.EmojiId == "" {
		return reply(errs.Param(ctx, "EmojiId"))
	}

	err := s.uc.CommentEmojiRemove(ctx, uid, req.CommentId, req.EmojiId)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.CommentEmojiRemoveReply{Result: nil}, nil
}
