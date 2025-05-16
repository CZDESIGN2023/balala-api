package data

import (
	"context"
	"testing"
)

func TestSpaceTagRepo_GetSpaceAllTags(t *testing.T) {
	// tags, err := SpaceTagRepo.GetSpaceTagList(context.Background(), 87)
	// if err != nil {
	// 	t.Error(err)
	// }

	// pprint.Println(tags)
}

func TestSpaceTagRepo_CheckTagNameIsExist(t *testing.T) {
	exist, err := SpaceTagRepo.CheckTagNameIsExist(context.Background(), 87, "12")
	if err != nil {
		t.Error(err)
	}

	t.Log(exist)
}
