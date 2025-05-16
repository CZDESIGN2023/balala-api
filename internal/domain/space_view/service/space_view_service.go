package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pb "go-cs/api/space_view/v1"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	member_repo "go-cs/internal/domain/space_member/repo"
	domain "go-cs/internal/domain/space_view"
	"go-cs/internal/domain/space_view/repo"
	"go-cs/internal/domain/work_item_status"
	shared "go-cs/internal/pkg/domain"
	"go-cs/pkg/stream"
)

type SpaceViewService struct {
	repo       repo.SpaceViewRepo
	memberRepo member_repo.SpaceMemberRepo
}

func NewSpaceViewService(
	repo repo.SpaceViewRepo,
	memberRepo member_repo.SpaceMemberRepo,
) *SpaceViewService {
	return &SpaceViewService{
		repo:       repo,
		memberRepo: memberRepo,
	}
}

func (s *SpaceViewService) CreateSpaceUserView(ctx context.Context, spaceId int64, name, key string, ranking int64, typ int64, outerId int64, queryConfig, tableConfig string, oper shared.Oper) (*domain.SpaceUserView, error) {
	if key == "" {
		key = uuid.NewString()
	}

	ins := domain.NewSpaceUserView(
		spaceId,
		key,
		name,
		ranking,
		typ,
		outerId,
		queryConfig,
		tableConfig,
		oper,
	)

	ins.AddMessage(oper, &domain_message.CreateSpaceView{
		SpaceId:  spaceId,
		ViewType: typ,
		ViewName: name,
	})
	return ins, nil
}

func (s *SpaceViewService) CreateSpaceGlobalView(ctx context.Context, spaceId int64, name string, ranking int64, typ int64, queryConfig, tableConfig string, oper shared.Oper) (*domain.SpaceGlobalView, error) {

	key := uuid.NewString()
	switch pb.SpaceViewType(typ) {
	case pb.SpaceViewType_SpaceViewType_Public:
		key = "pub_" + key
	default:
		return nil, errors.New("invalid type")
	}

	ins := domain.NewSpaceGlobalView(
		spaceId,
		key,
		name,
		ranking,
		typ,
		queryConfig,
		tableConfig,
		oper,
	)

	ins.AddMessage(oper, &domain_message.CreateSpaceView{
		SpaceId:  spaceId,
		ViewType: typ,
		ViewName: name,
	})
	return ins, nil
}

func (s *SpaceViewService) InitSpacePublicView(ctx context.Context, spaceId int64, status work_item_status.WorkItemStatusItems) ([]*domain.SpaceGlobalView, error) {
	statusIds := stream.Map(status.GetProcessingStatus(), func(item *work_item_status.WorkItemStatusItem) int64 {
		return item.Id
	})

	list := []*domain.SpaceGlobalView{
		s.newProcessingView(ctx, spaceId, statusIds),
		s.newExpiredView(ctx, spaceId, statusIds),
	}

	return list, nil
}

func (s *SpaceViewService) InitUserGlobalView(ctx context.Context, spaceId int64, userIds []int64) ([]*domain.SpaceUserView, error) {
	allSystemViewKeys := consts.GetAllSystemViewKeys()

	spaceGlobalViews, _ := s.repo.GetGlobalViewList(ctx, spaceId)
	spaceGlobalViewMap := stream.ToMap(spaceGlobalViews, func(k int, v *domain.SpaceGlobalView) (string, *domain.SpaceGlobalView) {
		return v.Key, v
	})

	spaceGlobalViewKeys := stream.Map(spaceGlobalViews, func(v *domain.SpaceGlobalView) string {
		return v.Key
	})

	allGlobalViewKeys := stream.Clone(allSystemViewKeys).Concat(spaceGlobalViewKeys...).Unique().List()

	var allUserViews []*domain.SpaceUserView
	if len(allGlobalViewKeys) != 0 {
		for _, userId := range userIds {
			userViews := stream.Map(allGlobalViewKeys, func(key string) *domain.SpaceUserView {
				var outerId int64
				var typ = int64(pb.SpaceViewType_SpaceViewType_System)

				if globalView := spaceGlobalViewMap[key]; globalView != nil {
					outerId = globalView.Id
					typ = globalView.Type
				}

				userView, _ := s.CreateSpaceUserView(ctx,
					spaceId,
					consts.GetSystemViewName(key),
					key,
					consts.GetViewRank(key),
					typ,
					outerId,
					"",
					"",
					shared.UserOper(userId),
				)

				return userView
			})

			allUserViews = append(allUserViews, userViews...)
		}
	}

	return allUserViews, nil
}

func (s *SpaceViewService) EnterpriseUserView(ctx context.Context, spaceId int64) ([]*domain.SpaceUserView, error) {
	allSystemViewKeys := consts.GetAllSystemViewKeys()
	spaceGlobalViews, _ := s.repo.GetGlobalViewList(ctx, spaceId)

	systemUserViews := stream.Map(allSystemViewKeys, func(key string) *domain.SpaceUserView {
		var outerId int64
		var typ = int64(pb.SpaceViewType_SpaceViewType_System)
		userView, _ := s.CreateSpaceUserView(ctx,
			spaceId,
			consts.GetSystemViewName(key),
			key,
			consts.GetViewRank(key),
			typ,
			outerId,
			"",
			"",
			shared.UserOper(0),
		)

		return userView
	})

	publicUserViews := stream.Map(spaceGlobalViews, func(v *domain.SpaceGlobalView) *domain.SpaceUserView {
		userView := v.CreateUserView(0)
		userView.Name = v.Name
		userView.TableConfig = v.TableConfig
		userView.QueryConfig = v.QueryConfig
		return userView
	})

	return append(systemUserViews, publicUserViews...), nil
}
