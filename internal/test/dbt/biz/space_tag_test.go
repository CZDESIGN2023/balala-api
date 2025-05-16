package biz

import (
	"context"
	"go-cs/internal/test/dbt"
	"testing"
)

func TestGetMySpaceTagListV2(t *testing.T) {
	res, err := dbt.UC.SpaceTagUsecase.GetMySpaceTagListV2(context.Background(), 42, 161)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func BenchmarkGetMySpaceTagListV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = dbt.UC.SpaceTagUsecase.GetMySpaceTagListV2(context.Background(), 42, 87)
	}
}
