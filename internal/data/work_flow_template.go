package data

import (
	"context"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	"go-cs/pkg/stream"
	"time"

	domain "go-cs/internal/domain/work_flow"

	"gorm.io/gorm"
)

func (r *workFlowRepo) CreateWorkFlowTemplate(ctx context.Context, template *domain.WorkFlowTemplate) error {
	po := convert.WorkFlowTemplateEntityToPo(template)
	err := r.data.db.Model(&db.WorkFlowTemplate{}).Create(po).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *workFlowRepo) CreateWorkFlowTemplates(ctx context.Context, templates []*domain.WorkFlowTemplate) error {

	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, template := range templates {
			po := convert.WorkFlowTemplateEntityToPo(template)
			err := r.data.db.Model(&db.WorkFlowTemplate{}).Create(po).Error
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

func (r *workFlowRepo) SaveWorkFlowTemplate(ctx context.Context, template *domain.WorkFlowTemplate) error {
	po := convert.WorkFlowTemplateEntityToPo(template)

	diffs := template.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.WorkFlowTemplate{}
	mColumns := m.Cloumns()

	columns := make(map[string]any)
	for _, v := range diffs {
		switch v {
		case domain.Diff_TemplateSetting:
			columns[mColumns.Setting] = po.Setting
		}
	}

	if len(columns) == 0 {
		return nil
	}

	columns[mColumns.UpdatedAt] = time.Now().Unix()
	err := r.data.DB(ctx).Model(m).Where("id=?", po.Id).UpdateColumns(columns).Error
	if err != nil {
		return err
	}

	// 清除缓存
	r.flowTpltCacheAccessObj.Delete(template.Id)

	return nil
}

func (r *workFlowRepo) DelWorkFlowTemplateBySpaceId(ctx context.Context, spaceId int64) error {
	var opValue = make(map[string]interface{})
	err := r.data.DB(ctx).Model(&db.WorkFlowTemplate{}).Unscoped().Where("space_id=?", spaceId).Delete(&opValue).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *workFlowRepo) GetFlowTemplate(ctx context.Context, id int64) (*domain.WorkFlowTemplate, error) {

	var row *db.WorkFlowTemplate
	err := r.data.DB(ctx).Model(&db.WorkFlowTemplate{}).Where("id=?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}

	ent := convert.WorkFlowTemplatePoToEntity(row)
	return ent, nil
}

func (r *workFlowRepo) GetWorkFlowTemplates(ctx context.Context, ids []*int64) ([]*domain.WorkFlowTemplate, error) {
	var rows []*db.WorkFlowTemplate
	err := r.data.DB(ctx).Model(&db.WorkFlowTemplate{}).Where("id in ?", ids).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.WorkFlowTemplatePoToEntities(rows), nil
}

func (r *workFlowRepo) GetFlowTemplateByIds(ctx context.Context, ids []int64) ([]*domain.WorkFlowTemplate, error) {

	var list []*db.WorkFlowTemplate
	err := r.data.DB(ctx).Model(&db.WorkFlowTemplate{}).Where("id in ?", ids).Find(&list).Error
	if err != nil {
		return nil, err
	}

	ent := convert.WorkFlowTemplatePoToEntities(list)
	return ent, nil
}

func (r *workFlowRepo) FlowTemplateMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkFlowTemplate, error) {
	list, err := r.GetFlowTemplateByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	return stream.ToMap(list, func(_ int, item *domain.WorkFlowTemplate) (int64, *domain.WorkFlowTemplate) {
		return item.Id, item
	}), nil
}

func (r *workFlowRepo) DelWorkFlowTemplateByFlowId(ctx context.Context, flowId int64) error {
	var opValue = make(map[string]any)
	err := r.data.DB(ctx).Model(&db.WorkFlowTemplate{}).Unscoped().Where("work_flow_id=?", flowId).Delete(&opValue).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *workFlowRepo) DelWorkFlowTemplate(ctx context.Context, id int64) error {
	var opValue = make(map[string]any)
	err := r.data.DB(ctx).Model(&db.WorkFlowTemplate{}).Unscoped().Where("id=?", id).Delete(&opValue).Error
	if err != nil {
		return err
	}
	return nil
}
