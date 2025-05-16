package service

import (
	"context"
	"go-cs/api/comm"
	v1 "go-cs/api/log/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"

	"github.com/go-kratos/kratos/v2/log"
)

type LogService struct {
	v1.UnimplementedLogServer
	log *log.Helper
	uc  *biz.LogUsecase
}

func NewLogService(stu *biz.LogUsecase, logger log.Logger) *LogService {
	moduleName := "LogService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &LogService{
		uc:  stu,
		log: hlog,
	}
}

func (s *LogService) LoginLogList(ctx context.Context, req *v1.LoginLogListRequest) (*v1.LoginLogListReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.LoginLogListReply, error) {
		return &v1.LoginLogListReply{Result: &v1.LoginLogListReply_Error{Error: err}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId
	if uid == 0 {
		return nil, errs.NotLogin(ctx)
	}

	data, err := s.uc.UserLoginLogList(ctx, uid, req)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.LoginLogListReply{Result: &v1.LoginLogListReply_Data{Data: data}}, nil
}

func (s *LogService) OpLogList(ctx context.Context, req *v1.OpLogListRequest) (*v1.OpLogListReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.OpLogListReply, error) {
		return &v1.OpLogListReply{Result: &v1.OpLogListReply_Error{Error: err}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId
	if uid == 0 {
		return nil, errs.NotLogin(ctx)
	}

	data, err := s.uc.OpLogList(ctx, uid, req)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.OpLogListReply{Result: &v1.OpLogListReply_Data{Data: data}}, nil
}

func (s *LogService) SpaceOpLogList(ctx context.Context, req *v1.SpaceOpLogListRequest) (*v1.SpaceOpLogListReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.SpaceOpLogListReply, error) {
		return &v1.SpaceOpLogListReply{Result: &v1.SpaceOpLogListReply_Error{Error: err}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId
	if uid == 0 {
		return nil, errs.NotLogin(ctx)
	}

	data, err := s.uc.SpaceOpLogList(ctx, uid, req)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.SpaceOpLogListReply{Result: &v1.SpaceOpLogListReply_Data{Data: data}}, nil
}

func (s *LogService) SystemOpLogList(ctx context.Context, req *v1.SystemOpLogListRequest) (*v1.SystemOpLogListReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.SystemOpLogListReply, error) {
		return &v1.SystemOpLogListReply{Error: err}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId
	if uid == 0 {
		return nil, errs.NotLogin(ctx)
	}

	data, err := s.uc.SystemOpLogList(ctx, uid, req)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.SystemOpLogListReply{Data: data}, nil
}

func (s *LogService) PersonalOpLogList(ctx context.Context, req *v1.PersonalOpLogListRequest) (*v1.PersonalOpLogListReply, error) {
	var reply = func(err *comm.ErrorInfo) (*v1.PersonalOpLogListReply, error) {
		return &v1.PersonalOpLogListReply{Error: err}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId
	if uid == 0 {
		return nil, errs.NotLogin(ctx)
	}

	data, err := s.uc.PersonalOpLogList(ctx, uid, req)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &v1.PersonalOpLogListReply{Data: data}, nil
}
