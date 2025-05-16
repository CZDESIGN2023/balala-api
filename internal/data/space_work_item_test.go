package data

import (
	"context"
	"encoding/json"
	"go-cs/internal/bean/vo"
	"go-cs/pkg/pprint"
	"testing"
)

func TestWorkItemDocMarshalJSON(t *testing.T) {
	marshal, err := json.Marshal(&vo.WorkItemDoc{})
	if err != nil {
		t.Error(err)
	}

	t.Log(string(marshal))
}

func TestSpaceWorkItemRepo_GetWorkItemByTagV2(t *testing.T) {
	res, err := SpaceWorkItemRepo.GetWorkItemByTagV2(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}

func TestSpaceTagRepo_RemoveTagFromAllWorkItem(t *testing.T) {
	update, err := SpaceWorkItemRepo.RemoveTagFromAllWorkItem(context.Background(), 87, 2)
	if err != nil {
		t.Error(err)
	}

	t.Log(update)
}

func TestSpaceWorkItemRepo_GetWorkItemByDirectorV2(t *testing.T) {
	res, err := SpaceWorkItemRepo.GetSpaceWorkItemByDirectorV2(context.Background(), 87, 42)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}

func TestGetSpaceWorkItemIdsByPidV2(t *testing.T) {
	v2, err := SpaceWorkItemRepo.GetSpaceWorkItemIdsByPid(context.Background(), 342)
	if err != nil {
		t.Error(err)
	}

	t.Log(v2)
}

func TestSpaceWorkItemRepo_UpdateSpaceAllWorkItemCreator(t *testing.T) {
	creator, err := SpaceWorkItemRepo.UpdateSpaceAllWorkItemCreator(context.Background(), 87, 42, 21)
	if err != nil {
		t.Error(err)
	}

	t.Log(creator)
}

func TestSpaceWorkItemRepo_GetSpaceWorkItemIdsByWorkObject(t *testing.T) {
	creator, err := SpaceWorkItemRepo.GetSpaceWorkItemIdsByWorkObject(context.Background(), 134)
	if err != nil {
		t.Error(err)
	}

	t.Log(creator)
}

func TestSpaceWorkItemRepo_CountUserRelatedSpaceWorkItem(t *testing.T) {
	creator, err := SpaceWorkItemRepo.CountUserRelatedSpaceWorkItem(context.Background(), 87, 42)
	if err != nil {
		t.Error(err)
	}

	t.Log(creator)
}

func TestSpaceWorkItemRepo_IncrCommentNum(t *testing.T) {
	creator, err := SpaceWorkItemRepo.IncrCommentNum(context.Background(), 358, 1)
	if err != nil {
		t.Error(err)
	}

	t.Log(creator)
}

func TestSpaceWorkItemRepo_ReplaceDirectorForWorkItemBySpace(t *testing.T) {
}

func TestSpaceWorkItemRepo_GetWorkItemByDirector(t *testing.T) {
	res, err := SpaceWorkItemRepo.HasWorkItemRelateFlow(context.Background(), 150, 8293)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}
