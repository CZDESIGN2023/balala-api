package data

// func TestSpaceWorkFlowRepo_SaveWorkFlowNodeV2(t *testing.T) {
// 	err := SpaceWorkFlowRepo.SaveWorkFlowNodeV2(context.Background(), &db.SpaceWorkItemFlowV2{
// 		WorkItemId:   1,
// 		FlowNodeUuid: "kuhufu",
// 		Directors:    "[]",
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestSpaceWorkFlowRepo_CreateWorkFlowNodeV2(t *testing.T) {
// 	var list = []*db.SpaceWorkItemFlowV2{
// 		{
// 			FlowNodeUuid: "kuhufu",
// 			Directors:    "[]",
// 		},
// 		{
// 			FlowNodeUuid: "kuhufu",
// 			Directors:    "[]",
// 		},
// 	}
// 	err := SpaceWorkFlowRepo.CreateWorkFlowNodeV2(context.Background(), list...)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(list)
// }

// func TestSpaceWorkFlowRepo_ReplaceAllWorkItemFlowDirectorV2(t *testing.T) {
// 	v2, err := SpaceWorkFlowRepo.ReplaceAllWorkItemFlowDirectorV2(context.Background(), 87, 4, 1)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(v2)
// }

// func TestSpaceWorkFlowRepo_AddDirectorForWorkItemFlow(t *testing.T) {
// 	v2, err := SpaceWorkFlowRepo.AddDirectorForWorkItemFlows(context.Background(), []int64{169}, 5)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(v2)
// }

// func TestSpaceWorkFlowRepo_RemoveDirectorForWorkItemFlow(t *testing.T) {
// 	v2, err := SpaceWorkFlowRepo.RemoveDirectorForWorkItemFlows(context.Background(), []int64{169}, 1)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(v2)
// }

// func TestSpaceWorkFlowRepo_SetDirectorForWorkItemFlow(t *testing.T) {
// 	v2, err := SpaceWorkFlowRepo.SetDirectorForWorkItemFlow(context.Background(), 33, []int64{1})
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(v2)
// }

// func TestGetSpaceWorkFlowIdsByDirector(t *testing.T) {
// 	v2, err := SpaceWorkFlowRepo.GetSpaceWorkItemFlowIdsByDirector(context.Background(), []int64{359}, 21)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(v2)
// }

// func TestGetWorkItemProgressingFlowNode(t *testing.T) {
// 	v2, err := SpaceWorkFlowRepo.GetWorkItemProgressingFlowNode(context.Background(), 358)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(v2)
// }

// func TestGetDirectorsByWorkItemIds(t *testing.T) {
// 	v2, err := SpaceWorkFlowRepo.NodeDirectorsMapByWorkItemIds(context.Background(), []int64{359})
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	pprint.Println(v2)
// }
