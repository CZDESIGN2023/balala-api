package command

import (
	"context"
	user_repo "go-cs/internal/domain/user/repo"
	login_log_domain "go-cs/internal/domain/user_login_log"
	login_log_repo "go-cs/internal/domain/user_login_log/repo"
	login_log_service "go-cs/internal/domain/user_login_log/service"
	"go-cs/internal/utils"
	"go-cs/internal/utils/user_agent"
	"go-cs/pkg/qqwry"
	"time"
)

type AddUserLoginLogCmd struct {
	userRepo       user_repo.UserRepo
	loginLogRepo   login_log_repo.UserLoginLogRepo
	loginLogServie *login_log_service.UserLgoinLogService
}

func NewAddUserLoginLogCommand(
	userRepo user_repo.UserRepo,
	loginLogRepo login_log_repo.UserLoginLogRepo,
	loginLogServie *login_log_service.UserLgoinLogService,
) *AddUserLoginLogCmd {
	return &AddUserLoginLogCmd{
		userRepo:       userRepo,
		loginLogRepo:   loginLogRepo,
		loginLogServie: loginLogServie,
	}
}

func (cmd *AddUserLoginLogCmd) Excute(ctx context.Context, userId int64, msg string, status int32) {
	user, _ := cmd.userRepo.GetUserByUserId(ctx, userId)
	if user == nil {
		return
	}

	ip := utils.GetIpFrom(ctx)
	//更新最后登录日志
	ipUtil := qqwry.NewQQwry()
	ipLocalInfo := ipUtil.Find(ip)
	ua := ""

	httpRequest := utils.GetRequestFromTransport(ctx)
	if httpRequest != nil {
		ua = httpRequest.UserAgent()
	}

	userAgent := user_agent.New(ua)
	// browser, _ := userAgent.Browser()
	os := userAgent.OS()
	loginStatus := login_log_domain.LoginStauts(status)
	loginMsg := msg
	loginLog := cmd.loginLogServie.NewLog(login_log_domain.LoginUser{
		LoginUserId:       user.Id,
		LoginUserName:     user.UserName,
		LoginUserNickname: user.UserNickname,
	}, loginStatus, login_log_domain.LoginIp{
		IpAddr:   ip,
		Location: ipLocalInfo.Area + " " + ipLocalInfo.Country,
	}, ua, os, loginMsg, time.Now())
	_ = cmd.loginLogRepo.CreateLoginInfo(ctx, loginLog)

}
