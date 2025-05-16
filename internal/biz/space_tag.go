package biz

import (
	"cmp"
	"context"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	perm_facade "go-cs/internal/domain/perm/facade"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	space_repo "go-cs/internal/domain/space/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	tag_repo "go-cs/internal/domain/space_tag/repo"
	tag_service "go-cs/internal/domain/space_tag/service"
	statics_repo "go-cs/internal/domain/statics/repo"
	user_repo "go-cs/internal/domain/user/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"slices"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceTagUsecase struct {
	tm              trans.Transaction
	repo            tag_repo.SpaceTagRepo
	userRepo        user_repo.UserRepo
	spaceRepo       space_repo.SpaceRepo
	spaceMemberRepo member_repo.SpaceMemberRepo
	workItemRepo    witem_repo.WorkItemRepo
	log             *log.Helper
	staticsRepo     statics_repo.StaticsRepo

	permService     *perm_service.PermService
	spaceTagService *tag_service.SpaceTagService

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewSpaceTagUsecase(
	repo tag_repo.SpaceTagRepo,
	tm trans.Transaction, staticsRepo statics_repo.StaticsRepo,

	workItemRepo witem_repo.WorkItemRepo,
	userRepo user_repo.UserRepo,
	spaceRepo space_repo.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,

	permService *perm_service.PermService,
	spaceTagService *tag_service.SpaceTagService,

	domainMessageProducer *domain_message.DomainMessageProducer,

	logger log.Logger) *SpaceTagUsecase {

	moduleName := "SpaceTagUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceTagUsecase{
		repo:                  repo,
		userRepo:              userRepo,
		spaceRepo:             spaceRepo,
		spaceMemberRepo:       spaceMemberRepo,
		workItemRepo:          workItemRepo,
		log:                   hlog,
		staticsRepo:           staticsRepo,
		tm:                    tm,
		permService:           permService,
		spaceTagService:       spaceTagService,
		domainMessageProducer: domainMessageProducer,
	}
}

func (s *SpaceTagUsecase) CreateMySpaceTag(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, tagName string) (*db.SpaceTag, error) {

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil {
		//成员不存在 不允许操作
		errInfo := errs.New(ctx, comm.ErrorCode_PERMISSION_INSUFFICIENT_DATA_PERMISSIONS)
		return nil, errInfo
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CREATE_SPACE_TAG,
	})
	if err != nil {
		return nil, err
	}

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		errInfo := errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
		return nil, errInfo
	}

	spaceTag, err := s.spaceTagService.CreateSpaceTag(ctx, space.Id, tagName, oper)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreateTag(ctx, spaceTag)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, spaceTag.GetMessages())

	return convert.SpaceTagEntityToPo(spaceTag), nil
}

func (s *SpaceTagUsecase) ModifyMySpaceTagName(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, tagId int64, tagName string) (*db.SpaceTag, error) {

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil {
		return nil, errs.NoPerm(ctx)
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_TAG,
	})
	if err != nil {
		return nil, err
	}

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	spaceTag, err := s.repo.GetSpaceTag(ctx, space.Id, tagId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_TAG_INFO_WRONG)
	}

	err = s.spaceTagService.UpdateTagName(ctx, spaceTag, tagName, oper)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveTag(ctx, spaceTag)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_TAG_CREATE_FAIL)
	}

	//添加日志
	s.domainMessageProducer.Send(ctx, spaceTag.GetMessages())

	return convert.SpaceTagEntityToPo(spaceTag), nil
}

func (s *SpaceTagUsecase) DelMySpaceTagV2(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, tagId int64) (*db.SpaceTag, error) {

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil {
		//成员不存在 不允许操作
		return nil, errs.NoPerm(ctx)
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_DELETE_SPACE_TAG,
	})
	if err != nil {
		return nil, err
	}

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		errInfo := errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
		return nil, errInfo
	}

	spaceTag, err := s.repo.GetSpaceTag(ctx, space.Id, tagId)
	if err != nil {
		errInfo := errs.New(ctx, comm.ErrorCode_SPACE_TAG_INFO_WRONG)
		return nil, errInfo
	}

	spaceTag.OnDelete(oper)

	//TAG删除事务
	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.repo.DelSpaceTag(ctx, tagId)
		if err != nil {
			return err
		}

		// 移除工作项中的标签
		_, err = s.workItemRepo.RemoveTagFromAllWorkItem(ctx, spaceId, tagId)
		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		errInfo := errs.New(ctx, comm.ErrorCode_SPACE_TAG_DEL_FAIL)
		return nil, errInfo
	}

	s.domainMessageProducer.Send(ctx, spaceTag.GetMessages())

	return convert.SpaceTagEntityToPo(spaceTag), nil
}

func (s *SpaceTagUsecase) GetMySpaceTagListV2(ctx context.Context, userId int64, spaceId int64) ([]*rsp.TagInfo, error) {

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	_, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, userId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	//countMap, err := s.staticsRepo.GetSpaceWorkItemTagCount(ctx, spaceId)
	//if err != nil {
	//	return nil, errs.Internal(ctx, err)
	//}

	list, _ := s.repo.QSpaceTagList(ctx, spaceId)
	var items []*rsp.TagInfo
	for _, v := range list {
		items = append(items, &rsp.TagInfo{
			Id:      v.Id,
			TagName: v.TagName,
		})
	}

	slices.SortFunc(items, func(i, j *rsp.TagInfo) int {
		return cmp.Compare(j.Id, i.Id)
	})

	return items, nil
}
