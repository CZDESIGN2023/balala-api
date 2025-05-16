package data

import (
	"context"
	"go-cs/internal/domain/space_member"
	"testing"
)

func TestSpaceMemberRepo_AllIsMember(t *testing.T) {
	member, err := SpaceMemberRepo.AllIsMember(context.Background(), 87, 42)
	if err != nil {
		t.Error(err)
	}

	t.Log(member)
}

func TestSpaceMemberRepo_GetUserSpaceMember(t *testing.T) {
	// 	member, err := SpaceMemberRepo.GetUserSpaceMember(context.Background(), 42, 87)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}

	// t.Log(member)
}

func TestSpaceMemberRepo_AnyOneIsMember(t *testing.T) {
	member, err := SpaceMemberRepo.AnyOneIsMember(context.Background(), 87, 42)
	if err != nil {
		t.Error(err)
	}

	t.Log(member)
}

func TestSpaceMemberRepo_UpdateUserSpaceRank(t *testing.T) {

	err := SpaceMemberRepo.UpdateUserSpaceOrder(context.Background(), 42, 2, 1)
	if err != nil {
		t.Error(err)
	}
}

func Test(t *testing.T) {

	err := SpaceMemberRepo.AddSpaceMember(context.Background(), &space_member.SpaceMember{
		Id:            1,
		UserId:        1,
		SpaceId:       1,
		RoleId:        1,
		Remark:        "",
		Ranking:       0,
		Notify:        0,
		CreatedAt:     0,
		UpdatedAt:     0,
		DeletedAt:     0,
		HistoryRoleId: 0,
	})
	if err != nil {
		t.Error(err)
	}
}
