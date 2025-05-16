package biz

import (
	"context"
	"go-cs/internal/bean/vo"
	"go-cs/internal/test/dbt"
	"go-cs/internal/utils"
	"testing"
	"time"
)

func MockUserInfo(uid int64) *utils.LoginUserInfo {
	return &utils.LoginUserInfo{
		UserId: uid,
	}
}

func TestGetWorkItemDetailV2(t *testing.T) {
	user := MockUserInfo(42)

	res, err := dbt.UC.SpaceWorkItemUsecase.GetWorkItemDetail(context.Background(), MockUserInfo(user.UserId), 1031, 2363)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func BenchmarkGetWorkItemDetailV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = dbt.UC.SpaceWorkItemUsecase.GetWorkItemDetail(context.Background(), MockUserInfo(42), 87, 358)
	}
}

func TestCreateWorkItem(t *testing.T) {
	_, err := dbt.UC.SpaceWorkItemUsecase.CreateTask(context.Background(), MockUserInfo(42), vo.CreateSpaceWorkItemVoV2{
		UserId:         42,
		SpaceId:        87,
		WorkObjectId:   5466,
		WorkVersionId:  247,
		WorkItemName:   "test",
		PlanStartAt:    time.Now().Unix(),
		PlanCompleteAt: time.Now().Add(time.Hour * 24).Unix(),
		Priority:       "P3",
		WorkFlowId:     13866,
		WorkItemTypeId: 5932,
		Owner: vo.CreateSpaceWorkItemOwnersV2{
			{
				OwnerRole: "12236",
				Directors: []int64{42},
			},
			{
				OwnerRole: "12237",
				Directors: []int64{42},
			},
			{
				OwnerRole: "12238",
				Directors: []int64{42},
			},
		},
	})

	if err != nil {
		t.Error(err)
	}
}

func Test_ConfirmStateFlowMainTaskState(t *testing.T) {
	const workItemId = 16954
	const nextStatusKey = "st_convert_to_story"

	err := dbt.UC.SpaceWorkItemUsecase.ConfirmStateFlowMain(context.Background(), MockUserInfo(42), workItemId, nextStatusKey, "test", "")
	if err != nil {
		t.Error(err)
	}
}
