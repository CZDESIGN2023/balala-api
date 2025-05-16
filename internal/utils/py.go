package utils

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

var hanRegx = regexp.MustCompile(`[\p{Han}]+`)

func Pinyin(userNickName string) string {
	if userNickName == "" {
		return ""
	}

	nickNamePy := userNickName
	nickNameFullPy := userNickName

	//转拼音规则 张x员-》zhangxyuan 中文部分替换成拼音，其它位置保留原来的内容
	hans := hanRegx.FindAllString(userNickName, -1)
	for _, v := range hans {
		//转首写字母拼音
		firstLetterPy := strings.Join(pinyin.LazyPinyin(v, pinyin.Args{
			Style: pinyin.FIRST_LETTER,
		}), "")

		if firstLetterPy != "" {
			nickNamePy = strings.ReplaceAll(nickNamePy, v, firstLetterPy)
		}

		//全拼
		fullPy := strings.Join(pinyin.LazyPinyin(v, pinyin.Args{}), "")
		if fullPy != "" {
			nickNameFullPy = strings.ReplaceAll(nickNameFullPy, v, fullPy)
		}
	}

	return "," + nickNamePy + "," + nickNameFullPy + ","
}

// PinyinByFirstChar 拼音首字母
func PinyinByFirstChar(str string) string {
	if str == "" {
		return ""
	}

	py := str

	hans := hanRegx.FindAllString(str, -1)
	for _, v := range hans {
		//转首写字母拼音
		firstLetterPy := strings.Join(pinyin.LazyPinyin(v, pinyin.Args{
			Style: pinyin.FIRST_LETTER,
		}), "")

		if firstLetterPy != "" {
			py = strings.ReplaceAll(py, v, firstLetterPy)
		}
	}

	return py
}
