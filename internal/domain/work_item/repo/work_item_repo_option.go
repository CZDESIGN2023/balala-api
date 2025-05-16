package repo

type WithDocOption struct {
	All           bool
	PlanTime      bool
	ProcessRate   bool
	Remark        bool
	Describe      bool
	Priority      bool
	Tags          bool
	Directors     bool
	Followers     bool
	Participators bool
}

type WithOption struct {
	FlowNodes bool
	FlowRoles bool
	FileInfos bool
}
