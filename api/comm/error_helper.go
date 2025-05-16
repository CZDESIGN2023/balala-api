package comm

// proto 无法自动生成 Error()string 方法
// 手动实现 error 接口
func (e *ErrorInfo) Error() string {
	return e.Message
}
