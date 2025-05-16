package service

import (
	"context"
	v1 "go-cs/api/search/v1"
	pb "go-cs/api/space_view/v1"
	"go-cs/internal/consts"
	"go-cs/internal/domain/search/search_es"
	domain "go-cs/internal/domain/space_view"
	"go-cs/internal/utils"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *SpaceViewService) newProcessingView(ctx context.Context, spaceId int64, processingStatusIds []int64) *domain.SpaceGlobalView {
	config := pb.QueryConfig{
		ConditionGroup: &v1.ConditionGroup{
			Conjunction: string(search_es.AND),
			Conditions: []*v1.Condition{
				{
					Field:    string(search_es.WorkItemStatusIdField),
					Operator: string(search_es.IN),
					Values:   utils.ToStrArray(processingStatusIds),
				},
			},
		},
	}

	marshal, _ := protojson.Marshal(&config)

	view, _ := s.CreateSpaceGlobalView(ctx, spaceId,
		consts.GetPublicViewName(consts.PublicViewKey_Processing),
		consts.GetViewRank(consts.PublicViewKey_Processing),
		int64(pb.SpaceViewType_SpaceViewType_Public),
		string(marshal),
		"",
		nil)

	view.Key = consts.PublicViewKey_Processing

	return view
}

func (s *SpaceViewService) newExpiredView(ctx context.Context, spaceId int64, processingStatusIds []int64) *domain.SpaceGlobalView {
	config := pb.QueryConfig{
		ConditionGroup: &v1.ConditionGroup{
			Conjunction: string(search_es.AND),
			Conditions: []*v1.Condition{
				{
					Field:    string(search_es.WorkItemStatusIdField),
					Operator: string(search_es.IN),
					Values:   utils.ToStrArray(processingStatusIds),
				},
				{
					Field:    string(search_es.PlanTimeField),
					Operator: string(search_es.LT),
					Values:   []string{"${TODAY}", "${TODAY}"},
				},
			},
		},
	}

	marshal, _ := protojson.Marshal(&config)

	view, _ := s.CreateSpaceGlobalView(ctx, spaceId,
		consts.GetPublicViewName(consts.PublicViewKey_Expired),
		consts.GetViewRank(consts.PublicViewKey_Expired),
		int64(pb.SpaceViewType_SpaceViewType_Public),
		string(marshal),
		"",
		nil)

	view.Key = consts.PublicViewKey_Expired

	return view
}
