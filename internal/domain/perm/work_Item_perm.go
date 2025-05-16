package checker

import (
	"go-cs/internal/consts"

	"github.com/casbin/casbin/v2"
	casbinM "github.com/casbin/casbin/v2/model"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
)

type WorkItemEditPermissionScene string

const (
	CreateWorkItemScene WorkItemEditPermissionScene = "create_scene"
	EditWorkItemScene   WorkItemEditPermissionScene = "edit_scene"
)

type workItemEditPerm struct {
	e       map[string]*casbin.CachedEnforcer
	setting map[string]interface{}
}

func newWorkItemPerm() *workItemEditPerm {

	var workItemPolicyModel = `
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

	var spaceSuperPolicy = `
	g, _creator, _creator
	g, _node_owner, _node_owner
	g, _no_relation, _no_relation

	p, _creator, PERM_CREATE_SPACE_WORK_OBJECT, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_OBJECT, *, allow
	p, _creator, PERM_CREATE_SPACE_WORK_VERSION, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_VERSION, *, allow
	p, _creator, PERM_CREATE_SPACE_TAG, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_RESUME_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, allow
	p, _creator, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _creator, PERM_CREATE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_UPGRADE_SPACE_WORK_ITEM_FLOW, *, allow

	p, _node_owner, PERM_CREATE_SPACE_WORK_OBJECT, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_OBJECT, *, allow
	p, _node_owner, PERM_CREATE_SPACE_WORK_VERSION, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_VERSION, *, allow
	p, _node_owner, PERM_CREATE_SPACE_TAG, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_RESUME_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, allow
	p, _node_owner, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _node_owner, PERM_CREATE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_UPGRADE_SPACE_WORK_ITEM_FLOW, *, allow

	p, _no_relation, PERM_CREATE_SPACE_WORK_OBJECT, *, allow
	p, _no_relation, PERM_MODIFY_SPACE_WORK_OBJECT, *, allow
	p, _no_relation, PERM_CREATE_SPACE_WORK_VERSION, *, allow
	p, _no_relation, PERM_MODIFY_SPACE_WORK_VERSION, *, allow
	p, _no_relation, PERM_CREATE_SPACE_TAG, *, allow
	p, _no_relation, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_RESUME_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, allow
	p, _no_relation, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _no_relation, PERM_CREATE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_CREATE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_UPGRADE_SPACE_WORK_ITEM_FLOW, *, allow

	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_CREATE_SPACE_WORK_OBJECT, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_MODIFY_SPACE_WORK_OBJECT, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_CREATE_SPACE_WORK_VERSION, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_MODIFY_SPACE_WORK_VERSION, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_CREATE_SPACE_TAG, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_CREATE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_UPGRADE_SPACE_WORK_ITEM_FLOW, *, allow
`

	var spaceManagerPolicy = `
	g, _creator, _creator
	g, _node_owner, _node_owner
	g, _no_relation, _no_relation

	p, _creator, PERM_CREATE_SPACE_WORK_OBJECT, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_OBJECT, *, allow
	p, _creator, PERM_CREATE_SPACE_WORK_VERSION, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_VERSION, *, allow
	p, _creator, PERM_CREATE_SPACE_TAG, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_RESUME_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, allow
	p, _creator, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _creator, PERM_CREATE_SPACE_WORK_ITEM, *, allow

	p, _node_owner, PERM_CREATE_SPACE_WORK_OBJECT, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_OBJECT, *, allow
	p, _node_owner, PERM_CREATE_SPACE_WORK_VERSION, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_VERSION, *, allow
	p, _node_owner, PERM_CREATE_SPACE_TAG, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_RESUME_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, allow
	p, _node_owner, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _node_owner, PERM_CREATE_SPACE_WORK_ITEM, *, allow

	p, _no_relation, PERM_CREATE_SPACE_WORK_OBJECT, *, allow
	p, _no_relation, PERM_MODIFY_SPACE_WORK_OBJECT, *, allow
	p, _no_relation, PERM_CREATE_SPACE_WORK_VERSION, *, allow
	p, _no_relation, PERM_MODIFY_SPACE_WORK_VERSION, *, allow
	p, _no_relation, PERM_CREATE_SPACE_TAG, *, allow
	p, _no_relation, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_RESUME_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, allow
	p, _no_relation, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _no_relation, PERM_CREATE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_CREATE_SPACE_WORK_ITEM, *, allow

	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_CREATE_SPACE_WORK_OBJECT, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_MODIFY_SPACE_WORK_OBJECT, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_CREATE_SPACE_WORK_VERSION, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_MODIFY_SPACE_WORK_VERSION, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_CREATE_SPACE_TAG, *, allow
	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_CREATE_SPACE_WORK_ITEM, *, allow
`

	var spaceEditorPolicy = `
	g, _creator, _creator
	g, _node_owner, _node_owner
	g, _no_relation, _no_relation

	p, _no_relation, Scene_WORK_ITEM_CREATE:PERM_CREATE_SPACE_TAG, *, allow

	p, _creator, PERM_CREATE_SPACE_TAG, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_RESUME_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, allow
	p, _creator, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _creator, PERM_CREATE_SPACE_WORK_ITEM, *, allow

	p, _node_owner, PERM_CREATE_SPACE_TAG, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, allow
	p, _node_owner, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _node_owner, PERM_CREATE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_DELETE_SPACE_WORK_ITEM, *, deny`

	var spaceWatcherPolicy = `
	g, _creator, _creator
	g, _node_owner, _node_owner
	g, _no_relation, _no_relation

	p, _creator, PERM_CREATE_SPACE_TAG, *, deny
	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_RESUME_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, allow
	p, _creator, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _creator, PERM_CREATE_SPACE_WORK_ITEM, *, deny

	p, _node_owner, PERM_CREATE_SPACE_TAG, *, deny
	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM, *, deny
	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE, *, deny
	p, _node_owner, PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE, *, allow
	p, _node_owner, PERM_CREATE_SPACE_WORK_ITEM, *, deny
	p, _node_owner, PERM_DELETE_SPACE_WORK_ITEM, *, deny`

	var taskSuperPolicy = `
	g, _creator, _creator
	g, _node_owner, _node_owner
	g, _no_relation, _no_relation

	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CREATE_SPACE_TAG, *, allow

	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_CREATE_SPACE_TAG, *, allow

	p, _no_relation, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_CREATE_SPACE_TAG, *, allow`

	var taskManagerPolicy = `
	g, _creator, _creator
	g, _node_owner, _node_owner
	g, _no_relation, _no_relation

	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CREATE_SPACE_TAG, *, allow

	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_CREATE_SPACE_TAG, *, allow

	p, _no_relation, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _no_relation, PERM_CREATE_SPACE_TAG, *, allow`

	var taskSpaceEditorPolicy = `
	g, _creator, _creator
	g, _node_owner, _node_owner
	g, _no_relation, _no_relation

	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CREATE_SPACE_TAG, *, allow

	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_DELETE_SPACE_WORK_ITEM, *, deny
	p, _node_owner, PERM_CREATE_SPACE_TAG, *, allow

	p, _no_relation, PERM_MODIFY_SPACE_WORK_ITEM, *, deny
	p, _no_relation, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, deny
	p, _no_relation, PERM_DELETE_SPACE_WORK_ITEM, *, deny
	p, _no_relation, PERM_CREATE_SPACE_TAG, *, deny`

	var taskSpaceWatcherPolicy = `
	g, _creator, _creator
	g, _node_owner, _node_owner
	g, _no_relation, _no_relation

	p, _creator, PERM_MODIFY_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_DELETE_SPACE_WORK_ITEM, *, allow
	p, _creator, PERM_CREATE_SPACE_TAG, *, deny

	p, _node_owner, PERM_MODIFY_SPACE_WORK_ITEM, *, deny
	p, _node_owner, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, allow
	p, _node_owner, PERM_DELETE_SPACE_WORK_ITEM, *, deny
	p, _node_owner, PERM_CREATE_SPACE_TAG, *, deny

	p, _no_relation, PERM_MODIFY_SPACE_WORK_ITEM, *, deny
	p, _no_relation, PERM_CHANGE_STATE_SPACE_WORK_ITEM, *, deny
	p, _no_relation, PERM_DELETE_SPACE_WORK_ITEM, *, deny
	p, _no_relation, PERM_CREATE_SPACE_TAG, *, deny`

	var workItemPermission = map[string]interface{}{

		"create_scene": map[string]interface{}{
			"work_item.create": "Scene_WORK_ITEM_CREATE:" + consts.PERM_CREATE_SPACE_WORK_ITEM,
			"object.create":    "Scene_WORK_ITEM_CREATE:" + consts.PERM_CREATE_SPACE_WORK_OBJECT,
			"object.modify":    "Scene_WORK_ITEM_CREATE:" + consts.PERM_MODIFY_SPACE_WORK_OBJECT,
			"version.create":   "Scene_WORK_ITEM_CREATE:" + consts.PERM_CREATE_SPACE_WORK_VERSION,
			"version.modify":   "Scene_WORK_ITEM_CREATE:" + consts.PERM_MODIFY_SPACE_WORK_VERSION,
			"tag.create":       "Scene_WORK_ITEM_CREATE:" + consts.PERM_CREATE_SPACE_TAG,
		},

		"edit_scene": map[string]interface{}{
			"work_item.modify":        consts.PERM_MODIFY_SPACE_WORK_ITEM,
			"work_item.delete":        consts.PERM_DELETE_SPACE_WORK_ITEM,
			"work_item.change_state":  consts.PERM_CHANGE_STATE_SPACE_WORK_ITEM,
			"flow_node.modify":        consts.PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE,
			"flow_node.confirm_state": consts.PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE,
			"object.create":           consts.PERM_CREATE_SPACE_WORK_OBJECT,
			"object.modify":           consts.PERM_MODIFY_SPACE_WORK_OBJECT,
			"version.create":          consts.PERM_CREATE_SPACE_WORK_VERSION,
			"version.modify":          consts.PERM_MODIFY_SPACE_WORK_VERSION,
			"tag.create":              consts.PERM_CREATE_SPACE_TAG,
			"comment":                 "*",
			"remind":                  "*",
			"task.create":             consts.PERM_CREATE_SPACE_WORK_ITEM,
			"flow.upgrade":            consts.PERM_UPGRADE_SPACE_WORK_ITEM_FLOW,
		},
		"task_edit_scene": map[string]interface{}{
			"work_item.modify":       consts.PERM_MODIFY_SPACE_WORK_ITEM,
			"work_item.delete":       consts.PERM_DELETE_SPACE_WORK_ITEM,
			"work_item.change_state": consts.PERM_CHANGE_STATE_SPACE_WORK_ITEM,
			"tag.create":             consts.PERM_CREATE_SPACE_TAG,
			"comment":                "*",
		},
	}

	// 父任务
	mSuper, _ := casbinM.NewModelFromString(workItemPolicyModel)
	aSuper := stringadapter.NewAdapter(spaceSuperPolicy)
	eSuper, err := casbin.NewCachedEnforcer(mSuper, aSuper)
	if err != nil {
		panic(err)
	}

	mManager, _ := casbinM.NewModelFromString(workItemPolicyModel)
	aManager := stringadapter.NewAdapter(spaceManagerPolicy)
	eManager, err := casbin.NewCachedEnforcer(mManager, aManager)
	if err != nil {
		panic(err)
	}

	mEditor, _ := casbinM.NewModelFromString(workItemPolicyModel)
	aEditor := stringadapter.NewAdapter(spaceEditorPolicy)
	eEditor, err := casbin.NewCachedEnforcer(mEditor, aEditor)
	if err != nil {
		panic(err)
	}

	mWatcher, _ := casbinM.NewModelFromString(workItemPolicyModel)
	aWatcher := stringadapter.NewAdapter(spaceWatcherPolicy)
	eWatcher, err := casbin.NewCachedEnforcer(mWatcher, aWatcher)
	if err != nil {
		panic(err)
	}

	// 子任务
	mTaskSuper, _ := casbinM.NewModelFromString(workItemPolicyModel)
	aTaskSuper := stringadapter.NewAdapter(taskSuperPolicy)
	eTaskSuper, err := casbin.NewCachedEnforcer(mTaskSuper, aTaskSuper)
	if err != nil {
		panic(err)
	}

	mTaskManager, _ := casbinM.NewModelFromString(workItemPolicyModel)
	aTaskManager := stringadapter.NewAdapter(taskManagerPolicy)
	eTaskManager, err := casbin.NewCachedEnforcer(mTaskManager, aTaskManager)
	if err != nil {
		panic(err)
	}

	mTaskEditor, _ := casbinM.NewModelFromString(workItemPolicyModel)
	aTaskEditor := stringadapter.NewAdapter(taskSpaceEditorPolicy)
	eTaskEditor, err := casbin.NewCachedEnforcer(mTaskEditor, aTaskEditor)
	if err != nil {
		panic(err)
	}

	mTaskWatcher, _ := casbinM.NewModelFromString(workItemPolicyModel)
	aTaskWatcher := stringadapter.NewAdapter(taskSpaceWatcherPolicy)
	eTaskWatcher, err := casbin.NewCachedEnforcer(mTaskWatcher, aTaskWatcher)
	if err != nil {
		panic(err)
	}

	return &workItemEditPerm{
		setting: workItemPermission,
		e: map[string]*casbin.CachedEnforcer{
			"supper":       eSuper,
			"manager":      eManager,
			"editor":       eEditor,
			"watcher":      eWatcher,
			"task_supper":  eTaskSuper,
			"task_manager": eTaskManager,
			"task_editor":  eTaskEditor,
			"task_watcher": eTaskWatcher,
		},
	}

}

func (p *workItemEditPerm) GetPermissionWithScene(role int64, itemRole string, scene WorkItemEditPermissionScene) map[string]interface{} {

	editorRole := "watcher"
	switch role {
	case consts.MEMBER_ROLE_SPACE_CREATOR:
		editorRole = "supper"
	case consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER:
		editorRole = "supper"
	case consts.MEMBER_ROLE_MANAGER:
		editorRole = "manager"
	case consts.MEMBER_ROLE_EDITOR:
		editorRole = "editor"
	case consts.MEMBER_ROLE_WATCHER:
		editorRole = "watcher"
	}

	if itemRole == "" {
		itemRole = "_no_relation"
	}

	settingVal := p.setting[string(scene)]
	if settingVal == nil {
		return nil
	}

	e := p.e[editorRole]
	if e == nil {
		return nil
	}

	result := make(map[string]interface{})
	setting := settingVal.(map[string]interface{})
	for k, v := range setting {
		if v == "*" {
			result[k] = true
			continue
		}

		hasPerm, _ := e.Enforce(itemRole, v, "*")
		result[k] = hasPerm
	}

	return result
}

// 判断任务表单的权限
func (p *workItemEditPerm) Check(role int64, itemRole string, method string) bool {

	editorRole := "watcher"
	switch role {
	case consts.MEMBER_ROLE_SPACE_CREATOR:
		editorRole = "supper"
	case consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER:
		editorRole = "supper"
	case consts.MEMBER_ROLE_MANAGER:
		editorRole = "manager"
	case consts.MEMBER_ROLE_EDITOR:
		editorRole = "editor"
	case consts.MEMBER_ROLE_WATCHER:
		editorRole = "watcher"
	}

	if itemRole == "" {
		itemRole = "_no_relation"
	}

	e := p.e[editorRole]
	if e == nil {
		return false
	}

	hasPerm, err := e.Enforce(itemRole, method, "*")
	//if hasPerm {
	//	fmt.Printf("%s CAN %s %s\n", editorRole, method, "*")
	//} else {
	//	fmt.Printf("%s CANNOT %s %s\n", editorRole, method, "*")
	//}
	if err != nil {
		return false
	}
	return hasPerm
}

// 判断任务表单的权限
func (p *workItemEditPerm) CheckTask(role int64, itemRole string, method string) bool {

	editorRole := "watcher"
	switch role {
	case consts.MEMBER_ROLE_SPACE_CREATOR:
		editorRole = "supper"
	case consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER:
		editorRole = "supper"
	case consts.MEMBER_ROLE_MANAGER:
		editorRole = "manager"
	case consts.MEMBER_ROLE_EDITOR:
		editorRole = "editor"
	case consts.MEMBER_ROLE_WATCHER:
		editorRole = "watcher"
	}

	if itemRole == "" {
		itemRole = "_no_relation"
	}

	e := p.e["task_"+editorRole]
	if e == nil {
		return false
	}

	hasPerm, err := e.Enforce(itemRole, method, "*")
	// if hasPerm {
	// 	fmt.Printf("%s CAN %s %s\n", itemRole, method, "*")
	// } else {
	// 	fmt.Printf("%s CANNOT %s %s\n", itemRole, method, "*")
	// }
	if err != nil {
		return false
	}
	return hasPerm
}

func (p *workItemEditPerm) GetTaskPermission(role int64, itemRole string) map[string]interface{} {

	editorRole := "watcher"
	switch role {
	case consts.MEMBER_ROLE_SPACE_CREATOR:
		editorRole = "supper"
	case consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER:
		editorRole = "supper"
	case consts.MEMBER_ROLE_MANAGER:
		editorRole = "manager"
	case consts.MEMBER_ROLE_EDITOR:
		editorRole = "editor"
	case consts.MEMBER_ROLE_WATCHER:
		editorRole = "watcher"
	}

	if itemRole == "" {
		itemRole = "_no_relation"
	}

	settingVal := p.setting["task_edit_scene"]
	if settingVal == nil {
		return nil
	}

	e := p.e["task_"+editorRole]
	if e == nil {
		return nil
	}

	result := make(map[string]interface{})
	setting := settingVal.(map[string]interface{})
	for k, v := range setting {
		if v == "*" {
			result[k] = true
			continue
		}

		hasPerm, _ := e.Enforce(itemRole, v, "*")
		result[k] = hasPerm

		// if hasPerm {
		// 	fmt.Printf("%s CAN %s %s\n", itemRole, v, "*")
		// } else {
		// 	fmt.Printf("%s CANNOT %s %s\n", itemRole, v, "*")
		// }
	}

	return result
}
