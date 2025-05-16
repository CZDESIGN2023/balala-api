package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
	"unicode/utf8"
)

func SubStrRunesWithEllipsis(s string, length int) string {
	ns, isSub := SubStrRunes(s, length)

	//regexp将所有空白字符替换为空格
	ns = ReplaceAllWhitespaceWithSpace(ns)

	if isSub {
		return ns + "..."
	}
	return ns
}

func ClearRichTextToPlanText(s string, escape bool) string {
	s = ClearAndReplaceHtmlTag(s)
	//regexp将所有空白字符替换为空格
	s = ReplaceAllWhitespaceWithSpace(s)

	// 将html中的特殊符号转换的html的转义文本
	if escape {
		s = EscapeHtmlSpecialCharacters(s)
	}

	return SubStrRunesWithEllipsis(s, 200)
}

func SubStrRunes(s string, length int) (string, bool) {
	if utf8.RuneCountInString(s) > length {
		rs := []rune(s)
		return string(rs[:length]), true
	}
	return s, false
}

func ParseWorkingDay(weekDays []int64) string {
	var weekDay = []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}

	var list []string
	for _, day := range weekDays {
		list = append(list, weekDay[day])
	}

	return strings.Join(list, "、")
}

func ClearAndReplaceHtmlTag(data string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	doc.Find(".file-video").ReplaceWithHtml("[视频]")
	doc.Find(".file-audio").ReplaceWithHtml("[音频]")
	doc.Find(".file-image").ReplaceWithHtml("[图片]")
	doc.Find(".inserted-file").ReplaceWithHtml("[文件]")
	doc.Find(".code-block-wrapper").ReplaceWithHtml("[代码块]")
	doc.Find("img").ReplaceWithHtml("[图片]")
	doc.Find("code").ReplaceWithHtml("[代码]")

	return doc.Text()
}

// EscapeSqlSpecialCharacters 转译sql特殊字符
func EscapeSqlSpecialCharacters(s string) string {
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")

	return s
}

var ReplaceAllWhitespaceWithSpaceReg = regexp.MustCompile(`\s+`)

// ReplaceAllWhitespaceWithSpace 使用正则表达式将所有空白字符替换为空格
func ReplaceAllWhitespaceWithSpace(s string) string {
	return ReplaceAllWhitespaceWithSpaceReg.ReplaceAllString(s, " ")
}

// 定义HTML特殊字符及其转义映射
var escapeMap = map[string]string{
	"<":  "&lt;",
	">":  "&gt;",
	"&":  "&amp;",
	"\"": "&quot;",
	"'":  "&apos;",
}
var regexpMatchHtmlSpecialCharacters = regexp.MustCompile(`[<>&"']`)

// EscapeHtmlSpecialCharacters 将HTML中的特殊符号转换为HTML转义文本
func EscapeHtmlSpecialCharacters(s string) string {
	// 使用正则表达式匹配特殊字符

	return regexpMatchHtmlSpecialCharacters.ReplaceAllStringFunc(s, func(match string) string {
		return escapeMap[match]
	})
}
