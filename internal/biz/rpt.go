package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
	v1 "go-cs/api/rpt/v1"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/consts"
	space_domain "go-cs/internal/domain/space"
	repo2 "go-cs/internal/domain/space/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	user_domain "go-cs/internal/domain/user"
	"go-cs/internal/domain/user/repo"
	login_log_repo "go-cs/internal/domain/user_login_log/repo"
	"go-cs/internal/dwh/model/aps"
	dwh_service "go-cs/internal/dwh/service"
	"go-cs/internal/utils"
	"go-cs/internal/utils/date"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/stream"
	"time"
)

type RptUsecase struct {
	dwhService      *dwh_service.DwhService
	log             *log.Helper
	userRepo        repo.UserRepo
	spaceRepo       repo2.SpaceRepo
	spaceMemberRepo member_repo.SpaceMemberRepo
	loginLogRepo    login_log_repo.UserLoginLogRepo
}

func NewRptUsecase(
	dwhService *dwh_service.DwhService,
	logger log.Logger,
	userRepo repo.UserRepo,
	spaceRepo repo2.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,
	loginLogRepo login_log_repo.UserLoginLogRepo,
) *RptUsecase {

	moduleName := "RptUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &RptUsecase{
		dwhService:      dwhService,
		log:             hlog,
		userRepo:        userRepo,
		spaceRepo:       spaceRepo,
		spaceMemberRepo: spaceMemberRepo,
		loginLogRepo:    loginLogRepo,
	}
}

func (uc *RptUsecase) SearchRptVersionWitem(ctx context.Context, oper *utils.LoginUserInfo, req *v1.SearchRptVersionWitemRequest) (*v1.SearchRptVersionWitemReplyData, error) {

	//查看是否有对应权限
	_, err := uc.spaceMemberRepo.GetSpaceMember(ctx, req.SpaceId, oper.UserId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	replyData := &v1.SearchRptVersionWitemReplyData{
		Total: 0,
		List:  make([]*v1.SearchRptVersionWitemReplyData_DataItem, 0),
	}

	startDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.StartDate)
	endDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.EndDate)

	result, err := uc.dwhService.Ads.SearchRptVersionWitem(ctx, &dwh_service.SearchRptVersionWitemReq{
		TimeSplitType: req.TimeSplitType,
		StartDate:     startDate,
		EndDate:       endDate,
		SpaceId:       req.SpaceId,
		VersionId:     req.VersionId,
		StatusId:      req.StatusId,
	})

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	result = ExpandDateRangeList(req.TimeSplitType, startDate, endDate, result)

	for i := 0; i < len(result); i++ {
		replyData.List = append(replyData.List, &v1.SearchRptVersionWitemReplyData_DataItem{
			SpaceId:     result[i].SpaceId,
			VersionId:   result[i].VersionId,
			ExpireNum:   result[i].ExpireNum,
			Num:         result[i].Num,
			TodoNum:     result[i].TodoNum,
			CloseNum:    result[i].CloseNum,
			AbortNum:    result[i].AbortNum,
			CompleteNum: result[i].CompleteNum,
			StartDate:   result[i].StartDate.Format("2006-01-02 15:04:05"),
			EndDate:     result[i].EndDate.Format("2006-01-02 15:04:05"),
		})
	}

	replyData.Total = cast.ToInt64(len(replyData.List))
	return replyData, nil
}

func (uc *RptUsecase) SearchRptMemberWitem(ctx context.Context, oper *utils.LoginUserInfo, req *v1.SearchRptMemberWitemRequest) (*v1.SearchRptMemberWitemReplyData, error) {

	//查看是否有对应权限
	_, err := uc.spaceMemberRepo.GetSpaceMember(ctx, req.SpaceId, oper.UserId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	replyData := &v1.SearchRptMemberWitemReplyData{
		Total: 0,
		List:  make([]*v1.SearchRptMemberWitemReplyData_DataItem, 0),
	}

	startDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.StartDate)
	endDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.EndDate)

	result, err := uc.dwhService.Ads.SearchRptMemberWitem(ctx, &dwh_service.SearchRptMemberWitemReq{
		TimeSplitType: req.TimeSplitType,
		StartDate:     startDate,
		EndDate:       endDate,
		SpaceIds:      []int64{req.SpaceId},
		UserId:        req.UserId,
	})

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	result = ExpandDateRangeList(req.TimeSplitType, startDate, endDate, result)

	for i := 0; i < len(result); i++ {
		replyData.List = append(replyData.List, &v1.SearchRptMemberWitemReplyData_DataItem{
			SpaceId:     result[i].SpaceId,
			UserId:      result[i].UserId,
			Num:         result[i].Num,
			ExpireNum:   result[i].ExpireNum,
			TodoNum:     result[i].TodoNum,
			CloseNum:    result[i].CloseNum,
			AbortNum:    result[i].AbortNum,
			CompleteNum: result[i].CompleteNum,
			StartDate:   result[i].StartDate.Format("2006-01-02 15:04:05"),
			EndDate:     result[i].EndDate.Format("2006-01-02 15:04:05"),
		})
	}

	replyData.Total = cast.ToInt64(len(replyData.List))
	return replyData, nil
}

func (uc *RptUsecase) DashboardRptSpaceWitem(ctx context.Context, oper *utils.LoginUserInfo, req *v1.DashboardRptSpaceWitemRequest) (*v1.DashboardRptSpaceWitemReplyData, error) {
	if len(req.SpaceIds) == 0 {
		return nil, nil
	}

	startDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.StartDate)
	endDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.EndDate)

	result, err := uc.dwhService.Ads.DashboardRptSpaceWitem(ctx, &dwh_service.DashboardRptSpaceWitemReq{
		TimeSplitType: req.TimeSplitType,
		StartDate:     startDate,
		EndDate:       GetEndDate(endDate),
		SpaceIds:      req.SpaceIds,
	})

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	groupBy := stream.GroupBy(result, func(item *aps.RptSpaceWitem) int64 {
		return item.SpaceId
	})

	finalResult := make([]*aps.RptSpaceWitem, 0)
	for _, v := range groupBy {
		list := ExpandDateRangeList(req.TimeSplitType, startDate, endDate, v)
		finalResult = append(finalResult, list...)
	}

	var items []*v1.DashboardRptSpaceWitemReplyData_DataItem
	for _, v := range finalResult {
		items = append(items, &v1.DashboardRptSpaceWitemReplyData_DataItem{
			SpaceId:     v.SpaceId,
			Num:         v.Num,
			ExpireNum:   v.ExpireNum,
			TodoNum:     v.TodoNum,
			CloseNum:    v.CloseNum,
			AbortNum:    v.AbortNum,
			CompleteNum: v.CompleteNum,
			StartDate:   v.StartDate.Format("2006-01-02 15:04:05"),
			EndDate:     v.EndDate.Format("2006-01-02 15:04:05"),
		})
	}

	return &v1.DashboardRptSpaceWitemReplyData{
		Total: int64(len(items)),
		List:  items,
	}, nil
}

func (uc *RptUsecase) DashboardRptMemberWitem(ctx context.Context, oper *utils.LoginUserInfo, req *v1.DashboardRptMemberWitemRequest) (*v1.DashboardRptMemberWitemReplyData, error) {
	if len(req.SpaceIds) == 0 {
		return nil, nil
	}

	startDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.StartDate)
	endDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.EndDate)

	result, err := uc.dwhService.Ads.SearchRptMemberWitem(ctx, &dwh_service.SearchRptMemberWitemReq{
		TimeSplitType: req.TimeSplitType,
		StartDate:     startDate,
		EndDate:       GetEndDate(endDate),
		SpaceIds:      req.SpaceIds,
	})

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	updatedAt := getUpdatedAt()

	members, err := uc.dwhService.Ads.SpaceMembers(ctx, req.SpaceIds, updatedAt)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	// 过滤不在项目的
	result = stream.Filter(result, func(item *aps.RptMemberWitem) bool {
		userIds := members[item.SpaceId]
		return stream.Contains(userIds, item.UserId)
	})

	groupBy := stream.GroupBy(result, func(item *aps.RptMemberWitem) string {
		return fmt.Sprintf("%v:%v", item.SpaceId, item.UserId)
	})

	finalResult := make([]*aps.RptMemberWitem, 0, len(result))
	for _, v := range groupBy {
		items := ExpandDateRangeList(req.TimeSplitType, startDate, endDate, v)
		finalResult = append(finalResult, items...)
	}

	var list []*v1.DashboardRptMemberWitemReplyData_DataItem
	for _, v := range finalResult {
		list = append(list, &v1.DashboardRptMemberWitemReplyData_DataItem{
			SpaceId:     v.SpaceId,
			UserId:      v.UserId,
			Num:         v.Num,
			ExpireNum:   v.ExpireNum,
			TodoNum:     v.TodoNum,
			CloseNum:    v.CloseNum,
			AbortNum:    v.AbortNum,
			CompleteNum: v.CompleteNum,
			StartDate:   v.StartDate.Format("2006-01-02 15:04:05"),
			EndDate:     v.EndDate.Format("2006-01-02 15:04:05"),
		})
	}

	return &v1.DashboardRptMemberWitemReplyData{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func (uc *RptUsecase) DashboardSpaceList(ctx context.Context, uid int64) (*v1.DashboardSpaceListReply_Data, error) {
	userInfo, err := uc.userRepo.GetUserByUserId(ctx, uid)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	var spaceList []*space_domain.Space

	switch userInfo.Role {
	case consts.SystemRole_Enterprise:
		spaceList, err = uc.spaceRepo.GetAllSpace(ctx)
	default:
		spaceIds, err := uc.spaceMemberRepo.GetUserSpaceIdList(ctx, uid)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		spaceList, err = uc.spaceRepo.GetSpaceByIds(ctx, spaceIds)
	}

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	updatedAt := getUpdatedAt()

	spaceIds := stream.Unique(stream.Map(spaceList, func(item *space_domain.Space) int64 {
		return item.Id
	}))

	members, err := uc.dwhService.Ads.SpaceMembers(ctx, spaceIds, updatedAt)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	userIds := stream.Unique(stream.Flat(stream.Values(members)))
	userMap, err := uc.userRepo.UserMap(ctx, userIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	userInfoMap := stream.MapValue(userMap, func(value *user_domain.User) *rsp.SimpleUserInfo {
		return &rsp.SimpleUserInfo{
			Id:           value.Id,
			UserId:       value.Id,
			UserNickname: value.UserNickname,
			UserName:     value.UserName,
			Avatar:       value.Avatar,
		}
	})

	var items []*v1.DashboardSpaceListReply_Data_Item
	for _, space := range spaceList {
		items = append(items, &v1.DashboardSpaceListReply_Data_Item{
			Id:        space.Id,
			Name:      space.SpaceName,
			CreatedAt: space.CreatedAt,
			Members:   members[space.Id],
		})
	}

	return &v1.DashboardSpaceListReply_Data{
		List:         items,
		Users:        stream.Values(userInfoMap),
		UpdatedAt:    updatedAt.Unix(),
		SystemInitAt: uc.loginLogRepo.GetLatestLoginTime(ctx),
	}, nil
}

func getUpdatedAt() time.Time {
	now := time.Now()
	var updateAt = date.HourBegin(now)
	if now.Sub(updateAt) < time.Minute*15 {
		updateAt = updateAt.Add(-time.Hour)
	}

	return updateAt
}

func (uc *RptUsecase) DashboardRptSpaceIncrWitem(ctx context.Context, uid int64, req *v1.DashboardSpaceIncrWitemRequest) (*v1.DashboardSpaceIncrWitemReply_Data, error) {
	if len(req.SpaceIds) == 0 {
		return nil, nil
	}

	startDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.StartDate)
	endDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.EndDate)

	var result []*aps.RptSpaceIncrWitem
	var err error

	r := &dwh_service.DashboardRptSpaceIncrWitemReq{
		TimeSplitType: req.TimeSplitType,
		StartDate:     startDate,
		EndDate:       GetEndDate(endDate),
		SpaceIds:      req.SpaceIds,
	}

	switch req.Type {
	case "completed":
		result, err = uc.dwhService.Ads.DashboardRptSpaceCompleteIncrWitem(ctx, r)
	case "created":
		result, err = uc.dwhService.Ads.DashboardRptSpaceNewIncrWitem(ctx, r)
	default:
		return nil, errs.Business(ctx, "type error")
	}

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var list []*v1.DashboardSpaceIncrWitemReply_Data_Item
	for _, v := range result {
		list = append(list, &v1.DashboardSpaceIncrWitemReply_Data_Item{
			SpaceId:   v.SpaceId,
			Num:       v.Num,
			StartDate: v.StartDate.Format("2006-01-02 15:04:05"),
		})
	}

	return &v1.DashboardSpaceIncrWitemReply_Data{
		List: list,
	}, nil
}

func (uc *RptUsecase) DashboardRptMemberIncrWitem(ctx context.Context, uid int64, req *v1.DashboardMemberIncrWitemRequest) (*v1.DashboardMemberIncrWitemReply_Data, error) {
	if len(req.SpaceIds) == 0 {
		return nil, nil
	}

	startDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.StartDate)
	endDate := date.ParseInLocation("2006-01-02 15:04:05", req.TimeRange.EndDate)

	var result []*aps.RptMemberIncrWitem
	var err error

	r := &dwh_service.DashboardRptMemberIncrWitemReq{
		TimeSplitType: req.TimeSplitType,
		StartDate:     startDate,
		EndDate:       GetEndDate(endDate),
		SpaceIds:      req.SpaceIds,
	}

	switch req.Type {
	case "completed":
		result, err = uc.dwhService.Ads.DashboardRptMemberCompleteIncrWitem(ctx, r)
	case "created":
		result, err = uc.dwhService.Ads.DashboardRptMemberNewIncrWitem(ctx, r)
	default:
		return nil, errs.Business(ctx, "type error")
	}

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var list []*v1.DashboardMemberIncrWitemReply_Data_Item
	for _, v := range result {
		list = append(list, &v1.DashboardMemberIncrWitemReply_Data_Item{
			SpaceId:   v.SpaceId,
			UserId:    v.UserId,
			Num:       v.Num,
			StartDate: v.StartDate.Format("2006-01-02 15:04:05"),
		})
	}

	return &v1.DashboardMemberIncrWitemReply_Data{
		List: list,
	}, nil
}
