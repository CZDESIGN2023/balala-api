package service

import (
	"context"
	"go-cs/api/comm"
	pb "go-cs/api/search/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"

	"github.com/go-kratos/kratos/v2/log"
)

type SearchService struct {
	pb.UnimplementedSearchServer

	uc  *biz.SearchUsecase
	log *log.Helper
}

func NewSearchService(SearchUsecase *biz.SearchUsecase, logger log.Logger) *SearchService {
	moduleName := "SearchService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SearchService{
		uc:  SearchUsecase,
		log: hlog,
	}
}

func clearCondGroup(req *pb.ConditionGroup) *pb.ConditionGroup {
	if req == nil {
		return nil
	}

	//req.Conditions = clearCondition(req.Conditions)

	var condGroups []*pb.ConditionGroup
	for i := 0; i < len(req.ConditionGroup); i++ {
		group := clearCondGroup(req.ConditionGroup[i])
		if group != nil {
			condGroups = append(condGroups, group)
		}
	}
	req.ConditionGroup = condGroups

	if len(req.Conditions) == 0 && req.ConditionGroup == nil {
		return nil
	}

	return req
}

func clearCondition(arr []*pb.Condition) []*pb.Condition {
	var list []*pb.Condition
	for _, v := range arr {
		if v.Field == "" || v.Operator == "" || len(v.Values) == 0 {
			continue
		}

		list = append(list, v)
	}
	return list
}

func (s *SearchService) SearchMySpaceWorkItemGroupInfoV2(ctx context.Context, req *pb.SearchSpaceWorkItemGroupInfoRequestV2) (*pb.SearchSpaceWorkItemGroupInfoReplyV2, error) {
	uid := utils.GetLoginUser(ctx).UserId

	req.ConditionGroup = clearCondGroup(req.ConditionGroup)

	data, err := s.uc.SearchGroupInfo(ctx, uid, req)
	if err != nil {
		return &pb.SearchSpaceWorkItemGroupInfoReplyV2{Result: &pb.SearchSpaceWorkItemGroupInfoReplyV2_Error{Error: errs.Cast(err)}}, nil
	}

	return &pb.SearchSpaceWorkItemGroupInfoReplyV2{Result: &pb.SearchSpaceWorkItemGroupInfoReplyV2_Data{Data: data}}, nil
}

func (s *SearchService) SearchMySpaceWorkItemsByIdV2(ctx context.Context, request *pb.SearchMySpaceWorkItemsByIdRequest) (*pb.SearchMySpaceWorkItemsByIdReplyV2, error) {

	var reply = func(err *comm.ErrorInfo) (*pb.SearchMySpaceWorkItemsByIdReplyV2, error) {
		return &pb.SearchMySpaceWorkItemsByIdReplyV2{Result: &pb.SearchMySpaceWorkItemsByIdReplyV2_Error{Error: err}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		return reply(errs.NotLogin(ctx))
	}

	out, err := s.uc.SearchMySpaceWorkItemsByIdV2(ctx, loginUser.UserId, request.Ids)
	if err != nil {
		return reply(errs.Cast(err))
	}

	okReply := &pb.SearchMySpaceWorkItemsByIdReplyV2{Result: &pb.SearchMySpaceWorkItemsByIdReplyV2_Data{Data: out}}
	return okReply, nil
}
