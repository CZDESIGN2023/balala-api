package biz

import (
	"context"
	"go-cs/api/comm"
	"go-cs/internal/consts"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	"go-cs/internal/domain/space"
	"go-cs/internal/domain/space_member"
	"go-cs/internal/domain/space_view"
	"go-cs/internal/domain/space_view/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/stream"
	"google.golang.org/protobuf/encoding/protojson"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	pb "go-cs/api/space_view/v1"

	space_repo "go-cs/internal/domain/space/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	view_service "go-cs/internal/domain/space_view/service"
	user_repo "go-cs/internal/domain/user/repo"
)

type SpaceViewUsecase struct {
	log *log.Helper
	tm  trans.Transaction

	repo            repo.SpaceViewRepo
	userRepo        user_repo.UserRepo
	spaceRepo       space_repo.SpaceRepo
	spaceMemberRepo member_repo.SpaceMemberRepo

	permService *perm_service.PermService
	service     *view_service.SpaceViewService

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewSpaceViewUsecase(

	logger log.Logger,
	tm trans.Transaction,

	repo repo.SpaceViewRepo,
	userRepo user_repo.UserRepo,
	spaceRepo space_repo.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,

	permService *perm_service.PermService,
	viewService *view_service.SpaceViewService,

	domainMessageProducer *domain_message.DomainMessageProducer,
) *SpaceViewUsecase {
	moduleName := "SpaceViewUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceViewUsecase{
		repo:                  repo,
		spaceMemberRepo:       spaceMemberRepo,
		userRepo:              userRepo,
		spaceRepo:             spaceRepo,
		log:                   hlog,
		tm:                    tm,
		permService:           permService,
		service:               viewService,
		domainMessageProducer: domainMessageProducer,
	}
}

func (s *SpaceViewUsecase) Create(ctx context.Context, uid int64, req *pb.CreateViewRequest) error {

	queryConfig := pb.QueryConfig{}
	err := protojson.Unmarshal([]byte(req.QueryConfig), &queryConfig)
	if err != nil {
		s.log.Error(err)
		return errs.Param(ctx, "QueryConfig")
	}

	oper := shared.UserOper(uid)

	space, err := s.spaceRepo.GetSpace(ctx, req.SpaceId)
	if err != nil {
		return errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, uid)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	switch req.Type {
	default:
		return errs.Param(ctx, "Type")
	case pb.SpaceViewType_SpaceViewType_Public:
		// 检查权限
		if member.RoleId < consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER {
			return errs.NoPerm(ctx)
		}

		globalView, err := s.service.CreateSpaceGlobalView(ctx,
			space.Id,
			req.Name,
			0,
			int64(req.Type),
			req.QueryConfig,
			req.TableConfig,
			oper,
		)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		members, err := s.spaceMemberRepo.GetSpaceMemberBySpaceId(ctx, space.Id)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		userViews := stream.Map(members, func(v *space_member.SpaceMember) *space_view.SpaceUserView {
			return globalView.CreateUserView(v.UserId)
		})

		err = s.tm.InTx(ctx, func(ctx context.Context) error {
			err = s.repo.CreateGlobalView(ctx, globalView)
			if err != nil {
				return errs.Internal(ctx, err)
			}

			for _, v := range userViews {
				v.OuterId = globalView.Id
			}

			err = s.repo.CreateUserViews(ctx, userViews)
			if err != nil {
				return errs.Internal(ctx, err)
			}

			return nil
		})
		if err != nil {
			return errs.Internal(ctx, err)
		}

		msg := &domain_message.CreateSpaceView{
			DomainMessageBase: shared.DomainMessageBase{
				Oper:     oper,
				OperTime: time.Now(),
			},
			SpaceId:  globalView.SpaceId,
			ViewId:   globalView.Id,
			ViewType: globalView.Type,
			ViewName: globalView.Name,
		}

		s.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	case pb.SpaceViewType_SpaceViewType_Personal:
		view, err := s.service.CreateSpaceUserView(
			ctx,
			space.Id,
			req.Name,
			"",
			0,
			int64(req.Type),
			0,
			req.QueryConfig,
			req.TableConfig,
			oper,
		)
		if err != nil {
			return err
		}

		err = s.tm.InTx(ctx, func(ctx context.Context) error {
			err = s.repo.CreateUserView(ctx, view)
			if err != nil {
				return errs.Internal(ctx, err)
			}
			return err
		})
		if err != nil {
			return errs.Internal(ctx, err)
		}

		msg := &domain_message.CreateSpaceView{
			DomainMessageBase: shared.DomainMessageBase{
				Oper:     oper,
				OperTime: time.Now(),
			},
			SpaceId:  view.SpaceId,
			ViewId:   view.Id,
			ViewType: view.Type,
			ViewName: view.Name,
		}

		s.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})
	}

	return nil
}

func (s *SpaceViewUsecase) ViewList(ctx context.Context, userId int64, spaceIds []int64, key string) (*pb.ViewListReply, error) {
	list, err := s.repo.UserViewList(ctx, userId, spaceIds, key)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var items []*pb.ViewListReply_Item
	for _, v := range list {
		item := &pb.ViewListReply_Item{
			Id:          v.Id,
			SpaceId:     v.SpaceId,
			Name:        v.Name,
			Ranking:     v.Ranking,
			Type:        v.Type,
			Key:         v.Key,
			OuterId:     v.OuterId,
			QueryConfig: v.QueryConfig,
			TableConfig: v.TableConfig,
			Status:      v.Status,
		}

		items = append(items, item)
	}

	// 处理企业管理员
	userInfo, err := s.userRepo.GetUserByUserId(ctx, userId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if userInfo.IsEnterpriseAdmin() {
		spaceList, err := s.spaceRepo.GetUserSpaceList(ctx, userId)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		userSpaceIds := stream.Map(spaceList, func(v *space.Space) int64 {
			return v.Id
		})

		notMemberSpaceIds := stream.Diff(spaceIds, userSpaceIds)

		for _, spaceId := range notMemberSpaceIds {
			view, err := s.service.EnterpriseUserView(ctx, spaceId)
			if err != nil {
				return nil, errs.Internal(ctx, err)
			}

			if key != "" {
				view = stream.Filter(view, func(v *space_view.SpaceUserView) bool {
					return v.Name == key
				})
			}

			items = append(items, stream.Map(view, func(v *space_view.SpaceUserView) *pb.ViewListReply_Item {
				return &pb.ViewListReply_Item{
					Id:          v.Id,
					SpaceId:     v.SpaceId,
					Name:        v.Name,
					Ranking:     v.Ranking,
					Type:        v.Type,
					Key:         v.Key,
					QueryConfig: v.QueryConfig,
					TableConfig: v.TableConfig,
					Status:      v.Status,
				}
			})...)
		}

	}

	return &pb.ViewListReply{Data: items}, nil
}

func (s *SpaceViewUsecase) Del(ctx context.Context, uid int64, id int64) error {
	view, err := s.repo.GetUserViewById(ctx, id)
	if err != nil {
		return errs.Business(ctx, "获取工作项版本信息失败")
	}

	oper := shared.UserOper(uid)

	// 检查创建空间的用户是否存在
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, view.SpaceId, uid)
	if err != nil || member == nil {
		return errs.NoPerm(ctx)
	}

	switch pb.SpaceViewType(view.Type) {
	case pb.SpaceViewType_SpaceViewType_System:
		return errs.Business(ctx, "系统视图不允许删除")
	case pb.SpaceViewType_SpaceViewType_Public:
		if member.RoleId < consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER {
			return errs.NoPerm(ctx)
		}

		err = s.tm.InTx(ctx, func(ctx context.Context) error {
			err = s.repo.DeleteGlobalViewById(ctx, view.OuterId)
			if err != nil {
				return errs.Internal(ctx, err)
			}

			err = s.repo.DeleteUserViewByOuterId(ctx, view.OuterId)
			if err != nil {
				return errs.Internal(ctx, err)
			}

			return nil
		})
		if err != nil {
			return errs.Internal(ctx, err)
		}
	case pb.SpaceViewType_SpaceViewType_Personal:
		if view.UserId != uid {
			return errs.NoPerm(ctx)
		}

		err = s.repo.DeleteUserViewById(ctx, id)
		if err != nil {
			return errs.Internal(ctx, err)
		}
	}

	view.OnDelete(oper)
	s.domainMessageProducer.Send(ctx, view.GetMessages())

	if err != nil {
		return err
	}

	return nil
}

func (s *SpaceViewUsecase) SetRanking(ctx context.Context, uid int64, spaceId int64, rankingList []*pb.SetViewRankingRequest_Item) error {

	// 判断当前用户是否在要查询的项目空间内
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return err
	}

	viewMap, err := s.repo.GetUserViewMap(ctx, uid, spaceId)
	if err != nil {
		return err
	}

	var updateList []*space_view.SpaceUserView
	for _, v := range rankingList {
		id := v.Id
		ranking := v.Ranking

		view := viewMap[id]
		if view == nil {
			continue
		}

		if view.SpaceId != spaceId {
			return errs.NoPerm(ctx)
		}

		if view.Ranking == ranking {
			continue
		}

		view.UpdateRanking(ranking, shared.UserOper(uid))
		updateList = append(updateList, view)
	}

	txErr := s.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range updateList {
			err = s.repo.SaveSpaceUserView(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return errs.Internal(ctx, txErr)
	}

	msg := &domain_message.ChangeSpaceViewOrder{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     shared.UserOper(uid),
			OperTime: time.Now(),
		},
		SpaceId: spaceId,
	}

	s.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}

func (s *SpaceViewUsecase) SetName(ctx context.Context, uid int64, req *pb.SetViewNameRequest) error {

	oper := shared.UserOper(uid)

	//判断这个空间里是不是有这个工作项
	view, err := s.repo.GetUserViewById(ctx, req.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, view.SpaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 判断用户是不是在这个空间里
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, uid)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	switch pb.SpaceViewType(view.Type) {
	case pb.SpaceViewType_SpaceViewType_System:
		return errs.Business(ctx, "系统视图不允许修改名称")

	case pb.SpaceViewType_SpaceViewType_Public:
		if member.RoleId < int64(consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER) {
			return errs.NoPerm(ctx)
		}

		globalView := view.GetGlobalView()

		globalView.SetName(req.Name, oper)
		err = s.repo.SaveSpaceGlobalView(ctx, globalView)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		s.domainMessageProducer.Send(ctx, globalView.GetMessages())
	case pb.SpaceViewType_SpaceViewType_Personal:
		if view.UserId != uid {
			return errs.NoPerm(ctx)
		}

		view.SetName(req.Name, oper)
		err := s.repo.SaveSpaceUserView(ctx, view)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		s.domainMessageProducer.Send(ctx, view.GetMessages())
	}

	return nil
}

func (s *SpaceViewUsecase) SetStatus(ctx context.Context, uid int64, req *pb.SetViewStatusRequest) error {
	oper := shared.UserOper(uid)

	//判断这个空间里是不是有这个工作项
	view, err := s.repo.GetUserViewById(ctx, req.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	//判断空间是否存在
	space, err := s.spaceRepo.GetSpace(ctx, view.SpaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 判断用户是不是在这个空间里
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, uid)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	view.SetStatus(req.Status, oper)

	err = s.tm.InTx(ctx, func(ctx context.Context) error {
		err = s.repo.SaveSpaceUserView(ctx, view)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		return nil
	})

	s.domainMessageProducer.Send(ctx, view.GetMessages())

	return nil
}

func (s *SpaceViewUsecase) SetQueryConfig(ctx context.Context, uid int64, req *pb.SetViewQueryConfigRequest) error {
	oper := shared.UserOper(uid)

	//判断这个空间里是不是有这个工作项
	view, err := s.repo.GetUserViewById(ctx, req.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	// 判断用户是不是在这个空间里
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, view.SpaceId, uid)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	switch pb.SpaceViewType(view.Type) {
	case pb.SpaceViewType_SpaceViewType_Public:
		if member.RoleId < int64(consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER) {
			return errs.NoPerm(ctx)
		}

		globalView := view.GetGlobalView()

		globalView.SetQueryConfig(req.Config, oper)
		err = s.repo.SaveSpaceGlobalView(ctx, globalView)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		s.domainMessageProducer.Send(ctx, globalView.GetMessages())
	case pb.SpaceViewType_SpaceViewType_Personal:
		if view.UserId != uid {
			return errs.NoPerm(ctx)
		}

		view.SetQueryConfig(req.Config)
		err := s.repo.SaveSpaceUserView(ctx, view)
		if err != nil {
			return errs.Internal(ctx, err)
		}
	}

	return nil
}

func (s *SpaceViewUsecase) SetTableConfig(ctx context.Context, uid int64, req *pb.SetViewTableConfigRequest) error {
	oper := shared.UserOper(uid)

	view, err := s.repo.GetUserViewById(ctx, req.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	spaceMember, err := s.spaceMemberRepo.GetSpaceMember(ctx, view.SpaceId, uid)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	if spaceMember.Id == 0 { //静默处理
		return nil
	}

	if view.UserId != uid {
		return errs.NoPerm(ctx)
	}

	if globalView := view.GetGlobalView(); globalView != nil &&
		spaceMember.RoleId >= int64(consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER) {

		globalView.SetTableConfig(req.Config, oper)
		err = s.repo.SaveSpaceGlobalView(ctx, globalView)
		if err != nil {
			return errs.Internal(ctx, err)
		}
	}

	view.SetTableConfig(req.Config)
	err = s.repo.SaveSpaceUserView(ctx, view)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	return nil
}
