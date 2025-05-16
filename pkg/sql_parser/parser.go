package sql_parser

import (
	"regexp"
	"strings"
)

// 单行注释
var singleLineComment = regexp.MustCompile(`(?m)^\s*--.*`)

// 多行注释
var multiLineComment = regexp.MustCompile(`(?m)^\s*/\*.*?\*/;?`)

// 语句分割
var split = regexp.MustCompile(`(?m);\s*$`)

// ToStatements 转换为SQL语句
func ToStatements(content string) []string {
	// 去除单行和多行注释
	cleanContent := singleLineComment.ReplaceAllString(content, "")
	cleanContent = multiLineComment.ReplaceAllString(cleanContent, "")

	// 分割SQL语句
	statements := split.Split(cleanContent, -1)

	// 去除空语句
	statements = filterEmptyStatements(statements)

	return statements
}

// filterEmptyStatements 去除空语句
func filterEmptyStatements(statements []string) []string {
	var filteredStatements []string
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement != "" {
			filteredStatements = append(filteredStatements, statement)
		}
	}
	return filteredStatements
}
