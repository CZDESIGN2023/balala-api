package service

import (
	"context"
	pb "go-cs/api/config/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"

	"github.com/go-kratos/kratos/v2/log"
)

type ConfigService struct {
	pb.UnimplementedConfigServer

	uc  *biz.ConfigUsecase
	log *log.Helper
}

func NewConfigService(configUsecase *biz.ConfigUsecase, logger log.Logger) *ConfigService {
	moduleName := "ConfigService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &ConfigService{
		uc:  configUsecase,
		log: hlog,
	}
}

func (s *ConfigService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListReply, error) {
	res, err := s.uc.List(ctx)
	if err != nil {
		return nil, err
	}

	okReply := &pb.ListReply{Result: &pb.ListReply_Data{Data: &pb.ListReplyData{Config: res}}}
	return okReply, nil
}

func (s *ConfigService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateReply, error) {
	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.UpdateByKey(ctx, uid, req.Key, req.Value)
	if err != nil {
		return &pb.UpdateReply{Result: &pb.UpdateReply_Error{Error: errs.Cast(err)}}, nil
	}

	return &pb.UpdateReply{Result: nil}, nil
}
