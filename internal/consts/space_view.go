package consts

const (
	SpaceViewType_System   = 1
	SpaceViewType_Public   = 2
	SpaceViewType_Personal = 3
)

const (
	SystemViewKey_All    = "all"
	SystemViewKey_Follow = "follow"
)

var all = []string{
	SystemViewKey_All,
	SystemViewKey_Follow,
}

var systemViewName = map[string]string{
	SystemViewKey_All:    "全部",
	SystemViewKey_Follow: "关注任务",
}

func GetAllSystemViewKeys() []string {
	return all
}

func GetSystemViewName(key string) string {
	return systemViewName[key]
}

const (
	PublicViewKey_Processing = "processing"
	PublicViewKey_Expired    = "expired"
)

var initPublicView = []string{
	PublicViewKey_Processing,
	PublicViewKey_Expired,
}

var initPublicViewName = map[string]string{
	PublicViewKey_Processing: "待办任务",
	PublicViewKey_Expired:    "逾期任务",
}

func GetAllPublicViewKeys() []string {
	return initPublicView
}

func GetPublicViewName(key string) string {
	return initPublicViewName[key]
}

func GetViewRank(key string) int64 {
	switch key {
	case SystemViewKey_All:
		return 1000
	case PublicViewKey_Processing:
		return 800
	case PublicViewKey_Expired:
		return 400
	case SystemViewKey_Follow:
		return 100
	}

	return 0
}
