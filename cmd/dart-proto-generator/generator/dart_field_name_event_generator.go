package generator

import (
	parser2 "go-cs/pkg/parser"
	"os"
	"path/filepath"
	"text/template"
)

const fieldNameEventDartTemplate = `// Auto-generated Dart code. DO NOT MODIFY.
// Source: {{.InputFilePath}}

{{ $field_type := "int" }} {{ $field_defaultvalue := "0" }}{{range .Messages}}{{ $messageName := .Name }}
class {{.Name}}ProtoField {
  {{range .Fields}}static const String {{ToHump .Name}} = "{{.Name}}";
  {{end}}
}

class {{.Name}}Event {
  {{range .Fields}}static const String {{ToHump .Name}} = "c.{{.Name}}";
  {{end}}
}
{{end}}
`

func GenerateDartFieldNameEventFile(tmplParams *parser2.TemplateParams, outputFilePath string) error {
	funcMap := template.FuncMap{
		"ToUpperWithUnderscores": parser2.ToUpperWithUnderscores,
		"ToPascalCase":           parser2.ToPascalCase,
		"ToHump":                 parser2.ToHump,
	}

	tmpl, err := template.New("applyChangeDart").Funcs(funcMap).Parse(fieldNameEventDartTemplate)
	if err != nil {
		return err
	}

	// 计算 EnumTypePrefix
	tmplParams.EnumTypePrefix = parser2.ComputeEnumTypePrefix(tmplParams.Messages)

	// 创建输出文件所在的目录
	dir := filepath.Dir(outputFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755) // Creates directory with permissions set to 0755
		if err != nil {
			return err
		}
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, tmplParams)
	if err != nil {
		return err
	}

	return nil
}
