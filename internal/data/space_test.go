package data

import (
	"context"
	"testing"
)

func TestSpaceRepo_GetSpace(t *testing.T) {
	space, err := SpaceRepo.GetSpace(context.Background(), 87)
	if err != nil {
		t.Error(err)
	}

	t.Log(space)
}

func TestSpaceRepo_GetUserSpace(t *testing.T) {
	space, err := SpaceRepo.GetSpaceByCreator(context.Background(), 42, 87)
	if err != nil {
		t.Error(err)
	}

	t.Log(space)
}

func TestSpaceRepo_GetSpaceConfig(t *testing.T) {
	config, err := SpaceRepo.GetSpaceConfig(context.Background(), 87)
	if err != nil {
		t.Error(err)
	}

	t.Log(config)
}
