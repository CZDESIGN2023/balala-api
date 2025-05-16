package errs

import (
	"context"
	"errors"
	"fmt"
	"go-cs/api/comm"
	"go-cs/pkg/i18n"
	"runtime"

	"gorm.io/gorm"
)

func Cast(err error) *comm.ErrorInfo {
	var v *comm.ErrorInfo
	if ok := errors.As(err, &v); !ok {
		return &comm.ErrorInfo{
			Code:    int32(comm.ErrorCode_ERROR_INTERNAL),
			Message: err.Error(),
		}
	}
	return v
}

func New(ctx context.Context, code comm.ErrorCode, fmtValues ...interface{}) *comm.ErrorInfo {
	var message = i18nErrMsg(ctx, code, fmtValues...)

	return &comm.ErrorInfo{
		Code:    int32(code),
		Message: message,
	}
}

func Custom(ctx context.Context, code any, message string) *comm.ErrorInfo {
	var code_ int32
	switch c := code.(type) {
	case comm.ErrorCode:
		code_ = int32(c)
	case int:
		code_ = int32(c)
	case int32:
		code_ = c
	case int64:
		code_ = int32(c)
	}

	return &comm.ErrorInfo{
		Code:    code_,
		Message: message,
	}
}

// Param 参数错误
// 快速创建参数错误
func Param(ctx context.Context, fmtValues ...interface{}) *comm.ErrorInfo {
	var code = comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT

	return &comm.ErrorInfo{
		Code:    int32(code),
		Message: i18nErrMsg(ctx, code, fmtValues...),
	}
}

func NotLogin(ctx context.Context) *comm.ErrorInfo {
	var code = comm.ErrorCode_LOGIN_USER_NOT_LOGIN

	return &comm.ErrorInfo{
		Code:    int32(code),
		Message: i18nErrMsg(ctx, code),
	}
}

func RepoErr(ctx context.Context, err error, fmtValues ...interface{}) *comm.ErrorInfo {
	code := comm.ErrorCode_ERROR_INTERNAL
	if errors.Is(err, gorm.ErrRecordNotFound) {
		code = comm.ErrorCode_ERROR_RES_NOT_EXIST
		return &comm.ErrorInfo{
			Code:    int32(code),
			Message: i18nErrMsg(ctx, code, fmtValues...),
		}
	} else {
		return Internal(ctx, err)
	}
}

func Internal(ctx context.Context, err error) *comm.ErrorInfo {
	code := comm.ErrorCode_ERROR_INTERNAL

	_, file, line, _ := runtime.Caller(1)
	logStr := fmt.Sprintf("error: %v, file: %v:%v", err, file, line)

	log.Error(logStr)

	return &comm.ErrorInfo{
		Code:    int32(code),
		Message: i18nErrMsg(ctx, code),
	}
}

func NoPerm(ctx context.Context) *comm.ErrorInfo {
	var code = comm.ErrorCode_PERMISSION_INSUFFICIENT_PERMISSIONS

	return &comm.ErrorInfo{
		Code:    int32(code),
		Message: i18nErrMsg(ctx, code),
	}
}

func Business(ctx context.Context, values ...interface{}) *comm.ErrorInfo {
	var code = comm.ErrorCode_ERROR_BUSINESS

	var message string
	if len(values) != 0 {
		first := values[0]
		switch v := first.(type) {
		case comm.ErrorCode: // 获取实际的ErrorCode
			code = v
			message = i18nErrMsg(ctx, code, values[1:]...)
		case string:
			message = v
		}
	}

	return &comm.ErrorInfo{
		Code:    int32(code),
		Message: message,
	}
}

func i18nErrMsg(ctx context.Context, errorCode comm.ErrorCode, fmtValues ...interface{}) string {
	errorCodeValue := int32(errorCode)
	var errMsg = errorCode.String()
	errCodeName, ok := comm.ErrorCode_name[errorCodeValue]

	if ok && ctx != nil {
		msgKey := "error." + errCodeName
		lang := i18n.GetLanguage(ctx)
		if lang != "" {
			if i18nMsg := i18n.GetMessage(lang, msgKey, fmtValues...); i18nMsg != "" {
				errMsg = i18nMsg
			}
		}
	}

	return errMsg
}

func IsDbRecordNotFoundErr(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
