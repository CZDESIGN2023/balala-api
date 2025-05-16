package data

import (
	"context"
	"testing"
)

func TestConfigRepo_UpdateByKey(t *testing.T) {
	err := CondfigRepo.UpdateByKey(context.Background(), "register", "false")
	if err != nil {
		t.Error(err)
	}
}
