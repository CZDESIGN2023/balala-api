package biz

import (
	"context"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	"go-cs/internal/test/dbt"
	"go-cs/internal/utils"
	"go-cs/internal/utils/rand"
	"go-cs/pkg/pprint"
	"testing"
)

func Test_GetMySpace(t *testing.T) {
	space, err := dbt.UC.SpaceUsecase.GetMySpace(context.Background(), 42, 87)
	if err != nil {
		t.Error(err)
	}

	t.Log(space)
}

func Test_GetMySpaceList(t *testing.T) {
	list, err := dbt.UC.SpaceUsecase.GetMySpaceList(context.Background(), 42)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(list)
}

func Test_CreateSpace(t *testing.T) {

	//随机获取几个用户
	users, err := dbt.R.UserRepo.SearchUser(context.Background(), "", []int64{1001})
	if err != nil {
		t.Fatal(err)
		return
	}

	maxUserCount := rand.N(1, len(users)%5+1)
	users = users[:maxUserCount]

	var members []*db.SpaceMember
	for _, v := range users {
		members = append(members, &db.SpaceMember{
			UserId: v.Id,
			RoleId: consts.MEMBER_ROLE_MANAGER,
		})
	}

	s, err := dbt.UC.SpaceUsecase.CreateMySpace(context.Background(), MockUserInfo(users[0].Id), "测试"+rand.Digits(3), "测试", members)

	if err != nil {
		t.Fatal(err)
		return
	}

	println(users[0].Id, s.Id, s.SpaceName)

}

func Test_DelSpace(t *testing.T) {

	s, err := dbt.UC.SpaceUsecase.DelSpace(context.Background(), MockUserInfo(185), 3470, "")
	if err != nil {
		t.Fatal(err)
		return
	}

	println(s)
}

func Test_CopySpace(t *testing.T) {

	s, err := dbt.UC.SpaceUsecase.Copy(context.Background(), &utils.LoginUserInfo{UserId: 42}, 87, "hmy1234", "")
	if err != nil {
		t.Fatal(err)
		return
	}

	println(s)
}

func Test_TransferSpaceOwnership(t *testing.T) {
	err := dbt.UC.SpaceUsecase.TransferSpaceOwnership(context.Background(), 185, 87, 152, 42)
	if err != nil {
		t.Fatal(err)
		return
	}
}
