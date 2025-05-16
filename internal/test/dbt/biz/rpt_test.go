package biz

import (
	"context"
	v1 "go-cs/api/rpt/v1"
	"go-cs/internal/test/dbt"
	"testing"
)

func TestDashboardRptSpaceWitem(t *testing.T) {
	result, err := dbt.UC.RptUsecase.DashboardRptSpaceWitem(context.Background(), MockUserInfo(21), &v1.DashboardRptSpaceWitemRequest{
		TimeSplitType: "Day",
		TimeRange: &v1.DashboardRptSpaceWitemRequest_TimeRange{
			StartDate: "2025-03-18 00:00:00",
			EndDate:   "2025-03-24 23:59:59",
		},
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}
