package service

import (
	"context"
	"encoding/json"
	"go-cs/api/comm"
	pb "go-cs/api/space_work_item/v1"
	"go-cs/internal/bean/vo"
	"go-cs/internal/biz"
	"go-cs/internal/consts"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/locker"
	"go-cs/pkg/stream"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceWorkItemService struct {
	pb.UnimplementedSpaceWorkItemServer
	log *log.Helper
	uc  *biz.SpaceWorkItemUsecase
}

func NewSpaceWorkItemService(uc *biz.SpaceWorkItemUsecase, logger log.Logger) *SpaceWorkItemService {
	moduleName := "SpaceWorkItemService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceWorkItemService{
		uc:  uc,
		log: hlog,
	}
}

func (s *SpaceWorkItemService) CreateWorkItemV2(ctx context.Context, req *pb.CreateWorkItemRequestV2) (*pb.CreateWorkItemReplyV2, error) {

	reply := func(err error) (*pb.CreateWorkItemReplyV2, error) {
		return &pb.CreateWorkItemReplyV2{Result: &pb.CreateWorkItemReplyV2_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkObjectId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkObjectId"))
	}

	if vaildErr = validate.Var(req.WorkVersionId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkVersionId"))
	}

	if vaildErr = validate.Var(req.FlowId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "FlowId"))
	}

	if vaildErr = validate.Var(strings.TrimSpace(req.ItemName), "required,utf8Len=2-200"); vaildErr != nil {
		return reply(errs.New(ctx, comm.ErrorCode_SPACE_WORK_ITEM_NAME_RULE_WRONG, "ItemName"))
	}

	if !consts.IsValidIconFlags(req.IconFlags...) {
		return reply(errs.Param(ctx, "icon_flags"))
	}

	if req.Tag != nil && len(req.Tag.New) > 0 {
		//创建的TAG分类
		for _, v := range req.Tag.New {
			if vaildErr = validate.Var(strings.TrimSpace(v), "utf8Len=2-10"); vaildErr != nil {
				return reply(errs.Param(ctx, "Tag.New"))
			}
		}
	}

	if req.File != nil && len(req.File.New) > 0 {
		//创建的TAG分类
		for _, v := range req.File.New {
			if vaildErr = validate.Var(v, "required,number,gt=0"); vaildErr != nil {
				return reply(errs.Param(ctx, "File.New"))
			}
		}
	}

	in := vo.CreateSpaceWorkItemVoV2{
		WorkItemName:   strings.TrimSpace(req.ItemName),
		WorkItemTypeId: req.WorkItemType,
		PlanStartAt:    0,
		PlanCompleteAt: 0,
		FileAdd:        make([]*vo.CreateSpaceWorkItemFileVoV2, 0),
		Owner:          make([]*vo.CreateSpaceWorkItemOwnerV2, 0),
		WorkVersionId:  req.WorkVersionId,
		WorkFlowId:     req.FlowId,
		ProcessRate:    req.ProgressRate,
		Describe:       strings.TrimSpace(req.Describe),
		Remark:         strings.TrimSpace(req.Remark),
		Priority:       strings.TrimSpace(req.Priority),
		SpaceId:        req.SpaceId,
		WorkObjectId:   req.WorkObjectId,
		UserId:         loginUser.UserId,
		IconFlags:      req.IconFlags,
		Followers:      stream.Unique(req.Followers),
	}

	reqNodes := req.GetOwner()
	for _, reqNode := range reqNodes {
		in.Owner = append(in.Owner, &vo.CreateSpaceWorkItemOwnerV2{
			OwnerRole: reqNode.Role,
			Directors: utils.ToInt64Array(reqNode.DirectorId),
		})
	}

	if req.Tag != nil {
		in.TagAdd = req.Tag.Add
	}

	if req.File != nil && len(req.File.Add) > 0 {
		for _, v := range req.File.Add {
			fileInfo := &vo.CreateSpaceWorkItemFileVoV2{}
			fileInfo.Id = v
			in.FileAdd = append(in.FileAdd, fileInfo)
		}
	}

	if req.PlanTimeAt != nil && req.PlanTimeAt.Start != "" {
		if planStartTime, err := time.ParseInLocation("2006/01/02 15:04:05", req.PlanTimeAt.Start, time.Local); err == nil {
			in.PlanStartAt = planStartTime.Unix()
		}

	}

	if req.PlanTimeAt != nil && req.PlanTimeAt.Complete != "" {
		if planCompleteTime, err := time.ParseInLocation("2006/01/02 15:04:05", req.PlanTimeAt.Complete, time.Local); err == nil {
			in.PlanCompleteAt = planCompleteTime.Unix()
		}
	}

	out, err := s.uc.CreateTask(ctx, loginUser, in)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.CreateWorkItemReplyV2{Result: &pb.CreateWorkItemReplyV2_Data{Data: cast.ToString(out)}}
	return okReply, nil

}

func (s *SpaceWorkItemService) ConfirmWorkFlowMain(ctx context.Context, req *pb.ConfirmWorkFlowMainRequest) (*pb.ConfirmWorkFlowMainReply, error) {

	reply := func(err error) (*pb.ConfirmWorkFlowMainReply, error) {
		return &pb.ConfirmWorkFlowMainReply{Result: &pb.ConfirmWorkFlowMainReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	if req.WorkItemId <= 0 {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if req.NodeState == "" {
		return reply(errs.Param(ctx, "NodeState"))
	}

	lock := locker.Lock(locker.NewWorkItemLockKey(req.WorkItemId))
	lock.Lock()
	defer lock.Unlock()

	out, err := s.uc.ConfirmWorkFlowMain(ctx, loginUser, req.WorkItemId, req.NodeState)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ConfirmWorkFlowMainReply{Result: &pb.ConfirmWorkFlowMainReply_Data{Data: &pb.ConfirmWorkItemNodeStateReplyV2Data{Ids: out}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) CreateWorkItemSubTask(ctx context.Context, req *pb.CreateWorkItemSubTaskRequest) (*pb.CreateWorkItemSubTaskReply, error) {

	reply := func(err error) (*pb.CreateWorkItemSubTaskReply, error) {
		return &pb.CreateWorkItemSubTaskReply{Result: &pb.CreateWorkItemSubTaskReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var err error
	validate := utils.NewValidator()

	if err = validate.Var(req.SpaceId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if err = validate.Var(strings.TrimSpace(req.ItemName), "required,utf8Len=2-200"); err != nil {
		return reply(errs.New(ctx, comm.ErrorCode_SPACE_WORK_ITEM_NAME_RULE_WRONG))
	}

	//必须有负责人
	if err = validate.Var(req.Director.Add, "required,gt=0,dive"); err != nil {
		return reply(errs.Param(ctx, "Director.Add"))
	}

	in := vo.CreateSpaceWorkItemTaskVoV2{}
	in.WorkItemName = strings.TrimSpace(req.ItemName)
	in.PlanStartAt = 0
	in.PlanCompleteAt = 0
	if req.Director != nil {
		in.DirectorAdd = req.Director.Add
	}

	if req.PlanTimeAt != nil && req.PlanTimeAt.Start != "" {
		if planStartTime, err := time.ParseInLocation("2006/01/02 15:04:05", req.PlanTimeAt.Start, time.Local); err == nil {
			in.PlanStartAt = planStartTime.Unix()
		}
	}

	if req.PlanTimeAt != nil && req.PlanTimeAt.Complete != "" {
		if planCompleteTime, err := time.ParseInLocation("2006/01/02 15:04:05", req.PlanTimeAt.Complete, time.Local); err == nil {
			in.PlanCompleteAt = planCompleteTime.Unix()
		}
	}

	in.ProcessRate = req.ProgressRate

	out, err := s.uc.CreateSubTask(ctx, loginUser, req.SpaceId, req.WorkItemId, in)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.CreateWorkItemSubTaskReply{Result: &pb.CreateWorkItemSubTaskReply_Data{Data: cast.ToString(out)}}
	return okReply, nil

}

func (s *SpaceWorkItemService) GetWorkItemDetailV2(ctx context.Context, req *pb.GetWorkItemDetailRequestV2) (*pb.GetWorkItemDetailReplyV2, error) {

	reply := func(err error) (*pb.GetWorkItemDetailReplyV2, error) {
		return &pb.GetWorkItemDetailReplyV2{Result: &pb.GetWorkItemDetailReplyV2_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	//以上数据检查完毕，由于入参过于复杂，直接把req传递到biz
	out, err := s.uc.GetWorkItemDetail(ctx, loginUser, req.SpaceId, req.WorkItemId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.GetWorkItemDetailReplyV2{Result: &pb.GetWorkItemDetailReplyV2_Data{Data: out}}
	return okReply, nil
}

func (s *SpaceWorkItemService) SetFlowMainDirector(ctx context.Context, req *pb.SetFlowMainDirectorRequest) (*pb.SetFlowMainDirectorReply, error) {

	var reply = func(err *comm.ErrorInfo) (*pb.SetFlowMainDirectorReply, error) {
		return &pb.SetFlowMainDirectorReply{Result: &pb.SetFlowMainDirectorReply_Error{Error: err}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var err error
	validate := utils.NewValidator()

	if err = validate.Var(req.WorkItemId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if err = validate.Var(req.Director, "required"); err != nil {
		return reply(errs.Param(ctx, "Director"))
	}

	if err = validate.Var(req.Role, "required"); err != nil {
		return reply(errs.Param(ctx, "Role"))
	}

	if len(req.Director.Remove) == 0 && len(req.Director.Add) == 0 {
		return reply(errs.Param(ctx, "Director.Remove/Add"))
	}

	lock := locker.Lock(locker.NewWorkItemLockKey(req.WorkItemId))
	lock.Lock()
	defer lock.Unlock()

	err = s.uc.SetFlowMainDirectorByRoleKey(ctx, loginUser, req.WorkItemId, req.Role, req.Director.Add, req.Director.Remove)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &pb.SetFlowMainDirectorReply{}, nil
}

func (s *SpaceWorkItemService) SetSubDirector(ctx context.Context, req *pb.SetSubDirectorRequest) (*pb.SetSubDirectorReply, error) {

	var reply = func(err *comm.ErrorInfo) (*pb.SetSubDirectorReply, error) {
		return &pb.SetSubDirectorReply{Result: &pb.SetSubDirectorReply_Error{Error: err}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var err error
	validate := utils.NewValidator()

	if err = validate.Var(req.WorkItemId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if err = validate.Var(req.Director, "required"); err != nil {
		return reply(errs.Param(ctx, "Director"))
	}

	if len(req.Director.Remove) == 0 && len(req.Director.Add) == 0 {
		return reply(errs.Param(ctx, "Director.Remove/Add"))
	}

	err = s.uc.SetSubDirector(ctx, loginUser, req.WorkItemId, req.Director.Add, req.Director.Remove)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &pb.SetSubDirectorReply{}, nil
}

func (s *SpaceWorkItemService) SetWorkItemToNewWorkObjectV2(ctx context.Context, req *pb.SetWorkItemToNewWorkObjectRequest) (*pb.SetWorkItemToNewWorkObjectReply, error) {

	reply := func(err error) (*pb.SetWorkItemToNewWorkObjectReply, error) {
		return &pb.SetWorkItemToNewWorkObjectReply{Result: &pb.SetWorkItemToNewWorkObjectReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if vaildErr = validate.Var(req.WorkObjectId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkObjectId"))
	}

	_, err := s.uc.ChangeWorkItemObject(ctx, loginUser, req.SpaceId, req.WorkItemId, req.WorkObjectId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetWorkItemToNewWorkObjectReply{Result: &pb.SetWorkItemToNewWorkObjectReply_Data{Data: &pb.SetWorkItemToNewWorkObjectReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) ModifyWorkItemPlanTimeV2(ctx context.Context, req *pb.ModifyWorkItemPlanTimeRequest) (*pb.ModifyWorkItemPlanTimeReply, error) {

	reply := func(err error) (*pb.ModifyWorkItemPlanTimeReply, error) {
		return &pb.ModifyWorkItemPlanTimeReply{Result: &pb.ModifyWorkItemPlanTimeReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	var planStartAt int64
	if req.PlanTimeAt != nil && req.PlanTimeAt.Start != "" {
		if planStartTime, err := time.ParseInLocation("2006/01/02 15:04:05", req.PlanTimeAt.Start, time.Local); err == nil {
			planStartAt = planStartTime.Unix()
		} else {
			return reply(errs.Param(ctx, "PlanTimeAt.Start"))
		}
	}

	var planCompleteAt int64
	if req.PlanTimeAt != nil && req.PlanTimeAt.Complete != "" {
		if planCompleteTime, err := time.ParseInLocation("2006/01/02 15:04:05", req.PlanTimeAt.Complete, time.Local); err == nil {
			planCompleteAt = planCompleteTime.Unix()
		} else {
			return reply(errs.Param(ctx, "PlanTimeAt.Complete"))
		}
	}

	err := s.uc.ModifyWorkItemPlanTime(ctx, loginUser, req.SpaceId, req.WorkItemId, planStartAt, planCompleteAt)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ModifyWorkItemPlanTimeReply{Result: &pb.ModifyWorkItemPlanTimeReply_Data{Data: &pb.ModifyWorkItemPlanTimeReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) ModifyWorkItemProcessRateV2(ctx context.Context, req *pb.ModifyWorkItemProcessRateRequest) (*pb.ModifyWorkItemProcessRateReply, error) {

	reply := func(err error) (*pb.ModifyWorkItemProcessRateReply, error) {
		return &pb.ModifyWorkItemProcessRateReply{Result: &pb.ModifyWorkItemProcessRateReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if vaildErr = validate.Var(req.ProcessRate, "gte=0,lte=100"); vaildErr != nil {
		return reply(errs.Param(ctx, "ProcessRate"))
	}

	err := s.uc.ModifyWorkItemProcessRate(ctx, loginUser, req.SpaceId, req.WorkItemId, int64(req.ProcessRate))
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ModifyWorkItemProcessRateReply{Result: &pb.ModifyWorkItemProcessRateReply_Data{Data: &pb.ModifyWorkItemProcessRateReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) ModifyWorkItemPriorityV2(ctx context.Context, req *pb.ModifyWorkItemPriorityRequest) (*pb.ModifyWorkItemPriorityReply, error) {

	reply := func(err error) (*pb.ModifyWorkItemPriorityReply, error) {
		return &pb.ModifyWorkItemPriorityReply{Result: &pb.ModifyWorkItemPriorityReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if vaildErr = validate.Var(req.Priority, "required"); vaildErr != nil {
		return reply(errs.Param(ctx, "Priority"))
	}

	err := s.uc.ModifyWorkItemPriority(ctx, loginUser, req.SpaceId, req.WorkItemId, req.Priority)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ModifyWorkItemPriorityReply{Result: &pb.ModifyWorkItemPriorityReply_Data{Data: &pb.ModifyWorkItemPriorityReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) ModifyWorkItemDescribeV2(ctx context.Context, req *pb.ModifyWorkItemDescribeRequest) (*pb.ModifyWorkItemDescribeReply, error) {

	reply := func(err error) (*pb.ModifyWorkItemDescribeReply, error) {
		return &pb.ModifyWorkItemDescribeReply{Result: &pb.ModifyWorkItemDescribeReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	err := s.uc.ModifyWorkItemDescribe(ctx, loginUser, req.SpaceId, req.WorkItemId, req.Describe, req.IconFlags)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ModifyWorkItemDescribeReply{Result: &pb.ModifyWorkItemDescribeReply_Data{Data: &pb.ModifyWorkItemDescribeReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) ModifyWorkItemNameV2(ctx context.Context, req *pb.ModifyWorkItemNameRequest) (*pb.ModifyWorkItemNameReply, error) {

	reply := func(err error) (*pb.ModifyWorkItemNameReply, error) {
		return &pb.ModifyWorkItemNameReply{Result: &pb.ModifyWorkItemNameReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if vaildErr = validate.Var(strings.TrimSpace(req.WorkItemName), "required,utf8Len=2-200"); vaildErr != nil {
		return reply(errs.Business(ctx, "请输入任务名称（2~200个字符）"))
	}

	err := s.uc.ModifyWorkItemName(ctx, loginUser, req.SpaceId, req.WorkItemId, req.WorkItemName)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ModifyWorkItemNameReply{Result: &pb.ModifyWorkItemNameReply_Data{Data: &pb.ModifyWorkItemNameReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) SetSpaceWorkItemFileInfoV2(ctx context.Context, req *pb.SetSpaceWorkItemFileInfoRequest) (*pb.SetSpaceWorkItemFileInfoReply, error) {

	reply := func(err error) (*pb.SetSpaceWorkItemFileInfoReply, error) {
		return &pb.SetSpaceWorkItemFileInfoReply{Result: &pb.SetSpaceWorkItemFileInfoReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)

	if req.WorkItemId <= 0 {
		return reply(errs.Param(ctx, "WorkItemId"))
	}
	if req.File == nil || len(req.File.Remove) == 0 && len(req.File.Add) == 0 {
		return reply(errs.Param(ctx, "File"))
	}

	in := vo.SetSpaceWorkItemFileInfoVoV2{
		WorkItemId:     req.WorkItemId,
		FileInfoAdd:    req.File.Add,
		FileInfoRemove: req.File.Remove,
	}

	err := s.uc.SetWorkItemFileInfo(ctx, loginUser, in)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetSpaceWorkItemFileInfoReply{Result: &pb.SetSpaceWorkItemFileInfoReply_Data{Data: &pb.SetSpaceWorkItemFileInfoReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) DelSpaceWorkItemV2(ctx context.Context, req *pb.DelSpaceWorkItemRequest) (*pb.DelSpaceWorkItemReply, error) {

	reply := func(err error) (*pb.DelSpaceWorkItemReply, error) {
		return &pb.DelSpaceWorkItemReply{Result: &pb.DelSpaceWorkItemReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	err := s.uc.DelWorkItem(ctx, loginUser, req.WorkItemId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.DelSpaceWorkItemReply{Result: &pb.DelSpaceWorkItemReply_Data{Data: &pb.DelSpaceWorkItemReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) SetWorkItemTagV2(ctx context.Context, req *pb.SetWorkItemTagRequest) (*pb.SetWorkItemTagReply, error) {

	reply := func(err error) (*pb.SetWorkItemTagReply, error) {
		return &pb.SetWorkItemTagReply{Result: &pb.SetWorkItemTagReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if vaildErr = validate.Var(req.Tag, "required"); vaildErr != nil {
		return reply(errs.Param(ctx, "Tag"))
	}

	//不允许空操作
	if len(req.Tag.Remove) == 0 && len(req.Tag.Add) == 0 && len(req.Tag.New) == 0 {
		return reply(errs.Param(ctx, "Director"))
	}

	if len(req.Tag.New) > 0 {
		//创建的TAG分类
		for _, v := range req.Tag.New {
			if vaildErr = validate.Var(v, "utf8Len=2-10"); vaildErr != nil {
				return reply(errs.Param(ctx, "Tag.New"))
			}
		}
	}

	in := vo.SetSpaceWorkItemTagVoV2{}
	in.SpaceId = req.SpaceId
	in.WorkItemId = req.WorkItemId
	in.TagAdd = req.Tag.Add
	in.TagNew = req.Tag.New
	in.TagRemove = req.Tag.Remove
	err := s.uc.SetWorkItemTag(ctx, loginUser, in)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetWorkItemTagReply{Result: &pb.SetWorkItemTagReply_Data{Data: &pb.SetWorkItemTagReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) TerminateSpaceWorkItemV2(ctx context.Context, req *pb.TerminateSpaceWorkItemRequest) (*pb.TerminateSpaceWorkItemReply, error) {

	reply := func(err error) (*pb.TerminateSpaceWorkItemReply, error) {
		return &pb.TerminateSpaceWorkItemReply{Result: &pb.TerminateSpaceWorkItemReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	out, err2 := s.uc.TerminateWorkItem(ctx, loginUser, req.WorkItemId, req.Reason)
	if err2 != nil {
		return reply(err2)
	}

	okReply := &pb.TerminateSpaceWorkItemReply{Result: &pb.TerminateSpaceWorkItemReply_Data{Data: &pb.TerminateSpaceWorkItemReplyData{Ids: out}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) RestartSpaceWorkItemV2(ctx context.Context, req *pb.RestartSpaceWorkItemRequest) (*pb.RestartSpaceWorkItemReply, error) {

	reply := func(err error) (*pb.RestartSpaceWorkItemReply, error) {
		return &pb.RestartSpaceWorkItemReply{Result: &pb.RestartSpaceWorkItemReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	out, err := s.uc.RestartTask(ctx, loginUser, req.SpaceId, req.WorkItemId, req.NodeCode, req.Reason)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.RestartSpaceWorkItemReply{Result: &pb.RestartSpaceWorkItemReply_Data{Data: &pb.RestartSpaceWorkItemReplyData{Ids: out}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) CloseSpaceWorkItemV2(ctx context.Context, req *pb.CloseSpaceWorkItemRequest) (*pb.CloseSpaceWorkItemReply, error) {
	reply := func(err *comm.ErrorInfo) (*pb.CloseSpaceWorkItemReply, error) {
		return &pb.CloseSpaceWorkItemReply{Result: &pb.CloseSpaceWorkItemReply_Error{Error: err}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var validErr error
	validate := utils.NewValidator()

	if validErr = validate.Var(req.SpaceId, "required,gt=0,number"); validErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if validErr = validate.Var(req.WorkItemId, "required,gt=0,number"); validErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if req.NodeCode == "" {
		return reply(errs.Param(ctx, "NodeCode"))
	}

	out, err := s.uc.CloseTask(ctx, loginUser, req.SpaceId, req.WorkItemId, req.NodeCode, req.Reason)
	if err != nil {
		return reply(errs.Cast(err))
	}

	okReply := &pb.CloseSpaceWorkItemReply{Result: &pb.CloseSpaceWorkItemReply_Data{Data: &pb.CloseSpaceWorkItemReplyData{Ids: out}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) RollbackWorkItemNodeStateV2(ctx context.Context, req *pb.RollbackSpaceWorkItemRequest) (*pb.RollbackSpaceWorkItemReply, error) {

	reply := func(err error) (*pb.RollbackSpaceWorkItemReply, error) {
		return &pb.RollbackSpaceWorkItemReply{Result: &pb.RollbackSpaceWorkItemReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	lock := locker.Lock(locker.NewWorkItemLockKey(req.WorkItemId))
	lock.Lock()
	defer lock.Unlock()

	out, err := s.uc.RollbackTask(ctx, loginUser, req.SpaceId, req.WorkItemId, req.NodeCode, req.Reason)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.RollbackSpaceWorkItemReply{Result: &pb.RollbackSpaceWorkItemReply_Data{Data: &pb.RollbackSpaceWorkItemReplyData{Ids: out}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) ResumeWorkItemV2(ctx context.Context, req *pb.ResumeSpaceWorkItemRequest) (*pb.ResumeSpaceWorkItemReply, error) {

	reply := func(err error) (*pb.ResumeSpaceWorkItemReply, error) {
		return &pb.ResumeSpaceWorkItemReply{Result: &pb.ResumeSpaceWorkItemReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,gt=0,number"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	lock := locker.Lock(locker.NewWorkItemLockKey(req.WorkItemId))
	lock.Lock()
	defer lock.Unlock()

	out, err := s.uc.ResumeTask(ctx, loginUser, req.SpaceId, req.WorkItemId, req.Reason)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ResumeSpaceWorkItemReply{Result: &pb.ResumeSpaceWorkItemReply_Data{Data: &pb.ResumeSpaceWorkItemReplyData{Ids: out}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) ConfirmSub(ctx context.Context, req *pb.ConfirmSubRequest) (*pb.ConfirmSubReply, error) {

	reply := func(err error) (*pb.ConfirmSubReply, error) {
		return &pb.ConfirmSubReply{Result: &pb.ConfirmSubReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)

	if req.WorkItemId <= 0 {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if req.State == "" {
		return reply(errs.Param(ctx, "State"))
	}

	err := s.uc.ConfirmSub(ctx, loginUser, req.WorkItemId, req.State, req.Reason)
	if err != nil {
		return reply(err)
	}

	return &pb.ConfirmSubReply{}, nil
}

func (s *SpaceWorkItemService) Follow(ctx context.Context, req *pb.FollowRequest) (*pb.FollowReply, error) {
	reply := func(err error) (*pb.FollowReply, error) {
		return &pb.FollowReply{Result: &pb.FollowReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.WorkItemId <= 0 {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	err := s.uc.FollowWorkItem(ctx, utils.GetLoginUser(ctx), req.WorkItemId, req.Unfollow)
	if err != nil {
		return reply(err)
	}

	return &pb.FollowReply{}, nil
}

func (s *SpaceWorkItemService) BatchDelSpaceWorkItemV2(ctx context.Context, req *pb.BatchDelSpaceWorkItemRequest) (*pb.BatchDelSpaceWorkItemReply, error) {

	reply := func(err error) (*pb.BatchDelSpaceWorkItemReply, error) {
		return &pb.BatchDelSpaceWorkItemReply{Result: &pb.BatchDelSpaceWorkItemReply_Error{Error: errs.Cast(err)}}, nil
	}

	return reply(errs.NoPerm(ctx))
}

func (s *SpaceWorkItemService) BatchTerminateSpaceWorkItemV2(ctx context.Context, req *pb.BatchTerminateSpaceWorkItemRequest) (*pb.BatchTerminateSpaceWorkItemReply, error) {

	reply := func(err error) (*pb.BatchTerminateSpaceWorkItemReply, error) {
		return &pb.BatchTerminateSpaceWorkItemReply{Result: &pb.BatchTerminateSpaceWorkItemReply_Error{Error: errs.Cast(err)}}, nil
	}

	return reply(errs.NoPerm(ctx))
}

func (s *SpaceWorkItemService) ChangeWorkItemVersion(ctx context.Context, req *pb.ChangeWorkItemVersionRequest) (*pb.ChangeWorkItemVersionReply, error) {

	reply := func(err error) (*pb.ChangeWorkItemVersionReply, error) {
		return &pb.ChangeWorkItemVersionReply{Result: &pb.ChangeWorkItemVersionReply_Error{Error: errs.Cast(err)}}, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.VersionId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "VersionId"))
	}

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	//进入逻辑
	effectIds, err := s.uc.ChangeTaskVersion(ctx, loginUser, req.SpaceId, req.WorkItemId, req.VersionId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ChangeWorkItemVersionReply{Result: &pb.ChangeWorkItemVersionReply_Data{Data: &pb.ChangeWorkItemVersionReplyData{Ids: effectIds}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) ModifyWorkItemRemark(ctx context.Context, req *pb.ModifyWorkItemRemarkRequest) (*pb.ModifyWorkItemRemarkReply, error) {

	reply := func(err error) (*pb.ModifyWorkItemRemarkReply, error) {
		return &pb.ModifyWorkItemRemarkReply{Result: &pb.ModifyWorkItemRemarkReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	err := s.uc.ModifyWorkItemRemark(ctx, loginUser, req.SpaceId, req.WorkItemId, req.Remark)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ModifyWorkItemRemarkReply{Result: &pb.ModifyWorkItemRemarkReply_Data{Data: &pb.ModifyWorkItemRemarkReplyData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemService) RemindWork(ctx context.Context, req *pb.RemindWorkRequest) (*pb.RemindWorkReply, error) {

	reply := func(err error) (*pb.RemindWorkReply, error) {
		return &pb.RemindWorkReply{Result: &pb.RemindWorkReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.WorkItemId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	err := s.uc.RemindWork(ctx, loginUser, req.SpaceId, req.WorkItemId, req.NodeCode)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.RemindWorkReply{Result: &pb.RemindWorkReply_Data{Data: ""}}
	return okReply, nil
}

func (s *SpaceWorkItemService) OperationPermissions(ctx context.Context, req *pb.OperationPermissionsRequest) (*pb.OperationPermissionsReply, error) {

	reply := func(err error) (*pb.OperationPermissionsReply, error) {
		return &pb.OperationPermissionsReply{Result: &pb.OperationPermissionsReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.Scene, "required,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "Scene"))
	}

	r, err := s.uc.OperationPermissions(ctx, loginUser.UserId, req.SpaceId, req.WorkItemId, req.Scene)
	if err != nil {
		return reply(err)
	}

	if r == nil {
		r = map[string]interface{}{}
	}

	data, _ := json.Marshal(r)
	okReply := &pb.OperationPermissionsReply{Result: &pb.OperationPermissionsReply_Data{Data: string(data)}}
	return okReply, nil
}

func (s *SpaceWorkItemService) SetWorkItemFollower(ctx context.Context, req *pb.SetWorkItemFollowerRequest) (*pb.SetWorkItemFollowerReply, error) {
	reply := func(err error) (*pb.SetWorkItemFollowerReply, error) {
		return &pb.SetWorkItemFollowerReply{Result: &pb.SetWorkItemFollowerReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)

	var err error
	validate := utils.NewValidator()

	if err = validate.Var(req.WorkItemId, "required,number,gt=0"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	err = s.uc.SetFollowers(ctx, loginUser, req.WorkItemId, req.UserIds)
	if err != nil {
		return reply(err)
	}

	return &pb.SetWorkItemFollowerReply{Result: nil}, nil
}

func (s *SpaceWorkItemService) ConfirmStateFlowMain(ctx context.Context, req *pb.ConfirmStateFlowMainRequest) (*pb.ConfirmStateFlowMainReply, error) {

	reply := func(err error) (*pb.ConfirmStateFlowMainReply, error) {
		return &pb.ConfirmStateFlowMainReply{Result: &pb.ConfirmStateFlowMainReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)

	if req.State == "" {
		return reply(errs.Param(ctx, "State"))
	}

	if req.WorkItemId <= 0 {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	lock := locker.Lock(locker.NewWorkItemLockKey(req.WorkItemId))
	lock.Lock()
	defer lock.Unlock()

	err := s.uc.ConfirmStateFlowMain(ctx, loginUser, req.WorkItemId, req.State, req.Reason, req.Remark)
	if err != nil {
		return reply(err)
	}

	return &pb.ConfirmStateFlowMainReply{}, nil
}
