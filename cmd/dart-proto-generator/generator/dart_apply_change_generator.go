package generator

import (
	parser2 "go-cs/pkg/parser"
	"os"
	"path/filepath"
	"text/template"
)

const applyChangeDartTemplate = `// Auto-generated Dart code. DO NOT MODIFY.
// Source: {{.InputFilePath}}

import 'package:common_base/common_base.dart';
import 'package:data_center/pb/bean/types.pb.dart' as types;
import 'package:data_center/pb/bean/types.pbenum.dart';

import 'request.pb.dart';
import 'types_field_enum.pb.dart';

abstract class _EventRowObjectMixin {
  // 设置键值
  void setValue(String key, dynamic value, {bool dispatcherEvent = true});
  // 根据key获得键值
  T getValue<T>(String? key, [dynamic def]);
}

{{ $field_type := "int" }} {{ $field_defaultvalue := "0" }}{{range .Messages}}{{ $messageName := .Name }}
abstract class {{.Name}}Mixin implements _EventRowObjectMixin {
  // 字段get set
  {{range .Fields}}{{- if ne .Name "id"}}
  {{- if or (eq .FieldType "TYPE_SINT64") (eq .FieldType "TYPE_INT64") (eq .FieldType "TYPE_UINT64")}}
  {{ $field_type = "int" }}{{ $field_defaultvalue = "0" }}
  {{- else if or (eq .FieldType "TYPE_SINT32") (eq .FieldType "TYPE_INT32") (eq .FieldType "TYPE_UINT32")}}
  {{ $field_type = "int" }}{{ $field_defaultvalue = "0" }}
  {{- else if or (eq .FieldType "TYPE_FLOAT") (eq .FieldType "TYPE_DOUBLE") }}
  {{ $field_type = "double" }}{{ $field_defaultvalue = "0" }}
  {{- else if eq .FieldType "TYPE_BOOL" }}
  {{ $field_type = "bool" }}{{ $field_defaultvalue = "false" }}
  {{- else if eq .FieldType "TYPE_STRING" }}
  {{ $field_type = "String" }}{{ $field_defaultvalue = "''" }}
  {{- end }}
  {{$field_type}} get {{ToHump .Name}} => getValue('{{ .Name }}', {{$field_defaultvalue}});
  set {{ToHump .Name}}({{$field_type}} value) {
    setValue('{{ .Name }}', value);
  }
  {{- end }}
  {{end}}
}
class {{.Name}} extends EventRowObject with {{.Name}}Mixin {}
{{end}}


// 表定义
class TableMapping {
  String name;
  Function createFunc;
  Map<String, dynamic> Function(UpdateData update) updateFunc;
  TableMapping._(this.name, this.createFunc, this.updateFunc);
}

// 表映射
Map<int, TableMapping> updateblokTableMappings = {
  {{range .Messages}}{{ $messageName := .Name }}
  Table.{{ ToUpperWithUnderscores .Name }}.value :TableMapping._('{{ ToDBName .Name }}', types.{{ .Name }}.create, {{ .Name }}BlockParser.parserUpdateblok),
  {{end}}
};

// 把有值的位放在一个数组里
List<int> _getSetBits(int num, {List<int>? setBits, int startIdx = 0}) {
  setBits ??= [];
  int index = startIdx;
  while (num > 0) {
    if ((num & 1) == 1) {
      setBits.add(index);
    }
    index++;
    num = num >> 1;
  }
  return setBits;
}

// 解析掩码
List<int> parseFieldMasks(List<int> masks) {
  List<int> setBits = [];
  for (var i = 0; i < masks.length; i++) {
    _getSetBits(masks[i], setBits: setBits, startIdx: i > 0 ? i * 32 - 1 : 0);
  }
  return setBits;
}

{{range .Messages}}
{{ $messageName := .Name }}
class {{.Name}}BlockParser {
  static Map<String, dynamic> parserUpdateblok(UpdateData update) {
    Map<String, dynamic>  data = {};
    final fieldMasks = parseFieldMasks(update.masks);
    final values = update.values;
    for (var i = 0; i < fieldMasks.length; i++) {
      final fieldIndex = fieldMasks[i];
      final fieldValue = values[i];
      _parserField(data, fieldIndex, fieldValue);
    }
    return data;
  }

  static void _parserField(Map<String, dynamic> target, int fieldIndex, FieldValue fieldValue) {
    switch ({{$messageName}}Field.values[fieldIndex]) {
      {{range .Fields}}{{- if ne .Name "id"}}
      case {{$messageName}}Field.{{ToUpperWithUnderscores $messageName}}_{{ToUpperWithUnderscores .Name}}:
        {{- if or (eq .FieldType "TYPE_SINT64") (eq .FieldType "TYPE_INT64") }}
        target['{{ .Name }}'] = fieldValue.int64Value.toInt();
        {{- else if eq .FieldType "TYPE_UINT64" }}
        target['{{ .Name }}'] = fieldValue.uint64Value.toInt();
        {{- else if or (eq .FieldType "TYPE_SINT32") (eq .FieldType "TYPE_INT32") }}
        target['{{ .Name }}'] = fieldValue.int32Value;
        {{- else if eq .FieldType "TYPE_UINT32" }}
        target['{{ .Name }}'] = fieldValue.uint32Value;
        {{- else if eq .FieldType "TYPE_FLOAT" }}
        target['{{ .Name }}'] = fieldValue.floatValue;
        {{- else if eq .FieldType "TYPE_DOUBLE" }}
        target['{{ .Name }}'] = fieldValue.doubleValue;
        {{- else if eq .FieldType "TYPE_BOOL" }}
        target['{{ .Name }}'] = fieldValue.boolValue;
        {{- else if eq .FieldType "TYPE_STRING" }}
        target['{{ .Name }}'] = fieldValue.str;
        {{- else }}
        throw ArgumentError('Unknown type [{{.FieldType}}] for field {{.Name}}');
        {{- end }}
        break;{{- end}}{{end}}
      default:
        throw ArgumentError('Unknown field index: $fieldIndex');
    }
  }
}
{{end}}

`

func GenerateDartApplyChangeFile(tmplParams *parser2.TemplateParams, outputFilePath string) error {
	funcMap := template.FuncMap{
		"ToUpperWithUnderscores": parser2.ToUpperWithUnderscores,
		"ToPascalCase":           parser2.ToPascalCase,
		"ToHump":                 parser2.ToHump,
		"ToDBName":               parser2.ToDBName,
	}

	tmpl, err := template.New("applyChangeDart").Funcs(funcMap).Parse(applyChangeDartTemplate)
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
