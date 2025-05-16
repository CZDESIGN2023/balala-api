package main

import (
	"fmt"
	"go-cs/cmd/generate_apply_update_block/generator"
	"go-cs/pkg/parser"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	InputProtoFilePath = "internal/bean/types.proto"
	//OutputEnumFilePath              = "internal/bean/types_field_enum.proto"
	//OutputGolangApplyChangeFilePath = "internal/bean/types_field_apply_change.go"
	//OutputDartApplyChangeFilePath   = ".output/generated/dart/pb/bean/types_field_apply_change.dart"
	//OutputCreateChangeFilePath      = "internal/bean/types_field_create_change.go"
	OutputGormTypeFilePath = "internal/bean/biz/types_gorm.go"
	// OutputGormQueryTypeFilePath     = "internal/bean/query/query_gorm.go"
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

	p := parser.NewParser()

	tmplParams, err := p.ParseProtoFile(InputProtoFilePath)
	if err != nil {
		fmt.Printf("解析Proto文件失败：%v\n", err)
		return
	}

	tmplParams.Time = time.Now()

	fmt.Printf("解析完成，共找到 %d 个消息。\n", len(tmplParams.Messages))

	//---------------------------------------------
	// fmt.Println("开始生成枚举文件...")
	// err = generator.GenerateEnumFile(tmplParams, OutputEnumFilePath)
	// if err != nil {
	// 	fmt.Printf("生成枚举文件失败：%v\n", err)
	// 	return
	// }
	// fmt.Printf("枚举文件生成完成，已保存至：%s。\n", OutputEnumFilePath)

	//---------------------------------------------
	// fmt.Println("开始生成CreateChange文件...")
	// err = generator.GenerateCreateChangeFile(tmplParams, OutputCreateChangeFilePath)
	// if err != nil {
	// 	fmt.Printf("生成CreateChange文件失败：%v\n", err)
	// 	return
	// }
	// fmtGoFile(OutputCreateChangeFilePath)
	// fmt.Printf("CreateChange文件生成完成，已保存至：%s。\n", OutputCreateChangeFilePath)

	//---------------------------------------------
	// fmt.Println("开始生成golang应用变更文件...")
	// err = generator.GenerateApplyChangeFile(tmplParams, OutputGolangApplyChangeFilePath)
	// if err != nil {
	// 	fmt.Printf("生成golang应用变更文件失败：%v\n", err)
	// 	return
	// }
	// fmtGoFile(OutputGolangApplyChangeFilePath)
	// fmt.Printf("golang应用变更文件生成完成，已保存至：%s。\n", OutputGolangApplyChangeFilePath)

	//---------------------------------------------
	fmt.Println("开始生成GORM类型文件...")
	err = generator.NewGormTypeGenerator().Generate(tmplParams, OutputGormTypeFilePath)
	if err != nil {
		fmt.Printf("生成GORM类型文件失败：%v\n", err)
		return
	}
	fmtGoFile(OutputGormTypeFilePath)
	fmt.Printf("GORM类型文件生成完成，已保存至：%s。\n", OutputGormTypeFilePath)
	//---------------------------------------------
	// fmt.Println("开始生成GORM QUERY类型文件...")
	// err = generator.NewGormQueryGenerator().Generate(tmplParams, OutputGormQueryTypeFilePath)
	// if err != nil {
	// 	fmt.Printf("生成GORM类型文件失败：%v\n", err)
	// 	return
	// }
	// fmtGoFile(OutputGormQueryTypeFilePath)
	// fmt.Printf("GORM类型文件生成完成，已保存至：%s。\n", OutputGormQueryTypeFilePath)
}

// 格式化go文件
func fmtGoFile(path string) {
	// 创建一个命令对象
	cmd := exec.Command("gofmt", "-w", path)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("gofmt error:%v\n", err)
	}
}
