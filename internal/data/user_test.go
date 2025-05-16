package data

import (
	"context"
	"testing"
)

func TestUserRepo_UpdateUserFields(t *testing.T) {
	// list, err := UserRepo.SearchSpaceUser(context.Background(), "", nil)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// for _, v := range list {
	// 	pinyin := utils.Pinyin(v.UserNickname)
	// 	if v.UserPinyin == pinyin {
	// 		continue
	// 	}

	// 	// err := UserRepo.UpdateUserFields(context.Background(), v.Id, map[string]any{
	// 	// 	"UserPinyin": pinyin,
	// 	// })
	// 	// if err != nil {
	// 	// 	t.Error(err)
	// 	// }
	// }
}

func TestUserRepo_GetUserByIdFromRedis(t *testing.T) {
	redis, err := UserRepo.GetUserByUserId(context.Background(), 42)
	if err != nil {
		t.Error(err)
	}

	t.Log(redis)
}

func TestUserRepo_UserMap(t *testing.T) {
	redis, err := UserRepo.UserMap(context.Background(), []int64{42, 1})
	if err != nil {
		t.Error(err)
	}

	t.Log(redis)
}

func TestUserRepo_GetUserByEmail(t *testing.T) {
	res, err := UserRepo.GetUserByEmail(context.Background(), "1428119544@qq.com")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestUserRepo_BindEmail(t *testing.T) {
	err := UserRepo.BindEmail(context.Background(), 42, "1428119544@qq.com")
	if err != nil {
		t.Error(err)
	}
}

func TestUserRepo_UnbindEmail(t *testing.T) {
	err := UserRepo.UnbindEmail(context.Background(), 42)
	if err != nil {
		t.Error(err)
	}
}

func TestUserRepo_BindMobile(t *testing.T) {
	err := UserRepo.BindMobile(context.Background(), 42, "12345678901")
	if err != nil {
		t.Error(err)
	}
}

func TestUserRepo_UnbindMobile(t *testing.T) {
	err := UserRepo.UnbindMobile(context.Background(), 42)
	if err != nil {
		t.Error(err)
	}
}

func TestUserRepo_BindThirdPlatform(t *testing.T) {
	// err := UserRepo.BindThirdPlatform(context.Background(), &db.ThirdPfAccount{
	// 	UserId:     42,
	// 	PfCode:     int32(comm.ThirdPlatformCode_pf_IM),
	// 	PfUserKey:  "1428119544",
	// 	PfUserName: "kuhufu",
	// })
	// if err != nil {
	// 	t.Error(err)
	// }
}

func TestUserRepo_UnbindThirdPlatform(t *testing.T) {
	// err := UserRepo.UnbindThirdPlatform(context.Background(), 42, int32(comm.ThirdPlatformCode_pf_IM))
	// if err != nil {
	// 	t.Error(err)
	// }
}

func TestUserRepo_SetEmailCode(t *testing.T) {
	// err := UserRepo.SetEmailCode(context.Background(), "1428119544@qq.com", int64(comm.VerifyType_SetPwd), "123456", "hmy123")
	// if err != nil {
	// 	t.Error(err)
	// }
}

func TestUserRepo_GetEmailCode(t *testing.T) {
	// res := UserRepo.GetEmailCode(context.Background(), "1428119544@qq.com", int64(comm.VerifyType_SetPwd), "hmy123")
	// t.Log(res)
}

func TestUserRepo_SetSmsCode(t *testing.T) {
	// err := UserRepo.SetSmsCode(context.Background(), "12345678901", int64(comm.VerifyType_SetPwd), "hmy123", "666666")
	// if err != nil {
	// 	t.Error(err)
	// }
}

func TestUserRepo_GetSmsCode(t *testing.T) {
	// res := UserRepo.GetSmsCode(context.Background(), "12345678901", int64(comm.VerifyType_SetPwd), "hmy123")
	// t.Log(res)
}

func TestUserRepo_AddUser(t *testing.T) {
	// res, err := UserRepo.AddUser(context.Background(), &db.User{})
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(res)
}

func TestUserRepo_SearchUserList(t *testing.T) {
	// res, err := UserRepo.SearchUserList(context.Background(), nil, "")
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(res)
}

func TestUserRepo_GetAllThirdPfAccount(t *testing.T) {
	res, err := UserRepo.GetAllThirdPfAccount(context.Background(), 21)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
