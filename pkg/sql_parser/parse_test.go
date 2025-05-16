package sql_parser

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed dml.sql
var initSql string

func Test(t *testing.T) {

	statements := ToStatements(initSql)

	// 打印提取的SQL语句
	for _, statement := range statements {
		fmt.Println(statement)
	}
}
