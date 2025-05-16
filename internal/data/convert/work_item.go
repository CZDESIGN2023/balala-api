package convert

import (
	"encoding/json"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_item"
)

func WorkItemEntityToEs(workItem *domain.WorkItem) map[string]any {

	mColumns := (&db.SpaceWorkItemV2{}).Cloumns()

	esDoc := make(map[string]any)
	esDoc[mColumns.Id] = workItem.Id
	esDoc[mColumns.Pid] = workItem.Pid
	esDoc[mColumns.SpaceId] = workItem.SpaceId
	esDoc[mColumns.UserId] = workItem.UserId
	esDoc[mColumns.WorkObjectId] = workItem.WorkObjectId
	esDoc[mColumns.WorkItemGuid] = workItem.WorkItemGuid
	esDoc[mColumns.WorkItemName] = workItem.WorkItemName
	esDoc[mColumns.IsRestart] = workItem.Restart.IsRestart
	esDoc[mColumns.RestartAt] = workItem.Restart.RestartAt
	esDoc[mColumns.RestartUserId] = workItem.Restart.RestartUserId
	esDoc[mColumns.IconFlags] = uint32(workItem.IconFlags)
	esDoc[mColumns.CommentNum] = workItem.CommentNum
	esDoc[mColumns.ResumeAt] = workItem.Resume.ResumeAt
	esDoc[mColumns.VersionId] = workItem.VersionId
	esDoc[mColumns.WorkItemTypeId] = workItem.WorkItemTypeId
	esDoc[mColumns.WorkItemTypeKey] = workItem.WorkItemTypeKey
	esDoc[mColumns.FlowTemplateId] = workItem.WorkFlowTemplateId
	esDoc[mColumns.FlowTemplateVersion] = workItem.WorkFlowTemplateVersion
	esDoc[mColumns.FlowId] = workItem.WorkFlowId
	esDoc[mColumns.FlowKey] = workItem.WorkFlowKey
	esDoc[mColumns.WorkItemFlowId] = workItem.WorkItemFlowId
	esDoc[mColumns.WorkItemFlowKey] = workItem.WorkItemFlowKey
	esDoc[mColumns.WorkItemStatus] = workItem.WorkItemStatus.Val
	esDoc[mColumns.WorkItemStatusKey] = workItem.WorkItemStatus.Key
	esDoc[mColumns.WorkItemStatusId] = workItem.WorkItemStatus.Id
	esDoc[mColumns.LastStatusAt] = workItem.LastWorkItemStatus.LastAt
	esDoc[mColumns.LastStatus] = workItem.LastWorkItemStatus.Val
	esDoc[mColumns.LastStatusKey] = workItem.LastWorkItemStatus.Key
	esDoc[mColumns.LastStatusId] = workItem.LastWorkItemStatus.Id
	esDoc[mColumns.CountAt] = workItem.CountAt
	esDoc[mColumns.CreatedAt] = workItem.CreatedAt
	esDoc[mColumns.UpdatedAt] = workItem.UpdatedAt
	esDoc[mColumns.DeletedAt] = workItem.DeletedAt

	esDoc[mColumns.Doc] = workItem.Doc
	esDoc[mColumns.Reason] = workItem.Reason

	return esDoc
}

func WorkItemEntityToPo(workItem *domain.WorkItem) *db.SpaceWorkItemV2 {
	dbWorkItem := &db.SpaceWorkItemV2{
		Id:           workItem.Id,
		Pid:          workItem.Pid,
		SpaceId:      workItem.SpaceId,
		UserId:       workItem.UserId,
		WorkObjectId: workItem.WorkObjectId,
		WorkItemGuid: workItem.WorkItemGuid,
		WorkItemName: workItem.WorkItemName,

		IsRestart:           workItem.Restart.IsRestart,
		RestartAt:           workItem.Restart.RestartAt,
		RestartUserId:       workItem.Restart.RestartUserId,
		IconFlags:           uint32(workItem.IconFlags),
		CommentNum:          workItem.CommentNum,
		ResumeAt:            workItem.Resume.ResumeAt,
		VersionId:           workItem.VersionId,
		WorkItemTypeId:      workItem.WorkItemTypeId,
		WorkItemTypeKey:     string(workItem.WorkItemTypeKey),
		FlowTemplateId:      workItem.WorkFlowTemplateId,
		FlowTemplateVersion: workItem.WorkFlowTemplateVersion,
		FlowId:              workItem.WorkFlowId,
		FlowKey:             workItem.WorkFlowKey,
		WorkItemFlowId:      workItem.WorkItemFlowId,
		WorkItemFlowKey:     workItem.WorkItemFlowKey,
		FlowMode:            string(workItem.FlowMode),
		WorkItemStatus:      workItem.WorkItemStatus.Val,
		WorkItemStatusKey:   workItem.WorkItemStatus.Key,
		WorkItemStatusId:    workItem.WorkItemStatus.Id,

		LastStatusAt:  workItem.LastWorkItemStatus.LastAt,
		LastStatus:    workItem.LastWorkItemStatus.Val,
		LastStatusKey: workItem.LastWorkItemStatus.Key,
		LastStatusId:  workItem.LastWorkItemStatus.Id,

		Reason:  workItem.Reason.ToJSON(),
		CountAt: workItem.CountAt,

		CreatedAt: workItem.CreatedAt,
		UpdatedAt: workItem.UpdatedAt,
		DeletedAt: workItem.DeletedAt,
	}

	docJson, _ := json.Marshal(workItem.Doc)
	dbWorkItem.Doc = string(docJson)

	return dbWorkItem
}

func WorkItemPoToEntity(workItem *db.DbSpaceWorkItem) *domain.WorkItem {
	ent := &domain.WorkItem{
		Id:           workItem.Id,
		Pid:          workItem.Pid,
		SpaceId:      workItem.SpaceId,
		UserId:       workItem.UserId,
		WorkObjectId: workItem.WorkObjectId,
		WorkItemGuid: workItem.WorkItemGuid,
		WorkItemName: workItem.WorkItemName,
		Restart: domain.Restart{
			IsRestart:     workItem.IsRestart,
			RestartAt:     workItem.RestartAt,
			RestartUserId: workItem.RestartUserId,
		},
		IconFlags:  domain.IconFlag(workItem.IconFlags),
		CommentNum: workItem.CommentNum,
		ChildNum:   workItem.ChildNum,
		Resume: domain.Resume{
			ResumeAt: workItem.ResumeAt,
		},
		VersionId:               workItem.VersionId,
		WorkItemTypeId:          workItem.WorkItemTypeId,
		WorkItemTypeKey:         consts.WorkItemTypeKey(workItem.WorkItemTypeKey),
		WorkFlowTemplateId:      workItem.FlowTemplateId,
		WorkFlowTemplateVersion: workItem.FlowTemplateVersion,
		WorkFlowId:              workItem.FlowId,
		WorkFlowKey:             workItem.FlowKey,
		WorkItemFlowId:          workItem.WorkItemFlowId,
		WorkItemFlowKey:         workItem.WorkItemFlowKey,
		FlowMode:                consts.WorkFlowMode(workItem.FlowMode),
		WorkItemStatus: domain.WorkItemStatus{
			Val: workItem.WorkItemStatus,
			Key: workItem.WorkItemStatusKey,
			Id:  workItem.WorkItemStatusId,
		},
		LastWorkItemStatus: domain.LastWorkItemStatus{
			LastAt: workItem.LastStatusAt,
			Val:    workItem.LastStatus,
			Key:    workItem.LastStatusKey,
			Id:     workItem.LastStatusId,
		},
		CountAt: workItem.CountAt,

		CreatedAt: workItem.CreatedAt,
		UpdatedAt: workItem.UpdatedAt,
		DeletedAt: workItem.DeletedAt,

		Doc: domain.WorkItemDoc{},
	}

	ent.Doc.PlanStartAt = workItem.PlanStartAt
	ent.Doc.PlanCompleteAt = workItem.PlanCompleteAt
	ent.Doc.ProcessRate = workItem.ProcessRate
	ent.Doc.Remark = workItem.Remark
	ent.Doc.Describe = workItem.Describe
	ent.Doc.Priority = workItem.Priority
	ent.Doc.Tags = workItem.TagSlice()
	ent.Doc.Directors = workItem.DirectorSlice()
	ent.Doc.Followers = workItem.FollowerSlice()
	ent.Doc.Participators = workItem.ParticipatorSlice()

	json.Unmarshal([]byte(workItem.Reason), &ent.Reason)

	return ent
}

func WorkItemPoToEntities(list []*db.DbSpaceWorkItem) []*domain.WorkItem {
	var items []*domain.WorkItem
	for _, v := range list {
		items = append(items, WorkItemPoToEntity(v))
	}
	return items
}

func WorkItemFLowNodeEntityToPo(workItemFlowNode *domain.WorkItemFlowNode) *db.SpaceWorkItemFlowV2 {
	return &db.SpaceWorkItemFlowV2{
		Id:             workItemFlowNode.Id,
		FlowId:         workItemFlowNode.FlowId,
		FlowTemplateId: workItemFlowNode.FlowTemplateId,
		SpaceId:        workItemFlowNode.SpaceId,
		WorkItemId:     workItemFlowNode.WorkItemId,

		WorkItemRoleId:  workItemFlowNode.WorkItemRoleId,
		WorkItemRoleKey: workItemFlowNode.WorkItemRoleKey,

		Directors: workItemFlowNode.Directors.ToJsonString(),

		FlowNodeUuid:    workItemFlowNode.FlowNodeUuid,
		FlowNodeStatus:  int32(workItemFlowNode.FlowNodeStatus),
		FlowNodeCode:    workItemFlowNode.FlowNodeCode,
		FlowNodePassed:  workItemFlowNode.FlowNodePassed,
		FlowNodeReached: workItemFlowNode.FlowNodeReached,
		StartAt:         workItemFlowNode.StartAt,
		FinishAt:        workItemFlowNode.FinishAt,
		PlanStartAt:     workItemFlowNode.PlanTime.StartAt,
		PlanCompleteAt:  workItemFlowNode.PlanTime.CompleteAt,

		CreatedAt: workItemFlowNode.CreatedAt,
		UpdatedAt: workItemFlowNode.UpdatedAt,
		DeletedAt: workItemFlowNode.DeletedAt,
	}
}

func WorkItemFLowNodeEntityToPos(workItemFlowNodes []*domain.WorkItemFlowNode) []*db.SpaceWorkItemFlowV2 {
	var workItemFlowNodesPos []*db.SpaceWorkItemFlowV2
	for _, workItemFlowNode := range workItemFlowNodes {
		workItemFlowNodesPos = append(workItemFlowNodesPos, WorkItemFLowNodeEntityToPo(workItemFlowNode))
	}
	return workItemFlowNodesPos
}

func WorkItemFlowNodePoToEntity(workItemFlowNode *db.SpaceWorkItemFlowV2) *domain.WorkItemFlowNode {

	flowNode := &domain.WorkItemFlowNode{
		Id:             workItemFlowNode.Id,
		FlowId:         workItemFlowNode.FlowId,
		FlowTemplateId: workItemFlowNode.FlowTemplateId,
		SpaceId:        workItemFlowNode.SpaceId,
		WorkItemId:     workItemFlowNode.WorkItemId,

		WorkItemRoleId:  workItemFlowNode.WorkItemRoleId,
		WorkItemRoleKey: workItemFlowNode.WorkItemRoleKey,

		FlowNodeUuid:    workItemFlowNode.FlowNodeUuid,
		FlowNodeStatus:  domain.FlowNodeStatus(workItemFlowNode.FlowNodeStatus),
		FlowNodeCode:    workItemFlowNode.FlowNodeCode,
		FlowNodePassed:  workItemFlowNode.FlowNodePassed,
		FlowNodeReached: workItemFlowNode.FlowNodeReached,
		StartAt:         workItemFlowNode.StartAt,
		FinishAt:        workItemFlowNode.FinishAt,
		PlanTime:        domain.PlanTime{StartAt: workItemFlowNode.PlanStartAt, CompleteAt: workItemFlowNode.PlanCompleteAt},
		CreatedAt:       workItemFlowNode.CreatedAt,
		UpdatedAt:       workItemFlowNode.UpdatedAt,
		DeletedAt:       workItemFlowNode.DeletedAt,
		Directors:       domain.Directors{}.FormJsonString(workItemFlowNode.Directors),
	}

	return flowNode
}

func WorkItemFlowNodePoToEntities(workItemFlowNodes []*db.SpaceWorkItemFlowV2) []*domain.WorkItemFlowNode {
	var workItemFlowNodesPos []*domain.WorkItemFlowNode
	for _, workItemFlowNode := range workItemFlowNodes {
		workItemFlowNodesPos = append(workItemFlowNodesPos, WorkItemFlowNodePoToEntity(workItemFlowNode))
	}
	return workItemFlowNodesPos
}

func WorkItemFLowRolePoToEntity(workItemFLowRole *db.SpaceWorkItemFlowRoleV2) *domain.WorkItemFlowRole {
	ent := &domain.WorkItemFlowRole{
		Id:              workItemFLowRole.Id,
		SpaceId:         workItemFLowRole.SpaceId,
		WorkItemId:      workItemFLowRole.WorkItemId,
		FlowId:          workItemFLowRole.FlowId,
		FlowTemplateId:  workItemFLowRole.FlowTemplateId,
		WorkItemRoleId:  workItemFLowRole.WorkItemRoleId,
		WorkItemRoleKey: workItemFLowRole.WorkItemRoleKey,
		CreatedAt:       workItemFLowRole.CreatedAt,
		UpdatedAt:       workItemFLowRole.UpdatedAt,
		DeletedAt:       workItemFLowRole.DeletedAt,
		Directors:       domain.Directors{}.FormJsonString(workItemFLowRole.Directors),
	}

	return ent
}

func WorkItemFLowRolePoToEntities(workItemFLowRoles []*db.SpaceWorkItemFlowRoleV2) []*domain.WorkItemFlowRole {
	var workItemFLowRolesPos []*domain.WorkItemFlowRole
	for _, workItemFLowRole := range workItemFLowRoles {
		workItemFLowRolesPos = append(workItemFLowRolesPos, WorkItemFLowRolePoToEntity(workItemFLowRole))
	}
	return workItemFLowRolesPos
}

func WorkItemFLowRoleEntityToPo(workItemFLowRole *domain.WorkItemFlowRole) *db.SpaceWorkItemFlowRoleV2 {
	return &db.SpaceWorkItemFlowRoleV2{
		Id:              workItemFLowRole.Id,
		SpaceId:         workItemFLowRole.SpaceId,
		WorkItemId:      workItemFLowRole.WorkItemId,
		FlowId:          workItemFLowRole.FlowId,
		FlowTemplateId:  workItemFLowRole.FlowTemplateId,
		WorkItemRoleId:  workItemFLowRole.WorkItemRoleId,
		WorkItemRoleKey: workItemFLowRole.WorkItemRoleKey,
		CreatedAt:       workItemFLowRole.CreatedAt,
		UpdatedAt:       workItemFLowRole.UpdatedAt,
		DeletedAt:       workItemFLowRole.DeletedAt,
		Directors:       workItemFLowRole.Directors.ToJsonString(),
	}
}

func WorkItemFLowRoleEntityToPos(workItemFLowRoles []*domain.WorkItemFlowRole) []*db.SpaceWorkItemFlowRoleV2 {
	var workItemFLowRolesPos []*db.SpaceWorkItemFlowRoleV2
	for _, workItemFLowRole := range workItemFLowRoles {
		workItemFLowRolesPos = append(workItemFLowRolesPos, WorkItemFLowRoleEntityToPo(workItemFLowRole))
	}
	return workItemFLowRolesPos
}

func WorkItemFilePoToEntity(fileInfo *db.SpaceFileInfo) *domain.WorkItemFile {
	// 转换逻辑
	return &domain.WorkItemFile{
		Id:         fileInfo.Id,
		SpaceId:    fileInfo.SpaceId,
		WorkItemId: fileInfo.SourceId,
		FileInfo: domain.FileInfo{
			FileInfoId: fileInfo.FileInfoId,
			FileName:   fileInfo.FileName,
			FileSize:   fileInfo.FileSize,
			FileUri:    fileInfo.FileUri,
		},
		CreatedAt: fileInfo.CreatedAt,
		UpdatedAt: fileInfo.UpdatedAt,
		DeletedAt: fileInfo.DeletedAt,
		Status:    fileInfo.Status,
	}
}

func WorkItemFilePoToEntities(fileInfos []*db.SpaceFileInfo) domain.WorkItemFiles {
	var fileInfosPos domain.WorkItemFiles
	for _, fileInfo := range fileInfos {
		fileInfosPos = append(fileInfosPos, WorkItemFilePoToEntity(fileInfo))
	}
	return fileInfosPos
}

func WorkItemFileEntityToPo(fileInfo *domain.WorkItemFile) *db.SpaceFileInfo {
	// 转换逻辑
	return &db.SpaceFileInfo{
		Id:         fileInfo.Id,
		SpaceId:    fileInfo.SpaceId,
		SourceId:   fileInfo.WorkItemId,
		SourceType: 1,
		FileInfoId: fileInfo.FileInfo.FileInfoId,
		FileName:   fileInfo.FileInfo.FileName,
		FileSize:   fileInfo.FileInfo.FileSize,
		FileUri:    fileInfo.FileInfo.FileUri,
		CreatedAt:  fileInfo.CreatedAt,
		UpdatedAt:  fileInfo.UpdatedAt,
		DeletedAt:  fileInfo.DeletedAt,
		Status:     fileInfo.Status,
	}
}
