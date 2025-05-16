package utils

import (
	"context"
	"fmt"
	"go-cs/api/comm"
	"go-cs/pkg/i18n"
	"strconv"
)

type ICode interface {
	Error() string
	Code() int32
}

// NewErrorInfo 傳入錯誤碼及格式化資料，產生錯誤回應訊息(i18n)
// 使用fmt.Sprintf格式化訊息
// 訊息內容範例: "Hello %s, you are %d years old."
func NewErrorInfo(ctx context.Context, errorCode comm.ErrorCode, fmtValues ...interface{}) *comm.ErrorInfo {
	errMsg := ""
	errorCodeValue := int32(errorCode)
	errCodeName, ok := comm.ErrorCode_name[errorCodeValue]
	if ok {
		msgKey := "error." + errCodeName
		lang := i18n.GetLanguage(ctx)
		errMsg = i18n.GetMessage(lang, msgKey, fmtValues...)
		if errMsg == "" {
			// 找不到錯誤訊息時，回傳錯誤碼
			errMsg = strconv.Itoa(int(errorCodeValue))
		}
	} else {
		// 找不到錯誤碼定義名稱時，回傳錯誤碼
		errMsg = strconv.Itoa(int(errorCodeValue))
	}

	return &comm.ErrorInfo{
		Code:    errorCodeValue,
		Message: errMsg,
	}
}

// NewErrorInfoFmtData 傳入錯誤碼及格式化樣板資料，產生錯誤回應訊息(i18n)
// 使用i18n內建的資料樣板格式化訊息
// 訊息內容範例: "Hello {{.Name}}, you are {{.Age}} years old."
// 用法請查看i18n.GetMessageFmtData
func NewErrorInfoFmtData(ctx context.Context, errorCode comm.ErrorCode, fmtData interface{}) *comm.ErrorInfo {
	errMsg := ""
	errorCodeValue := int32(errorCode)
	errCodeName, ok := comm.ErrorCode_name[errorCodeValue]
	if ok {
		msgKey := "error." + errCodeName
		lang := i18n.GetLanguage(ctx)
		errMsg = i18n.GetMessageFmtData(lang, msgKey, fmtData)
		if errMsg == "" {
			// 找不到錯誤訊息時，回傳錯誤碼
			errMsg = strconv.Itoa(int(errorCodeValue))
		}
	} else {
		// 找不到錯誤碼定義名稱時，回傳錯誤碼
		errMsg = strconv.Itoa(int(errorCodeValue))
	}

	return &comm.ErrorInfo{
		Code:    errorCodeValue,
		Message: errMsg,
	}
}

func NewCommonErrorReply(ctx context.Context, errorCode comm.ErrorCode, fmtValues ...interface{}) *comm.CommonReply {
	errorReply := &comm.CommonReply{Result: &comm.CommonReply_Error{Error: NewErrorInfo(ctx, errorCode, fmtValues...)}}
	return errorReply
}

func NewCommonOkReply(ctx context.Context) *comm.CommonReply {
	okReply := &comm.CommonReply{Result: &comm.CommonReply_Data{Data: ""}}
	return okReply
}

// NewAppErrorCode 傳入錯誤碼定義，產生完整錯誤碼
// 总的6位 xyyzzz, app-service-num
func NewAppErrorCode(appCode int, serviceCode int, num int) (int, error) {
	strX := fmt.Sprintf("%01d", appCode)
	strYY := fmt.Sprintf("%02d", serviceCode)
	strZZZ := fmt.Sprintf("%03d", num)

	merge := strX + strYY + strZZZ
	appErrorCode, err := strconv.Atoi(merge)
	if err != nil {
		return 0, err
	}
	return appErrorCode, nil
}
