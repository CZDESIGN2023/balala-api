package service

import (
	"context"
	pb "go-cs/api/space_work_object/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceWorkObjectService struct {
	pb.UnimplementedSpaceWorkObjectServer

	uc  *biz.SpaceWorkObjectUsecase
	log *log.Helper
}

func NewSpaceWorkObjectService(spaceWorkObjectUsecase *biz.SpaceWorkObjectUsecase, logger log.Logger) *SpaceWorkObjectService {
	moduleName := "SpaceWorkObjectService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceWorkObjectService{
		uc:  spaceWorkObjectUsecase,
		log: hlog,
	}
}

func (s *SpaceWorkObjectService) CreateSpaceWorkObject(ctx context.Context, req *pb.CreateSpaceWorkObjectRequest) (*pb.CreateSpaceWorkObjectReply, error) {

	reply := func(err error) (*pb.CreateSpaceWorkObjectReply, error) {
		return &pb.CreateSpaceWorkObjectReply{Result: &pb.CreateSpaceWorkObjectReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(strings.TrimSpace(req.WorkObjectName), "required,utf8Len=1-20"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkObjectName"))
	}

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	//进入逻辑
	out, err := s.uc.CreateMySpaceWorkObject(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.CreateSpaceWorkObjectReply{Result: &pb.CreateSpaceWorkObjectReply_Data{
		Data: out.ToProto(),
	}}
	return okReply, nil
}

func (s *SpaceWorkObjectService) SpaceWorkObjectList(ctx context.Context, req *pb.SpaceWorkObjectListRequest) (*pb.SpaceWorkObjectListReply, error) {

	reply := func(err error) (*pb.SpaceWorkObjectListReply, error) {
		return &pb.SpaceWorkObjectListReply{Result: &pb.SpaceWorkObjectListReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	//进入逻辑
	out, err := s.uc.QSpaceWorkObjectList(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceWorkObjectListReply{Result: &pb.SpaceWorkObjectListReply_Data{
		Data: out,
	}}
	return okReply, nil
}

func (s *SpaceWorkObjectService) SpaceWorkObjectById(ctx context.Context, req *pb.SpaceWorkObjectByIdRequest) (*pb.SpaceWorkObjectByIdReply, error) {

	reply := func(err error) (*pb.SpaceWorkObjectByIdReply, error) {
		return &pb.SpaceWorkObjectByIdReply{Result: &pb.SpaceWorkObjectByIdReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if len(req.Ids) == 0 {
		return reply(errs.Param(ctx, "Ids"))
	}

	//进入逻辑
	out, err := s.uc.QSpaceWorkObjectById(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceWorkObjectByIdReply{Result: &pb.SpaceWorkObjectByIdReply_Data{
		Data: out,
	}}
	return okReply, nil
}

func (s *SpaceWorkObjectService) ModifySpaceWorkObjectName(ctx context.Context, req *pb.ModifySpaceWorkObjectNameRequest) (*pb.ModifySpaceWorkObjectNameReply, error) {

	reply := func(err error) (*pb.ModifySpaceWorkObjectNameReply, error) {
		return &pb.ModifySpaceWorkObjectNameReply{Result: &pb.ModifySpaceWorkObjectNameReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkObjectId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkObjectId"))
	}

	if vaildErr = validate.Var(strings.TrimSpace(req.WorkObjectName), "required,utf8Len=1-20"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkObjectName"))
	}

	//进入逻辑
	_, err := s.uc.SetMySpaceWorkObjectName(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ModifySpaceWorkObjectNameReply{Result: &pb.ModifySpaceWorkObjectNameReply_Data{Data: &pb.ModifySpaceWorkObjectNameReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkObjectService) DelSpaceWorkObject(ctx context.Context, req *pb.DelSpaceWorkObjectRequest) (*pb.DelSpaceWorkObjectReply, error) {

	reply := func(err error) (*pb.DelSpaceWorkObjectReply, error) {
		return &pb.DelSpaceWorkObjectReply{Result: &pb.DelSpaceWorkObjectReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkObjectId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkObjectId"))
	}

	//进入逻辑
	_, err := s.uc.DelMySpaceWorkObject(ctx, loginUser, req.SpaceId, req.WorkObjectId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.DelSpaceWorkObjectReply{Result: &pb.DelSpaceWorkObjectReply_Data{Data: &pb.DelSpaceWorkObjectReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkObjectService) DelSpaceWorkObject2(ctx context.Context, req *pb.DelSpaceWorkObjectRequest2) (*pb.DelSpaceWorkObjectReply, error) {

	reply := func(err error) (*pb.DelSpaceWorkObjectReply, error) {
		return &pb.DelSpaceWorkObjectReply{Result: &pb.DelSpaceWorkObjectReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkObjectId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkObjectId"))
	}

	//进入逻辑
	_, err := s.uc.DelAndTransferWorkItem(ctx, loginUser, req.SpaceId, req.WorkObjectId, req.ToWorkObjectId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.DelSpaceWorkObjectReply{Result: &pb.DelSpaceWorkObjectReply_Data{Data: &pb.DelSpaceWorkObjectReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkObjectService) SetOrder(ctx context.Context, req *pb.SetWorkObjectOrderRequest) (*pb.SetWorkObjectOrderReply, error) {

	reply := func(err error) (*pb.SetWorkObjectOrderReply, error) {
		return &pb.SetWorkObjectOrderReply{Result: &pb.SetWorkObjectOrderReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	if req.SpaceId < 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.WorkObjectId < 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.FromIdx < 0 {
		return reply(errs.Param(ctx, "FromIdx"))
	}

	if req.ToIdx < 0 {
		return reply(errs.Param(ctx, "ToIdx"))
	}

	err := s.uc.SetOrder(ctx, loginUser, req.SpaceId, req.WorkObjectId, req.FromIdx, req.ToIdx)
	if err != nil {
		return reply(err)
	}

	return &pb.SetWorkObjectOrderReply{}, nil
}

func (s *SpaceWorkObjectService) SetSpaceWorkObjectRanking(ctx context.Context, req *pb.SetSpaceWorkObjectRankingRequest) (*pb.SetSpaceWorkObjectRankingReply, error) {
	reply := func(err error) (*pb.SetSpaceWorkObjectRankingReply, error) {
		return &pb.SetSpaceWorkObjectRankingReply{Result: &pb.SetSpaceWorkObjectRankingReply_Error{Error: errs.Cast(err)}}, nil
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
	err := s.uc.SetSpaceWorkWorkObjectRanking(ctx, loginUser, req.SpaceId, rankList)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetSpaceWorkObjectRankingReply{Result: &pb.SetSpaceWorkObjectRankingReply_Data{Data: ""}}
	return okReply, nil
}
