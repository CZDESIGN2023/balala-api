package data

import (
	"context"
	"go-cs/pkg/pprint"
	"testing"
)

func TestStaticsRepo_GetWorkbenchCount(t *testing.T) {
	var userId int64 = 42

	ids, err := SpaceRepo.GetUserSpaceIds(context.Background(), userId)
	if err != nil {
		t.Error(err)
		return
	}
	v2, err := StaticsRepo.GetWorkbenchCount(context.Background(), userId, ids)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(v2)
}

func TestStaticsRepo_GetSpaceWorkbenchCount(t *testing.T) {
	v2, err := StaticsRepo.GetSpaceWorkbenchCount(context.Background(), 42, 87)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(v2)
}

func TestStaticsRepo_GetSpaceWorkItemTagCount(t *testing.T) {
	v2, err := StaticsRepo.GetSpaceWorkItemTagCount(context.Background(), 87)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(v2)
}

func TestStaticsRepo_GetTagRelationCounts2(t *testing.T) {
	v2, err := StaticsRepo.GetSpaceTagRelationCounts(context.Background(), 87, []int64{489, 568})
	if err != nil {
		t.Error(err)
	}

	pprint.Println(v2)
}

func TestStaticsRepo_GetTagRelationCounts(t *testing.T) {
	v2, err := StaticsRepo.GetTagRelationCounts(context.Background(), []int64{489, 568})
	if err != nil {
		t.Error(err)
	}

	pprint.Println(v2)
}

func TestStaticsRepo_TagCountMap(t *testing.T) {
	v2, err := StaticsRepo.TagCountMap(context.Background(), []int64{489, 568})
	if err != nil {
		t.Error(err)
	}

	pprint.Println(v2)
}

func Benchmark_GetSpaceWorkItemTagCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = StaticsRepo.GetSpaceWorkItemTagCount(context.Background(), 87)
	}
}

func TestStaticsRepo_GetSpaceWorkObjectCountByIds(t *testing.T) {
	res, err := StaticsRepo.GetSpaceWorkObjectCountByIds(context.Background(), 87, nil, 0, 0)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}

func TestStaticsRepo_GetSpaceWorkVersionCountByIds(t *testing.T) {
	res, err := StaticsRepo.GetSpaceWorkVersionCountByIds(context.Background(), 87, nil, 0, 0)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}

func TestStaticsRepo_GetSpaceUserCountByIds(t *testing.T) {
	res, err := StaticsRepo.GetSpaceUserCount(context.Background(), 87, 0, 0)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(res)
}
