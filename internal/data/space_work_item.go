package data

import (
	"context"
	"errors"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/data/convert"
	search22 "go-cs/internal/domain/search/search2"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	domain "go-cs/internal/domain/work_item"
	repo "go-cs/internal/domain/work_item/repo"
)

type spaceWorkItemRepo struct {
	baseRepo
	cache bool
	index string
}

func NewSpaceWorkItemRepo(data *Data, logger log.Logger) repo.WorkItemRepo {
	moduleName := "spaceWorkItemRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &spaceWorkItemRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		cache: true,
		index: data.conf.Es.Index,
	}

	return repo
}

func buildWithDocGormFields(opt *repo.WithDocOption) []string {
	var fields []string
	if opt == nil {
		return fields
	}

	if opt.PlanTime || opt.All {
		fields = append(fields, "doc->>'$.plan_start_at' 'plan_start_at'")
		fields = append(fields, "doc->>'$.plan_complete_at' 'plan_complete_at'")
	}
	if opt.ProcessRate || opt.All {
		fields = append(fields, "doc->>'$.process_rate' 'process_rate'")
	}
	if opt.Remark || opt.All {
		fields = append(fields, "doc->>'$.remark' 'remark'")
	}
	if opt.Describe || opt.All {
		fields = append(fields, "doc->>'$.describe' 'describe'")
	}
	if opt.Priority || opt.All {
		fields = append(fields, "doc->>'$.priority' 'priority'")
	}
	if opt.Tags || opt.All {
		fields = append(fields, "doc->>'$.tags' 'tags'")
	}
	if opt.Directors || opt.All {
		fields = append(fields, "doc->>'$.directors' 'directors'")
	}
	if opt.Followers || opt.All {
		fields = append(fields, "doc->>'$.followers' 'followers'")
	}
	if opt.Participators || opt.All {
		fields = append(fields, "doc->>'$.participators' 'participators'")
	}

	return fields
}

func (r *spaceWorkItemRepo) CreateWorkItem(ctx context.Context, workItem *domain.WorkItem) error {
	po := convert.WorkItemEntityToPo(workItem)
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).Create(&po).Error
	return err
}

func (r *spaceWorkItemRepo) CreateWorkItemFiles(ctx context.Context, workItemFiles domain.WorkItemFiles) error {

	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {

		for _, v := range workItemFiles {
			err := tx.Model(&db.SpaceFileInfo{}).Create(convert.WorkItemFileEntityToPo(v)).Error
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

//func (p *spaceWorkItemRepo) GetNotArchiveChildItems(ctx context.Context, workItemPid int64) (*domain.WorkItem, error) {
//
//	var list []*search2.Model
//
//	//更新状态
//	err := p.data.DB(ctx).Model(&search2.Model{}).
//		Where("pid=?", workItemPid).
//		Where("work_item_status not in ?", consts.WorkItemEndStatusList()).
//		Find(&list).Error
//	if err != nil {
//		return nil, err
//	}
//
//	return list, nil
//}
//
//func (p *spaceWorkItemRepo) GetNotCompletedChildItems(ctx context.Context, workItemPid int64) (*domain.WorkItem, error) {
//
//	var list []*search2.Model
//
//	//更新状态
//	err := p.data.DB(ctx).Model(&search2.Model{}).
//		Where("pid=?", workItemPid).
//		Where("work_item_status != ?", consts.WORKITEM_STATUS_COMPLETED).
//		Find(&list).Error
//	if err != nil {
//		return nil, err
//	}
//
//	return list, nil
//}

func (r *spaceWorkItemRepo) GetWorkItem(ctx context.Context, workItemId int64, withDocFiledOpt *repo.WithDocOption, withOption *repo.WithOption) (*domain.WorkItem, error) {

	m := &db.DbSpaceWorkItem{}
	mColumns := (&db.SpaceWorkItemV2{}).Cloumns()

	var selectFields []string
	docFields := buildWithDocGormFields(withDocFiledOpt)
	selectFields = mColumns.SelectEx("doc")
	if len(docFields) > 0 {
		selectFields = append(selectFields, buildWithDocGormFields(withDocFiledOpt)...)
	}

	var row *db.DbSpaceWorkItem
	err := r.data.DB(ctx).Model(m).Select(selectFields).Where("id = ?", workItemId).Take(&row).Error
	if err != nil {
		return nil, err
	}

	entity := convert.WorkItemPoToEntity(row)
	err = r.fillWorkItemByWithOption(ctx, entity, withOption)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *spaceWorkItemRepo) GetWorkItemByIds(ctx context.Context, workItemIds []int64, withDocFiledOpt *repo.WithDocOption, withOption *repo.WithOption) (domain.WorkItems, error) {

	m := &db.DbSpaceWorkItem{}
	columns := (&db.SpaceWorkItemV2{}).Cloumns()

	var selectFields []string
	docFields := buildWithDocGormFields(withDocFiledOpt)
	selectFields = columns.SelectEx("doc")
	if len(docFields) > 0 {
		selectFields = append(selectFields, buildWithDocGormFields(withDocFiledOpt)...)
	}

	var rows []*db.DbSpaceWorkItem
	err := r.data.DB(ctx).Model(m).Select(selectFields).Where("id in ?", workItemIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	entities := convert.WorkItemPoToEntities(rows)
	err = r.batchFillWorkItemByWithOption(ctx, entities, withOption)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *spaceWorkItemRepo) GetWorkItemByPid(ctx context.Context, workItemPid int64, withDocFiledOpt *repo.WithDocOption, withOption *repo.WithOption) (domain.WorkItems, error) {

	m := &db.SpaceWorkItemV2{}
	mColumns := m.Cloumns()

	selectFields := mColumns.SelectEx("doc")
	selectFields = append(selectFields, buildWithDocGormFields(withDocFiledOpt)...)

	var rows []*db.DbSpaceWorkItem
	err := r.data.DB(ctx).Model(m).Select(selectFields).Where("pid = ?", workItemPid).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var list []*domain.WorkItem
	for _, row := range rows {

		workItemEnt := convert.WorkItemPoToEntity(row)

		err = r.fillWorkItemByWithOption(ctx, workItemEnt, withOption)
		if err != nil {
			return nil, err
		}

		list = append(list, workItemEnt)
	}

	return list, nil
}

func (r *spaceWorkItemRepo) GetWorkItemByPids(ctx context.Context, workItemPids []int64, withDocFiledOpt *repo.WithDocOption, withOption *repo.WithOption) (domain.WorkItems, error) {

	m := &db.SpaceWorkItemV2{}
	mColumns := m.Cloumns()

	selectFields := mColumns.SelectEx("doc")
	selectFields = append(selectFields, buildWithDocGormFields(withDocFiledOpt)...)

	var rows []*db.DbSpaceWorkItem
	err := r.data.DB(ctx).
		Select(selectFields).
		Where("pid in ?", workItemPids).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var list []*domain.WorkItem
	for _, row := range rows {

		workItemEnt := convert.WorkItemPoToEntity(row)

		err = r.fillWorkItemByWithOption(ctx, workItemEnt, withOption)
		if err != nil {
			return nil, err
		}

		list = append(list, workItemEnt)
	}

	return list, nil
}

func (r *spaceWorkItemRepo) fillWorkItemByWithOption(ctx context.Context, workItem *domain.WorkItem, withOption *repo.WithOption) error {

	var err error

	if withOption != nil {
		var itemFlowNodes domain.WorkItemFlowNodes
		if withOption.FlowNodes {
			itemFlowNodes, err = r.GetWorkItemFlowNodes(ctx, workItem.Id)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		var itemFlowRoles domain.WorkItemFlowRoles
		if withOption.FlowRoles {
			itemFlowRoles, err = r.GetWorkItemFlowRoles(ctx, workItem.Id)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		var itemFileInfos domain.WorkItemFiles
		if withOption.FileInfos {
			itemFileInfos, err = r.GetWorkItemFiles(ctx, workItem.Id)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		workItem.WorkItemFlowNodes = itemFlowNodes
		workItem.WorkItemFlowRoles = itemFlowRoles
		workItem.WorkItemFiles = itemFileInfos
	}

	return nil
}

func (r *spaceWorkItemRepo) batchFillWorkItemByWithOption(ctx context.Context, items []*domain.WorkItem, withOption *repo.WithOption) error {
	if withOption == nil {
		return nil
	}

	ids := stream.Map(items, func(item *domain.WorkItem) int64 {
		return item.Id
	})

	// 流程节点
	if withOption.FlowNodes {
		m, err := r.GetWorkItemFlowNodesMap(ctx, ids)
		if err != nil {
			return err
		}

		for _, v := range items {
			v.WorkItemFlowNodes = m[v.Id]
		}
	}

	// 流程节点角色
	if withOption.FlowRoles {
		m, err := r.GetWorkItemFlowRolesMap(ctx, ids)
		if err != nil {
			return err
		}

		for _, v := range items {
			v.WorkItemFlowRoles = m[v.Id]
		}
	}

	// 任务附件
	if withOption.FileInfos {
		m, err := r.GetWorkItemFilesMap(ctx, ids)
		if err != nil {
			return err
		}

		for _, v := range items {
			v.WorkItemFiles = m[v.Id]
		}
	}

	return nil
}

func (r *spaceWorkItemRepo) GetWorkItemFiles(ctx context.Context, workItemId int64) (domain.WorkItemFiles, error) {

	m := &db.SpaceFileInfo{}

	var rows []*db.SpaceFileInfo
	err := r.data.DB(ctx).Model(m).Where("source_type = 1 and source_id = ? and status = 1", workItemId).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.WorkItemFilePoToEntities(rows), nil
}

func (r *spaceWorkItemRepo) GetWorkItemFilesMap(ctx context.Context, workItemIds []int64) (map[int64]domain.WorkItemFiles, error) {
	var rows []*db.SpaceFileInfo
	err := r.data.DB(ctx).Model(&db.SpaceFileInfo{}).
		Where("source_type = 1 and source_id = ? and status = 1", workItemIds).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	entities := convert.WorkItemFilePoToEntities(rows)

	m := stream.GroupBy(entities, func(entity *domain.WorkItemFile) int64 {
		return entity.WorkItemId
	})

	return stream.MapValue(m, func(v []*domain.WorkItemFile) domain.WorkItemFiles {
		return v
	}), nil
}

func (r *spaceWorkItemRepo) SaveWorkItem(ctx context.Context, workItem *domain.WorkItem) error {

	diffs := workItem.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceWorkItemV2{}
	mColumns := m.Cloumns()

	mdoc := &db.SpaceWorkItemDocV2{}
	mdocCloums := mdoc.Cloumns()

	mdocJSONSet := datatypes.JSONSet("doc")
	hasDocDiff := false

	cloums := make(map[string]interface{})
	for _, diff := range diffs {
		switch diff {
		case domain.Diff_Name:
			cloums[mColumns.WorkItemName] = workItem.WorkItemName
		case domain.Diff_Restart:
			cloums[mColumns.RestartAt] = workItem.Restart.RestartAt
			cloums[mColumns.IsRestart] = workItem.Restart.IsRestart
			cloums[mColumns.RestartUserId] = workItem.Restart.RestartUserId

		case domain.Diff_IconFlags:
			cloums[mColumns.IconFlags] = workItem.IconFlags
		case domain.Diff_Resume:
			cloums[mColumns.ResumeAt] = workItem.Resume.ResumeAt
		case domain.Diff_CommentNum:
			cloums[mColumns.CommentNum] = workItem.CommentNum
		case domain.Diff_Status:
			cloums[mColumns.WorkItemStatusKey] = workItem.WorkItemStatus.Key
			cloums[mColumns.WorkItemStatusId] = workItem.WorkItemStatus.Id
			cloums[mColumns.WorkItemStatus] = workItem.WorkItemStatus.Val
			cloums[mColumns.LastStatusAt] = workItem.LastWorkItemStatus.LastAt
			cloums[mColumns.LastStatusKey] = workItem.LastWorkItemStatus.Key
			cloums[mColumns.LastStatusId] = workItem.LastWorkItemStatus.Id
			cloums[mColumns.LastStatus] = workItem.LastWorkItemStatus.Val
		case domain.Diff_WorkItemType:
			cloums[mColumns.WorkItemTypeId] = workItem.WorkItemTypeId
			cloums[mColumns.WorkItemTypeKey] = workItem.WorkItemTypeKey
		case domain.Diff_VersionId:
			cloums[mColumns.VersionId] = workItem.VersionId
		case domain.Diff_ObjectId:
			cloums[mColumns.WorkObjectId] = workItem.WorkObjectId

		case domain.Diff_WorkFlowTemplate:
			cloums[mColumns.FlowTemplateId] = workItem.WorkFlowTemplateId
			cloums[mColumns.FlowTemplateVersion] = workItem.WorkFlowTemplateVersion
		case domain.Diff_Reason:
			cloums[mColumns.Reason] = workItem.Reason.ToJSON()
		case domain.Diff_CountAt:
			cloums[mColumns.CountAt] = workItem.CountAt

		// ----- doc json 部分 ----------
		case domain.Diff_PlanTime:
			mdocJSONSet.Set(mdocCloums.PlanStartAt, workItem.Doc.PlanStartAt)
			mdocJSONSet.Set(mdocCloums.PlanCompleteAt, workItem.Doc.PlanCompleteAt)
			hasDocDiff = true
		case domain.Diff_ProcessRate:
			mdocJSONSet.Set(mdocCloums.ProcessRate, workItem.Doc.ProcessRate)
			hasDocDiff = true
		case domain.Diff_Priority:
			mdocJSONSet.Set(mdocCloums.Priority, workItem.Doc.Priority)
			hasDocDiff = true
		case domain.Diff_Remark:
			mdocJSONSet.Set(mdocCloums.Remark, workItem.Doc.Remark)
			hasDocDiff = true
		case domain.Diff_Describe:
			mdocJSONSet.Set(mdocCloums.Describe, workItem.Doc.Describe)
			hasDocDiff = true
		case domain.Diff_Directors:
			mdocJSONSet.Set(mdocCloums.Directors, workItem.Doc.Directors)
			hasDocDiff = true
		case domain.Diff_Followers:
			mdocJSONSet.Set(mdocCloums.Followers, workItem.Doc.Followers)
			hasDocDiff = true
		case domain.Diff_Participators:
			mdocJSONSet.Set(mdocCloums.Participators, workItem.Doc.Participators)
			hasDocDiff = true
		case domain.Diff_NodeDirectors:
			mdocJSONSet.Set(mdocCloums.NodeDirectors, workItem.Doc.NodeDirectors)
			hasDocDiff = true
		case domain.Diff_Tags:
			mdocJSONSet.Set(mdocCloums.Tags, workItem.Doc.Tags)
			hasDocDiff = true
		}
	}

	if hasDocDiff {
		cloums[mColumns.Doc] = mdocJSONSet
	}

	cloums[mColumns.UpdatedAt] = time.Now().Unix()
	err := r.data.DB(ctx).Model(m).Where("id=?", workItem.Id).Updates(cloums).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *spaceWorkItemRepo) ResetChildTaskNum(ctx context.Context, workItemId int64) error {

	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		var totalCount int64
		err := tx.Model(&db.SpaceWorkItemV2{}).
			Where("pid = ? ", workItemId).Count(&totalCount).Error
		if err != nil {
			return err
		}

		err = tx.Model(&db.SpaceWorkItemV2{}).Where("id = ?", workItemId).Update("child_num", totalCount).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func (r *spaceWorkItemRepo) DelSpaceWorkItem(ctx context.Context, workItemIds ...int64) (int64, error) {
	res := r.data.DB(ctx).Where("id in ?", workItemIds).Unscoped().Delete(&db.SpaceWorkItemV2{})
	if res.Error != nil {
		r.log.Error(res.Error)
		return 0, res.Error
	}

	return res.RowsAffected, res.Error
}

func (r *spaceWorkItemRepo) DelSpaceWorkItemByWorkObjectId(ctx context.Context, workObjectId int64) (int64, error) {

	res := r.data.DB(ctx).Where("work_object_id = ?", workObjectId).Delete(&db.SpaceWorkItemV2{})
	err := res.Error
	if err != nil {
		return 0, err
	}

	return res.RowsAffected, nil
}

func (r *spaceWorkItemRepo) GetSpaceWorkItemIdsByParticipators(ctx context.Context, spaceId, userId int64) ([]int64, error) {
	if userId == 0 {
		return nil, nil
	}

	idStr := strconv.FormatInt(userId, 10)

	var list []int64

	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ?", spaceId).
		Where("? MEMBER OF(doc->'$.participators')", idStr).
		Pluck("id", &list).Error

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *spaceWorkItemRepo) AllCommentNumMap(ctx context.Context, spaceId int64) (map[int64]int64, error) {

	type entity struct {
		Id         int64
		CommentNum int64
	}

	m := &db.SpaceWorkItemV2{}
	mColumn := m.Cloumns()

	var list []entity
	err := r.data.RoDB(ctx).Model(m).Select(mColumn.Id, mColumn.CommentNum).
		Where("space_id = ? AND comment_num > 0", spaceId).
		Find(&list).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	commentNumMap := stream.ToMap(list, func(i int, t entity) (int64, int64) {
		return t.Id, t.CommentNum
	})

	return commentNumMap, nil
}

func (r *spaceWorkItemRepo) GetWorkItemIdsByFollower(ctx context.Context, userId int64, spaceId int64) ([]int64, error) {
	if userId == 0 {
		return nil, nil
	}

	idStr := strconv.FormatInt(userId, 10)

	var list []int64

	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ?", spaceId).
		Where("? MEMBER OF(doc->'$.followers')", idStr).
		Pluck("id", &list).Error

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *spaceWorkItemRepo) GetSpaceWorkItemIdsByPid(ctx context.Context, workItemPid int64) ([]int64, error) {
	var list []int64
	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("pid = ?", workItemPid).
		Pluck("id", &list).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *spaceWorkItemRepo) QSpaceWorkItemByPid(ctx context.Context, workItemPid int64) (*rsp.QSpaceWorkItemByPidResult, error) {
	var rows []*db.SpaceWorkItemV2
	err := r.data.RoDB(ctx).Select("id", "space_id", "work_item_name", "user_id", "work_item_status").Model(&rows).
		Where("pid = ?", workItemPid).
		Pluck("id", &rows).Error

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	result := &rsp.QSpaceWorkItemByPidResult{}
	result.List = make([]*rsp.QSpaceWorkItemByPidResult_Item, 0)
	for _, v := range rows {
		result.List = append(result.List, &rsp.QSpaceWorkItemByPidResult_Item{
			Id:           v.Id,
			SpaceId:      v.SpaceId,
			UserId:       v.UserId,
			WorkItemName: v.WorkItemName,
		})
	}

	return result, nil
}

func (r *spaceWorkItemRepo) GetWorkItemByTagV2(ctx context.Context, tagId int64) ([]*search22.Model, error) {
	if tagId == 0 {
		return nil, nil
	}

	tagIdStr := strconv.FormatInt(tagId, 10)

	var list []*search22.Model

	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Select(search22.SelectAll()).
		Where("? MEMBER OF(doc->'$.tags')", tagIdStr).
		Find(&list).Error

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *spaceWorkItemRepo) GetSpaceWorkItemByDirectorV2(ctx context.Context, spaceId, userId int64) ([]*search22.Model, error) {
	if userId == 0 {
		return nil, nil
	}

	idStr := strconv.FormatInt(userId, 10)

	var list []*search22.Model

	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Select(search22.SelectAll()).
		Where("space_id = ?", spaceId).
		Where("? MEMBER OF(doc->'$.directors')", idStr).
		Find(&list).Error

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *spaceWorkItemRepo) UpdateSpaceAllWorkItemCreator(ctx context.Context, spaceId, oldUserId, newUserId int64) (int64, error) {
	if spaceId == 0 || oldUserId == 0 || newUserId == 0 {
		return 0, nil
	}

	//itemIds, _ := p.getWorkItemIdsByCreator(ctx, spaceId, oldUserId)

	res := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ? AND user_id = ?", spaceId, oldUserId).
		Update("user_id", newUserId)
	err := res.Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	//p.clearCache(ctx, itemIds...)
	return res.RowsAffected, nil
}

func (r *spaceWorkItemRepo) CountUserRelatedSpaceWorkItem(ctx context.Context, spaceId, userId int64) (int64, error) {
	var count int64
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ?", spaceId).
		Where("? MEMBER OF(doc->'$.participators')", cast.ToString(userId)).
		Count(&count).Error

	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return count, err
}

func (r *spaceWorkItemRepo) CountUserRelatedSpaceWorkItemBySpaceIds(ctx context.Context, userId int64, spaceIds []int64) (map[int64]int64, error) {
	type res struct {
		SpaceId int64
		Count   int64
	}

	var results []res
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Select("space_id, COUNT(*) as count").
		Where("space_id in ?", spaceIds).
		Where("? MEMBER OF(doc->'$.participators')", cast.ToString(userId)).
		Group("space_id").
		Find(&results).Error

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	count := stream.ToMap(results, func(i int, t res) (int64, int64) {
		return t.SpaceId, t.Count
	})

	return count, err
}

func (r *spaceWorkItemRepo) GetSpaceWorkItemIdsByWorkObject(ctx context.Context, workObjectId int64) ([]int64, error) {
	var list []int64
	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("work_object_id = ?", workObjectId).
		Pluck("id", &list).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *spaceWorkItemRepo) UpdateWorkItemWorkObjectIdByIds(ctx context.Context, workItemIds []int64, newWorkObjectId int64) error {

	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).Where("id in ?", workItemIds).UpdateColumn("work_object_id", newWorkObjectId).Error
	if err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

func (r *spaceWorkItemRepo) RemoveTagFromAllWorkItem(ctx context.Context, spaceId, tagId int64) (int64, error) {
	tagStr := cast.ToString(tagId)

	//itemIds, _ := p.getWorkItemIdsByTag(ctx, tagId)

	res := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ?", spaceId).
		Where("? MEMBER OF(doc->'$.tags')", tagStr).
		Update("doc", datatypes.JSONSet("doc").Set("tags", gorm.Expr("JSON_REMOVE(doc->'$.tags', JSON_UNQUOTE(JSON_SEARCH(doc->'$.tags', 'one', ?)))", tagStr)))

	err := res.Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	//p.clearCache(ctx, itemIds...)
	return res.RowsAffected, nil
}

func (r *spaceWorkItemRepo) ResetVersion(ctx context.Context, oldVersionId int64, newVersionId int64) error {
	model := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{})
	model = model.Where("version_id = ?", oldVersionId).Update("version_id", newVersionId)
	err := model.Error
	if err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

func (r *spaceWorkItemRepo) IncrCommentNum(ctx context.Context, workItemId int64, num int) (int64, error) {
	res := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("id = ?", workItemId).
		Update("comment_num", gorm.Expr("comment_num + ?", num))
	err := res.Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return res.RowsAffected, nil
}

func (r *spaceWorkItemRepo) UpdateNewWorkObject(ctx context.Context, workItemId int64, objectId int64) error {
	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("id = ?", workItemId).
		Or("pid = ?", workItemId).
		Update("work_object_id", objectId).Error
	if err != nil {
		r.log.Error(err)
	}

	// 修改要删缓存
	//workItemIdsByPid := p.getWorkItemIdsByPid(ctx, workItemId)
	//p.clearCache(ctx, append(workItemIdsByPid, workItemId)...)
	return err
}

func (r *spaceWorkItemRepo) CountWorkItemStatusRelatedSpaceWorkItem(ctx context.Context, spaceId, workItemStatusId int64) (int64, error) {

	var count int64

	err := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).Select("id").
		Where("space_id = ? and  work_item_status_id=?", spaceId, workItemStatusId).
		Count(&count).Error

	if err != nil {
		return count, err
	}

	return count, nil
}

func (r *spaceWorkItemRepo) Unfollow(ctx context.Context, userId int64, workItemIds []int64) error {

	if len(workItemIds) == 0 {
		return nil
	}

	s := cast.ToString(userId)

	res := r.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("id IN ?", workItemIds).
		Where("? MEMBER OF(doc->'$.followers')", s).
		Update("doc", datatypes.JSONSet("doc").Set("followers", gorm.Expr("JSON_REMOVE(doc->'$.followers', JSON_UNQUOTE(JSON_SEARCH(doc->'$.followers', 'one', ?)))", s)))

	err := res.Error
	if err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

func (r *spaceWorkItemRepo) GetSpaceAllWorkItemIdsHasComment(ctx context.Context, spaceId int64) ([]int64, error) {
	var list []int64
	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ? AND comment_num > 0", spaceId).
		Pluck("id", &list).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *spaceWorkItemRepo) ReplaceDirectorForWorkItemBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error) {

	idStr := cast.ToString(userId)
	toIdStr := cast.ToString(newUserId)

	var effectRow int64
	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {

		res := tx.Model(&db.SpaceWorkItemV2{}).
			Where("space_id = ?", spaceId).
			Where("? MEMBER OF(doc->'$.directors')", idStr).
			Where("? MEMBER OF(doc->'$.directors')", toIdStr).
			Update("doc", datatypes.JSONSet("doc").Set("directors", gorm.Expr("JSON_REMOVE(doc->'$.directors', JSON_UNQUOTE(JSON_SEARCH(doc->'$.directors', 'one', ?)))", idStr)))

		if res.Error != nil {
			return res.Error
		}

		res = tx.Model(&db.SpaceWorkItemV2{}).
			Where("space_id=?", spaceId).
			Where("? MEMBER OF(doc->>'$.directors')", idStr).
			Not("? MEMBER OF(doc->>'$.directors')", toIdStr).
			Update("doc", datatypes.JSONSet("doc").Set("directors", gorm.Expr("JSON_REPLACE(doc->>'$.directors', JSON_UNQUOTE(JSON_SEARCH(doc->>'$.directors', 'one', ?)), ?)", idStr, toIdStr)))
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

func (r *spaceWorkItemRepo) ReplaceParticipatorsForWorkItemBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error) {

	idStr := cast.ToString(userId)
	toIdStr := cast.ToString(newUserId)

	var effectRow int64
	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {

		res := tx.Model(&db.SpaceWorkItemV2{}).
			Where("space_id = ?", spaceId).
			Where("? MEMBER OF(doc->'$.participators')", idStr).
			Where("? MEMBER OF(doc->'$.participators')", toIdStr).
			Update("doc", datatypes.JSONSet("doc").Set("participators", gorm.Expr("JSON_REMOVE(doc->'$.participators', JSON_UNQUOTE(JSON_SEARCH(doc->'$.participators', 'one', ?)))", idStr)))

		if res.Error != nil {
			return res.Error
		}

		res = tx.Model(&db.SpaceWorkItemV2{}).
			Where("space_id=?", spaceId).
			Where("? MEMBER OF(doc->>'$.participators')", idStr).
			Not("? MEMBER OF(doc->>'$.participators')", toIdStr).
			Update("doc", datatypes.JSONSet("doc").Set("participators", gorm.Expr("JSON_REPLACE(doc->>'$.participators', JSON_UNQUOTE(JSON_SEARCH(doc->>'$.participators', 'one', ?)), ?)", idStr, toIdStr)))
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

func (r *spaceWorkItemRepo) ReplaceNodeDirectorsForWorkItemBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error) {

	idStr := cast.ToString(userId)
	toIdStr := cast.ToString(newUserId)

	var effectRow int64
	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {

		res := tx.Model(&db.SpaceWorkItemV2{}).
			Where("space_id = ?", spaceId).
			Where("? MEMBER OF(doc->'$.node_directors')", idStr).
			Where("? MEMBER OF(doc->'$.node_directors')", toIdStr).
			Update("doc", datatypes.JSONSet("doc").Set("node_directors", gorm.Expr("JSON_REMOVE(doc->'$.node_directors', JSON_UNQUOTE(JSON_SEARCH(doc->'$.node_directors', 'one', ?)))", idStr)))

		if res.Error != nil {
			return res.Error
		}

		res = tx.Model(&db.SpaceWorkItemV2{}).
			Where("space_id=?", spaceId).
			Where("? MEMBER OF(doc->>'$.node_directors')", idStr).
			Not("? MEMBER OF(doc->>'$.node_directors')", toIdStr).
			Update("doc", datatypes.JSONSet("doc").Set("node_directors", gorm.Expr("JSON_REPLACE(doc->>'$.node_directors', JSON_UNQUOTE(JSON_SEARCH(doc->>'$.node_directors', 'one', ?)), ?)", idStr, toIdStr)))
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

func (r *spaceWorkItemRepo) DelSpaceWorkItemBySpaceId(ctx context.Context, spaceId int64) (int64, error) {
	res := r.data.DB(ctx).Where("space_id = ?", spaceId).Unscoped().Delete(&db.SpaceWorkItemV2{})
	if res.Error != nil {
		r.log.Error(res.Error)
		return 0, res.Error
	}

	return res.RowsAffected, res.Error
}

func (r *spaceWorkItemRepo) FilterChildWorkItemIds(ctx context.Context, workItemIds []int64) ([]int64, error) {
	if len(workItemIds) == 0 {
		return nil, nil
	}

	var ids []int64
	res := r.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("id in ? AND pid != 0", workItemIds).
		Pluck("id", &ids)

	if res.Error != nil {
		r.log.Error(res.Error)
		return nil, res.Error
	}

	return ids, nil
}
