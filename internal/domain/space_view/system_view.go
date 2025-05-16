package space_view

import (
	pb "go-cs/api/search/v1"
	"go-cs/internal/utils"
)

type QueryConfig struct {
	ConditionGroup *pb.ConditionGroup `json:"conditionGroup"`
	Sorts          []*pb.Sort         `json:"sorts"`
	Groups         []*pb.GroupBy      `json:"groups"`
}

var (
	queryAll        string
	queryProcessing string
	queryExpired    string
	queryFollow     string
)

func init() {
	queryAll = utils.ToJSON(&QueryConfig{})
	queryProcessing = utils.ToJSON(&QueryConfig{
		ConditionGroup: &pb.ConditionGroup{
			ConditionGroup: []*pb.ConditionGroup{
				{
					Conditions: []*pb.Condition{
						{
							Field:    "work_item_status_id",
							Operator: "IN",
							Values:   []string{"processing"},
						},
					},
				},
			},
		},
	})
}

func SystemViews() []*SpaceUserView {
	all := &SpaceUserView{}

	processing := &SpaceUserView{
		QueryConfig: "",
		TableConfig: "",
	}

	expired := &SpaceUserView{
		QueryConfig: "",
		TableConfig: "",
	}

	follow := &SpaceUserView{
		QueryConfig: "",
		TableConfig: "",
	}

	return []*SpaceUserView{
		all,
		processing,
		expired,
		follow,
	}
}
