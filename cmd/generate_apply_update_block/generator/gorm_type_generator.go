package generator

import (
	"fmt"
	parser2 "go-cs/pkg/parser"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type GormTypeGenerator struct {
}

func NewGormTypeGenerator() *GormTypeGenerator {
	return &GormTypeGenerator{}
}

func (g *GormTypeGenerator) Generate(params *parser2.TemplateParams, fileName string) error {
	// 跟输出文件同级目录的, 除了扩展名不一样以外
	tmplFilePath := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".tmpl"
	templateFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("unable to create template file: %v", err)
	}
	defer templateFile.Close()

	// 生成模板需要的函数
	funcMap := template.FuncMap{
		"ToSnakeCase":            parser2.ToSnakeCase,
		"ShortName":              parser2.ShortName,
		"ToDBName":               parser2.ToDBName,
		"ToUpperWithUnderscores": parser2.ToUpperWithUnderscores,
		"ToHump":                 parser2.ToHump,
	}

	// 绑定函数和输入变量
	t := template.Must(template.New(filepath.Base(tmplFilePath)).Funcs(funcMap).ParseFiles(tmplFilePath))
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}
	defer f.Close()

	// Filter out messages without an 'id' field and convert proto types to Go types
	originalMessages := params.Messages
	params.Messages = []parser2.Message{}
	for _, msg := range originalMessages {
		//if msg.DbTmplType != "gorm" {
		//	continue
		//}

		// hasIDField := false
		// for _, field := range msg.Fields {
		// 	if strings.ToLower(field.Name) == "id" {
		// 		hasIDField = true
		// 	}
		// }
		// if hasIDField {
		params.Messages = append(params.Messages, msg)
		// }
	}

	err = t.Execute(f, params)
	if err != nil {
		return fmt.Errorf("unable to execute template: %v", err)
	}

	// Restore the original messages
	params.Messages = originalMessages

	return nil
}
