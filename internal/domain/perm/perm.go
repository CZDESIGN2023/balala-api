package checker

import (
	"fmt"
	"go-cs/internal/consts"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v2"
	casbinM "github.com/casbin/casbin/v2/model"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
)

var model string = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`

// 项目权限
var policy string = `
g, 999, space_creator
g, 99, space_supper_manager
g, 1, space_manager
g, 2, space_editor
g, 3, space_watcher

g, space_creator, space_supper_manager
g, space_supper_manager, space_manager
g, space_manager, space_editor
g, space_editor, space_watcher

p, space_creator, PERM_CREATE_SPACE, *, allow
p, space_creator, PERM_DELETE_SPACE, *, allow
p, space_creator, PERM_QUITE_SPACE, *, deny

p, space_supper_manager, PERM_SET_SPACE_MEMBER_ROLE_TO_MANAGER, *, allow
p, space_supper_manager, PERM_CreateMemberCategory, *, allow
p, space_supper_manager, PERM_ModifyMemberCategory, *, allow
p, space_supper_manager, PERM_DeleteMemberCategory, *, allow
p, space_supper_manager, PERM_DeleteComment, *, allow
p, space_supper_manager, PERM_UPGRADE_SPACE_WORK_ITEM_FLOW, *, allow
p, space_supper_manager, PERM_COPY_SPACE, *, allow

p, space_manager, PERM_CREATE_SPACE_WORK_VERSION, *, allow
p, space_manager, PERM_MODIFY_SPACE_WORK_VERSION, *, allow
p, space_manager, PERM_DELETE_SPACE_WORK_VERSION, *, allow
p, space_manager, PERM_CREATE_SPACE_WORK_OBJECT, *, allow
p, space_manager, PERM_MODIFY_SPACE_WORK_OBJECT, *, allow
p, space_manager, PERM_DELETE_SPACE_WORK_OBJECT, *, allow
p, space_manager, PERM_DELETE_SPACE_WORK_OBJECT2, *, allow
p, space_manager, PERM_ADD_SPACE_MEMBER, *, allow
p, space_manager, PERM_REMOVE_SPACE_MEMBER, *, allow
p, space_manager, PERM_SET_SPACE_MEMBER_ROLE, *, allow
p, space_manager, PERM_MODIFY_SPACE, *, allow
p, space_manager, PERM_CREATE_WORK_FLOW_ROLE, *, allow
p, space_manager, PERM_MODIFY_WORK_FLOW_ROLE, *, allow
p, space_manager, PERM_DELETE_WORK_FLOW_ROLE, *, allow
p, space_manager, PERM_CREATE_WORK_FLOW_STATUS, *, allow
p, space_manager, PERM_MODIFY_WORK_FLOW_STATUS, *, allow
p, space_manager, PERM_DELETE_WORK_FLOW_STATUS, *, allow
p, space_manager, CREATE_SPACE_WORK_FLOW, *, allow
p, space_manager, MODIFY_SPACE_WORK_FLOW, *, allow
p, space_manager, DELETE_SPACE_WORK_FLOW, *, allow

p, space_editor, PERM_CREATE_SPACE_WORK_ITEM, *, allow
p, space_editor, PERM_SET_SPACE_WORK_ITEM_DIRECTOR, *, allow
p, space_editor, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
p, space_editor, PERM_DELETE_SPACE_WORK_ITEM, *, allow
p, space_editor, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
p, space_editor, PERM_UPLOAD_SPACE_FILE, *, allow
p, space_editor, PERM_MOVE_SPACE_WORK_ITEM_TO_NEW_WORK_OBJECT, *, allow
p, space_editor, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
p, space_editor, PERM_CREATE_SPACE_WORK_TASK_ITEM, *, allow
p, space_editor, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
p, space_editor, PERM_RESUME_SPACE_WORK_ITEM, *, allow
p, space_editor, PERM_CREATE_SPACE_TAG, *, allow
p, space_editor, PERM_MODIFY_SPACE_TAG, *, allow
p, space_editor, PERM_DELETE_SPACE_TAG, *, allow
p, space_editor, PERM_SET_SPACE_WORK_ITEM_TAG, *, allow

p, space_watcher, PERM_QUITE_SPACE, *, allow
`

var LevelPolicy string = `
g, 999, space_creator
g, 99, space_supper_manager
g, 1, space_manager
g, 2, space_editor
g, 3, space_watcher

p, space_creator, space_supper_manager, *, allow
p, space_creator, space_editor, *, allow
p, space_creator, space_watcher, *, allow
p, space_creator, space_manager, *, allow

p, space_supper_manager, space_manager, *, allow
p, space_supper_manager, space_editor, *, allow
p, space_supper_manager, space_watcher, *, allow

p, space_manager, space_watcher, *, allow
p, space_manager, space_editor, *, allow

`

type perm struct {
	WorkItemEditPerm *workItemEditPerm

	e  *casbin.CachedEnforcer
	el *casbin.CachedEnforcer
}

var instance *perm
var instanceOnce sync.Once

func Instance() *perm {
	instanceOnce.Do(func() {

		workItemEditPerm := newWorkItemPerm()

		m, _ := casbinM.NewModelFromString(model)
		a := stringadapter.NewAdapter(policy)
		e, err := casbin.NewCachedEnforcer(m, a)
		if err != nil {
			panic(err)
		}

		m3, _ := casbinM.NewModelFromString(model)
		a3 := stringadapter.NewAdapter(LevelPolicy)
		e3, err3 := casbin.NewCachedEnforcer(m3, a3)
		if err3 != nil {
			panic(err3)
		}

		instance = &perm{
			e:                e,
			el:               e3,
			WorkItemEditPerm: workItemEditPerm,
		}
	})
	return instance
}

// 用这个判断权限
func (p *perm) Check(role int64, method string) bool {
	hasPerm, err := p.e.Enforce(strconv.FormatInt(role, 10), method, "*")
	if hasPerm {
		fmt.Printf("%v CAN %s %v\n", role, method, "*")
	} else {
		fmt.Printf("%v CANNOT %v %v\n", role, method, "*")
	}
	if err != nil {
		return false
	}
	return hasPerm
}

func (p *perm) CheckLevel(roleA int64, roleB int64) bool {

	roleBTag := consts.GetSpaceMemberRoleTag(roleB)
	hasPerm, err := p.el.Enforce(strconv.FormatInt(roleA, 10), roleBTag, "*")
	if err != nil {
		return false
	}
	return hasPerm
}
