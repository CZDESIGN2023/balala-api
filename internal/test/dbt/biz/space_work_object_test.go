package biz

import (
	"context"
	v1 "go-cs/api/space_work_object/v1"
	"go-cs/internal/test/dbt"
	"testing"
)

func TestGetMySpaceWorkObjectList(t *testing.T) {
	space, err := dbt.UC.SpaceWorkObjectUsecase.QSpaceWorkObjectList(context.Background(), MockUserInfo(42), &v1.SpaceWorkObjectListRequest{})
	if err != nil {
		t.Error(err)
	}

	t.Log(space)
}

func TestDelMySpaceWorkObject2(t *testing.T) {

}
