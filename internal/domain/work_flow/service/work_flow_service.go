package service

import (
	"cmp"
	"context"
	"go-cs/internal/consts"
	wf_domain "go-cs/internal/domain/work_flow"
	"go-cs/internal/domain/work_flow/facade"
	wf_repo "go-cs/internal/domain/work_flow/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"
	"go-cs/pkg/stream"
	"math"
	"slices"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type WorkFlowService struct {
	repo      wf_repo.WorkFlowRepo
	log       *log.Helper
	idService *biz_id.BusinessIdService
}

func NewWorkFlowService(
	repo wf_repo.WorkFlowRepo,
	idService *biz_id.BusinessIdService,
	logger log.Logger,
) *WorkFlowService {

	moduleName := "WorkFlowService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkFlowService{
		log:       hlog,
		idService: idService,
		repo:      repo,
	}
}

type GenerateWorkFlowReq struct {
	SpaceId            int64
	WorkItemTypeId     int64
	Ranking            int64
	Uid                int64
	WorkItemStatusInfo *facade.WorkItemStatusInfo
	WorkItemRoleInfo   *facade.WorkItemRoleInfo
}

type GenerateWorkFlowReqResult struct {
	WorkFlow         *wf_domain.WorkFlow
	WorkFlowTemplate *wf_domain.WorkFlowTemplate
}

func (s *WorkFlowService) NewXuQiuWorkFlow(ctx context.Context, req *GenerateWorkFlowReq) *GenerateWorkFlowReqResult {
	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow)
	if bizId == nil {
		return nil
	}

	wfTplt := s.newXuQiuWorkFlowTemplate(ctx, &GenerateWorkFlowTemplateReq{
		SpaceId:            req.SpaceId,
		WorkItemTypeId:     req.WorkItemTypeId,
		WorkFlowId:         bizId.Id,
		WorkItemStatusInfo: req.WorkItemStatusInfo,
		WorkItemRoleInfo:   req.WorkItemRoleInfo,
		UserId:             req.Uid,
	})

	if wfTplt == nil {
		return nil
	}

	wf := wf_domain.NewWorkFlow(bizId.Id, req.SpaceId, req.WorkItemTypeId, "需求", "xuqiu", req.Ranking, consts.FlowMode_WorkFlow, wfTplt, wf_domain.WorkFlowStatus_Enable, 1, req.Uid, nil)
	return &GenerateWorkFlowReqResult{
		WorkFlow:         wf,
		WorkFlowTemplate: wfTplt,
	}
}

func (s *WorkFlowService) NewBugWorkFlow(ctx context.Context, req *GenerateWorkFlowReq) *GenerateWorkFlowReqResult {
	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow)
	if bizId == nil {
		return nil
	}

	wfTplt := s.newBugWorkFlowTemplate(ctx, &GenerateWorkFlowTemplateReq{
		SpaceId:            req.SpaceId,
		WorkItemTypeId:     req.WorkItemTypeId,
		WorkFlowId:         bizId.Id,
		WorkItemStatusInfo: req.WorkItemStatusInfo,
		WorkItemRoleInfo:   req.WorkItemRoleInfo,
		UserId:             req.Uid,
	})

	if wfTplt == nil {
		return nil
	}

	wf := wf_domain.NewWorkFlow(bizId.Id, req.SpaceId, req.WorkItemTypeId, "BUG", "bug", req.Ranking, consts.FlowMode_WorkFlow, wfTplt, wf_domain.WorkFlowStatus_Disable, 1, req.Uid, nil)
	return &GenerateWorkFlowReqResult{
		WorkFlow:         wf,
		WorkFlowTemplate: wfTplt,
	}

}

func (s *WorkFlowService) NewZouChaWorkFlow(ctx context.Context, req *GenerateWorkFlowReq) *GenerateWorkFlowReqResult {

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow)
	if bizId == nil {
		return nil
	}

	wfTplt := s.newZouChaWorkFlowTemplate(ctx, &GenerateWorkFlowTemplateReq{
		SpaceId:            req.SpaceId,
		WorkItemTypeId:     req.WorkItemTypeId,
		WorkFlowId:         bizId.Id,
		WorkItemStatusInfo: req.WorkItemStatusInfo,
		WorkItemRoleInfo:   req.WorkItemRoleInfo,
		UserId:             req.Uid,
	})

	if wfTplt == nil {
		return nil
	}

	wf := wf_domain.NewWorkFlow(bizId.Id, req.SpaceId, req.WorkItemTypeId, "走查", "zoucha", req.Ranking, consts.FlowMode_WorkFlow, wfTplt, wf_domain.WorkFlowStatus_Enable, 1, req.Uid, nil)
	return &GenerateWorkFlowReqResult{
		WorkFlow:         wf,
		WorkFlowTemplate: wfTplt,
	}
}

func (s *WorkFlowService) NewDesignWorkFlow(ctx context.Context, req *GenerateWorkFlowReq) *GenerateWorkFlowReqResult {
	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow)
	if bizId == nil {
		return nil
	}

	wfTplt := s.newDesignWorkFlowTemplate(ctx, &GenerateWorkFlowTemplateReq{
		SpaceId:            req.SpaceId,
		WorkItemTypeId:     req.WorkItemTypeId,
		WorkFlowId:         bizId.Id,
		WorkItemStatusInfo: req.WorkItemStatusInfo,
		WorkItemRoleInfo:   req.WorkItemRoleInfo,
		UserId:             req.Uid,
	})

	if wfTplt == nil {
		return nil
	}

	wf := wf_domain.NewWorkFlow(bizId.Id, req.SpaceId, req.WorkItemTypeId, "设计", "design", req.Ranking, consts.FlowMode_WorkFlow, wfTplt, wf_domain.WorkFlowStatus_Enable, 1, req.Uid, nil)
	return &GenerateWorkFlowReqResult{
		WorkFlow:         wf,
		WorkFlowTemplate: wfTplt,
	}
}

func (s *WorkFlowService) NewSubTaskFlow(ctx context.Context, req *GenerateWorkFlowReq) *GenerateWorkFlowReqResult {
	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow)
	if bizId == nil {
		return nil
	}

	wfTplt := s.newSubTaskStateFlowTemplate(ctx, &GenerateWorkFlowTemplateReq{
		SpaceId:            req.SpaceId,
		WorkItemTypeId:     req.WorkItemTypeId,
		WorkFlowId:         bizId.Id,
		WorkItemStatusInfo: req.WorkItemStatusInfo,
		WorkItemRoleInfo:   req.WorkItemRoleInfo,
		UserId:             req.Uid,
	})

	if wfTplt == nil {
		return nil
	}

	wf := wf_domain.NewWorkFlow(bizId.Id, req.SpaceId, req.WorkItemTypeId, "子任务", "sub_task", req.Ranking, consts.FlowMode_StateFlow, wfTplt, wf_domain.WorkFlowStatus_Enable, 1, req.Uid, nil)
	return &GenerateWorkFlowReqResult{
		WorkFlow:         wf,
		WorkFlowTemplate: wfTplt,
	}
}

func (s *WorkFlowService) NewIssueStateFlow(ctx context.Context, req *GenerateWorkFlowReq) *GenerateWorkFlowReqResult {
	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow)
	if bizId == nil {
		return nil
	}

	wfTplt := s.newIssueStateFlowTemplate(ctx, &GenerateWorkFlowTemplateReq{
		SpaceId:            req.SpaceId,
		WorkItemTypeId:     req.WorkItemTypeId,
		WorkFlowId:         bizId.Id,
		WorkItemStatusInfo: req.WorkItemStatusInfo,
		WorkItemRoleInfo:   req.WorkItemRoleInfo,
		UserId:             req.Uid,
	})

	if wfTplt == nil {
		return nil
	}

	wf := wf_domain.NewWorkFlow(bizId.Id, req.SpaceId, req.WorkItemTypeId, "缺陷", "issue", req.Ranking, consts.FlowMode_StateFlow, wfTplt, wf_domain.WorkFlowStatus_Enable, 1, req.Uid, nil)
	return &GenerateWorkFlowReqResult{
		WorkFlow:         wf,
		WorkFlowTemplate: wfTplt,
	}
}

type GenerateDefaultWorkFlowReq struct {
	SpaceId            int64
	WorkItemTypeId     int64
	Uid                int64
	FlowName           string
	FlowMode           string
	WorkItemStatusInfo *facade.WorkItemStatusInfo
	WorkItemRoleInfo   *facade.WorkItemRoleInfo
	WorkFlowStatus     wf_domain.WorkFlowStatus
}

type GenerateDefaultWorkFlowReqResult struct {
	WorkFlow         *wf_domain.WorkFlow
	WorkFlowTemplate *wf_domain.WorkFlowTemplate
}

// 创建一个包含默认流程模版的 工作流程
func (s *WorkFlowService) NewDefaultWorkFlow(ctx context.Context, req GenerateDefaultWorkFlowReq, oper shared.Oper) (*GenerateDefaultWorkFlowReqResult, error) {

	//检查名称是否重复
	flowName := req.FlowName
	if flowName == "" {
		flowName = "未命名流程"
	}

	//flowName, err := s.generateWorkFlowName(ctx, req.SpaceId, flowName)
	//if err != nil {
	//	return nil, errs.Business(ctx, "检查工作流程名称失败")
	//}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow)
	if bizId == nil {
		return nil, errs.Business(ctx, "生成业务id失败")
	}

	//获取最大排序值
	maxRanking, _ := s.repo.GetMaxRanking(ctx, req.SpaceId)

	var wfTplt *wf_domain.WorkFlowTemplate
	var flowMode consts.WorkFlowMode
	var flowKey string

	switch consts.WorkFlowMode(req.FlowMode) {
	case consts.FlowMode_WorkFlow:
		flowMode = consts.FlowMode_WorkFlow
		flowKey = "flow_st_" + rand.S(5)
		wfTplt = s.newDefaultWorkFlowTemplate(ctx, &GenerateWorkFlowTemplateReq{
			SpaceId:            req.SpaceId,
			WorkItemTypeId:     req.WorkItemTypeId,
			WorkFlowId:         bizId.Id,
			UserId:             req.Uid,
			WorkItemStatusInfo: req.WorkItemStatusInfo,
			WorkItemRoleInfo:   req.WorkItemRoleInfo,
		})

	case consts.FlowMode_StateFlow:
		flowMode = consts.FlowMode_StateFlow
		flowKey = "flow_" + rand.S(5)
		wfTplt = s.newDefaultStateFlowTemplate(ctx, &GenerateWorkFlowTemplateReq{
			SpaceId:            req.SpaceId,
			WorkItemTypeId:     req.WorkItemTypeId,
			WorkFlowId:         bizId.Id,
			UserId:             req.Uid,
			WorkItemStatusInfo: req.WorkItemStatusInfo,
			WorkItemRoleInfo:   req.WorkItemRoleInfo,
		})
	}

	if wfTplt == nil {
		return nil, errs.Business(ctx, "生成默认流程模版失败")
	}

	wf := wf_domain.NewWorkFlow(bizId.Id, req.SpaceId, req.WorkItemTypeId, flowName, flowKey, cast.ToInt64(maxRanking)+100, flowMode, wfTplt, req.WorkFlowStatus, 0, req.Uid, oper)
	return &GenerateDefaultWorkFlowReqResult{
		WorkFlow:         wf,
		WorkFlowTemplate: wfTplt,
	}, nil
}

type CopyWorkFlowResult struct {
	WorkFlow         *wf_domain.WorkFlow
	WorkFlowTemplate *wf_domain.WorkFlowTemplate
}

func (s *WorkFlowService) CopyWorkFlow(ctx context.Context, workFlow *wf_domain.WorkFlow, newFlowName string, status wf_domain.WorkFlowStatus, oper shared.Oper) (*CopyWorkFlowResult, error) {

	flowTplt, err := s.repo.GetFlowTemplate(ctx, workFlow.LastTemplateId)
	if err != nil {
		return nil, errs.Business(ctx, "获取流程模版失败")
	}

	var cpyName = workFlow.Name + "-副本"
	if newFlowName != "" {
		cpyName = newFlowName
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow)
	if bizId == nil {
		return nil, errs.Business(ctx, "生成流程id失败")
	}

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil, errs.Business(ctx, "生成流程模版id失败")
	}

	//获取最大排序值
	maxRanking, _ := s.repo.GetMaxRanking(ctx, workFlow.SpaceId)

	cpyFlowTplt := wf_domain.NewWorkFlowTemplate(
		tpltId.Id,
		workFlow.SpaceId,
		workFlow.WorkItemTypeId,
		bizId.Id,
		1,
		workFlow.FlowMode,
		flowTplt.WorkFlowConf(),
		flowTplt.StateFlowConf(),
		wf_domain.WorkFlowTemplateStatus_Disable,
		oper.GetId(),
		oper,
	)

	cpyFlowKey := "flow_" + rand.S(5)
	cpyFlow := wf_domain.NewWorkFlow(
		bizId.Id,
		workFlow.SpaceId,
		workFlow.WorkItemTypeId,
		cpyName,
		cpyFlowKey,
		cast.ToInt64(maxRanking)+100,
		workFlow.FlowMode,
		cpyFlowTplt,
		status,
		0,
		oper.GetId(),
		oper,
	)

	return &CopyWorkFlowResult{
		WorkFlow:         cpyFlow,
		WorkFlowTemplate: cpyFlowTplt,
	}, nil
}

func (s *WorkFlowService) generateWorkFlowName(ctx context.Context, spaceId int64, flowName string) (string, error) {

	if flowName == "" {
		flowName = "未命名流程"
	}

	names, err := s.repo.GetAllWorkFlowNameBySpaceId(ctx, spaceId)
	if err != nil {
		return "", err
	}

	newFlowName := utils.GenerateName(flowName, names)
	return newFlowName, nil
}

func (s *WorkFlowService) ChangeWorkFlowName(ctx context.Context, workFlow *wf_domain.WorkFlow, newName string, oper shared.Oper) error {
	workFlow.ChangeName(newName, oper)
	return nil
}

func (s *WorkFlowService) GetWorkFlowTemplate(ctx context.Context, templateId int64) (*wf_domain.WorkFlowTemplate, error) {
	return s.repo.GetWorkFlowTemplateFormMemoryCache(ctx, templateId)
}

func (s *WorkFlowService) GetWorkFlow(ctx context.Context, flowId int64) (*wf_domain.WorkFlow, error) {
	return s.repo.GetWorkFlow(ctx, flowId)
}

func (s *WorkFlowService) ResetRank(ctx context.Context, spaceId int64) ([]*wf_domain.WorkFlow, error) {
	// 全量调整排序值
	list, err := s.repo.GetWorkFlowBySpaceId(ctx, spaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	m := stream.GroupBy(list, func(v *wf_domain.WorkFlow) wf_domain.WorkFlowStatus {
		return v.Status
	})

	slices.SortFunc(m[wf_domain.WorkFlowStatus_Enable], func(a, b *wf_domain.WorkFlow) int {
		return cmp.Compare(b.Ranking, a.Ranking)
	})

	slices.SortFunc(m[wf_domain.WorkFlowStatus_Disable], func(a, b *wf_domain.WorkFlow) int {
		return cmp.Compare(b.UpdatedAt, a.UpdatedAt)
	})

	list = append(m[wf_domain.WorkFlowStatus_Enable], m[wf_domain.WorkFlowStatus_Disable]...)
	listLen := len(list)
	for i, v := range list {
		v.UpdateRanking(int64(listLen-i) * 100)
	}

	return list, nil
}

func (s *WorkFlowService) ChangeWorkFlowStatus(ctx context.Context, workFlow *wf_domain.WorkFlow, status wf_domain.WorkFlowStatus, oper shared.Oper) error {

	err := workFlow.ChangeStatus(int64(status), oper)
	if err != nil {
		return err
	}

	if status == wf_domain.WorkFlowStatus_Enable {
		workFlow.UpdateRanking(math.MaxInt) // 启用置顶
	}

	return nil
}
