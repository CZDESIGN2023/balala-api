package facade

import (
	witem_role_domain "go-cs/internal/domain/work_item_role"
)

type WorkItemRoleInfo struct {
	roles    []*witem_role_domain.WorkItemRole
	rolesMap map[string]*witem_role_domain.WorkItemRole
}

func (w *WorkItemRoleInfo) GetRoleByKey(roleKey string) *witem_role_domain.WorkItemRole {
	return w.rolesMap[roleKey]
}

func (w *WorkItemRoleInfo) GetRoles() []*witem_role_domain.WorkItemRole {
	return w.roles
}

func BuildWorkItemRoleInfo(roles []*witem_role_domain.WorkItemRole) *WorkItemRoleInfo {

	rolesMap := make(map[string]*witem_role_domain.WorkItemRole)
	for _, v := range roles {
		rolesMap[v.Key] = v
	}

	roleInfo := &WorkItemRoleInfo{
		roles:    roles,
		rolesMap: rolesMap,
	}

	return roleInfo
}
