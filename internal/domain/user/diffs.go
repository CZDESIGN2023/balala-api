package user

import shared "go-cs/internal/pkg/domain"

const (
	Diff_Avatar    shared.PropDiff = "avatar"
	Diff_Name      shared.PropDiff = "name"
	Diff_NickName  shared.PropDiff = "nick_name"
	Diff_Email     shared.PropDiff = "email"
	Diff_Password  shared.PropDiff = "user_password"
	Diff_Role      shared.PropDiff = "role"
	Diff_Status    shared.PropDiff = "status"
	Diff_LastLogin shared.PropDiff = "last_login"
)

const (
	Diff_ThirdPlatform_Notify shared.PropDiff = "notify"
)
