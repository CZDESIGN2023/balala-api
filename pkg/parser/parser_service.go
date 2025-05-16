package parser

import googleapi "google.golang.org/genproto/googleapis/api/annotations"

// MethodPath 表示方法的路径规则
type MethodPath struct {
	Method       string
	FormatString string
	Params       map[string]string
}

func NewMethodPath(method string, fs string) {

}

// Method 表示服务中的方法
type Method struct {
	Name       string
	InputType  string
	OutputType string
	PostPath   []*MethodPath
}

// Service 表示协议服务
type Service struct {
	Name    string
	Methods []Method
}

func makeHttpPath(httpOptions *googleapi.HttpRule) []*MethodPath {
	var paths []*MethodPath

	// 处理 put 路径
	putPath := httpOptions.GetPut()
	if putPath != "" {
		paths = append(paths, &MethodPath{
			Method:       "PUT",
			FormatString: putPath,
		})
	}

	// 处理 post 路径
	postPath := httpOptions.GetPost()
	if postPath != "" {
		paths = append(paths, &MethodPath{
			Method:       "POST",
			FormatString: postPath,
		})
	}

	// 处理 delete 路径
	deletePath := httpOptions.GetDelete()
	if deletePath != "" {
		paths = append(paths, &MethodPath{
			Method:       "DELETE",
			FormatString: deletePath,
		})
	}

	// 处理 get 路径
	getPath := httpOptions.GetGet()
	if getPath != "" {
		paths = append(paths, &MethodPath{
			Method:       "GET",
			FormatString: getPath,
		})
	}

	// 处理 patch 路径
	patchPath := httpOptions.GetPatch()
	if patchPath != "" {
		paths = append(paths, &MethodPath{
			Method:       "PATCH",
			FormatString: patchPath,
		})
	}

	// 处理 custom 路径
	customPath := httpOptions.GetCustom()
	if customPath != nil {
		paths = append(paths, &MethodPath{
			Method:       customPath.Kind,
			FormatString: customPath.Path,
		})
	}

	// 进一步处理路径值，提取路径参数
	for _, methodPath := range paths {
		methodPath.Params = extractParams(methodPath.FormatString)
	}

	return paths
}

// extractParams 提取路径中的参数
func extractParams(formatString string) map[string]string {
	params := make(map[string]string)

	// 解析格式字符串中的路径参数
	// 这里假设路径参数以花括号 {} 包围，例如 "/chat/{typ}/{user_id}/{friend_id}/{group_id}"
	// 可根据实际情况进行修改
	// 这里只是一个示例实现，可能需要根据实际需求进行修改
	paramStart := false
	paramName := ""
	for _, char := range formatString {
		switch char {
		case '{':
			paramStart = true
		case '}':
			if paramName != "" {
				params[paramName] = ""
				paramName = ""
			}
			paramStart = false
		default:
			if paramStart {
				paramName += string(char)
			}
		}
	}

	return params
}
