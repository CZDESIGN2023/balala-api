package data

import (
	"context"
	"testing"
)

func TestSpaceFileInfoRepo_SoftDelSpaceWorkItemFileInfo(t *testing.T) {
	err := SpaceFileInfoRepo.SoftDelSpaceWorkItemFileInfo(context.Background(), 995, 87, 343)
	if err != nil {
		t.Error(err)
	}
}

func TestSpaceFileInfoRepo_HardDelSpaceWorkItemFileInfo(t *testing.T) {
	err := SpaceFileInfoRepo.HardDelSpaceWorkItemFileInfo(context.Background(), 995, 87, 343)
	if err != nil {
		t.Error(err)
	}
}

func TestSpaceFileInfoRepo_CountWorkItemFileNum(t *testing.T) {
	num, err := SpaceFileInfoRepo.CountWorkItemFileNum(context.Background(), 358)
	if err != nil {
		t.Error(err)
	}

	t.Log(num)
}

func TestSpaceFileInfoRepo_SaveSpaceFileInfo(t *testing.T) {
	// err := SpaceFileInfoRepo.SaveSpaceFileInfo(context.Background(), &db.SpaceFileInfo{
	// 	SpaceId:  87,
	// 	FileName: "12323",
	// })
	// if err != nil {
	// 	t.Error(err)
	// }
}
