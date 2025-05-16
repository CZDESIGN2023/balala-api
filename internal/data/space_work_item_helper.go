package data

//func (p *spaceWorkItemRepo) clearCache(ctx context.Context, workItemId ...int64) {
//	keys := stream.Map(workItemId, func(v int64) string {
//		return NewWorkItemKey(v).Key()
//	})
//
//	cmdRes := p.data.rdb.Del(ctx, keys...)
//	if r, err := cmdRes.Result(); err != nil {
//		p.log.Error(r, err)
//	}
//}

// func (p *spaceWorkItemRepo) clearCache2(workItemIds []int64) {
// 	keys := stream.Map(workItemIds, func(v int64) string {
// 		return NewWorkItemKey(v).Key()
// 	})

// 	cmdRes := p.data.rdb.Del(context.Background(), keys...)
// 	if r, err := cmdRes.Result(); err != nil {
// 		p.log.Error(r, err)
// 	}
// }

// func getFieldValue(field string, value any) (string, any) {
// 	model := search2.FieldByColumn(field)

// 	if model.IsJSONField() {
// 		subField := field
// 		field = "doc"
// 		value = datatypes.JSONSet("doc").Set(subField, value)
// 	}

// 	return field, value
// }

// func (p *spaceWorkItemRepo) updateWorkItemFieldByPidV2(ctx context.Context, workItemPid int64, field string, value any) error {

// 	m := &db.SpaceWorkItemV2{}
// 	mColumns := m.Cloumns()

// 	if field == mColumns.IconFlags {
// 		return errors.New("icon_flags 不能通过此方法修改")
// 	}

// 	field, value = getFieldValue(field, value)

// 	err := p.data.DB(ctx).Model(m).Where("pid = ?", workItemPid).Update(field, value).Error
// 	if err != nil {
// 		p.log.Error(err)
// 		return err
// 	}

// 	return nil
// }

// func (p *spaceWorkItemRepo) updateWorkItemsFieldV2(ctx context.Context, workItemIds []int64, field string, value any) error {

// 	m := &db.SpaceWorkItemV2{}
// 	mColumns := m.Cloumns()

// 	if field == mColumns.IconFlags {
// 		return errors.New("icon_flags 不能通过此方法修改")
// 	}

// 	field, value = getFieldValue(field, value)

// 	err := p.data.DB(ctx).Model(m).Where("id in ?", workItemIds).Update(field, value).Error
// 	if err != nil {
// 		p.log.Error(err)
// 		return err
// 	}

// 	return nil
// }

// func (p *spaceWorkItemRepo) getWorkItemIdsByCreator(ctx context.Context, spaceId int64, creator int64) ([]int64, error) {
// 	var ids []int64
// 	err := p.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).Where("space_id = ? AND user_id = ?", spaceId, creator).Pluck("id", &ids).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ids, nil
// }

// func (p *spaceWorkItemRepo) getWorkItemIdsByPid(ctx context.Context, pid int64) []int64 {
// 	if pid == 0 {
// 		return nil
// 	}
// 	var ids []int64
// 	p.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).Where("pid = ?", pid).Pluck("id", &ids)
// 	return ids
// }

// func (p *spaceWorkItemRepo) getWorkItemIdsByTag(ctx context.Context, tagId int64) ([]int64, error) {
// 	if tagId == 0 {
// 		return nil, nil
// 	}

// 	tagIdStr := strconv.FormatInt(tagId, 10)

// 	var list []int64

// 	err := p.data.DB(ctx).Model(&db.SpaceWorkItemV2{}).
// 		Where("? MEMBER OF(doc->'$.tags')", tagIdStr).
// 		Pluck("", &list).Error

// 	if err != nil {
// 		p.log.Error(err)
// 		return nil, err
// 	}

// 	return list, nil
// }
