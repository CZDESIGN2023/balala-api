package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go-cs/internal/utils"
	"reflect"
)

type OpenAPIService struct {
	//之后看需要被转发的有哪些服务再增加 先留一个范例
	//Test   *service.TestService
	log *log.Helper
}

func NewOpenAPIService(logger log.Logger) *OpenAPIService {
	moduleName := "OpenAPIService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &OpenAPIService{
		log: hlog,
	}
}

type Result struct {
	//推送過來固定的格式
	Ctl  string      `json:"ctl"`
	Act  string      `json:"act"`
	Data interface{} `json:"data"`
}

func (api *OpenAPIService) RouteToOther(ctx context.Context, args interface{}) {
	//先解析傳來的參數
	values := reflect.ValueOf(args).String()
	var resJson Result
	//塞进固定的格式内
	err := json.Unmarshal([]byte(values), &resJson)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		return
	}
	//獲取對應服務
	serviceValue := reflect.ValueOf(*api).FieldByName(resJson.Ctl)
	if !serviceValue.IsValid() {
		fmt.Println(serviceValue)
		fmt.Println("Service not found")
		return
	}
	//獲取對應函數
	funcValue := serviceValue.MethodByName(resJson.Act)
	if !funcValue.IsValid() {
		fmt.Println("Function not found")
		return
	}
	//获取要调用函数的参数类型
	funcType := funcValue.Type()
	paramType := funcType.In(1)
	//创建一个空的参数值
	argValue := reflect.New(paramType.Elem())
	//将 resJson.Data 转换为byte
	dataBytes, err := json.Marshal(resJson.Data)
	if err != nil {
		fmt.Println("Failed to convert data to []byte:", err)
		return
	}
	//将 resJson.Data 转换为参数类型
	err = json.Unmarshal(dataBytes, argValue.Interface())
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		return
	}
	// 构造参数列表
	var value []reflect.Value
	value = append(value, reflect.ValueOf(ctx))
	value = append(value, argValue)
	// 调用函数
	funcValue.Call(value)
	//TODO 错误及回传处理
}
