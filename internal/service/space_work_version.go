package service

import (
	"context"
	"go-cs/api/comm"
	pb "go-cs/api/space_work_version/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceWorkVersionService struct {
	pb.UnimplementedSpaceWorkVersionServer

	uc  *biz.SpaceWorkVersionUsecase
	log *log.Helper
}

func NewSpaceWorkVersionService(uc *biz.SpaceWorkVersionUsecase, logger log.Logger) *SpaceWorkVersionService {
	moduleName := "SpaceWorkVersionService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceWorkVersionService{
		uc:  uc,
		log: hlog,
	}
}

func (s *SpaceWorkVersionService) CreateSpaceWorkVersion(ctx context.Context, req *pb.CreateSpaceWorkVersionRequest) (*pb.CreateSpaceWorkVersionReply, error) {

	reply := func(err error) (*pb.CreateSpaceWorkVersionReply, error) {
		return &pb.CreateSpaceWorkVersionReply{Result: &pb.CreateSpaceWorkVersionReply_Error{Error: errs.Cast(err)}}, nil
	}

	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(strings.TrimSpace(req.VersionName), "required,utf8Len=2-20"); vaildErr != nil {
		return reply(errs.Param(ctx, "请输入有效格式（2 ～ 20个字符）"))
	}

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	//进入逻辑
	out, err := s.uc.CreateMySpaceWorkVersion(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.CreateSpaceWorkVersionReply{Result: &pb.CreateSpaceWorkVersionReply_Data{
		Data: out.ToProto(),
	}}
	return okReply, nil
}

func (s *SpaceWorkVersionService) SpaceWorkVersionList(ctx context.Context, req *pb.SpaceWorkVersionListRequest) (*pb.SpaceWorkVersionListReply, error) {

	reply := func(err error) (*pb.SpaceWorkVersionListReply, error) {
		return &pb.SpaceWorkVersionListReply{Result: &pb.SpaceWorkVersionListReply_Error{Error: errs.Cast(err)}}, nil
	}

	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)

	//进入逻辑
	out, err := s.uc.QSpaceWorkVersionList(ctx, loginUser.UserId, req.SpaceId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceWorkVersionListReply{Result: &pb.SpaceWorkVersionListReply_Data{
		Data: out,
	}}
	return okReply, nil
}

func (s *SpaceWorkVersionService) SpaceWorkVersionById(ctx context.Context, req *pb.SpaceWorkVersionByIdRequest) (*pb.SpaceWorkVersionByIdReply, error) {

	reply := func(err error) (*pb.SpaceWorkVersionByIdReply, error) {
		return &pb.SpaceWorkVersionByIdReply{Result: &pb.SpaceWorkVersionByIdReply_Error{Error: errs.Cast(err)}}, nil
	}

	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if len(req.Ids) == 0 {
		return reply(errs.Param(ctx, "Ids"))
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	out, err := s.uc.QSpaceWorkVersionById(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceWorkVersionByIdReply{Result: &pb.SpaceWorkVersionByIdReply_Data{
		Data: out,
	}}
	return okReply, nil
}

func (s *SpaceWorkVersionService) ModifySpaceWorkVersionName(ctx context.Context, req *pb.ModifySpaceWorkVersionNameRequest) (*pb.ModifySpaceWorkVersionNameReply, error) {
	var vaildErr error

	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.VersionId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
		errReply := &pb.ModifySpaceWorkVersionNameReply{Result: &pb.ModifySpaceWorkVersionNameReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if vaildErr = validate.Var(strings.TrimSpace(req.VersionName), "required,utf8Len=1-20"); vaildErr != nil {
		//参数检查失败
		errInfo := errs.Business(ctx, "名称长度必须为1-20个字符")
		errReply := &pb.ModifySpaceWorkVersionNameReply{Result: &pb.ModifySpaceWorkVersionNameReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.ModifySpaceWorkVersionNameReply{Result: &pb.ModifySpaceWorkVersionNameReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//进入逻辑
	_, err := s.uc.SetMySpaceWorkVersionName(ctx, loginUser, req)
	if err != nil {
		errReply := &pb.ModifySpaceWorkVersionNameReply{Result: &pb.ModifySpaceWorkVersionNameReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.ModifySpaceWorkVersionNameReply{Result: &pb.ModifySpaceWorkVersionNameReply_Data{Data: ""}}
	return okReply, nil
}

func (s *SpaceWorkVersionService) DelSpaceWorkVersion(ctx context.Context, req *pb.DelSpaceWorkVersionRequest) (*pb.DelSpaceWorkVersionReply, error) {
	var vaildErr error

	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.VersionId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "VersionId")
		errReply := &pb.DelSpaceWorkVersionReply{Result: &pb.DelSpaceWorkVersionReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.DelSpaceWorkVersionReply{Result: &pb.DelSpaceWorkVersionReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//进入逻辑
	_, err := s.uc.DelMySpaceWorkVersion(ctx, loginUser, req.VersionId, req.ToVersionId)
	if err != nil {
		errReply := &pb.DelSpaceWorkVersionReply{Result: &pb.DelSpaceWorkVersionReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.DelSpaceWorkVersionReply{Result: &pb.DelSpaceWorkVersionReply_Data{Data: ""}}
	return okReply, nil
}

func (s *SpaceWorkVersionService) SetSpaceWorkVersionOrder(ctx context.Context, req *pb.SetWorkVersionOrderRequest) (*pb.SetWorkVersionOrderReply, error) {
	return &pb.SetWorkVersionOrderReply{}, nil
}

func (s *SpaceWorkVersionService) SpaceWorkVersionRelationCount(ctx context.Context, req *pb.SpaceWorkVersionRelationCountRequest) (*pb.SpaceWorkVersionRelationCountReply, error) {
	reply := func(err error) (*pb.SpaceWorkVersionRelationCountReply, error) {
		return &pb.SpaceWorkVersionRelationCountReply{Result: &pb.SpaceWorkVersionRelationCountReply_Error{Error: errs.Cast(err)}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	if req.SpaceId < 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.VersionId < 0 {
		return reply(errs.Param(ctx, "VersionId"))
	}

	count, err := s.uc.GetMySpaceWorkVersionRelationCount(ctx, uid, req.SpaceId, req.VersionId)
	if err != nil {
		return reply(err)
	}

	return &pb.SpaceWorkVersionRelationCountReply{Result: &pb.SpaceWorkVersionRelationCountReply_Data{Data: count}}, nil
}

func (s *SpaceWorkVersionService) SetSpaceWorkVersionRanking(ctx context.Context, req *pb.SetSpaceWorkVersionRankingRequest) (*pb.SetSpaceWorkVersionRankingReply, error) {
	reply := func(err error) (*pb.SetSpaceWorkVersionRankingReply, error) {
		return &pb.SetSpaceWorkVersionRankingReply{Result: &pb.SetSpaceWorkVersionRankingReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if len(req.List) == 0 {
		return reply(errs.Param(ctx, "List"))
	}

	rankList := make([]map[string]int64, 0)
	for _, v := range req.List {
		rankList = append(rankList, map[string]int64{
			"id":      v.Id,
			"ranking": v.Ranking,
		})
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.SetSpaceWorkWorkVersionRanking(ctx, loginUser, req.SpaceId, rankList)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetSpaceWorkVersionRankingReply{Result: &pb.SetSpaceWorkVersionRankingReply_Data{Data: ""}}
	return okReply, nil
}
