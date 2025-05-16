package data

import (
	"context"
	"testing"
)

func TestSpaceWorkObjectRepo_GetSpaceWorkObjectByIds(t *testing.T) {
	SpaceWorkObjectRepo.GetSpaceWorkObjectByIds(context.Background(), []int64{1})
}

func TestSpaceWorkObjectRepo_IsEmpty(t *testing.T) {
	empty, err := SpaceWorkObjectRepo.IsEmpty(context.Background(), 368)
	if err != nil {
		t.Error(err)
	}

	t.Log(empty)
}
