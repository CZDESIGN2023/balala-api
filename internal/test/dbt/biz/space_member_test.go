package biz

import (
	"context"
	"go-cs/internal/test/dbt"
	"go-cs/pkg/pprint"
	"testing"
)

func TestRemoveMySpaceMember(t *testing.T) {
	// member, err := dbt.UC.SpaceMemberUsecase.removeMySpaceMember(context.Background(), "remove_member", 42, 87, 21, 42)
	// if err != nil {
	// 	t.Error(err)
	// }

	// pprint.Print(member)
}

func TestGetSpaceMemberList(t *testing.T) {
	list, err := dbt.UC.SpaceMemberUsecase.GetMySpaceMemberList(context.Background(), 42, 87, "hmy")
	if err != nil {
		t.Error(err)
	}
	pprint.Println(list)
}
