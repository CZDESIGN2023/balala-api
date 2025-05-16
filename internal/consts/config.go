package consts

const (
	CONFIG_SPACE_FILE_DOMAIN        = "space.file.domain"
	CONFIG_USER_AVATAR_DOMAIN       = "user.avatar.domain"
	CONFIG_NOTIFY_REDIRECT_DOMAIN   = "notify.redirect.domain"
	CONFIG_BALALA_ASSECT_DOMIN      = "balala.assect.domain"
	CONFIG_BALALA_LOGO              = "balala.logo"
	CONFIG_BALALA_TITLE             = "balala.title"
	CONFIG_BALALA_REGISTER_ENTRY    = "balala.register.entry"
	CONFIG_BALALA_BG                = "balala.bg"
	CONFIG_BALALA_ATTACH            = "balala.attach"
	CONFIG_BALALA_VERSION           = "balala.version"
	CONFIG_BALALA_THIRD_IM_CODE     = "balala.third.im.code"
	CONFIG_BALALA_THIRD_QL_CODE     = "balala.third.ql.code"
	CONFIG_BALALA_THIRD_HALALA_CODE = "balala.third.halala.code"
)

func MutableConfigKeyList() []string {
	return []string{
		CONFIG_NOTIFY_REDIRECT_DOMAIN,
		CONFIG_BALALA_LOGO,
		CONFIG_BALALA_TITLE,
		CONFIG_BALALA_REGISTER_ENTRY,
		CONFIG_BALALA_BG,
		CONFIG_BALALA_ATTACH,
	}
}

func AllConfigKeyList() []string {
	return []string{
		CONFIG_SPACE_FILE_DOMAIN,
		CONFIG_USER_AVATAR_DOMAIN,
		CONFIG_NOTIFY_REDIRECT_DOMAIN,
		CONFIG_BALALA_ASSECT_DOMIN,
		CONFIG_BALALA_LOGO,
		CONFIG_BALALA_TITLE,
		CONFIG_BALALA_REGISTER_ENTRY,
		CONFIG_BALALA_BG,
		CONFIG_BALALA_ATTACH,
		CONFIG_BALALA_VERSION,
	}
}
