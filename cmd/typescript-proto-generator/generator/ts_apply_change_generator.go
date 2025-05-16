package generator

import (
	parser2 "go-cs/pkg/parser"
	"os"
	"path/filepath"
	"text/template"
)

const applyChangeDartTemplate = `// Auto-generated Typescript code. DO NOT MODIFY.
// Source: {{.InputFilePath}}

import {UpdateData, FieldValue, AnyObj} from "./request";
import * as module from "./types_field_enum";
import * as BaseType from "./types";
import BaseObject from "../base_object";

// 把有值的位放在一个数组里
function _getSetBits(num:number, setBits?:Array<number>,  startIdx?:number):Array<number> {
  setBits ??= [];
  let index:number = startIdx??0;
  while (num > 0) {
    if ((num & 1) == 1) {
      setBits.push(index)
    }
    index++;
    num = num >> 1;
  }
  return setBits;
}

// 解析掩码
function parseFieldMasks(masks:Uint8Array):Array<number> {
  let setBits:Array<number> = [];
  for (var i = 0; i < masks.length; i++) {
    _getSetBits(masks[i], setBits, i > 0 ? i * 8 : 0);
  }
  return setBits;
}

{{range .Messages}}
{{ $messageName := .Name }}
export class {{.Name}}Extend extends BaseObject<BaseType.{{.Name}}>{

  constructor(u8a?:Uint8Array, json?:string) {
    if(u8a != undefined)
      super(BaseType.{{.Name}}.fromBinary(u8a));
    else if(json != undefined)
      super(BaseType.{{.Name}}.fromJsonString(json));
  }

  applyUpdateblok(update:UpdateData) {
    // let data:Map<String, any> = new Map<String, any>();
    let fieldMasks = parseFieldMasks(update.masks);
    const values = update.values;
    for (var i = 0; i < fieldMasks.length; i++) {
      const fieldIndex = fieldMasks[i];
      const fieldValue = values[i];
      this._applyField(fieldIndex, fieldValue);
    }
  }

  _applyField(fieldIndex:number, fieldValue:FieldValue) {
    switch (fieldIndex) {
      {{range .Fields}}{{- if ne .Name "id"}}
      case module.{{$messageName}}Field.{{ToUpperWithUnderscores $messageName}}_{{ToUpperWithUnderscores .Name}}:
        {{- if or (eq .FieldType "TYPE_SINT64") (eq .FieldType "TYPE_INT64") }}
        this.setValue('{{ToHump .Name }}', (fieldValue.value as any).int64Value);
        {{- else if eq .FieldType "TYPE_UINT64" }}
        this.setValue('{{ToHump .Name }}', (fieldValue.value as any).uint64Value);
        {{- else if or (eq .FieldType "TYPE_SINT32") (eq .FieldType "TYPE_INT32") }}
        this.setValue('{{ToHump .Name }}', (fieldValue.value as any).int32Value);
        {{- else if eq .FieldType "TYPE_UINT32" }}
        this.setValue('{{ToHump .Name }}', (fieldValue.value as any).uint32Value);
        {{- else if eq .FieldType "TYPE_FLOAT" }}
        this.setValue('{{ToHump .Name }}', (fieldValue.value as any).floatValue);
        {{- else if eq .FieldType "TYPE_DOUBLE" }}
        this.setValue('{{ToHump .Name }}', (fieldValue.value as any).doubleValue);
        {{- else if eq .FieldType "TYPE_BOOL" }}
        this.setValue('{{ToHump .Name }}', (fieldValue.value as any).boolValue);
        {{- else if eq .FieldType "TYPE_STRING" }}
        this.setValue('{{ToHump .Name }}', (fieldValue.value as any).str);
        {{- else }}
        throw Error('Unknown type [{{.FieldType}}] for field {{.Name}}');
        {{- end }}
        break;{{- end}}{{end}}
      default:
        throw Error('Unknown field index: ' + fieldIndex);
    }
  }

}


{{end}}

// export function parseBlock(classIndex:number, update:UpdateData):Map<String, any>{
//   switch(classIndex){
//     {{range .Messages}}
//     {{- $messageName := .Name }}
//     case BaseType.Table.{{ToUpperWithUnderscores $messageName}}:
//       return {{$messageName}}BlockParser.parserUpdateblok(update)
//     {{- end}}
//   default:
//     return new Map<String, any>();
//   }
// }

export function createBlock(classIndex:number, update:AnyObj):any{
  switch(classIndex){
    {{range .Messages}}
    {{- $messageName := .Name }}
    case BaseType.Table.{{ToUpperWithUnderscores $messageName}}:
      return BaseType.{{$messageName}}.fromBinary(update.obj)
    {{- end}}
  default:
    return new Map<String, any>();
  }
}


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
