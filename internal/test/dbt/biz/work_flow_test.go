package biz

import (
	"context"
	v1 "go-cs/api/work_flow/v1"
	"go-cs/internal/test/dbt"
	"testing"
)

func TestSearchTpltByOwnerRule(t *testing.T) {
	result, err := dbt.UC.WorkFlowUsecase.GetOwnerRuleRelationTemplateRequest(context.Background(), MockUserInfo(42), &v1.GetOwnerRuleRelationTemplateRequest{
		SpaceId:  1481,
		OwnerUid: 42,
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}

func Test_CopyWorkFlowFromOtherSpace(t *testing.T) {
	result, err := dbt.UC.WorkFlowUsecase.CopyWorkFlowFromOtherSpace(context.Background(), MockUserInfo(42), 87, 10032, "newName", 1)
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}
