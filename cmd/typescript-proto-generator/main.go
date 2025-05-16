package main

import (
	"fmt"
	"go-cs/cmd/typescript-proto-generator/generator"
	parser2 "go-cs/pkg/parser"
	"os"
	"path/filepath"
	"time"
)

const (
	InputProtoFilePath            = "internal/bean/types.proto"
	OutputDartFieldEventFilePath  = ".output/generated/typescript/pb/bean/types_field_event.ts"
	OutputDartApplyChangeFilePath = ".output/generated/typescript/pb/bean/types_field_apply_change.ts"
)

func printPwd() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to retrieve executable path:", err)
		return
	}

	toolDir := filepath.Dir(exePath)
	fmt.Println("Tool directory:", toolDir)
}

func main() {
	fmt.Println("开始解析过程...")

	printPwd()

	p := parser2.NewParser()

	tmplParams, err := p.ParseProtoFile(InputProtoFilePath)
	if err != nil {
		fmt.Printf("解析Proto文件失败：%v\n", err)
		return
	}

	tmplParams.Time = time.Now()

	fmt.Printf("解析完成，共找到 %d 个消息。\n", len(tmplParams.Messages))

	//---------------------------------------------
	fmt.Println("开始生成dart表事件文件...")
	err = generator.GenerateDartFieldNameEventFile(tmplParams, OutputDartFieldEventFilePath)
	if err != nil {
		fmt.Printf("生成dart表事件文件失败：%v\n", err)
		return
	}
	fmt.Printf("dart表事件文件生成完成，已保存至：%s。\n", OutputDartFieldEventFilePath)

	//---------------------------------------------
	fmt.Println("开始生成dart应用变更文件...")
	err = generator.GenerateDartApplyChangeFile(tmplParams, OutputDartApplyChangeFilePath)
	if err != nil {
		fmt.Printf("生成dart应用变更文件失败：%v\n", err)
		return
	}
	fmt.Printf("dart应用变更文件生成完成，已保存至：%s。\n", OutputDartApplyChangeFilePath)

	//---------------------------------------------
	// fmt.Println("开始生成dart api文件...")
	// 搜索api目录的proto文件
	// err = filepath.Walk("./api/", func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		fmt.Printf("阻止访问路径 %q: %v\n", path, err)
	// 		return err
	// 	}
	// 	// if !info.IsDir() && strings.HasSuffix(info.Name(), ".proto") {
	// 	// 	generatorApi(path)
	// 	// }
	// 	return nil
	// })
	if err != nil {
		fmt.Printf("搜索api目录的proto文件: %v\n", err)
	}
}

func generatorApi(protoFilePath string) {
	fmt.Printf("开始生成dart api: %v\n", protoFilePath)
	p := parser2.NewParser()
	tmplParams, err := p.ParseProtoFile(protoFilePath)
	if err != nil {
		fmt.Printf("解析Proto文件失败：%v\n", err)
		return
	}
	fmt.Printf("解析完成，共找到 %d 个api。\n", len(tmplParams.Services))
	if len(tmplParams.Services) == 0 {
		return
	}

	outFilePath := ".output/generated/dart/http/" + parser2.GetNameByPath(protoFilePath) + ".dart"
	err = generator.GenerateDartHttpApiFile(tmplParams, outFilePath)
	if err != nil {
		fmt.Printf("生成dart api文件失败：%v\n", err)
		return
	}
	fmt.Printf("dart api文件生成完成，已保存至：%s。\n", outFilePath)
}
