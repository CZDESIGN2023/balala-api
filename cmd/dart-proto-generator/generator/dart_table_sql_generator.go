package generator

import (
	parser2 "go-cs/pkg/parser"
	"os"
	"path/filepath"
	"text/template"
)

const tableSqlDartTemplate = `// Auto-generated Dart code. DO NOT MODIFY.
// Source: {{.InputFilePath}}
class TableSql {
{{ $field_type := "int" }} {{ $field_defaultvalue := "0" }}{{range .Messages}}{{ $messageName := .Name }}
  static const String {{ToHump .Name}}TableName = '{{ ToDBName .Name }}';
  static const String {{ToHump .Name}}Sql = '''
  	CREATE TABLE IF NOT EXISTS {{ToHump .Name}} (
		id INTEGER PRIMARY KEY,
		{{range .Fields}}{{- if ne .Name "id"}} 
		{{- if or (eq .FieldType "TYPE_SINT64") (eq .FieldType "TYPE_INT64") (eq .FieldType "TYPE_UINT64")}}
		{{ $field_type = "INTEGER" }}
		{{- else if or (eq .FieldType "TYPE_SINT32") (eq .FieldType "TYPE_INT32") (eq .FieldType "TYPE_UINT32")}}
		{{ $field_type = "INTEGER" }}
		{{- else if or (eq .FieldType "TYPE_FLOAT") (eq .FieldType "TYPE_DOUBLE") }}
		{{ $field_type = "REAL" }}
		{{- else if eq .FieldType "TYPE_BOOL" }}
		{{ $field_type = "INTEGER" }}
		{{- else if eq .FieldType "TYPE_STRING" }}
		{{ $field_type = "TEXT" }}
		{{- end }}
		{{.Name}} {{ $field_type }},
		{{- end }}
		{{end}}
		__add_index INTEGER
	);
  ''';
{{end}}
}
`

func GenerateDartTableSqlFile(tmplParams *parser2.TemplateParams, outputFilePath string) error {
	funcMap := template.FuncMap{
		"ToUpperWithUnderscores": parser2.ToUpperWithUnderscores,
		"ToPascalCase":           parser2.ToPascalCase,
		"ToHump":                 parser2.ToHump,
		"ToDBName":               parser2.ToDBName,
	}

	tmpl, err := template.New("applyChangeDart").Funcs(funcMap).Parse(tableSqlDartTemplate)
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
