package data

import (
	"context"
	"fmt"
	"go-cs/internal/dwh/model/aps"
	"go-cs/internal/utils/date"
	"go-cs/pkg/pprint"
	"go-cs/pkg/stream"
	"testing"
	"time"
)

func TestSpace(t *testing.T) {
	startDate := date.Parse("2025/03/01 00:00:00")
	endDate := date.Parse("2025/03/01 23:59:59")

	h, err := repo.DashboardRptSpaceWitem1h(context.Background(), &aps.RptSpaceWitemQuery{
		SpaceIds:  []int64{87},
		UserId:    42,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	h = ExpandDateRangeList("Hour", startDate, endDate, h)

	pprint.Print(h)
}

func TestSpaceDay(t *testing.T) {
	startDate := date.Parse("2025/02/01 00:00:00")
	endDate := date.Parse("2025/03/31 23:59:59")

	list, err := repo.DashboardRptSpaceWitem1day(context.Background(), &aps.RptSpaceWitemQuery{
		SpaceIds:  []int64{87},
		UserId:    42,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	list = ExpandDateRangeList("Day", startDate, endDate, list)

	pprint.Print(list)
}

func TestSpaceMonth(t *testing.T) {
	startDate := date.Parse("2025/02/28 00:00:00")
	endDate := date.Parse("2025/04/01 00:00:00")

	list, err := repo.DashboardRptSpaceWitem1month(context.Background(), &aps.RptSpaceWitemQuery{
		SpaceIds:  []int64{87},
		UserId:    42,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	list = ExpandDateRangeList("Month", startDate, endDate, list)

	pprint.Print(list)
}

func TestMemberDay(t *testing.T) {
	startDate := date.Parse2("2025-04-01 00:00:00")
	endDate := date.Parse2("2025-04-01 16:00:00")

	list, err := repo.SearchRptMemberWitem1day(context.Background(), &aps.RptMemberWitemQuery{
		SpaceIds:  []int64{8510},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	groupBy := stream.GroupBy(list, func(item *aps.RptMemberWitem) string {
		return fmt.Sprintf("%v:%v", item.SpaceId, item.UserId)
	})

	finalResult := make([]*aps.RptMemberWitem, 0, len(list))
	for _, v := range groupBy {
		items := ExpandDateRangeList("Day", startDate, endDate, v)
		finalResult = append(finalResult, items...)
	}

	pprint.Print(list)
}

func GetEndDate(endDate time.Time) time.Time {
	now := time.Now()

	if endDate.After(now) {
		endDate = now
	}

	nowHourBegin := date.HourBegin(now)

	dur := endDate.Sub(nowHourBegin)

	if dur >= time.Minute*15 {
		return nowHourBegin
	}

	if dur > 0 {
		return nowHourBegin.Add(-time.Hour)
	}

	if dur > -time.Minute*60 {
		return date.HourBegin(endDate)
	}

	return endDate
}

func TestSpaceIncrDay(t *testing.T) {
	startDate := date.Parse2("2024-03-31 00:00:00")
	endDate := GetEndDate(date.Parse2("2025-03-31 23:59:59"))
	list, err := repo.DashboardSpaceCompleteIncrDay(context.Background(), &aps.RptSpaceIncrWitemQuery{
		SpaceIds:  []int64{8270},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	pprint.Println(list)
}

func TestSpaceIncrHour(t *testing.T) {
	startDate := date.Parse("2024/03/24 00:00:00")
	endDate := date.Parse("2025/03/24 23:59:59")
	list, err := repo.DashboardSpaceCompleteIncrHour(context.Background(), &aps.RptSpaceIncrWitemQuery{
		SpaceIds:  []int64{87},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(list)
}

func TestRptRepo_DashboardSpaceNewIncrMonth(t *testing.T) {
	startDate := date.Parse("2024/03/24 00:00:00")
	endDate := date.Parse("2025/03/24 23:59:59")
	list, err := repo.DashboardSpaceNewIncrMonth(context.Background(), &aps.RptSpaceIncrWitemQuery{
		SpaceIds:  []int64{87},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	pprint.Println(list)
}

func TestRptRepo_DashboardSpaceNewIncrDay(t *testing.T) {
	startDate := date.Parse2("2025-03-26 00:00:00")
	endDate := GetEndDate(date.Parse2("2025-04-01 13:57:59"))
	list, err := repo.DashboardSpaceNewIncrDay(context.Background(), &aps.RptSpaceIncrWitemQuery{
		SpaceIds:  []int64{8480},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	pprint.Println(list)
}

func TestRptRepo_DashboardSpaceNewIncrHour(t *testing.T) {
	startDate := date.Parse("2024/03/24 00:00:00")
	endDate := date.Parse("2025/03/24 23:59:59")
	list, err := repo.DashboardSpaceNewIncrHour(context.Background(), &aps.RptSpaceIncrWitemQuery{
		SpaceIds:  []int64{87},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	pprint.Println(list)
}

func TestRptRepo_DashboardMemberCompleteIncrDay(t *testing.T) {
	startDate := date.Parse("2025/04/08 00:00:00")
	endDate := date.Parse("2025/04/08 14:20:00")
	list, err := repo.DashboardMemberCompleteIncrDay(context.Background(), &aps.RptMemberIncrWitemQuery{
		SpaceIds:  []int64{4790},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	pprint.Println(list)
}

func TestRptRepo_DashboardMemberCompleteIncrMonth(t *testing.T) {
	startDate := date.Parse("2025/04/07 00:00:00")
	endDate := date.Parse("2025/04/07 14:20:00")
	list, err := repo.DashboardMemberCompleteIncrMonth(context.Background(), &aps.RptMemberIncrWitemQuery{
		SpaceIds:  []int64{8270},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	pprint.Println(list)
}

func TestRptRepo_DashboardMemberNewIncr(t *testing.T) {
	startDate := date.Parse2("2025-04-02 00:00:00")
	endDate := date.Parse2("2025-04-02 15:00:00")
	list, err := repo.DashboardMemberNewIncr(context.Background(), &aps.RptMemberIncrWitemQuery{
		SpaceIds:  []int64{8691},
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		t.Fatal(err)
	}

	pprint.Println(list)
}

func TestGetEndDate(t *testing.T) {
	t.Log(GetEndDate(date.Parse2("2025-04-01 23:59:59")))
}
