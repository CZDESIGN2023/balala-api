package facade

import (
	"context"
	"go-cs/internal/domain/work_flow"
	flow_service "go-cs/internal/domain/work_flow/service"
)

type WorkFlowTemplateFacade struct {
	template *work_flow.WorkFlowTemplate
}

func (w *WorkFlowTemplateFacade) Template() *work_flow.WorkFlowTemplate {
	return w.template
}

func BuildWorkFlowTemplateFacade(template *work_flow.WorkFlowTemplate) *WorkFlowTemplateFacade {
	return &WorkFlowTemplateFacade{
		template: template,
	}
}

type WorkFlowTemplateServiceFacade struct {
	flowService *flow_service.WorkFlowService
}

type GetWorkFlowTemplateResult struct {
	Template *work_flow.WorkFlowTemplate
}

func (w *WorkFlowTemplateServiceFacade) GetWorkFlowTemplate(ctx context.Context, id int64) (*GetWorkFlowTemplateResult, error) {
	// 获取模板
	template, err := w.flowService.GetWorkFlowTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetWorkFlowTemplateResult{
		Template: template,
	}, nil
}

func (w *WorkFlowTemplateServiceFacade) GetWorkFlow(ctx context.Context, id int64) (*work_flow.WorkFlow, error) {
	// 获取模板
	flow, err := w.flowService.GetWorkFlow(ctx, id)
	if err != nil {
		return nil, err
	}

	return flow, nil
}

func BuildWorkFlowTemplateServiceFacade(flowService *flow_service.WorkFlowService) *WorkFlowTemplateServiceFacade {
	return &WorkFlowTemplateServiceFacade{
		flowService: flowService,
	}
}
