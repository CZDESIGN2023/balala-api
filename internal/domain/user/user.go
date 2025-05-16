package user

import (
	"errors"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/rand"
	"time"
)

type User struct {
	shared.AggregateRoot

	Id            int64             ` bson:"_id" json:"id"`
	UserName      string            ` bson:"user_name" json:"user_name"`
	Mobile        string            ` bson:"mobile" json:"mobile"`
	UserNickname  string            ` bson:"user_nickname" json:"user_nickname"`
	UserPinyin    string            ` bson:"user_pinyin" json:"user_pinyin"`
	UserPassword  string            ` bson:"user_password" json:"user_password"`
	UserSalt      string            ` bson:"user_salt" json:"user_salt"`
	UserStatus    int32             ` bson:"user_status" json:"user_status"`
	UserEmail     string            ` bson:"user_email" json:"user_email"`
	Sex           int32             ` bson:"sex" json:"sex"`
	Avatar        string            ` bson:"avatar" json:"avatar"`
	Remark        string            ` bson:"remark" json:"remark"`
	Describe      string            ` bson:"describe" json:"describe"`
	LastLoginIp   string            ` bson:"last_login_ip" json:"last_login_ip"`
	LastLoginTime int64             ` bson:"last_login_time" json:"last_login_time"`
	Role          consts.SystemRole ` bson:"role" json:"role"`
	CreatedAt     int64             ` bson:"created_at" json:"created_at"`
	UpdatedAt     int64             ` bson:"updated_at" json:"updated_at"`
	DeletedAt     int64             ` bson:"deleted_at" json:"deleted_at"`
}

func (u *User) ChangeNickName(newName string, oper shared.Oper) error {
	if u.UserNickname == newName {
		return nil
	}

	oldValue := u.UserNickname

	newNamePy := utils.Pinyin(newName)
	u.UserNickname = newName
	u.UserPinyin = newNamePy

	u.AddDiff(Diff_NickName)

	u.AddMessage(oper, &domain_message.PersonalChangeNickName{
		UserId:      u.Id,
		OldNickName: oldValue,
		NewNickName: newName,
	})

	return nil
}

func (u *User) UpdateLastLogin(ip string, tm time.Time) {
	u.LastLoginIp = ip
	u.LastLoginTime = tm.Unix()

	u.AddDiff(Diff_LastLogin)
}

func (u *User) ValidatePwd(pwd string) bool {
	return u.UserPassword == utils.EncryptUserPassword(pwd, u.UserSalt)
}

// 重置密码
func (u *User) ResetPwd(pwd string, oper shared.Oper) {

	salt := rand.S(10)
	password := utils.EncryptUserPassword(pwd, salt)

	u.UserSalt = salt
	u.UserPassword = password

	u.AddDiff(Diff_Password)

	u.AddMessage(oper, &domain_message.PersonalResetPwd{
		UserId:       u.Id,
		UserNickname: u.UserNickname,
	})
}

// 改密码
func (u *User) ChangePwd(oldPwd string, newPwd string, oper shared.Oper) error {

	if utils.EncryptUserPassword(oldPwd, u.UserSalt) != u.UserPassword {
		return errors.New("原密码不正确")
	}

	salt := rand.S(10)
	password := utils.EncryptUserPassword(newPwd, salt)

	u.UserSalt = salt
	u.UserPassword = password

	u.AddDiff(Diff_Password)

	u.AddMessage(oper, &domain_message.PersonalResetPwd{
		UserId:       u.Id,
		UserNickname: u.UserNickname,
	})

	return nil
}

func (u *User) UpdateRole(newRoleId consts.SystemRole) {
	if u.Role == newRoleId {
		return
	}
	u.Role = newRoleId

	u.AddDiff(Diff_Role)
}

func (u *User) UpdateAvatar(newAvatar string, oper shared.Oper) {
	u.Avatar = newAvatar

	u.AddDiff(Diff_Avatar)

	u.AddMessage(oper, &domain_message.PersonalChangeAvatar{
		UserId: u.Id,
	})
}

// 注销账号
func (u *User) Cancel() {
	u.UserName = ""
	u.UserStatus = 3

	u.AddDiff(Diff_Status)
	u.AddDiff(Diff_Name)
}

func (u *User) IsSystemAdmin() bool {
	return u.Role >= consts.SystemRole_Admin
}

func (u *User) IsEnterpriseAdmin() bool {
	return u.Role == consts.SystemRole_Enterprise
}

func (u *User) IsSystemSuperAdmin() bool {
	return u.Role == consts.SystemRole_SuperAdmin
}

// 用户角色权限是否大于
func (u *User) RoleGreaterThan(role consts.SystemRole) bool {
	return u.Role > role
}
