package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/user_login_log"
)

func UserLoginLogEntityToPo(log *domain.UserLoginLog) *db.UserLoginLog {
	return &db.UserLoginLog{
		Id:                log.Id,
		LoginUserId:       log.LoginUser.LoginUserId,
		LoginUserName:     log.LoginUser.LoginUserName,
		LoginUserNickname: log.LoginUser.LoginUserNickname,
		Ipaddr:            log.LoginIp.IpAddr,
		LoginLocation:     log.LoginIp.Location,
		Browser:           log.Browser,
		Os:                log.Os,
		Status:            int32(log.LoginStauts),
		Msg:               log.Msg,
		LoginAt:           log.LoginAt,
		CreatedAt:         log.CreatedAt,
		UpdatedAt:         log.UpdatedAt,
		DeletedAt:         log.DeletedAt,
	}
}
