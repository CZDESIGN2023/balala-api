package utils

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

/*
logger使用前言說明:
kratos的log分為三層
1. log lib (可彈性選擇, 目前選用zap)
2. kratos log.logger封裝log lib, 這一層可以附加額外的欄位訊息記錄到log內容中
3. kratos log.Helper封裝log.logger, 提供更簡化的Info, Warn...等方便的func

為了簡化說明, 上面1=zlog, 2=klog, 3=hlog
目前Service初始化後將klog封裝成hlog, 往後傳遞的log是hlog, 因為方便使用

但因為有附加訊息的需求, 所以Service之後要額外傳遞klog
以ChatService為例:
type ChatService struct {
	log *log.Helper
    klog  *log.Logger  // 新傳遞的klog
}
klog有提供With func, 將目前階段的klog及想附加的資訊傳入, 可取回新的klog(已附加訊息)
為了使用方便的Info...等func, 又需要將新的klog用log.Helper封裝, 這個過程就不太方便
所以有附加訊息需求時, 可透過以下提供的func簡化操作
--------------------------------------------------------------------------

附加自訂欄位補充說明:
main初始化會將目前專案的service.id及service.name寫入klog, 然後透過wire傳遞到各Service
此時klog有兩個欄位了

服務模塊初始化, 透過utils.InitModuleLogger, 傳入模組名稱
此時klog有三個欄位了service.id、service.name、module

再加上log鏈路追蹤的關鍵字是trace_id、span_id
所以klog中如果有自訂欄位和訊息要加入, 以下五個關鍵字不能使用
service.id、service.name、module、trace_id、span_id

另外還有UI工具定義的一些特殊欄位要避開
為了避免記憶困難, 傳入的自訂欄位時func會自動幫key加上module_前綴
例如:傳入自訂欄位data, 會被調整成module_data

*/

// InitModuleLogger 初始化各服務模塊klog, 附加module資訊
// 使用情境可參考:service.NewChatService
func InitModuleLogger(klog log.Logger, moduleName string) (*log.Logger, *log.Helper) {
	newKlog := log.With(klog, "module", moduleName)
	hlog := log.NewHelper(newKlog)
	return &newKlog, hlog
}

// GetLogger 將目前klog附加資訊, kv變數兩個一組, 可多組
// 例如:utils.GetLogger(klog, "field1", "value1", "field2", "value2")
//
// 重要業務log需求請優先使用GetTraceLogger取得log Helper, 以利錯誤追蹤
func GetLogger(klog log.Logger, kv ...interface{}) *log.Helper {
	prefixedKv := addPrefix(kv...)
	newKlog := log.With(klog, prefixedKv...)
	hlog := log.NewHelper(newKlog)
	return hlog
}

// GetTraceLogger 將目前klog附加資訊, kv變數兩個一組, 可多組
// 例如:utils.GetTraceLogger(ctx, klog, "field1", "value1", "field2", "value2")
//
// 會額外將trace_id和span_id附加到klog, 以利錯誤追蹤
func GetTraceLogger(ctx context.Context, klog log.Logger, kv ...interface{}) *log.Helper {
	prefixedKv := addPrefix(kv...)
	newKv := make([]interface{}, 0, len(prefixedKv)+4)
	newKv = append(newKv, prefixedKv...)
	newKv = append(newKv,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
	newKlog := log.With(klog, newKv...)
	hlog := log.NewHelper(newKlog)
	return hlog.WithContext(ctx)
}

// 幫傳入的自訂欄位key加入module_前綴
func addPrefix(kv ...interface{}) []interface{} {
	if len(kv) < 2 {
		return []interface{}{}
	}

	// 檢查 kv 長度是否為偶數，如果不是，則去掉最後一個元素
	if len(kv)%2 != 0 {
		kv = kv[:len(kv)-1]
	}
	result := make([]interface{}, 0, len(kv))
	// 為每個鍵添加 "module_" 前綴
	for i := 0; i < len(kv); i += 2 {
		if key, ok := kv[i].(string); ok {
			result = append(result, "module_"+key, kv[i+1])
		}
	}

	return result
}
