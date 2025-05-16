package data

import (
	"context"
	"errors"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/query"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/local_cache"
	"go-cs/pkg/stream"
	"time"

	domain "go-cs/internal/domain/work_flow"
	repo "go-cs/internal/domain/work_flow/repo"

	goCache "github.com/Code-Hex/go-generics-cache"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type workFlowRepo struct {
	baseRepo
	cacheConfig            *local_cache.Cache[string, *db.Config] // 配置信息缓存
	flowTpltCacheAccessObj *goCache.Cache[int64, *domain.WorkFlowTemplate]
	flowCache              *goCache.Cache[string, any]
}

func NewWorkFlowRepo(data *Data, logger log.Logger) repo.WorkFlowRepo {
	moduleName := "workFlowRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)
	r := &workFlowRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		cacheConfig:            local_cache.NewCache[string, *db.Config](-1),
		flowTpltCacheAccessObj: goCache.New(goCache.AsFIFO[int64, *domain.WorkFlowTemplate]()),
		flowCache:              goCache.New(goCache.AsFIFO[string, any]()),
	}
	return r
}

func (r *workFlowRepo) CreateWorkFlow(ctx context.Context, flow *domain.WorkFlow) error {
	po := convert.WorkFlowEntityToPo(flow)
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Create(po).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *workFlowRepo) CreateWorkFlows(ctx context.Context, flows []*domain.WorkFlow) error {

	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, flow := range flows {
			po := convert.WorkFlowEntityToPo(flow)
			err := r.data.DB(ctx).Model(&db.WorkFlow{}).Create(po).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *workFlowRepo) SaveWorkFlow(ctx context.Context, workFlow *domain.WorkFlow) error {

	diffs := workFlow.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.WorkFlow{}
	mColumns := m.Cloumns()

	columns := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_Ranking:
			columns[mColumns.Ranking] = workFlow.Ranking
		case domain.Diff_Name:
			columns[mColumns.Name] = workFlow.Name
		case domain.Diff_Status:
			columns[mColumns.Status] = workFlow.Status
		case domain.Diff_Version:
			columns[mColumns.Version] = workFlow.Version
		case domain.Diff_LastTemplate:
			columns[mColumns.LastTemplateId] = workFlow.LastTemplateId
		}
	}

	if len(columns) == 0 {
		return nil
	}

	columns[mColumns.UpdatedAt] = time.Now().Unix()
	err := r.data.DB(ctx).Model(m).Where("id=?", workFlow.Id).UpdateColumns(columns).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *workFlowRepo) DelWorkFlowBySpaceId(ctx context.Context, spaceId int64) error {
	var opValue = make(map[string]interface{})
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Unscoped().Where("space_id=?", spaceId).Delete(&opValue).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *workFlowRepo) DelWorkFlow(ctx context.Context, id int64) error {
	var opValue = make(map[string]interface{})
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Unscoped().Where("id=?", id).Delete(&opValue).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *workFlowRepo) GetWorkFlow(ctx context.Context, id int64) (*domain.WorkFlow, error) {
	var row *db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("id=?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	ent := convert.WorkFlowPoToEntity(row)
	return ent, nil
}

func (r *workFlowRepo) GetWorkFlowBySpaceId(ctx context.Context, spaceId int64) ([]*domain.WorkFlow, error) {
	var rows []*db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("space_id=?", spaceId).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return convert.WorkFlowPoToEntities(rows), nil
}

func (r *workFlowRepo) GetWorkFlowBySpaceIds(ctx context.Context, spaceIds []int64) ([]*domain.WorkFlow, error) {
	var rows []*db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("space_id in ?", spaceIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return convert.WorkFlowPoToEntities(rows), nil
}

func (r *workFlowRepo) GetWorkFlowBySpaceWorkItemTypeId(ctx context.Context, spaceId int64, workItemTypeId int64) ([]*domain.WorkFlow, error) {
	var rows []*db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("space_id=? and work_item_type_id=?", spaceId, workItemTypeId).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.WorkFlowPoToEntities(rows), nil
}

func (r *workFlowRepo) GetWorkFlowMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkFlow, error) {
	var list []*db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("id in ", ids).Find(&list).Error
	if err != nil {
		return nil, err
	}
	ent := convert.WorkFlowPoToEntities(list)

	m := stream.ToMap(ent, func(_ int, v *domain.WorkFlow) (int64, *domain.WorkFlow) {
		return v.Id, v
	})
	return m, nil
}

func (r *workFlowRepo) QTaskWorkFlowViewList(ctx context.Context, req *query.TaskWorkFlowListQuery) (*query.TaskWorkFlowListQueryResult, error) {

	//查询工作流
	var rows []*db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).
		Where("space_id=? AND work_item_type_id IN ?", req.SpaceId, req.WorkItemTypeIds).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var tpltIds []int64
	for _, v := range rows {
		tpltIds = append(tpltIds, v.LastTemplateId)
	}

	//查询工作流对应的模版
	var tpltRows []*db.WorkFlowTemplate
	err = r.data.DB(ctx).Model(&db.WorkFlowTemplate{}).Where("id in ?", tpltIds).Find(&tpltRows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	flows := convert.WorkFlowPoToEntities(rows)
	flowTplts := convert.WorkFlowTemplatePoToEntities(tpltRows)
	flowTpltsMap := stream.ToMap(flowTplts, func(idx int, flowTplt *domain.WorkFlowTemplate) (int64, *domain.WorkFlowTemplate) {
		return flowTplt.Id, flowTplt
	})

	result := &query.TaskWorkFlowListQueryResult{}
	result.List = make([]*query.TaskWorkFlowListQueryResult_Item, 0)
	for _, flow := range flows {
		flowTplt := flowTpltsMap[flow.LastTemplateId]
		if flowTplt == nil {
			continue
		}

		item := &query.TaskWorkFlowListQueryResult_Item{
			WorkFlow:         flow,
			WorkFlowTemplate: flowTplt,
		}

		result.List = append(result.List, item)
	}

	return result, nil
}

func (r *workFlowRepo) QTaskWorkFlowList(ctx context.Context, spaceId int64, workItemTypeIds []int64, FlowMode consts.WorkFlowMode) ([]*domain.WorkFlow, error) {

	//查询工作流
	var rows []*db.WorkFlow
	tx := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("space_id=? AND work_item_type_id IN ?", spaceId, workItemTypeIds)
	if FlowMode != "" {
		tx = tx.Where("flow_mode=?", FlowMode)
	}

	err := tx.Find(&rows).Order("ranking desc, id desc").Error
	if err != nil {
		return nil, err
	}

	return convert.WorkFlowPoToEntities(rows), nil
}

func (r *workFlowRepo) QTaskWorkFlowById(ctx context.Context, spaceId int64, workFlowIds []int64) ([]*domain.WorkFlow, error) {

	//查询工作流
	var rows []*db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).
		Where("id in ? and space_id=? and status=1 and flow_mode='work_flow'", spaceId, workFlowIds).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.WorkFlowPoToEntities(rows), nil
}

func (r *workFlowRepo) GetAllWorkFlowNameBySpaceId(ctx context.Context, spaceId int64) ([]string, error) {
	rows := make([]map[string]interface{}, 0)
	err := r.data.RoDB(ctx).Model(&db.WorkFlow{}).Select("name").Where("space_id = ?", spaceId).Find(&rows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var names []string
	for _, v := range rows {
		names = append(names, cast.ToString(v["name"]))
	}

	return names, nil
}

func (r *workFlowRepo) QWorkFlowInfo(ctx context.Context, req *query.WorkFlowInfoQuery) (*query.WorkFlowInfoQueryResult, error) {
	var row *db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("id=?", req.FlowId).Take(&row).Error
	if err != nil {
		return nil, err
	}

	var tpltRow *db.WorkFlowTemplate
	err = r.data.DB(ctx).Model(&db.WorkFlowTemplate{}).Where("id = ?", row.LastTemplateId).Take(&tpltRow).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	result := &query.WorkFlowInfoQueryResult{
		WorkFlow:         convert.WorkFlowPoToEntity(row),
		WorkFlowTemplate: convert.WorkFlowTemplatePoToEntity(tpltRow),
	}

	return result, nil
}

func (r *workFlowRepo) IsExistByWorkFlowName(ctx context.Context, spaceId int64, workFlowName string) (bool, error) {
	var rowNum int64
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("space_id=? and BINARY name=?", spaceId, workFlowName).Count(&rowNum).Error
	if err != nil {
		return false, err
	}

	if rowNum > 0 {
		return true, nil
	}
	return false, nil
}

func (r *workFlowRepo) GetWorkFlowByIds(ctx context.Context, ids []int64) ([]*domain.WorkFlow, error) {
	var list []*db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("id in ? ", ids).Find(&list).Error
	if err != nil {
		return nil, err
	}

	ret := convert.WorkFlowPoToEntities(list)
	return ret, err
}

func (r *workFlowRepo) WorkFlowMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkFlow, error) {
	list, err := r.GetWorkFlowByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	return stream.ToMap(list, func(i int, t *domain.WorkFlow) (int64, *domain.WorkFlow) {
		return t.Id, t
	}), nil
}

func (r *workFlowRepo) SearchHistoryTaskWorkFlowTemplateByOwnerRule(ctx context.Context, spaceId int64, owner string) ([]*domain.WorkFlowTemplate, error) {
	//查询工作流模板
	var tpltRows []*db.WorkFlowTemplate
	err := r.data.RoDB(ctx).Model(&db.WorkFlowTemplate{}).
		Where("space_id = ?", spaceId).
		Where("JSON_SEARCH(setting, 'one' , ?, null, '$.nodes[*].owner.value.*[*].value') is not null OR JSON_SEARCH(setting, 'one' , ?, null, '$.stateFlowNodes[*].owner.value.*[*].value') is not null", owner, owner).
		Find(&tpltRows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return convert.WorkFlowTemplatePoToEntities(tpltRows), nil
}

func (r *workFlowRepo) SearchTaskWorkFlowLastTemplateByOwnerRule(ctx context.Context, spaceId int64, owner string) ([]int64, error) {
	//查询工作流
	var rows []*db.WorkFlow
	err := r.data.RoDB(ctx).Model(&db.WorkFlow{}).Select("id,last_template_id").
		Where("space_id=?", spaceId).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var tpltIds []int64
	for _, row := range rows {
		tpltIds = append(tpltIds, row.LastTemplateId)
	}

	//JSON_SEARCH(setting, 'one' , '42', null, '$.nodes[*].owner.value.*[*].value') is not null

	//查询工作流模板
	var tpltRows []*db.WorkFlowTemplate
	err = r.data.RoDB(ctx).Model(&db.WorkFlowTemplate{}).Select("id").
		Where("id in ?", tpltIds).
		Where("JSON_SEARCH(setting, 'one' , ?, null, '$.nodes[*].owner.value.*[*].value') is not null", owner).
		Find(&tpltRows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var ids []int64
	for _, row := range tpltRows {
		ids = append(ids, row.Id)
	}

	return ids, nil
}

func (r *workFlowRepo) SearchTaskWorkFlowTemplateByNodeStateEvent(ctx context.Context, spaceId int64, subStateId string) ([]int64, error) {
	//查询工作流
	var rows []*db.WorkFlow
	err := r.data.RoDB(ctx).Model(&db.WorkFlow{}).Select("id,last_template_id").
		Where("space_id=?", spaceId).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var tpltIds []int64
	for _, row := range rows {
		tpltIds = append(tpltIds, row.LastTemplateId)
	}

	//JSON_SEARCH(setting, 'one' , '42', null, '$.nodes[*].owner.value.*[*].value') is not null

	//查询工作流模板
	var tpltRows []*db.WorkFlowTemplate
	err = r.data.RoDB(ctx).Model(&db.WorkFlowTemplate{}).Select("id").
		Where("id in ?", tpltIds).
		Where(`JSON_SEARCH(setting, 'one' , ?, null, '$.nodes[*].*[*].targetSubState.id') is not null OR JSON_SEARCH(setting, 'one' , ?, null, '$.stateFlowNodes[*].subStateId') is not null`, subStateId, subStateId).
		Find(&tpltRows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var ids []int64
	for _, row := range tpltRows {
		ids = append(ids, row.Id)
	}

	return ids, nil
}

func (r *workFlowRepo) SearchTaskWorkFlowTemplateByOwnerRoleRule(ctx context.Context, spaceId int64, roleId string) ([]int64, error) {
	//查询工作流
	var rows []*db.WorkFlow
	err := r.data.RoDB(ctx).Model(&db.WorkFlow{}).Select("id,last_template_id").
		Where("space_id=?", spaceId).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var tpltIds []int64
	for _, row := range rows {
		tpltIds = append(tpltIds, row.LastTemplateId)
	}

	//JSON_SEARCH(setting, 'one' , '42', null, '$.nodes[*].owner.value.*[*].value') is not null

	//查询工作流模板
	var tpltRows []*db.WorkFlowTemplate
	err = r.data.RoDB(ctx).Model(&db.WorkFlowTemplate{}).Select("id").
		Where("id in ?", tpltIds).
		Where("JSON_SEARCH(setting, 'one' , ?, null, '$.nodes[*].owner.ownerRole[*].id') is not null OR JSON_SEARCH(setting, 'one' , ?, null, '$.stateFlowNodes[*].owner.ownerRole[*].id') is not null", roleId, roleId).
		Find(&tpltRows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var ids []int64
	for _, row := range tpltRows {
		ids = append(ids, row.Id)
	}

	return ids, nil
}

func (r *workFlowRepo) GetMaxRanking(ctx context.Context, spaceId int64) (int64, error) {
	var max int64
	row := r.data.DB(ctx).Model(&db.WorkFlow{}).Select("MAX(ranking)").Where("space_id=?", spaceId).Row()
	err := row.Scan(&max)
	if err != nil {
		return max, err
	}
	return max, nil
}

func (r *workFlowRepo) ClearHistoryTemplate(ctx context.Context, flowId int64) error {
	var flow db.WorkFlow
	err := r.data.DB(ctx).Model(&db.WorkFlow{}).Where("id=?", flowId).Take(&flow).Error
	if err != nil {
		return err
	}

	var templateIds []int64
	err = r.data.DB(ctx).Model(&db.WorkFlowTemplate{}).Where("work_flow_id=? AND id != ?", flowId, flow.LastTemplateId).Pluck("id", &templateIds).Error
	if err != nil {
		return err
	}

	var needClearTemplateIds []int64
	for _, templateId := range templateIds {
		var spaceWorkItem db.SpaceWorkItemV2
		err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).Where("flow_template_id = ?", templateId).Take(&spaceWorkItem).Error
		if spaceWorkItem.Id == 0 && errs.IsDbRecordNotFoundErr(err) {
			needClearTemplateIds = append(needClearTemplateIds, templateId)
		}
	}

	if len(needClearTemplateIds) > 0 {
		err := r.data.DB(ctx).Where("id in ?", needClearTemplateIds).Delete(&db.WorkFlowTemplate{}).Error
		if err != nil {
			return err
		}
	}

	return nil
}
