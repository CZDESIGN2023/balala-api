package consts

const (
	UserConfigKey_NotifySwitchGlobal        = "notify_switch_global" // 全局通知开关
	UserConfigKey_NotifySwitchThirdPlatform = "notify_switch_third_platform"
	UserConfigKey_NotifySwitchSpace         = "notify_switch_space"
)

// key 名称
var UserConfigKeyName = map[string]string{
	UserConfigKey_NotifySwitchGlobal:        "通知",
	UserConfigKey_NotifySwitchThirdPlatform: "通道消息通知",
	UserConfigKey_NotifySwitchSpace:         "项目通知",
}

// 根据key获取名称
func GetUserConfigKeyName(key string) string {
	if name, ok := UserConfigKeyName[key]; ok {
		return name
	}
	return ""
}
