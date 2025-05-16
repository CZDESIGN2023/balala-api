package biz

import (
	"context"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	perm_facade "go-cs/internal/domain/perm/facade"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/stream"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"

	pb "go-cs/api/space_work_version/v1"

	space_repo "go-cs/internal/domain/space/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	workVersion_domain "go-cs/internal/domain/space_work_version"
	workVersion_repo "go-cs/internal/domain/space_work_version/repo"
	workVersion_service "go-cs/internal/domain/space_work_version/service"
	user_repo "go-cs/internal/domain/user/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"
)

type SpaceWorkVersionUsecase struct {
	log *log.Helper
	tm  trans.Transaction

	repo              workVersion_repo.SpaceWorkVersionRepo
	userRepo          user_repo.UserRepo
	spaceRepo         space_repo.SpaceRepo
	spaceMemberRepo   member_repo.SpaceMemberRepo
	spaceWorkItemRepo witem_repo.WorkItemRepo

	permService        *perm_service.PermService
	workVersionService *workVersion_service.SpaceWorkVersionService

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewSpaceWorkVersionUsecase(

	logger log.Logger,
	tm trans.Transaction,

	repo workVersion_repo.SpaceWorkVersionRepo,
	spaceWorkItemRepo witem_repo.WorkItemRepo,
	userRepo user_repo.UserRepo,
	spaceRepo space_repo.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,

	permService *perm_service.PermService,
	workVersionService *workVersion_service.SpaceWorkVersionService,

	domainMessageProducer *domain_message.DomainMessageProducer,
) *SpaceWorkVersionUsecase {
	moduleName := "SpaceWorkVersionUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceWorkVersionUsecase{
		repo:                  repo,
		spaceMemberRepo:       spaceMemberRepo,
		userRepo:              userRepo,
		spaceRepo:             spaceRepo,
		spaceWorkItemRepo:     spaceWorkItemRepo,
		log:                   hlog,
		tm:                    tm,
		permService:           permService,
		workVersionService:    workVersionService,
		domainMessageProducer: domainMessageProducer,
	}
}

func (s *SpaceWorkVersionUsecase) CreateMySpaceWorkVersion(ctx context.Context, oper *utils.LoginUserInfo, in *pb.CreateSpaceWorkVersionRequest) (*db.SpaceWorkVersion, error) {

	ownerId := oper.UserId
	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, in.SpaceId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	//判断成员信息是否存在
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, ownerId)
	if err != nil {
		//创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_MEMBER_WRONG)
		return nil, err
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CREATE_SPACE_WORK_VERSION,
	})

	if err != nil {
		return nil, err
	}

	// 检查空间名称是否存在
	workVersion, err := s.workVersionService.CreateSpaceWorkVersion(ctx, in.SpaceId, in.VersionName, 0, oper)
	if err != nil {
		return nil, err
	}

	// 填充数据实体
	err2 := s.repo.CreateSpaceWorkVersion(ctx, workVersion)
	if err2 != nil {
		//创建空间失败
		err := errs.Business(ctx, "创建失败")
		return nil, err
	}

	s.domainMessageProducer.Send(ctx, workVersion.GetMessages())

	return convert.SpaceWorkVersionEntityToPo(workVersion), nil
}

func (s *SpaceWorkVersionUsecase) QSpaceWorkVersionList(ctx context.Context, userId int64, spaceId int64) (*pb.SpaceWorkVersionListReplyData, error) {
	// 判断空间信息是否存在
	_, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, userId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	list, err := s.repo.QSpaceWorkVersionList(ctx, spaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	result := &pb.SpaceWorkVersionListReplyData{}
	result.List = make([]*pb.SpaceWorkVersionListReplyData_Item, 0)
	for _, v := range list {
		result.List = append(result.List, &pb.SpaceWorkVersionListReplyData_Item{
			Id:          v.Id,
			SpaceId:     v.SpaceId,
			VersionKey:  v.VersionKey,
			VersionName: v.VersionName,
			Ranking:     v.Ranking,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return result, nil
}

func (s *SpaceWorkVersionUsecase) QSpaceWorkVersionById(ctx context.Context, oper *utils.LoginUserInfo, req *pb.SpaceWorkVersionByIdRequest) (*pb.SpaceWorkVersionByIdReplyData, error) {
	list, err := s.repo.QSpaceWorkVersionById(ctx, req.SpaceId, req.Ids)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	spaceIds := stream.Map(list, func(item *workVersion_domain.SpaceWorkVersion) int64 {
		return item.SpaceId
	})

	isMember, err := s.spaceMemberRepo.UserInAllSpace(ctx, oper.UserId, spaceIds...)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if !isMember && !s.userRepo.IsEnterpriseAdmin(ctx, oper.UserId) {
		return nil, errs.NoPerm(ctx)
	}

	result := &pb.SpaceWorkVersionByIdReplyData{}
	result.List = make([]*pb.SpaceWorkVersionByIdReplyData_Item, 0)
	for _, v := range list {
		result.List = append(result.List, &pb.SpaceWorkVersionByIdReplyData_Item{
			Id:          v.Id,
			SpaceId:     v.SpaceId,
			VersionKey:  v.VersionKey,
			VersionName: v.VersionName,
			Ranking:     v.Ranking,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return result, nil
}

func (s *SpaceWorkVersionUsecase) SetMySpaceWorkVersionName(ctx context.Context, oper *utils.LoginUserInfo, in *pb.ModifySpaceWorkVersionNameRequest) (*db.SpaceWorkVersion, error) {

	//判断这个空间里是不是有这个工作项
	workVersion, err := s.repo.GetSpaceWorkVersion(ctx, in.VersionId)
	if err != nil {
		return nil, errs.Business(ctx, "获取工作项版本信息失败")
	}

	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, workVersion.SpaceId)
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
		Perm:              consts.PERM_MODIFY_SPACE_WORK_VERSION,
	})

	if err != nil {
		return nil, err
	}

	//判断名字是否重复
	err = s.workVersionService.UpdateSpaceWorkVersionName(ctx, workVersion, in.VersionName, oper)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveSpaceWorkVersion(ctx, workVersion)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, workVersion.GetMessages())

	return convert.SpaceWorkVersionEntityToPo(workVersion), nil
}

func (s *SpaceWorkVersionUsecase) DelMySpaceWorkVersion(ctx context.Context, oper *utils.LoginUserInfo, workVersionId int64, toWorkVersionId int64) (*db.SpaceWorkVersion, error) {

	// 判断这个空间里是不是有这个工作项
	workVersion, err := s.repo.GetSpaceWorkVersion(ctx, workVersionId)
	if err != nil {
		return nil, errs.Business(ctx, "获取工作项版本信息失败")
	}

	count, err := s.repo.GetSpaceWorkVersionCount(ctx, workVersion.SpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if count <= 1 {
		return nil, errs.Business(ctx, "最后一个版本不可删除")
	}

	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, workVersion.SpaceId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	// 检查创建空间的用户是否存在
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_MEMBER_WRONG)
	}

	relationCount, err := s.repo.GetVersionRelationCount(ctx, workVersion.SpaceId, workVersion.Id)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if relationCount > 0 && toWorkVersionId == 0 {
		return nil, errs.Business(ctx, "存在关联任务，不可删除")
	}

	var toWorkVersion *workVersion_domain.SpaceWorkVersion
	if toWorkVersionId > 0 {
		toWorkVersion, err = s.repo.GetSpaceWorkVersion(ctx, toWorkVersionId)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		//不是同一个空间下的工作项版本不能互相转移
		if workVersion.SpaceId != toWorkVersion.SpaceId {
			return nil, errs.Business(ctx, "不是同一个空间下的工作项版本不能互相转移")
		}
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_DELETE_SPACE_WORK_VERSION,
	})

	if err != nil {
		return nil, err
	}

	workVersion.OnDelete(oper)

	err = s.tm.InTx(ctx, func(ctx context.Context) error {
		_, err := s.repo.DelWorkVersion(ctx, workVersionId)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		if toWorkVersion != nil {
			err = s.spaceWorkItemRepo.ResetVersion(ctx, workVersionId, toWorkVersion.Id)
			if err != nil {
				return errs.Internal(ctx, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.domainMessageProducer.Send(ctx, workVersion.GetMessages())

	return convert.SpaceWorkVersionEntityToPo(workVersion), nil
}

func (s *SpaceWorkVersionUsecase) GetMySpaceWorkVersionRelationCount(ctx context.Context, userId int64, spaceId int64, workVersionId int64) (int64, error) {
	// 判断是否当前空间成员
	_, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, userId)
	if err != nil {
		//不是的话无权查看
		return 0, errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	count, err := s.repo.GetVersionRelationCount(ctx, spaceId, workVersionId)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	return count, nil
}

func (uc *SpaceWorkVersionUsecase) SetSpaceWorkWorkVersionRanking(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, rankingList []map[string]int64) error {

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
		Perm:              consts.PERM_MODIFY_SPACE_WORK_VERSION,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	var workVersions workVersion_domain.SpaceWorkVersions
	for _, v := range rankingList {
		versionId := cast.ToInt64(v["id"])
		newRanking := cast.ToInt64(v["ranking"])

		//调整排序
		workVersion, err := uc.repo.GetSpaceWorkVersion(ctx, versionId)
		if err != nil {
			return err
		}

		if !workVersion.IsSameSpace(spaceId) {
			return errs.NoPerm(ctx)
		}

		workVersion.UpdateRanking(newRanking, oper)

		workVersions = append(workVersions, workVersion)
	}

	txErr := uc.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range workVersions {
			err = uc.repo.SaveSpaceWorkVersion(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return errs.Internal(ctx, txErr)
	}

	msg := &domain_message.ChangeVersionOrder{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     oper,
			OperTime: time.Now(),
		},
		SpaceId: spaceId,
	}

	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}
