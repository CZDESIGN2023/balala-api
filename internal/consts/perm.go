package consts

const (
	//权限相关

	PERM_CREATE_SPACE_TAG = "PERM_CREATE_SPACE_TAG" //创建TAG
	PERM_MODIFY_SPACE_TAG = "PERM_MODIFY_SPACE_TAG" //修改TAG信息
	PERM_DELETE_SPACE_TAG = "PERM_DELETE_SPACE_TAG" //删除TAG

	PERM_CREATE_SPACE_WORK_ITEM                  = "PERM_CREATE_SPACE_WORK_ITEM"                  //创建任务
	PERM_SET_SPACE_WORK_ITEM_DIRECTOR            = "PERM_SET_SPACE_WORK_ITEM_DIRECTOR"            //设置任务负责人
	PERM_MODIFY_SPACE_WORK_ITEM                  = "PERM_MODIFY_SPACE_WORK_ITEM"                  //修改任务基本信息
	PERM_CHANGE_STATE_SPACE_WORK_ITEM            = "PERM_CHANGE_STATE_SPACE_WORK_ITEM"            //变更任务状态
	PERM_DELETE_SPACE_WORK_ITEM                  = "PERM_DELETE_SPACE_WORK_ITEM"                  //删除任务
	PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE = "PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE" //确认任务节点状态
	PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE        = "PERM_MODIFY_SPACE_WORK_ITEM_FLOW_NODE"        //修改任务节点信息
	PERM_UPGRADE_SPACE_WORK_ITEM_FLOW            = "PERM_UPGRADE_SPACE_WORK_ITEM_FLOW"            //升级任务流程

	PERM_RESUME_SPACE_WORK_ITEM = "PERM_RESUME_SPACE_WORK_ITEM" //恢复任务

	PERM_CREATE_SPACE_WORK_OBJECT  = "PERM_CREATE_SPACE_WORK_OBJECT"  //创建模块
	PERM_DELETE_SPACE_WORK_OBJECT  = "PERM_DELETE_SPACE_WORK_OBJECT"  //移除模块
	PERM_DELETE_SPACE_WORK_OBJECT2 = "PERM_DELETE_SPACE_WORK_OBJECT2" //移除模块
	PERM_MODIFY_SPACE_WORK_OBJECT  = "PERM_MODIFY_SPACE_WORK_OBJECT"  //修改模块基本信息

	PERM_MODIFY_SPACE = "PERM_MODIFY_SPACE" //修改空间
	PERM_CREATE_SPACE = "PERM_CREATE_SPACE" //创建空间
	PERM_DELETE_SPACE = "PERM_DELETE_SPACE" //删除空间
	PERM_QUITE_SPACE  = "PERM_QUITE_SPACE"  //退出空间
	PERM_COPY_SPACE   = "PERM_COPY_SPACE"   //复制空间

	PERM_ADD_SPACE_MEMBER                 = "PERM_ADD_SPACE_MEMBER"                 //添加空间成员
	PERM_REMOVE_SPACE_MEMBER              = "PERM_REMOVE_SPACE_MEMBER"              //移除空间成员
	PERM_SET_SPACE_MEMBER_ROLE            = "PERM_SET_SPACE_MEMBER_ROLE"            //设置空间成员权限
	PERM_SET_SPACE_MEMBER_ROLE_TO_MANAGER = "PERM_SET_SPACE_MEMBER_ROLE_TO_MANAGER" //调整成员角色至管理员权限

	PERM_UPLOAD_SPACE_FILE = "PERM_UPLOAD_SPACE_FILE" //上传空间附件

	PERM_CREATE_SPACE_WORK_VERSION = "PERM_CREATE_SPACE_WORK_VERSION" //创建空间任务项版本
	PERM_MODIFY_SPACE_WORK_VERSION = "PERM_MODIFY_SPACE_WORK_VERSION" //修改空间任务项版本
	PERM_DELETE_SPACE_WORK_VERSION = "PERM_DELETE_SPACE_WORK_VERSION" //删除空间任务项版本

	PERM_ADD_SPACE_MANAGER    = "PERM_ADD_SPACE_MANAGER"    //添加空间管理员
	PERM_REMOVE_SPACE_MANAGER = "PERM_REMOVE_SPACE_MANAGER" //移除空间管理员

	PERM_CreateMemberCategory = "PERM_CreateMemberCategory" //添加成员分类
	PERM_ModifyMemberCategory = "PERM_ModifyMemberCategory" //修改成员分类
	PERM_DeleteMemberCategory = "PERM_DeleteMemberCategory" //删除成员分类
	PERM_DeleteComment        = "PERM_DeleteComment"        //删除成员分类

	PERM_CreateWorkFlowRole = "PERM_CREATE_WORK_FLOW_ROLE" // 添加角色
	PERM_ModifyWorkFlowRole = "PERM_MODIFY_WORK_FLOW_ROLE" // 编辑角色
	PERM_DeleteWorkFlowRole = "PERM_DELETE_WORK_FLOW_ROLE" // 删除角色

	PERM_CreateWorkFlowStatus = "PERM_CREATE_WORK_FLOW_STATUS" // 添加状态
	PERM_ModifyWorkFlowStatus = "PERM_MODIFY_WORK_FLOW_STATUS" // 编辑状态
	PERM_DeleteWorkFlowStatus = "PERM_DELETE_WORK_FLOW_STATUS" // 删除状态

	PERM_CreateSpaceWorkFlow = "CREATE_SPACE_WORK_FLOW" // 新建流程
	PERM_ModifySpaceWorkFlow = "MODIFY_SPACE_WORK_FLOW" // 编辑流程
	PERM_DeleteSpaceWorkFlow = "DELETE_SPACE_WORK_FLOW" // 流程删除

	WORK_ITEM_ROLE_CREATOR    = "_creator"
	WORK_ITEM_ROLE_NODE_OWNER = "_node_owner"

	WORK_ITEM_SUB_ROLE_NODE_ITEM_TERMINATOR  = "_item_terminator"
	WORK_ITEM_SUB_ROLE_NODE_NODE_EDITOR      = "_node_editor"
	WORK_ITEM_SUB_ROLE_NODE_ITEM_EDITOR      = "_item_editor"
	WORK_ITEM_SUB_ROLE_NODE_ITEM_COMMENTATOR = "_item_commentator"

	MEMBER_ROLE_SPACE_CREATOR        = 999 //超级管理员=空间创建者
	MEMBER_ROLE_SPACE_SUPPER_MANAGER = 99  //空间管理员
	MEMBER_ROLE_MANAGER              = 1   //普通管理员
	MEMBER_ROLE_EDITOR               = 2   //普通编辑
	MEMBER_ROLE_WATCHER              = 3   //普通观察者

	MEMBER_ROLE_SPACE_CREATOR_TAG        = "space_creator"
	MEMBER_ROLE_SPACE_SUPPER_MANAGER_TAG = "space_supper_manager"
	MEMBER_ROLE_SPACE_MANAGER_TAG        = "space_manager"
	MEMBER_ROLE_SPACE_EDITOR_TAG         = "space_editor"
	MEMBER_ROLE_SPACE_WATCHER_TAG        = "space_watcher"

	Scene_WORK_ITEM_CREATE = "Scene_WORK_ITEM_CREATE"
	Scene_WORK_ITEM_EDIT   = "Scene_WORK_ITEM_EDIT"
)

func GetSpaceMemberRoleTag(roleId int64) string {
	switch roleId {
	case MEMBER_ROLE_SPACE_CREATOR:
		return MEMBER_ROLE_SPACE_CREATOR_TAG
	case MEMBER_ROLE_SPACE_SUPPER_MANAGER:
		return MEMBER_ROLE_SPACE_SUPPER_MANAGER_TAG
	case MEMBER_ROLE_MANAGER:
		return MEMBER_ROLE_SPACE_MANAGER_TAG
	case MEMBER_ROLE_EDITOR:
		return MEMBER_ROLE_SPACE_EDITOR_TAG
	case MEMBER_ROLE_WATCHER:
		return MEMBER_ROLE_SPACE_WATCHER_TAG
	default:
		return ""
	}
}

var allMemberRole = []int64{
	MEMBER_ROLE_SPACE_CREATOR,
	MEMBER_ROLE_SPACE_SUPPER_MANAGER,
	MEMBER_ROLE_MANAGER,
	MEMBER_ROLE_EDITOR,
	MEMBER_ROLE_WATCHER,
}

var memberRoleRank = map[int64]int64{
	MEMBER_ROLE_SPACE_CREATOR:        1,
	MEMBER_ROLE_SPACE_SUPPER_MANAGER: 2,
	MEMBER_ROLE_MANAGER:              3,
	MEMBER_ROLE_EDITOR:               4,
	MEMBER_ROLE_WATCHER:              5,
}

func GetAllMemberRole() []int64 {
	return allMemberRole
}

func GetMemberRoleRank(roleId int32) int64 {
	return memberRoleRank[int64(roleId)]
}

func GetSpaceRoleName(roleId int) string {
	switch roleId {
	case MEMBER_ROLE_WATCHER:
		return "可查看"
	case MEMBER_ROLE_MANAGER:
		return "可管理"
	case MEMBER_ROLE_EDITOR:
		return "可编辑"
	case MEMBER_ROLE_SPACE_SUPPER_MANAGER:
		return "项目管理员"
	case MEMBER_ROLE_SPACE_CREATOR:
		return "项目创建者"
	default:
		return "可查看"
	}
}
