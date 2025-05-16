package biz_utils

import (
	"go-cs/internal/utils"
	"regexp"
	"testing"
)

func TestPinyin(t *testing.T) {
	pinyin := utils.Pinyin("湖h北")

	t.Log(pinyin)
}

func TestPinyinByFirstChar(t *testing.T) {
	pinyin := utils.PinyinByFirstChar("湖H北")

	t.Log(pinyin)
}

func Test1(t *testing.T) {
	findString := regexp.MustCompile(`[^!-~]`).FindString("asd@@986986")
	match := regexp.MustCompile(`[^!-~]`).MatchString("asd@@986986")

	t.Log(findString)
	t.Log(match)
}
