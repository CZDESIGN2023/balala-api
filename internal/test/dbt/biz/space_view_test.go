package biz

import (
	"context"
	pb "go-cs/api/space_view/v1"
	"go-cs/internal/test/dbt"
	"testing"
)

func Test_ViewList(t *testing.T) {
	result, err := dbt.UC.SpaceViewUsecase.ViewList(context.Background(), 2025, []int64{1940, 467}, "all")
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}

func Test_Create(t *testing.T) {
	err := dbt.UC.SpaceViewUsecase.Create(context.Background(), 42, &pb.CreateViewRequest{
		Name:        "123",
		SpaceId:     87,
		QueryConfig: "",
		TableConfig: "",
		Type:        2,
	})
	if err != nil {
		t.Error(err)
	}
}
func Test_SetName(t *testing.T) {
	err := dbt.UC.SpaceViewUsecase.SetName(context.Background(), 42, &pb.SetViewNameRequest{
		Id:   1,
		Name: "12",
	})
	if err != nil {
		t.Error(err)
	}
}
