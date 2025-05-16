package biz

import (
	"context"
	v1 "go-cs/api/log/v1"
	db "go-cs/internal/bean/biz"
	vo "go-cs/internal/bean/vo"
	space_repo "go-cs/internal/domain/space/repo"
	user_repo "go-cs/internal/domain/user/repo"
	work_item_repo "go-cs/internal/domain/work_item/repo"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/oper"
	"go-cs/pkg/stream"
	"math"

	"github.com/go-kratos/kratos/v2/log"
)

type LogRepo interface {
	UserLoginLogPagination(ctx context.Context, userId int64, pos, size int) ([]*db.UserLoginLog, error)
	OpLogPagination(ctx context.Context, searchVo *vo.OpLogPaginationSearchVo) ([]*db.OperLog, error)
}

type LogUsecase struct {
	tm           trans.Transaction
	log          *log.Helper
	repo         LogRepo
	userRepo     user_repo.UserRepo
	spaceRepo    space_repo.SpaceRepo
	workItemRepo work_item_repo.WorkItemRepo
}

func NewLogUsecase(repo LogRepo, tm trans.Transaction, spaceRepo space_repo.SpaceRepo, userRepo user_repo.UserRepo, workItemRepo work_item_repo.WorkItemRepo,
	logger log.Logger) *LogUsecase {

	moduleName := "LogUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &LogUsecase{
		tm:           tm,
		log:          hlog,
		repo:         repo,
		userRepo:     userRepo,
		spaceRepo:    spaceRepo,
		workItemRepo: workItemRepo,
	}
}

func (s *LogUsecase) UserLoginLogList(ctx context.Context, uid int64, req *v1.LoginLogListRequest) (*v1.LoginLogListReplyData, error) {

	if req.Pos == 0 {
		req.Pos = math.MaxInt64
	}

	if req.Size == 0 {
		req.Size = 20
	}

	var pos = int(req.Pos)
	var size = int(req.Size)

	list, err := s.repo.UserLoginLogPagination(ctx, uid, pos, size+1)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var hasNext bool
	var nextPos int64
	if len(list) > size {
		hasNext = true
		nextPos = list[size].Id
		list = list[:size]
	}

	var items []*v1.LoginLogListReplyData_Item
	for _, v := range list {
		items = append(items, &v1.LoginLogListReplyData_Item{
			Info: v.ToProto(),
		})
	}

	return &v1.LoginLogListReplyData{
		Items:   items,
		HasNext: hasNext,
		NextPos: nextPos,
	}, nil
}

func (s *LogUsecase) SpaceOpLogList(ctx context.Context, uid int64, req *v1.SpaceOpLogListRequest) (*v1.SpaceOpLogListReplyData, error) {

	//获取当前用户可以查看的操作日志空间范围
	spaceIds, _ := s.spaceRepo.GetUserSpaceIds(ctx, uid)
	if req.SpaceId > 0 {
		if !stream.Contains(spaceIds, req.SpaceId) {
			return nil, errs.NoPerm(ctx)
		}
		spaceIds = []int64{req.SpaceId}
	}

	if req.Size == 0 {
		req.Size = 20
	}

	var pos = int(req.Pos)
	var size = int(req.Size)

	var IncludeModuleType []int
	if req.Scene == "space_sys_log" {
		IncludeModuleType = []int{
			int(oper.ModuleTypeSpace),
			int(oper.ModuleTypeSpaceMember),
			int(oper.ModuleTypeSpaceTag),
			int(oper.ModuleTypeSpaceWorkObject),
			int(oper.ModuleTypeSpaceWorkVersion),
			int(oper.ModuleTypeSpaceMemberCategory),
			int(oper.ModuleTypeWorkFlow),
			int(oper.ModuleTypeWorkItemRole),
			int(oper.ModuleTypeWorkItemStatus),
			int(oper.ModuleTypeSpaceGlobalView),
		}
	}

	list, err := s.repo.OpLogPagination(ctx, &vo.OpLogPaginationSearchVo{
		Size:              size + 1,
		Pos:               int64(pos),
		SpaceIds:          spaceIds,
		ModuleType:        int(req.ModuleType),
		ModuleId:          req.ModuleId,
		OperId:            req.OperId,
		IncludeModuleType: IncludeModuleType,
	})
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var hasNext bool
	var nextPos int64
	if len(list) > size {
		hasNext = true
		nextPos = list[size].Id
		list = list[:size]
	}

	var items []*v1.SpaceOpLogListReplyData_Item
	for _, v := range list {
		items = append(items, &v1.SpaceOpLogListReplyData_Item{
			Info: v.ToProto(),
		})
	}

	return &v1.SpaceOpLogListReplyData{
		Items:   items,
		HasNext: hasNext,
		NextPos: nextPos,
	}, nil
}

func (s *LogUsecase) OpLogList(ctx context.Context, uid int64, req *v1.OpLogListRequest) (*v1.OpLogListReplyData, error) {

	if req.Size == 0 {
		req.Size = 20
	}

	var pos = int(req.Pos)
	var size = int(req.Size)

	list, err := s.repo.OpLogPagination(ctx, &vo.OpLogPaginationSearchVo{
		Size:       size + 1,
		Pos:        int64(pos),
		ModuleId:   req.ModuleId,
		ModuleType: int(req.ModuleType),
		OperId:     req.OperId,
	})
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var hasNext bool
	var nextPos int64
	if len(list) > size {
		hasNext = true
		nextPos = list[size].Id
		list = list[:size]
	}

	spaceIds := stream.Map(list, func(e *db.OperLog) int64 {
		return e.SpaceId
	})
	spaceMap, err := s.spaceRepo.SpaceMap(ctx, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	for _, v := range list {
		if v.SpaceId != 0 && spaceMap[v.SpaceId] != nil {
			v.SpaceName = spaceMap[v.SpaceId].SpaceName
		}
	}

	list = stream.Filter(list, func(e *db.OperLog) bool {
		if e.SpaceId != 0 && spaceMap[e.SpaceId] == nil {
			return false
		}

		return true
	})

	var items []*v1.OpLogListReplyData_Item
	for _, v := range list {
		items = append(items, &v1.OpLogListReplyData_Item{
			Info: v.ToProto(),
		})
	}

	return &v1.OpLogListReplyData{
		Items:   items,
		HasNext: hasNext,
		NextPos: nextPos,
	}, nil
}

func (s *LogUsecase) SystemOpLogList(ctx context.Context, uid int64, req *v1.SystemOpLogListRequest) (*v1.SystemOpLogListReplyData, error) {
	user, err := s.userRepo.GetUserByUserId(ctx, uid)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if !user.IsSystemAdmin() {
		return nil, errs.NoPerm(ctx)
	}

	if req.Size == 0 {
		req.Size = 20
	}

	var pos = int(req.Pos)
	var size = int(req.Size)

	list, err := s.repo.OpLogPagination(ctx, &vo.OpLogPaginationSearchVo{
		Size:         size + 1,
		Pos:          int64(pos),
		OperatorType: int(oper.OperatorTypeSys),
	})
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var hasNext bool
	var nextPos int64
	if len(list) > size {
		hasNext = true
		nextPos = list[size].Id
		list = list[:size]
	}

	var items []*v1.SystemOpLogListReplyData_Item
	for _, v := range list {
		items = append(items, &v1.SystemOpLogListReplyData_Item{
			Info: v.ToProto(),
		})
	}

	return &v1.SystemOpLogListReplyData{
		Items:   items,
		HasNext: hasNext,
		NextPos: nextPos,
	}, nil
}

func (s *LogUsecase) PersonalOpLogList(ctx context.Context, uid int64, req *v1.PersonalOpLogListRequest) (*v1.PersonalOpLogListReplyData, error) {

	if req.Size == 0 {
		req.Size = 20
	}

	var pos = int(req.Pos)
	var size = int(req.Size)

	list, err := s.repo.OpLogPagination(ctx, &vo.OpLogPaginationSearchVo{
		Size:   size + 1,
		Pos:    int64(pos),
		OperId: uid,
	})
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var hasNext bool
	var nextPos int64
	if len(list) > size {
		hasNext = true
		nextPos = list[size].Id
		list = list[:size]
	}

	// 过滤子任务创建日志
	list = stream.Filter(list, func(e *db.OperLog) bool {
		return !(e.ModuleType == int32(oper.ModuleTypeSpaceWorkItem) &&
			e.BusinessType == int32(oper.BusinessTypeAdd) &&
			e.ModuleFlag&int64(oper.ModuleFlag_subWorkItem) == 1)
	})

	spaceIds := stream.Map(stream.Filter(list, func(e *db.OperLog) bool {
		return e.SpaceId > 0 && e.SpaceName == ""
	}), func(e *db.OperLog) int64 {
		return e.SpaceId
	})

	spaceMap, err := s.spaceRepo.SpaceMap(ctx, stream.Unique(spaceIds))
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	for _, v := range list {
		if v.SpaceId != 0 && spaceMap[v.SpaceId] != nil {
			v.SpaceName = spaceMap[v.SpaceId].SpaceName
		}
	}

	list = stream.Filter(list, func(e *db.OperLog) bool {
		if e.SpaceId != 0 && e.SpaceName == "" {
			return false
		}

		return true
	})

	var items []*v1.PersonalOpLogListReplyData_Item
	for _, v := range list {
		items = append(items, &v1.PersonalOpLogListReplyData_Item{
			Info: v.ToProto(),
		})
	}

	return &v1.PersonalOpLogListReplyData{
		Items:   items,
		HasNext: hasNext,
		NextPos: nextPos,
	}, nil
}
