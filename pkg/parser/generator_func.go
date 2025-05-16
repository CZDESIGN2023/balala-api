package parser

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// ToPascalCase 将字符串转换为 PascalCase
func ToPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool { return r == '_' })
	for i := 0; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

// ToHump 将字符串转换为 驼峰
func ToHump(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool { return r == '_' })
	for i := 0; i < len(words); i++ {
		if i > 0 {
			words[i] = strings.Title(words[i])
		} else {
			words[i] = lowercaseFirstLetter(words[i])
		}
	}
	return strings.Join(words, "")
}

func lowercaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	// 获取字符串的第一个字符
	firstChar := s[0]
	// 将第一个字符转换为小写
	lowerFirstChar := unicode.ToLower(rune(firstChar))
	// 构建转换后的字符串
	lowercasedStr := string(lowerFirstChar) + s[1:]
	return lowercasedStr
}

// ToUpperWithUnderscores 将字符串转换为大写，并用下划线分隔单词
func ToUpperWithUnderscores(s string) string {
	var buffer bytes.Buffer
	for i, char := range s {
		if i > 0 && (unicode.IsUpper(char) || unicode.IsDigit(char)) {
			buffer.WriteRune('_')
		}
		buffer.WriteRune(unicode.ToUpper(char))
	}
	return buffer.String()
}

func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func ShortName(name string) string {
	// 这里根据需要实现缩写逻辑，例如提取首字母或使用其他规则
	// 这里只是一个示例
	return strings.ToLower(string(name[0]))
}

func ToDBName(name string) string {
	return strings.ToLower(ToSnakeCase(name))
}

func FormatEnumName(s string) string {
	s = strings.ToUpper(s)
	return strings.Replace(s, ".", "_", -1)
}

// 计算 EnumTypePrefix
func ComputeEnumTypePrefix(messages []Message) string {
	// 假设枚举类型前缀与第一个消息的名称相同
	if len(messages) > 0 {
		return messages[0].Name
	}
	return ""
}

// GoType 将 Protobuf 类型转换为 Go 类型
func GoType(s string) string {
	switch s {
	case "TYPE_INT32", "TYPE_SINT32":
		return "int32"
	case "TYPE_UINT32":
		return "uint32"
	case "TYPE_INT64", "TYPE_SINT64":
		return "int64"
	case "TYPE_UINT64":
		return "uint64"
	case "TYPE_FLOAT":
		return "float32"
	case "TYPE_DOUBLE":
		return "float64"
	case "TYPE_BOOL":
		return "bool"
	case "TYPE_STRING":
		return "string"
	default:
		return "interface{}"
	}
}

// FieldValueName 返回对应的 FieldValue 中的字段名
func FieldValueName(fieldType string) string {
	switch fieldType {
	case "TYPE_INT32", "TYPE_SINT32":
		return "Int32Value"
	case "TYPE_INT64", "TYPE_SINT64":
		return "Int64Value"
	case "TYPE_UINT32":
		return "Uint32Value"
	case "TYPE_UINT64":
		return "Uint64Value"
	case "TYPE_FLOAT":
		return "FloatValue"
	case "TYPE_DOUBLE":
		return "DoubleValue"
	case "TYPE_BOOL":
		return "BoolValue"
	case "TYPE_STRING":
		return "Str"
	default:
		return "Unknown"
	}
}

// FieldValueStructName 返回对应的 FieldValue 中的结构体名
func FieldValueStructName(fieldType string) string {
	switch fieldType {
	case "TYPE_INT32", "TYPE_SINT32":
		return "FieldValue_Int32Value"
	case "TYPE_INT64", "TYPE_SINT64":
		return "FieldValue_Int64Value"
	case "TYPE_UINT32":
		return "FieldValue_Uint32Value"
	case "TYPE_UINT64":
		return "FieldValue_Uint64Value"
	case "TYPE_FLOAT":
		return "FieldValue_FloatValue"
	case "TYPE_DOUBLE":
		return "FieldValue_DoubleValue"
	case "TYPE_BOOL":
		return "FieldValue_BoolValue"
	case "TYPE_STRING":
		return "FieldValue_Str"
	default:
		return "FieldValue_Unknown"
	}
}

func GetNameByPath(path string) string {
	fileNameWithExt := filepath.Base(path)
	extension := filepath.Ext(fileNameWithExt)
	fileName := strings.TrimSuffix(fileNameWithExt, extension)
	return fileName
}

func GetTitleNameByPath(path string) string {
	return strings.Title(ToHump(GetNameByPath(path)))
}

func GetImportGrpcPathByPorto(path string) string {
	path = strings.Replace(path, "api/", "../pb/", 1)
	path = strings.Replace(path, ".proto", ".pbgrpc.dart", 1)
	return path
}

func GetDartApiPath(path string) string {
	re := regexp.MustCompile(`\{[^{}]+\}`)

	result := re.ReplaceAllStringFunc(path, func(match string) string {
		// 获取 {} 内部的内容
		content := match[1 : len(match)-1]
		// 进行相应的处理
		// 这里可以根据需要进行替换或其他操作
		// 这里将 {} 内部的内容转换为大写
		return "${data." + ToHump(content) + "}"
	})
	return result
}
