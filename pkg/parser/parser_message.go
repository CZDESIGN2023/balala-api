package parser

import (
	"strings"

	"github.com/jhump/protoreflect/desc"
)

// Field 表示消息中的字段
type Field struct {
	Name       string
	EnumValue  int
	FieldType  string
	CamelName  string // 添加 CamelName 字段，表示转换为驼峰命名的字段名
	GoType     string // 添加 GoType 字段，表示转换为 Go 类型的字段类型
	IsRepeated bool
	Comments   string //注释信息
}

// Message 表示协议消息
type Message struct {
	Name       string
	Fields     []Field
	Comments   string //注释信息
	DbTmplType string // 記錄DB的類型應該使用哪一種tmpl輸出(gorm、bson)
}

// protoTypeToGoType 将 Protobuf 类型转换为 Go 类型
func protoTypeToGoType(protoType string) string {
	switch protoType {
	case "TYPE_DOUBLE":
		return "float64"
	case "TYPE_FLOAT":
		return "float32"
	case "TYPE_INT64", "TYPE_SINT64":
		return "int64"
	case "TYPE_UINT64":
		return "uint64"
	case "TYPE_UINT32":
		return "uint32"
	case "", "TYPE_INT32", "TYPE_SINT32", "TYPE_FIXED32", "TYPE_SFIXED32":
		return "int32"
	case "TYPE_BOOL":
		return "bool"
	case "TYPE_STRING":
		return "string"
	case "TYPE_BYTES":
		return "[]byte"
	default:
		return "string"
	}
}

// toCamel 将下划线命名转换为驼峰命名
func toCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func parseTmplType(md *desc.MessageDescriptor) string {
	sourceInfo := md.GetSourceInfo()
	if sourceInfo == nil {
		return "mysql"
	}
	comment := sourceInfo.GetLeadingComments()
	if comment != "" {
		lowerComment := strings.ToLower(comment)
		if strings.Index(lowerComment, "[mongodb]") != -1 {
			return DbTmplBson
		}
	}
	return DbTmplGorm
}

func parseTmplComment(md *desc.MessageDescriptor) string {
	sourceInfo := md.GetSourceInfo()
	if sourceInfo == nil {
		return ""
	}
	comment := sourceInfo.GetLeadingComments()
	if comment != "" {
		lowerComment := strings.ToLower(comment)
		return lowerComment
	}
	return comment
}
