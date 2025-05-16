package service

import (
	"context"
	"go-cs/api/comm"
	pb "go-cs/api/space/v1"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/stream"
	"slices"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceService struct {
	pb.UnimplementedSpaceServer

	uc *biz.SpaceUsecase
	sm *biz.SpaceMemberUsecase

	log *log.Helper
}

func NewSpaceService(spaceUsecase *biz.SpaceUsecase, sm *biz.SpaceMemberUsecase, logger log.Logger) *SpaceService {
	moduleName := "SpaceService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceService{
		uc:  spaceUsecase,
		sm:  sm,
		log: hlog,
	}
}

func (s *SpaceService) CreateSpace(ctx context.Context, req *pb.CreateSpaceRequest) (*pb.CreateSpaceReply, error) {

	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceName, "required,runeLen=1-255"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
		errReply := &pb.CreateSpaceReply{Result: &pb.CreateSpaceReply_Error{Error: errInfo}}
		return errReply, nil
	}

	// if vaildErr = validate.Var(req.Describe, "max=500"); vaildErr != nil {
	// 	//参数检查失败
	// 	errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
	// 	errReply := &pb.CreateSpaceReply{Result: &pb.CreateSpaceReply_Error{Error: errInfo}}
	// 	return errReply, nil
	// }

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.CreateSpaceReply{Result: &pb.CreateSpaceReply_Error{Error: errInfo}}
		return errReply, nil
	}

	var members []*db.SpaceMember
	for i := 0; i < len(req.Users); i++ {

		members = append(members, &db.SpaceMember{
			UserId: req.Users[i].UserId,
			RoleId: int64(req.Users[i].RoleId),
		})
	}

	//进入逻辑
	out, err := s.uc.CreateMySpace(ctx, loginUser, strings.TrimSpace(req.SpaceName), strings.TrimSpace(req.Describe), members)
	if err != nil {
		errReply := &pb.CreateSpaceReply{Result: &pb.CreateSpaceReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.CreateSpaceReply{Result: &pb.CreateSpaceReply_Data{
		Data: out.ToProto(),
	}}
	return okReply, nil
}

func (s *SpaceService) GetSpaceInfo(ctx context.Context, req *pb.GetSpaceInfoRequest) (*pb.GetSpaceInfoReply, error) {

	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
		errReply := &pb.GetSpaceInfoReply{Result: &pb.GetSpaceInfoReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.GetSpaceInfoReply{Result: &pb.GetSpaceInfoReply_Error{Error: errInfo}}
		return errReply, nil
	}

	out, err := s.uc.GetMySpace(ctx, loginUser.UserId, req.SpaceId)
	if err != nil {
		errReply := &pb.GetSpaceInfoReply{Result: &pb.GetSpaceInfoReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.GetSpaceInfoReply{Result: &pb.GetSpaceInfoReply_Data{Data: out}}
	return okReply, nil
}

func (s *SpaceService) SetSpaceDescribe(ctx context.Context, req *pb.SetSpaceDescribeRequest) (*pb.SetSpaceDescribeReply, error) {

	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
		errReply := &pb.SetSpaceDescribeReply{Result: &pb.SetSpaceDescribeReply_Error{Error: errInfo}}
		return errReply, nil
	}

	// if vaildErr = validate.Var(req.Describe, "max=500"); vaildErr != nil {
	// 	//参数检查失败
	// 	errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
	// 	errReply := &pb.SetSpaceDescribeReply{Result: &pb.SetSpaceDescribeReply_Error{Error: errInfo}}
	// 	return errReply, nil
	// }

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.SetSpaceDescribeReply{Result: &pb.SetSpaceDescribeReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//进入逻辑
	_, err := s.uc.SetMySpaceDescribe(ctx, loginUser, req.SpaceId, strings.TrimSpace(req.Describe))

	if err != nil {
		errReply := &pb.SetSpaceDescribeReply{Result: &pb.SetSpaceDescribeReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.SetSpaceDescribeReply{Result: &pb.SetSpaceDescribeReply_Data{
		Data: &pb.SetSpaceDescribeReplyData{},
	}}
	return okReply, nil
}

func (s *SpaceService) SpaceList(ctx context.Context, req *pb.SpaceListRequest) (*pb.SpaceListReply, error) {

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.SpaceListReply{Result: &pb.SpaceListReply_Error{Error: errInfo}}
		return errReply, nil
	}

	data, _ := s.uc.GetMySpaceList(ctx, loginUser.UserId)

	return &pb.SpaceListReply{Result: &pb.SpaceListReply_Data{Data: data}}, nil
}

func (s *SpaceService) SetSpaceBaseInfo(ctx context.Context, req *pb.SetSpaceBaseInfoRequest) (*pb.SetSpaceBaseInfoReply, error) {

	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "SpaceId")
		errReply := &pb.SetSpaceBaseInfoReply{Result: &pb.SetSpaceBaseInfoReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if vaildErr = validate.Var(req.SpaceName, "required,utf8Len=2-20"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "SpaceName")
		errReply := &pb.SetSpaceBaseInfoReply{Result: &pb.SetSpaceBaseInfoReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.SetSpaceBaseInfoReply{Result: &pb.SetSpaceBaseInfoReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//进入逻辑
	_, err := s.uc.SetMySpaceBaseInfo(ctx, loginUser, req.SpaceId, strings.TrimSpace(req.SpaceName), strings.TrimSpace(req.Describe))

	if err != nil {
		errReply := &pb.SetSpaceBaseInfoReply{Result: &pb.SetSpaceBaseInfoReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.SetSpaceBaseInfoReply{Result: &pb.SetSpaceBaseInfoReply_Data{
		Data: &pb.SetSpaceBaseInfoReplyData{},
	}}
	return okReply, nil
}

func (s *SpaceService) DelSpace(ctx context.Context, req *pb.DelSpaceRequest) (*pb.DelSpaceReply, error) {

	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "SpaceId")
		errReply := &pb.DelSpaceReply{Result: &pb.DelSpaceReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.DelSpaceReply{Result: &pb.DelSpaceReply_Error{Error: errInfo}}
		return errReply, nil
	}

	out, err := s.uc.DelSpace(ctx, loginUser, req.SpaceId, req.Scene)
	if err != nil {
		errReply := &pb.DelSpaceReply{Result: &pb.DelSpaceReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.DelSpaceReply{Result: &pb.DelSpaceReply_Data{Data: out.ToProto()}}
	return okReply, nil
}

func (s *SpaceService) QuitMySpace(ctx context.Context, req *pb.QuitMySpaceRequest) (*pb.QuitMySpaceReply, error) {
	loginUser := utils.GetLoginUser(ctx)
	if loginUser.UserId == 0 { //用户信息获取失败
		err := errs.NotLogin(ctx)
		errReply := &pb.QuitMySpaceReply{Result: &pb.QuitMySpaceReply_Error{Error: err}}
		return errReply, nil
	}

	uid := loginUser.UserId
	if req.SpaceId <= 0 {
		err := errs.Param(ctx, "spaceId")
		errReply := &pb.QuitMySpaceReply{Result: &pb.QuitMySpaceReply_Error{Error: err}}
		return errReply, nil
	}

	if req.UserId == 0 {
		req.UserId = uid
	}

	if req.TargetUserId < 0 || req.TargetUserId == req.UserId {
		err := errs.Param(ctx, "targetUserId")
		errReply := &pb.QuitMySpaceReply{Result: &pb.QuitMySpaceReply_Error{Error: err}}
		return errReply, nil
	}

	err := s.sm.QuitMySpaceV2(ctx, uid, req.SpaceId, req.UserId, req.TargetUserId)
	if err != nil {
		errReply := &pb.QuitMySpaceReply{Result: &pb.QuitMySpaceReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.QuitMySpaceReply{Result: &pb.QuitMySpaceReply_Data{Data: &pb.QuitMySpaceReplyData{}}}
	return okReply, nil
}

func (s *SpaceService) TransferSpaceOwnership(ctx context.Context, req *pb.TransferSpaceOwnershipRequest) (*pb.TransferSpaceOwnershipReply, error) {
	reply := func(err *comm.ErrorInfo) (*pb.TransferSpaceOwnershipReply, error) {
		return &pb.TransferSpaceOwnershipReply{Result: &pb.TransferSpaceOwnershipReply_Error{Error: err}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser.UserId == 0 { //用户信息获取失败
		return reply(errs.NotLogin(ctx))
	}

	uid := loginUser.UserId
	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "spaceId"))
	}

	if req.UserId <= 0 || req.UserId == uid {
		return reply(errs.Param(ctx, "UserId"))
	}

	if req.SrcUserId == 0 {
		req.SrcUserId = uid
	}

	err := s.uc.TransferSpaceOwnership(ctx, uid, req.SpaceId, req.SrcUserId, req.UserId)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &pb.TransferSpaceOwnershipReply{Result: nil}, nil
}

func (s *SpaceService) SetNotify(ctx context.Context, req *pb.SetNotifyRequest) (*pb.SetNotifyReply, error) {
	reply := func(err error) (*pb.SetNotifyReply, error) {
		return &pb.SetNotifyReply{Result: &pb.SetNotifyReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if !slices.Contains([]int64{0, 1}, req.Notify) {
		return reply(errs.Param(ctx, "Notify"))
	}

	err := s.uc.SetNotify(ctx, utils.GetLoginUser(ctx), req.SpaceId, req.Notify)
	if err != nil {
		return reply(err)
	}

	return &pb.SetNotifyReply{}, nil
}

func (s *SpaceService) SetWorkingDay(ctx context.Context, req *pb.SetWorkingDayRequest) (*pb.SetWorkingDayReply, error) {
	reply := func(err error) (*pb.SetWorkingDayReply, error) {
		return &pb.SetWorkingDayReply{Result: &pb.SetWorkingDayReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if len(req.WeekDays) == 0 || len(stream.Diff(req.WeekDays, []int64{0, 1, 2, 3, 4, 5, 6})) != 0 {
		return reply(errs.Param(ctx, "WeekDays"))
	}

	slices.SortFunc(req.WeekDays, func(a, b int64) int {
		if a == 0 {
			a = 7
		}
		if b == 0 {
			b = 7
		}

		return int(a - b)
	})

	err := s.uc.SetWorkingDay(ctx, utils.GetLoginUser(ctx), req.SpaceId, req.WeekDays)
	if err != nil {
		return reply(err)
	}

	return &pb.SetWorkingDayReply{}, nil
}

func (s *SpaceService) SetName(ctx context.Context, req *pb.SetNameRequest) (*pb.SetNameReply, error) {
	reply := func(err error) (*pb.SetNameReply, error) {
		return &pb.SetNameReply{Result: &pb.SetNameReply_Error{Error: errs.Cast(err)}}, nil
	}

	validate := utils.NewValidator()
	if err := validate.Var(req.SpaceId, "required,number,gt=0"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if err := validate.Var(req.SpaceName, "required,runeLen=1-255"); err != nil {
		return reply(errs.Param(ctx, "SpaceName"))
	}

	//进入逻辑
	_, err := s.uc.SetName(ctx, utils.GetLoginUser(ctx), req.SpaceId, req.SpaceName)
	if err != nil {
		return reply(err)
	}
	return &pb.SetNameReply{}, nil
}

func (s *SpaceService) SearchWorkItem(ctx context.Context, req *pb.SearchWorkItemRequest) (*pb.SearchWorkItemReply, error) {
	reply := func(err error) (*pb.SearchWorkItemReply, error) {
		return &pb.SearchWorkItemReply{Result: &pb.SearchWorkItemReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.Keyword == "" {
		return reply(errs.Param(ctx, "Keyword"))
	}

	//进入逻辑
	data, err := s.uc.SearchWorkItem(ctx, utils.GetLoginUser(ctx), req.SpaceId, req.Keyword)
	if err != nil {
		return reply(err)
	}
	return &pb.SearchWorkItemReply{Result: &pb.SearchWorkItemReply_Data{
		Data: data,
	}}, nil
}

func (s *SpaceService) GetWorkItemTypes(ctx context.Context, req *pb.GetWorkItemTypesRequest) (*pb.GetWorkItemTypesReply, error) {
	reply := func(err error) (*pb.GetWorkItemTypesReply, error) {
		return &pb.GetWorkItemTypesReply{Result: &pb.GetWorkItemTypesReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	//进入逻辑
	items, err := s.uc.GetWorkItemTypes(ctx, utils.GetLoginUser(ctx), req)
	if err != nil {
		return reply(err)
	}

	return &pb.GetWorkItemTypesReply{Result: &pb.GetWorkItemTypesReply_Data{
		Data: &pb.GetWorkItemTypesReplyData{List: items},
	}}, nil
}

func (s *SpaceService) SetCommentDeletable(ctx context.Context, req *pb.SetCommentDeletableRequest) (*pb.SetCommentDeletableReply, error) {
	reply := func(err error) (*pb.SetCommentDeletableReply, error) {
		return &pb.SetCommentDeletableReply{Result: &pb.SetCommentDeletableReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.CommentDeletable != 0 && req.CommentDeletable != 1 {
		return reply(errs.Param(ctx, "CommentDeletable"))
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.SetCommentDeletable(ctx, uid, req.SpaceId, req.CommentDeletable)
	if err != nil {
		return reply(err)
	}

	return &pb.SetCommentDeletableReply{}, nil
}

func (s *SpaceService) GetTempConfig(ctx context.Context, req *pb.GetTempConfigRequest) (*pb.GetTempConfigReply, error) {
	reply := func(err error) (*pb.GetTempConfigReply, error) {
		return &pb.GetTempConfigReply{Error: errs.Cast(err)}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	data, err := s.uc.GetTempConfig(ctx, uid, req.SpaceId, req.GetKeys())
	if err != nil {
		return reply(err)
	}

	return &pb.GetTempConfigReply{
		Data: data,
	}, nil
}

func (s *SpaceService) SetTempConfig(ctx context.Context, req *pb.SetTempConfigRequest) (*pb.SetTempConfigReply, error) {
	reply := func(err error) (*pb.SetTempConfigReply, error) {
		return &pb.SetTempConfigReply{Error: errs.Cast(err)}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.SetTempConfig(ctx, uid, req.SpaceId, req.GetConfigs())
	if err != nil {
		return reply(err)
	}

	return &pb.SetTempConfigReply{}, nil
}

func (s *SpaceService) DelTempConfig(ctx context.Context, req *pb.DelTempConfigRequest) (*pb.DelTempConfigReply, error) {
	reply := func(err error) (*pb.DelTempConfigReply, error) {
		return &pb.DelTempConfigReply{Error: errs.Cast(err)}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.DelTempConfig(ctx, uid, req.SpaceId, req.Keys)
	if err != nil {
		return reply(err)
	}

	return &pb.DelTempConfigReply{}, nil
}

func (s *SpaceService) SetCommentDeletableWhenArchived(ctx context.Context, req *pb.SetCommentDeletableWhenArchivedRequest) (*pb.SetCommentDeletableWhenArchivedReply, error) {
	reply := func(err error) (*pb.SetCommentDeletableWhenArchivedReply, error) {
		return &pb.SetCommentDeletableWhenArchivedReply{Error: errs.Cast(err)}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.Value != 0 && req.Value != 1 {
		return reply(errs.Param(ctx, "Value"))
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.SetCommentDeletableWhenArchived(ctx, uid, req.SpaceId, req.Value)
	if err != nil {
		return reply(err)
	}

	return &pb.SetCommentDeletableWhenArchivedReply{}, nil
}

func (s *SpaceService) SetCommentShowPos(ctx context.Context, req *pb.SetCommentShowPosRequest) (*pb.SetCommentShowPosReply, error) {
	reply := func(err error) (*pb.SetCommentShowPosReply, error) {
		return &pb.SetCommentShowPosReply{Error: errs.Cast(err)}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.Value != 0 && req.Value != 1 {
		return reply(errs.Param(ctx, "Value"))
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.SetCommentShowPos(ctx, uid, req.SpaceId, req.Value)
	if err != nil {
		return reply(err)
	}

	return &pb.SetCommentShowPosReply{}, nil
}

func (s *SpaceService) Copy(ctx context.Context, req *pb.CopyRequest) (*pb.CopyReply, error) {
	reply := func(err error) (*pb.CopyReply, error) {
		return &pb.CopyReply{Error: errs.Cast(err)}, nil
	}

	if req.SrcSpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	oper := utils.GetLoginUser(ctx)

	data, err := s.uc.Copy(ctx, oper, req.SrcSpaceId, req.SpaceName, req.SpaceDescribe)
	if err != nil {
		return reply(err)
	}

	return &pb.CopyReply{Data: data}, nil
}
