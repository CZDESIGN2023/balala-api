package biz

import (
	"context"
	"encoding/json"
	"go-cs/internal/utils"
	"go-cs/internal/utils/oper"
	"go-cs/pkg/qqwry"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
)

type MiddleareUsecase struct {
	userBiz *UserUsecase
	log     *log.Helper
}

func NewMiddlearesecase(userBiz *UserUsecase, logger log.Logger) *MiddleareUsecase {
	moduleName := "MiddleareUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	uc := &MiddleareUsecase{
		userBiz: userBiz,
		log:     hlog,
	}

	return uc
}

func (uc *MiddleareUsecase) OpLoggerMiddleware() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {

		return func(ctx context.Context, req interface{}) (interface{}, error) {

			// loginUserInfo, err := utils.GetLoginUserInfo(ctx)
			// if err != nil {
			// 	return h(ctx, req)
			// }

			// userInfo, err := uc.userBiz.MyInfo(ctx, loginUserInfo.UserId)
			// if err != nil {
			// 	return h(ctx, req)
			// }

			newCtx, opLogger := oper.NewOperLoggerWithCtx(ctx)

			// opLogger.Operator = &oper.OperUser{
			// 	OperType:         1,
			// 	OperUid:          userInfo.Id,
			// 	OperUname:        userInfo.UserName,
			// 	OperUserNickName: userInfo.UserNickname,
			// }

			if operParam, err := json.Marshal(req); err == nil {
				opLogger.RequestInfo.OperParam = string(operParam)
			}

			r := utils.GetRequestFromTransport(newCtx)
			if r != nil {
				opLogger.RequestInfo.RequestMethod = r.Method
				rawQuery := r.URL.RawQuery
				if rawQuery != "" {
					rawQuery = "?" + rawQuery
				}
				opLogger.RequestInfo.OperUrl = r.URL.Path + rawQuery
				opLogger.RequestInfo.OperIp = utils.GetIpFrom(newCtx)
				ipUtils := qqwry.NewQQwry()
				ipInfo := ipUtils.Find(opLogger.RequestInfo.OperIp)
				opLogger.RequestInfo.OperLocation = ipInfo.Area + " " + ipInfo.Country
			}

			rsp, err := h(newCtx, req)

			return rsp, err
		}

	}
}
