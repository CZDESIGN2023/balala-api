package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "go-cs/api/space_view/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
)

type SpaceViewService struct {
	pb.UnimplementedSpaceViewServer
	uc  *biz.SpaceViewUsecase
	log *log.Helper
}

func NewSpaceViewService(uc *biz.SpaceViewUsecase, logger log.Logger) *SpaceViewService {
	moduleName := "SpaceViewService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceViewService{
		uc:  uc,
		log: hlog,
	}
}

func (s *SpaceViewService) CreateView(ctx context.Context, req *pb.CreateViewRequest) (*pb.CreateViewReply, error) {
	reply := func(err error) (*pb.CreateViewReply, error) {
		return &pb.CreateViewReply{Error: errs.Cast(err)}, err
	}

	uid := utils.GetLoginUser(ctx).UserId

	if req.Type != pb.SpaceViewType_SpaceViewType_Public && req.Type != pb.SpaceViewType_SpaceViewType_Personal {
		return reply(errs.Param(ctx, "Type"))
	}

	err := s.uc.Create(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.CreateViewReply{}, nil
}

func (s *SpaceViewService) DelView(ctx context.Context, req *pb.DelViewRequest) (*pb.DelViewReply, error) {
	reply := func(err error) (*pb.DelViewReply, error) {
		return &pb.DelViewReply{Error: errs.Cast(err)}, err
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.Del(ctx, uid, req.Id)
	if err != nil {
		return reply(err)
	}

	return &pb.DelViewReply{}, nil
}

func (s *SpaceViewService) ViewList(ctx context.Context, req *pb.ViewListRequest) (*pb.ViewListReply, error) {
	reply := func(err error) (*pb.ViewListReply, error) {
		return &pb.ViewListReply{Error: errs.Cast(err)}, err
	}

	uid := utils.GetLoginUser(ctx).UserId

	res, err := s.uc.ViewList(ctx, uid, req.SpaceIds, req.Key)
	if err != nil {
		return reply(err)
	}

	return res, nil
}

func (s *SpaceViewService) SetViewName(ctx context.Context, req *pb.SetViewNameRequest) (*pb.SetViewNameReply, error) {
	reply := func(err error) (*pb.SetViewNameReply, error) {
		return &pb.SetViewNameReply{Error: errs.Cast(err)}, err
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.SetName(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SetViewNameReply{}, nil
}

func (s *SpaceViewService) SetViewStatus(ctx context.Context, req *pb.SetViewStatusRequest) (*pb.SetViewStatusReply, error) {
	reply := func(err error) (*pb.SetViewStatusReply, error) {
		return &pb.SetViewStatusReply{Error: errs.Cast(err)}, err
	}

	uid := utils.GetLoginUser(ctx).UserId

	if req.Status < 0 || req.Status > 1 {
		return reply(errs.Param(ctx, "Status"))
	}

	err := s.uc.SetStatus(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SetViewStatusReply{}, nil
}

func (s *SpaceViewService) SetViewRanking(ctx context.Context, req *pb.SetViewRankingRequest) (*pb.SetViewRankingReply, error) {
	reply := func(err error) (*pb.SetViewRankingReply, error) {
		return &pb.SetViewRankingReply{Error: errs.Cast(err)}, err
	}

	uid := utils.GetLoginUser(ctx).UserId

	if len(req.List) == 0 {
		return reply(errs.Param(ctx, "ViewList"))
	}

	err := s.uc.SetRanking(ctx, uid, req.SpaceId, req.List)
	if err != nil {
		return reply(err)
	}

	return &pb.SetViewRankingReply{}, nil
}

func (s *SpaceViewService) SetViewQueryConfig(ctx context.Context, req *pb.SetViewQueryConfigRequest) (*pb.SetViewQueryConfigReply, error) {
	reply := func(err error) (*pb.SetViewQueryConfigReply, error) {
		return &pb.SetViewQueryConfigReply{Error: errs.Cast(err)}, err
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.SetQueryConfig(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SetViewQueryConfigReply{}, nil
}

func (s *SpaceViewService) SetViewTableConfig(ctx context.Context, req *pb.SetViewTableConfigRequest) (*pb.SetViewTableConfigReply, error) {
	reply := func(err error) (*pb.SetViewTableConfigReply, error) {
		return &pb.SetViewTableConfigReply{Error: errs.Cast(err)}, err
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.SetTableConfig(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SetViewTableConfigReply{}, nil
}
