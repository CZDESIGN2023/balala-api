package parser

import (
	"fmt"
	"time"

	"github.com/jhump/protoreflect/desc/protoparse"

	"github.com/golang/protobuf/proto"
	googleapi "google.golang.org/genproto/googleapis/api/annotations"
)

const DbTmplGorm = "gorm"
const DbTmplBson = "bson"

// TemplateParams 表示模板参数
type TemplateParams struct {
	InputFilePath  string
	Time           time.Time
	Messages       []Message
	Services       []Service
	GoPackageName  string
	EnumTypePrefix string // 添加 EnumTypePrefix 字段，表示枚举类型的前缀
}

// Parser 是解析器
type Parser struct {
}

// NewParser 创建一个新的解析器
func NewParser() *Parser {
	return &Parser{}
}

// ParseProtoFile 解析 Proto 文件
func (p *Parser) ParseProtoFile(fileName string) (*TemplateParams, error) {
	parser := protoparse.Parser{
		ImportPaths: []string{
			"",
			"api",
			"internal",
			"third_party",
		},
		IncludeSourceCodeInfo: true,
	}
	fds, err := parser.ParseFiles(fileName)
	if err != nil {
		return nil, err
	}

	params := &TemplateParams{
		InputFilePath: fileName,
		Time:          time.Now(),
		GoPackageName: "bean",
	}

	for _, fd := range fds {
		for _, sd := range fd.GetServices() {
			service := Service{Name: sd.GetName()}
			for _, smd := range sd.GetMethods() {
				method := Method{Name: smd.GetName(), InputType: smd.GetInputType().GetName(), OutputType: smd.GetOutputType().GetName()}

				if ext, err := proto.GetExtension(smd.GetOptions(), googleapi.E_Http); err == nil {
					// 从 httpOptions 中提取所需的值
					method.PostPath = makeHttpPath(ext.(*googleapi.HttpRule))
				} else {
					// 处理无法解析扩展的错误
					fmt.Printf("方法：%v 处理无法解析扩展的错误: %v ", method.Name, err)
				}

				service.Methods = append(service.Methods, method)
			}
			params.Services = append(params.Services, service)
		}
		for _, md := range fd.GetMessageTypes() {
			msg := Message{
				Name:       md.GetName(),
				DbTmplType: parseTmplType(md),
				Comments:   parseTmplComment(md),
			}

			for i, fd := range md.GetFields() {
				// 因为 id 肯定不更新，并且是第一个，可以支持跳过
				enumValue := i //+ 1
				//msg.Fields[i].FieldType = fd.GetType()
				goType := protoTypeToGoType(fd.GetType().String()) // 转换为 Go 类型
				camelName := toCamel(fd.GetName())
				field := Field{
					Name:       fd.GetName(),
					CamelName:  camelName,
					GoType:     goType,
					EnumValue:  enumValue,
					FieldType:  fd.GetType().String(),
					IsRepeated: fd.IsRepeated(),
					Comments:   fd.GetSourceInfo().GetTrailingComments(),
				}
				fd.GetJSONName()
				msg.Fields = append(msg.Fields, field)
			}
			params.Messages = append(params.Messages, msg)
		}
	}

	return params, nil
}
