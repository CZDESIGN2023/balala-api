package service

import (
	"context"
	"go-cs/internal/dwh/data"
	"go-cs/internal/dwh/model/aps"
	"time"
)

type DwhAds struct {
	repo *data.RptRepo
}

func NewDwhAds(
	repo *data.RptRepo,
) *DwhAds {
	return &DwhAds{
		repo: repo,
	}
}

type SearchRptVersionWitemReq struct {
	// Hour, Day, Week, Month
	TimeSplitType string
	StartDate     time.Time
	EndDate       time.Time
	SpaceId       int64
	VersionId     int64
	StatusId      int64
}

func (d *DwhAds) SearchRptVersionWitem(ctx context.Context, req *SearchRptVersionWitemReq) ([]*aps.RptVersionWitem, error) {

	searchReq := &aps.RptVersionWitemQuery{
		SpaceIds:  []int64{req.SpaceId},
		VersionId: req.VersionId,
		StatusId:  req.StatusId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	var err error
	var result []*aps.RptVersionWitem

	switch req.TimeSplitType {
	case "Hour":
		result, err = d.repo.SearchRptVersionWitem1h(ctx, searchReq)
	case "Day":
		result, err = d.repo.SearchRptVersionWitem1day(ctx, searchReq)
	case "Month":
		result, err = d.repo.SearchRptVersionWitem1month(ctx, searchReq)
	}

	return result, err
}

type SearchRptMemberWitemReq struct {
	// Hour, Day, Week, Month
	TimeSplitType string
	StartDate     time.Time
	EndDate       time.Time
	SpaceIds      []int64
	UserId        int64
}

func (d *DwhAds) SearchRptMemberWitem(ctx context.Context, req *SearchRptMemberWitemReq) ([]*aps.RptMemberWitem, error) {

	searchReq := &aps.RptMemberWitemQuery{
		SpaceIds:  req.SpaceIds,
		UserId:    req.UserId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	var err error
	var result []*aps.RptMemberWitem

	switch req.TimeSplitType {
	case "Hour":
		if req.EndDate.Sub(req.StartDate) > 24*time.Hour {
			return nil, nil
		}
		result, err = d.repo.SearchRptMemberWitem1h(ctx, searchReq)
	case "Day":
		result, err = d.repo.SearchRptMemberWitem1day(ctx, searchReq)
	case "Month":
		result, err = d.repo.SearchRptMemberWitem1month(ctx, searchReq)
	}

	return result, err
}

type DashboardRptSpaceWitemReq struct {
	// Hour, Day, Month
	TimeSplitType string
	StartDate     time.Time
	EndDate       time.Time
	SpaceIds      []int64
	UserId        int64
}

func (d *DwhAds) DashboardRptSpaceWitem(ctx context.Context, req *DashboardRptSpaceWitemReq) ([]*aps.RptSpaceWitem, error) {

	searchReq := &aps.RptSpaceWitemQuery{
		SpaceIds:  req.SpaceIds,
		UserId:    req.UserId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	var err error
	var result []*aps.RptSpaceWitem

	switch req.TimeSplitType {
	case "Hour":
		result, err = d.repo.DashboardRptSpaceWitem1h(ctx, searchReq)
	case "Day":
		result, err = d.repo.DashboardRptSpaceWitem1day(ctx, searchReq)
	case "Month":
		result, err = d.repo.DashboardRptSpaceWitem1month(ctx, searchReq)
	}

	return result, err
}

func (d *DwhAds) SpaceMembers(ctx context.Context, spaceIds []int64, endTime time.Time) (map[int64][]int64, error) {
	return d.repo.SpaceMemberMap(spaceIds, endTime)
}

type DashboardRptSpaceIncrWitemReq struct {
	// Hour, Day, Month
	TimeSplitType string
	StartDate     time.Time
	EndDate       time.Time
	SpaceIds      []int64
}

func (d *DwhAds) DashboardRptSpaceCompleteIncrWitem(ctx context.Context, req *DashboardRptSpaceIncrWitemReq) ([]*aps.RptSpaceIncrWitem, error) {
	searchReq := &aps.RptSpaceIncrWitemQuery{
		SpaceIds:  req.SpaceIds,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	var err error
	var result []*aps.RptSpaceIncrWitem

	switch req.TimeSplitType {
	case "Hour":
		result, err = d.repo.DashboardSpaceCompleteIncrHour(ctx, searchReq)
	case "Day":
		result, err = d.repo.DashboardSpaceCompleteIncrDay(ctx, searchReq)
	case "Month":
		result, err = d.repo.DashboardSpaceCompleteIncrMonth(ctx, searchReq)
	}

	return result, err
}

func (d *DwhAds) DashboardRptSpaceNewIncrWitem(ctx context.Context, req *DashboardRptSpaceIncrWitemReq) ([]*aps.RptSpaceIncrWitem, error) {
	searchReq := &aps.RptSpaceIncrWitemQuery{
		SpaceIds:  req.SpaceIds,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	var err error
	var result []*aps.RptSpaceIncrWitem

	switch req.TimeSplitType {
	case "Hour":
		result, err = d.repo.DashboardSpaceNewIncrHour(ctx, searchReq)
	case "Day":
		result, err = d.repo.DashboardSpaceNewIncrDay(ctx, searchReq)
	case "Month":
		result, err = d.repo.DashboardSpaceNewIncrMonth(ctx, searchReq)
	}

	return result, err
}

type DashboardRptMemberIncrWitemReq struct {
	// Hour, Day, Month
	TimeSplitType string
	StartDate     time.Time
	EndDate       time.Time
	SpaceIds      []int64
}

func (d *DwhAds) DashboardRptMemberCompleteIncrWitem(ctx context.Context, req *DashboardRptMemberIncrWitemReq) ([]*aps.RptMemberIncrWitem, error) {
	searchReq := &aps.RptMemberIncrWitemQuery{
		SpaceIds:  req.SpaceIds,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	var err error
	var result []*aps.RptMemberIncrWitem

	switch req.TimeSplitType {
	case "Hour":
		result, err = d.repo.DashboardMemberCompleteIncrHour(ctx, searchReq)
	case "Day":
		result, err = d.repo.DashboardMemberCompleteIncrDay(ctx, searchReq)
	case "Month":
		result, err = d.repo.DashboardMemberCompleteIncrMonth(ctx, searchReq)
	}

	return result, err
}

func (d *DwhAds) DashboardRptMemberNewIncrWitem(ctx context.Context, req *DashboardRptMemberIncrWitemReq) ([]*aps.RptMemberIncrWitem, error) {
	searchReq := &aps.RptMemberIncrWitemQuery{
		SpaceIds:  req.SpaceIds,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	result, err := d.repo.DashboardMemberNewIncr(ctx, searchReq)

	return result, err
}
