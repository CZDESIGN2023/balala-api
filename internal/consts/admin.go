package consts

import (
	"slices"
)

type SystemRole int

const (
	SystemRole_Normal     SystemRole = 0
	SystemRole_Admin      SystemRole = 50
	SystemRole_Enterprise SystemRole = 70
	SystemRole_SuperAdmin SystemRole = 100
)

var systemRoles = []SystemRole{SystemRole_Admin, SystemRole_Enterprise, SystemRole_SuperAdmin, SystemRole_Normal}

func GetAllSystemRoles() []SystemRole {
	return slices.Clone(systemRoles)
}

func (r SystemRole) String() string {

	switch r {
	case SystemRole_Admin:
		return "系统管理员"
	case SystemRole_Enterprise:
		return "企业管理员"
	case SystemRole_SuperAdmin:
		return "超级管理员"
	case SystemRole_Normal:
		return "普通用户"
	}

	return ""
}
