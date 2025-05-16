package oper

type ModuleType int
type BusinessType int
type OperatorType int
type ShowType int
type ModuleFlag uint64

const (
	//操作日志相关

	//功能模块ID

	ModuleTypeSpace               ModuleType = 1  //空间模块 关联 空间表ID
	ModuleTypeSpaceMember         ModuleType = 2  //空间成员模块 关联 空间成员表ID
	ModuleTypeSpaceTag            ModuleType = 3  //TAG模块 关联 TAG表ID
	ModuleTypeSpaceWorkObject     ModuleType = 4  //工作项模块 关联 工作项表ID
	ModuleTypeSpaceWorkItem       ModuleType = 5  //工作任务模块 关联 工作任务表ID
	ModuleTypeSpaceWorkVersion    ModuleType = 6  //工作项版本模块 关联 工作项版本表ID
	ModuleTypeSpaceMemberCategory ModuleType = 7  //工作项版本模块 关联 工作项版本表ID
	ModuleTypeWorkFlow            ModuleType = 8  //工作流
	ModuleTypeWorkItemRole        ModuleType = 9  //工作流角色
	ModuleTypeWorkItemStatus      ModuleType = 10 //工作状态
	ModuleTypeUser                ModuleType = 11 //用户
	ModuleTypeSystem              ModuleType = 12 //系统配置
	ModuleTypeSpaceGlobalView     ModuleType = 13 //全局视图
	ModuleTypeSpaceUserView       ModuleType = 14 //个人视图

	BusinessTypeNone   BusinessType = 0 //未定义
	BusinessTypeAdd    BusinessType = 1 //新增操作
	BusinessTypeDel    BusinessType = 2 //删除操作
	BusinessTypeModify BusinessType = 3 //修改操作

	OperatorTypeOther OperatorType = 0 //其它
	OperatorTypeUser  OperatorType = 1 //用户
	OperatorTypeSys   OperatorType = 2 //系统

	ShowTypeNone               ShowType = 0 //默认
	ShowTypeWorkItemNodeChange ShowType = 1 //任务节点变化
)

const (
	ModuleFlag_subWorkItem ModuleFlag = 1 << iota
)
