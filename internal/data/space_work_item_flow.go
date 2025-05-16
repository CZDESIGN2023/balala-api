package data

import (
	"context"
	"fmt"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	"go-cs/pkg/stream"
	"time"

	domain "go-cs/internal/domain/work_item"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func (r *spaceWorkItemRepo) CreateWorkItemFlowNode(ctx context.Context, workItemFlowNode *domain.WorkItemFlowNode) error {
	po := convert.WorkItemFLowNodeEntityToPo(workItemFlowNode)
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowV2{}).Create(&po).Error
	return err
}

func (r *spaceWorkItemRepo) CreateWorkItemFlowNodes(ctx context.Context, workItemFlowNode ...*domain.WorkItemFlowNode) error {
	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, workItemFlowNode := range workItemFlowNode {
			po := convert.WorkItemFLowNodeEntityToPo(workItemFlowNode)
			err := tx.Model(&db.SpaceWorkItemFlowV2{}).Create(&po).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *spaceWorkItemRepo) SaveWorkItemFlowNode(ctx context.Context, workItemFlowNode *domain.WorkItemFlowNode) error {

	diffs := workItemFlowNode.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceWorkItemFlowV2{}
	mColumns := m.Cloumns()

	columns := make(map[string]interface{})
	for _, diff := range diffs {
		switch diff {
		case domain.Diff_WorkItemFlowNode_Status:
			columns[mColumns.FlowNodeStatus] = workItemFlowNode.FlowNodeStatus
			columns[mColumns.FlowNodePassed] = workItemFlowNode.FlowNodePassed
			columns[mColumns.FlowNodeReached] = workItemFlowNode.FlowNodeReached
			columns[mColumns.StartAt] = workItemFlowNode.StartAt
			columns[mColumns.FinishAt] = workItemFlowNode.FinishAt

		case domain.Diff_WorkItemFlowNode_PlanTime:
			columns[mColumns.PlanStartAt] = workItemFlowNode.PlanTime.StartAt
			columns[mColumns.PlanCompleteAt] = workItemFlowNode.PlanTime.CompleteAt

		case domain.Diff_WorkItemFlowNode_Directors:
			columns[mColumns.Directors] = workItemFlowNode.Directors.ToJsonString()
		case domain.Diff_WorkItemFlowNode_UpdatedAt:
			columns[mColumns.UpdatedAt] = workItemFlowNode.UpdatedAt
		}
	}

	var err error
	if len(columns) == 0 {
		return nil
	}

	if _, ok := columns[mColumns.UpdatedAt]; !ok {
		columns[mColumns.UpdatedAt] = time.Now().Unix()
	}

	err = r.data.DB(ctx).Model(m).Where("id=?", workItemFlowNode.Id).UpdateColumns(columns).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *spaceWorkItemRepo) GetWorkItemFlowNodes(ctx context.Context, workItemId int64) (domain.WorkItemFlowNodes, error) {
	var rows []*db.SpaceWorkItemFlowV2
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowV2{}).Where("work_item_id=?", workItemId).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.WorkItemFlowNodePoToEntities(rows), nil
}

func (r *spaceWorkItemRepo) GetWorkItemFlowNodesMap(ctx context.Context, workItemIds []int64) (map[int64]domain.WorkItemFlowNodes, error) {
	var rows []*db.SpaceWorkItemFlowV2
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowV2{}).Where("work_item_id in ?", workItemIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	entities := convert.WorkItemFlowNodePoToEntities(rows)
	m := stream.GroupBy(entities, func(entity *domain.WorkItemFlowNode) int64 {
		return entity.WorkItemId
	})

	return stream.MapValue(m, func(v []*domain.WorkItemFlowNode) domain.WorkItemFlowNodes {
		return v
	}), nil
}

func (r *spaceWorkItemRepo) CreateWorkItemFlowRole(ctx context.Context, workItemFlowRole *domain.WorkItemFlowRole) error {
	po := convert.WorkItemFLowRoleEntityToPo(workItemFlowRole)
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowRoleV2{}).Create(&po).Error
	return err
}

func (r *spaceWorkItemRepo) CreateWorkItemFlowRoles(ctx context.Context, workItemFlowRoles ...*domain.WorkItemFlowRole) error {
	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, workItemFlowNode := range workItemFlowRoles {
			po := convert.WorkItemFLowRoleEntityToPo(workItemFlowNode)
			err := tx.Model(&db.SpaceWorkItemFlowRoleV2{}).Create(&po).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *spaceWorkItemRepo) SaveWorkItemFlowRole(ctx context.Context, workItemFlowRole *domain.WorkItemFlowRole) error {

	diffs := workItemFlowRole.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceWorkItemFlowRoleV2{}
	mColumns := m.Cloumns()

	columns := make(map[string]interface{})
	for _, diff := range diffs {
		switch diff {
		case domain.Diff_WorkItemFlowRole_Directors:
			columns[mColumns.Directors] = workItemFlowRole.Directors.ToJsonString()
		}
	}

	var err error
	if len(columns) == 0 {
		return nil
	}

	err = r.data.DB(ctx).Model(m).Where("id=?", workItemFlowRole.Id).Updates(columns).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *spaceWorkItemRepo) GetWorkItemFlowRoles(ctx context.Context, workItemId int64) (domain.WorkItemFlowRoles, error) {
	var rows []*db.SpaceWorkItemFlowRoleV2
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowRoleV2{}).Where("work_item_id=?", workItemId).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.WorkItemFLowRolePoToEntities(rows), nil
}

func (r *spaceWorkItemRepo) GetWorkItemFlowRolesMap(ctx context.Context, workItemIds []int64) (map[int64]domain.WorkItemFlowRoles, error) {
	var rows []*db.SpaceWorkItemFlowRoleV2
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowRoleV2{}).Where("work_item_id in ?", workItemIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	entities := convert.WorkItemFLowRolePoToEntities(rows)

	m := stream.GroupBy(entities, func(entity *domain.WorkItemFlowRole) int64 {
		return entity.WorkItemId
	})

	return stream.MapValue(m, func(v []*domain.WorkItemFlowRole) domain.WorkItemFlowRoles {
		return v
	}), nil
}

func (r *spaceWorkItemRepo) DelWorkItemFlowNodeByWorkItemIds(ctx context.Context, workItemIds ...int64) (int64, error) {
	res := r.data.DB(ctx).Unscoped().
		Where("work_item_id in ?", workItemIds).
		Delete(&db.SpaceWorkItemFlowV2{})
	err := res.Error
	if err != nil {
		return 0, err
	}
	return res.RowsAffected, err
}

func (r *spaceWorkItemRepo) DelWorkItemFlowNodeByIds(ctx context.Context, ids ...int64) (int64, error) {
	res := r.data.DB(ctx).Unscoped().
		Where("id in ?", ids).
		Delete(&db.SpaceWorkItemFlowV2{})
	err := res.Error
	if err != nil {
		return 0, err
	}
	return res.RowsAffected, err
}

func (r *spaceWorkItemRepo) DelWorkItemFlowRoleByWorkItemIds(ctx context.Context, workItemIds ...int64) (int64, error) {
	res := r.data.DB(ctx).Unscoped().
		Where("work_item_id in ?", workItemIds).
		Delete(&db.SpaceWorkItemFlowRoleV2{})
	err := res.Error
	if err != nil {
		return 0, err
	}
	return res.RowsAffected, err
}

func (r *spaceWorkItemRepo) DelWorkItemFlowRoleByIds(ctx context.Context, ids ...int64) (int64, error) {
	res := r.data.DB(ctx).Unscoped().
		Where("id in ?", ids).
		Delete(&db.SpaceWorkItemFlowRoleV2{})
	err := res.Error
	if err != nil {
		return 0, err
	}
	return res.RowsAffected, err
}

func (r *spaceWorkItemRepo) GetWorkItemFlowNodeByNodeCode(ctx context.Context, workItemId int64, workFlowNodeCode string) (*domain.WorkItemFlowNode, error) {
	var row *db.SpaceWorkItemFlowV2
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowV2{}).Where("work_item_id=? and flow_node_code=?", workItemId, workFlowNodeCode).First(&row).Error
	if err != nil {
		return nil, err
	}
	return convert.WorkItemFlowNodePoToEntity(row), nil
}

func (r *spaceWorkItemRepo) ProcessingNodeMapByWorkItemIds(ctx context.Context, workItemIds []int64) (map[int64][]*domain.WorkItemFlowNode, error) {
	var row []*db.SpaceWorkItemFlowV2
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowV2{}).Where("work_item_id in ? AND flow_node_status = 2", workItemIds).Find(&row).Error
	if err != nil {
		return nil, err
	}

	m := stream.GroupBy(convert.WorkItemFlowNodePoToEntities(row), func(entity *domain.WorkItemFlowNode) int64 {
		return entity.WorkItemId
	})

	return m, nil
}

func (r *spaceWorkItemRepo) AddDirectorForWorkItemFlows(ctx context.Context, workFlowIds []int64, userId int64) (int64, error) {

	idStr := cast.ToString(userId)

	res := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowV2{}).
		Where("id IN (?)", workFlowIds).
		Not("? MEMBER OF(directors)", idStr).
		Update("directors", gorm.Expr("JSON_ARRAY_APPEND(directors, '$', ?)", idStr))

	err := res.Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return res.RowsAffected, nil
}

func (r *spaceWorkItemRepo) RemoveDirectorForWorkItemFlows(ctx context.Context, workFlowIds []int64, userId int64) (int64, error) {
	if len(workFlowIds) == 0 {
		return 0, nil
	}

	idStr := cast.ToString(userId)

	res := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowV2{}).
		Where("id IN (?)", workFlowIds).
		Where("? MEMBER OF(directors)", idStr).
		Update("directors", gorm.Expr("JSON_REMOVE(directors, JSON_UNQUOTE(JSON_SEARCH(directors, 'one', ?)))", idStr))

	err := res.Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return res.RowsAffected, nil
}

func (r *spaceWorkItemRepo) ReplaceDirectorForWorkItemFlowRolesBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error) {

	idStr := cast.ToString(userId)
	toIdStr := cast.ToString(newUserId)

	var effectRow int64
	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {

		res := tx.Model(&db.SpaceWorkItemFlowRoleV2{}).
			Where("space_id = ?", spaceId).
			Where("? MEMBER OF(directors)", idStr).
			Where("? MEMBER OF(directors)", toIdStr).
			Update("directors", gorm.Expr("JSON_REMOVE(directors, JSON_UNQUOTE(JSON_SEARCH(directors, 'one', ?)))", idStr))

		if res.Error != nil {
			return res.Error
		}

		res = tx.Model(&db.SpaceWorkItemFlowRoleV2{}).
			Where("space_id=?", spaceId).
			Where("? MEMBER OF(directors)", idStr).
			Not("? MEMBER OF(directors)", toIdStr).
			Update("directors", gorm.Expr("JSON_REPLACE(directors, JSON_UNQUOTE(JSON_SEARCH(directors, 'one', ?)), ?)", idStr, toIdStr))
		if res.Error != nil {
			return res.Error
		}

		effectRow = res.RowsAffected
		return nil
	})

	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return effectRow, nil
}

func (r *spaceWorkItemRepo) ReplaceDirectorForWorkItemFlowBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error) {

	idStr := cast.ToString(userId)
	toIdStr := cast.ToString(newUserId)

	var effectRow int64
	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {

		res := tx.Model(&db.SpaceWorkItemFlowV2{}).
			Where("space_id = ?", spaceId).
			Where("? MEMBER OF(directors)", idStr).
			Where("? MEMBER OF(directors)", toIdStr).
			Update("directors", gorm.Expr("JSON_REMOVE(directors, JSON_UNQUOTE(JSON_SEARCH(directors, 'one', ?)))", idStr))

		if res.Error != nil {
			return res.Error
		}

		res = tx.Model(&db.SpaceWorkItemFlowV2{}).
			Where("space_id=?", spaceId).
			Where("? MEMBER OF(directors)", idStr).
			Not("? MEMBER OF(directors)", toIdStr).
			Update("directors", gorm.Expr("JSON_REPLACE(directors, JSON_UNQUOTE(JSON_SEARCH(directors, 'one', ?)), ?)", idStr, toIdStr))
		if res.Error != nil {
			return res.Error
		}

		effectRow = res.RowsAffected
		return nil
	})

	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return effectRow, nil
}

func (r *spaceWorkItemRepo) ReplaceDirectorForWorkItemById(ctx context.Context, workItemIds []int64, userId int64, newUserId int64) (int64, error) {

	idStr := cast.ToString(userId)
	toIdStr := cast.ToString(newUserId)

	res := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("id IN (?)", workItemIds).
		Where("? MEMBER OF(directors)", idStr).
		Update("directors", gorm.Expr("IF(JSON_SEARCH(directors, 'one', ?) is NULL ,directors ,JSON_REPLACE(directors, JSON_UNQUOTE(JSON_SEARCH(directors, 'one', ?)), ?)) ", idStr, toIdStr))

	err := res.Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return res.RowsAffected, nil
}

func (r *spaceWorkItemRepo) AddDirectorForWorkItemFlow(ctx context.Context, workFlowId, userId int64) (int64, error) {
	return r.AddDirectorForWorkItemFlows(ctx, []int64{workFlowId}, userId)
}

func (r *spaceWorkItemRepo) RemoveDirectorForWorkItemFlow(ctx context.Context, workFlowId, userId int64) (int64, error) {
	return r.RemoveDirectorForWorkItemFlows(ctx, []int64{workFlowId}, userId)
}

func (r *spaceWorkItemRepo) GetSpaceWorkItemFlowIdsBySpaceUserId(ctx context.Context, spaceId int64, userId int64) ([]int64, error) {
	if userId == 0 {
		return nil, nil
	}

	idStr := cast.ToString(userId)

	var list []int64

	err := r.data.DB(ctx).Model(&db.SpaceWorkItemFlowV2{}).
		Where("space_id = ?", spaceId).
		Where("? MEMBER OF(directors)", idStr).
		Pluck("id", &list).Error

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *spaceWorkItemRepo) GetSpaceWorkItemIdsForUpgradeFlow(ctx context.Context, spaceId int64, flowId int64, lastVersion int64) ([]int64, error) {

	var list []int64

	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Select("id").
		Where("space_id = ? and  flow_id=? and  flow_template_version <> ? ", spaceId, flowId, lastVersion).
		Pluck("id", &list).Error

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *spaceWorkItemRepo) CountWorkFlowRelatedSpaceWorkItem(ctx context.Context, spaceId, workFlowId int64, excludeStatusKeys []string) (int64, error) {

	var count int64

	model := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ? and  work_item_flow_id=?", spaceId, workFlowId)
	if len(excludeStatusKeys) > 0 {
		model = model.Where("work_item_status_key NOT IN (?)", excludeStatusKeys)
	}

	err := model.
		Count(&count).Error

	if err != nil {
		return count, err
	}

	return count, nil
}

func (r *spaceWorkItemRepo) CountWorkFlowRoleRelatedSpaceWorkItem(ctx context.Context, spaceId, flowRoleId int64) (int64, error) {
	var count int64

	m := &db.SpaceWorkItemFlowRoleV2{}

	sql := fmt.Sprintf(` SELECT COUNT(*) FROM ( 
		select DISTINCT work_item_id from %v where space_id = ? and  work_item_role_id=?
	 )   AS subquery `, m.TableName())

	err := r.data.DB(ctx).Raw(sql, spaceId, flowRoleId).Count(&count).Error
	if err != nil {
		return count, err
	}

	return count, nil
}

func (r *spaceWorkItemRepo) DelSpaceWorkItemFlowBySpaceId(ctx context.Context, spaceId int64) (int64, error) {
	res := r.data.DB(ctx).Where("space_id = ?", spaceId).Unscoped().Delete(&db.SpaceWorkItemFlowV2{})
	if res.Error != nil {
		r.log.Error(res.Error)
		return 0, res.Error
	}
	return res.RowsAffected, res.Error
}

func (r *spaceWorkItemRepo) DelSpaceWorkItemFlowRoleBySpaceId(ctx context.Context, spaceId int64) (int64, error) {
	res := r.data.DB(ctx).Where("space_id = ?", spaceId).Unscoped().Delete(&db.SpaceWorkItemFlowRoleV2{})
	if res.Error != nil {
		r.log.Error(res.Error)
		return 0, res.Error
	}
	return res.RowsAffected, res.Error
}

func (r *spaceWorkItemRepo) HasWorkItemRelateFlow(ctx context.Context, spaceId, workFlowId int64) (bool, error) {

	var id int64
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ? AND work_item_flow_id=?", spaceId, workFlowId).
		Limit(1).
		Pluck("id", &id).Error
	if err != nil {
		return false, err
	}

	return id > 0, nil
}
