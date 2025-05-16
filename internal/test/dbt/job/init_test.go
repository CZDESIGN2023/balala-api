package main

import "go-cs/internal/test/dbt"

func init() {
	// 测试执行目录为当前目录，所以项目配置文件路径这么奇怪
	dbt.Init("../../../../configs/config.yaml", "../../../../.env.local", true)
}
