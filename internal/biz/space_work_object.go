package biz

import (
	"context"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	shared "go-cs/internal/pkg/domain"
	"go-cs/pkg/stream"
	"time"

	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"

	pb "go-cs/api/space_work_object/v1"
	perm_facade "go-cs/internal/domain/perm/facade"
	space_repo "go-cs/internal/domain/space/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	workObj_domain "go-cs/internal/domain/space_work_object"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	workObj_service "go-cs/internal/domain/space_work_object/service"
	user_repo "go-cs/internal/domain/user/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"
)

type SpaceWorkObjectUsecase struct {
	repo              workObj_repo.SpaceWorkObjectRepo
	userRepo          user_repo.UserRepo
	spaceRepo         space_repo.SpaceRepo
	spaceMemberRepo   member_repo.SpaceMemberRepo
	spaceWorkItemRepo witem_repo.WorkItemRepo
	log               *log.Helper
	tm                trans.Transaction

	permService         *perm_service.PermService
	spaceWorkObjService *workObj_service.SpaceWorkObjectService

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewSpaceWorkObjectUsecase(
	repo workObj_repo.SpaceWorkObjectRepo,
	spaceWorkItemRepo witem_repo.WorkItemRepo,
	tm trans.Transaction,
	userRepo user_repo.UserRepo,
	spaceRepo space_repo.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,
	permService *perm_service.PermService,
	spaceWorkObjService *workObj_service.SpaceWorkObjectService,
	domainMessageProducer *domain_message.DomainMessageProducer,
	logger log.Logger,
) *SpaceWorkObjectUsecase {
	moduleName := "SpaceWorkObjectUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceWorkObjectUsecase{
		repo:                  repo,
		spaceMemberRepo:       spaceMemberRepo,
		userRepo:              userRepo,
		spaceRepo:             spaceRepo,
		spaceWorkItemRepo:     spaceWorkItemRepo,
		log:                   hlog,
		tm:                    tm,
		permService:           permService,
		spaceWorkObjService:   spaceWorkObjService,
		domainMessageProducer: domainMessageProducer,
	}
}

func (s *SpaceWorkObjectUsecase) CreateMySpaceWorkObject(ctx context.Context, oper *utils.LoginUserInfo, in *pb.CreateSpaceWorkObjectRequest) (*db.SpaceWorkObject, error) {
	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, in.SpaceId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	//判断成员信息是否存在
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if err != nil {
		//创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_MEMBER_WRONG)
		return nil, err
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CREATE_SPACE_WORK_OBJECT,
	})
	if err != nil {
		return nil, err
	}

	workObject, err := s.spaceWorkObjService.CreateSpaceWorkObject(ctx, in.SpaceId, in.WorkObjectName, 0, oper.UserId, oper)
	if err != nil {
		return nil, err
	}

	// 填充数据实体
	err = s.repo.CreateSpaceWorkObject(ctx, workObject)
	if err != nil {
		//创建空间失败
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_WORK_OBJECT_CREATE_FAIL)
	}

	s.domainMessageProducer.Send(ctx, workObject.GetMessages())

	return convert.SpaceWorkObjectEntityToPo(workObject), nil
}

func (s *SpaceWorkObjectUsecase) QSpaceWorkObjectList(ctx context.Context, oper *utils.LoginUserInfo, req *pb.SpaceWorkObjectListRequest) (*pb.SpaceWorkObjectListReplyData, error) {
	// 判断空间信息是否存在
	_, err := s.spaceMemberRepo.GetSpaceMember(ctx, req.SpaceId, oper.UserId)
	if err != nil { //创建人的信息不存在
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	list, err := s.repo.QSpaceWorkObjectList(ctx, req.SpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	result := &pb.SpaceWorkObjectListReplyData{}
	result.List = make([]*pb.SpaceWorkObjectListReplyData_Item, 0)

	for _, item := range list {
		result.List = append(result.List, &pb.SpaceWorkObjectListReplyData_Item{
			Id:               item.Id,
			SpaceId:          item.SpaceId,
			UserId:           item.UserId,
			WorkObjectGuid:   item.WorkObjectGuid,
			WorkObjectName:   item.WorkObjectName,
			WorkObjectStatus: int32(item.WorkObjectStatus),
			Ranking:          item.Ranking,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
		})
	}

	return result, nil
}

func (s *SpaceWorkObjectUsecase) QSpaceWorkObjectById(ctx context.Context, oper *utils.LoginUserInfo, req *pb.SpaceWorkObjectByIdRequest) (*pb.SpaceWorkObjectByIdReplyData, error) {
	list, err := s.repo.QSpaceWorkObjectById(ctx, req.SpaceId, req.Ids)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	spaceIds := stream.Map(list, func(item *workObj_domain.SpaceWorkObject) int64 {
		return item.SpaceId
	})

	isMember, err := s.spaceMemberRepo.UserInAllSpace(ctx, oper.UserId, spaceIds...)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if !isMember && !s.userRepo.IsEnterpriseAdmin(ctx, oper.UserId) {
		return nil, errs.NoPerm(ctx)
	}

	result := &pb.SpaceWorkObjectByIdReplyData{}
	result.List = make([]*pb.SpaceWorkObjectByIdReplyData_Item, 0)

	for _, item := range list {
		result.List = append(result.List, &pb.SpaceWorkObjectByIdReplyData_Item{
			Id:               item.Id,
			SpaceId:          item.SpaceId,
			UserId:           item.UserId,
			WorkObjectGuid:   item.WorkObjectGuid,
			WorkObjectName:   item.WorkObjectName,
			WorkObjectStatus: int32(item.WorkObjectStatus),
			Ranking:          item.Ranking,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
		})
	}

	return result, nil
}

func (s *SpaceWorkObjectUsecase) SetMySpaceWorkObjectName(ctx context.Context, oper *utils.LoginUserInfo, in *pb.ModifySpaceWorkObjectNameRequest) (*db.SpaceWorkObject, error) {

	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, in.SpaceId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	// 判断用户是不是在这个空间里
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return nil, errs.NoPerm(ctx)
	}

	// 验证权限

	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_OBJECT,
	})

	if err != nil {
		return nil, err
	}

	//判断这个空间里是不是有这个工作项
	workObject, err := s.repo.GetSpaceWorkObject(ctx, in.SpaceId, in.WorkObjectId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_WORK_OBJECT_INFO_WRONG)
	}

	err = s.spaceWorkObjService.UpdateSpaceWorkObjectName(ctx, workObject, in.WorkObjectName, oper)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveSpaceWorkObject(ctx, workObject)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	//添加日志
	s.domainMessageProducer.Send(ctx, workObject.GetMessages())

	return convert.SpaceWorkObjectEntityToPo(workObject), nil
}

func (s *SpaceWorkObjectUsecase) DelMySpaceWorkObject(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workObjectId int64) (*db.SpaceWorkObject, error) {

	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		//创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
		return nil, err
	}

	//是不是有权限操作
	// 检查创建空间的用户是否存在
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if err != nil {
		//创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_MEMBER_WRONG)
		return nil, err
	}

	// 验证权限

	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_DELETE_SPACE_WORK_OBJECT,
	})

	if err != nil {
		return nil, err
	}

	// 判断这个空间里是不是有这个工作项
	workObject, err := s.repo.GetSpaceWorkObject(ctx, spaceId, workObjectId)
	if err != nil {
		//创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_WORK_OBJECT_INFO_WRONG)
		return nil, err
	}

	workObject.OnDelete(oper)

	err = s.tm.InTx(ctx, func(ctx context.Context) error {
		_, err := s.repo.DelWorkObject(ctx, workObjectId)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		workItemIds, err := s.spaceWorkItemRepo.GetSpaceWorkItemIdsByWorkObject(ctx, workObjectId)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		if len(workItemIds) != 0 {
			return errs.Business(ctx, "模块存在任务，不能删除")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.domainMessageProducer.Send(ctx, workObject.GetMessages())

	return convert.SpaceWorkObjectEntityToPo(workObject), nil
}

func (s *SpaceWorkObjectUsecase) DelAndTransferWorkItem(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workObjectId int64, toWorkObjectId int64) (*db.SpaceWorkObject, error) {

	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		err := errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
		return nil, err
	}

	//是不是有权限操作
	// 检查创建空间的用户是否存在
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if err != nil {
		//创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_MEMBER_WRONG)
		return nil, err
	}

	// 验证权限

	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_DELETE_SPACE_WORK_OBJECT2,
	})

	if err != nil {
		return nil, err
	}

	// 判断这个空间里是不是有这个工作项
	workObject, err := s.repo.GetSpaceWorkObject(ctx, spaceId, workObjectId)
	if err != nil {
		//创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_WORK_OBJECT_INFO_WRONG)
		return nil, err
	}

	err = s.spaceWorkObjService.CheckTransfer(ctx, workObject, toWorkObjectId, oper)
	if err != nil {
		return nil, err
	}

	workObject.OnDelete(oper)

	err = s.tm.InTx(ctx, func(ctx context.Context) error {
		// 删除旧模块
		update, err := s.repo.DelWorkObject(ctx, workObjectId)
		if err != nil {
			return err
		}
		if update == 0 {
			return errs.Business(ctx, "模块存在任务，不能删除")
		}

		// 获取旧模块关联的任务
		workItemIds, err := s.spaceWorkItemRepo.GetSpaceWorkItemIdsByWorkObject(ctx, workObjectId)
		if err != nil {
			return err
		}

		if len(workItemIds) > 0 && toWorkObjectId == 0 {
			return errs.Business(ctx, "模块存在任务，不能删除")
		}

		// 更新任务的模块id
		err = s.spaceWorkItemRepo.UpdateWorkItemWorkObjectIdByIds(ctx, workItemIds, toWorkObjectId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.domainMessageProducer.Send(ctx, workObject.GetMessages())

	return convert.SpaceWorkObjectEntityToPo(workObject), nil
}

func (s *SpaceWorkObjectUsecase) SetOrder(ctx context.Context, oper *utils.LoginUserInfo, spaceId, workObjectId, fromIdx, toIdx int64) error {
	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		err := errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
		return err
	}

	//是不是有权限操作
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_OBJECT,
	})
	if err != nil {
		return err
	}

	// 判断这个空间里是不是有这个工作项
	workObject, err := s.repo.GetSpaceWorkObject(ctx, spaceId, workObjectId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	workObject.UpdateRanking(time.Now().Unix(), oper)

	err = s.repo.SetOrder(ctx, spaceId, fromIdx, toIdx)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, workObject.GetMessages())

	return nil
}

func (uc *SpaceWorkObjectUsecase) SetSpaceWorkWorkObjectRanking(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, rankingList []map[string]int64) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.spaceMemberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_OBJECT,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	var workObjects workObj_domain.SpaceWorkObjects
	for _, v := range rankingList {
		objectId := cast.ToInt64(v["id"])
		newRanking := cast.ToInt64(v["ranking"])

		//调整排序
		workObject, err := uc.repo.GetSpaceWorkObject(ctx, spaceId, objectId)
		if err != nil {
			return err
		}

		workObject.UpdateRanking(newRanking, oper)

		workObjects = append(workObjects, workObject)
	}

	txErr := uc.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range workObjects {
			err = uc.repo.SaveSpaceWorkObject(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return errs.Internal(ctx, txErr)
	}

	msg := &domain_message.ChangeWorkObjectOrder{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     oper,
			OperTime: time.Now(),
		},
		SpaceId: spaceId,
	}

	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}
