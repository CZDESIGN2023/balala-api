package data

import (
	"context"
	"go-cs/pkg/pprint"
	"testing"
)

func TestStaticsEsRepo_GetWorkbenchCount(t *testing.T) {
	var userId int64 = 42

	ids, err := SpaceRepo.GetUserSpaceIds(context.Background(), userId)
	if err != nil {
		t.Error(err)
		return
	}
	v2, err := StaticsEsRepo.GetWorkbenchCount(context.Background(), userId, ids)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(v2)
}

func TestStaticsEsRepo_GetSpaceWorkbenchCount(t *testing.T) {
	v2, err := StaticsEsRepo.GetSpaceWorkbenchCount(context.Background(), 42, 87)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(v2)
}

func TestStaticsEsRepo_GetSpaceWorkObjectCountByIds(t *testing.T) {
	res, err := StaticsEsRepo.GetSpaceWorkObjectCountByIds(context.Background(), 1271, nil, 0, 0)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}

func TestStaticsEsRepo_GetSpaceWorkVersionCountByIds(t *testing.T) {
	res, err := StaticsEsRepo.GetSpaceWorkVersionCountByIds(context.Background(), 87, nil, 0, 0)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}

func TestStaticsEsRepo_GetSpaceUserCountByIds(t *testing.T) {
	res, err := StaticsEsRepo.GetSpaceUserCount(context.Background(), 87, 0, 0)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}
