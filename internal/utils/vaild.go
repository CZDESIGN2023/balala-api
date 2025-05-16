package utils

import (
	"context"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var punctuationChars = []rune(`·”“‘’'!"#\$%&'\(\)\*\+,-\./:;<=>\?@\[\]\^_` + "`" + `\{\|\}~`)

var validate *validator.Validate

func init() {
	validate = newValidator()
}

func NewValidator() *validator.Validate {
	return validate
}

func newValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	//任意 中文 英文 下划线 数字
	validate.RegisterValidationCtx("username", func(ctx context.Context, fl validator.FieldLevel) bool {
		val := fl.Field().String()
		regx := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
		isMatch := regx.MatchString(val)
		return !isMatch
	})

	//任意 中文 英文 下划线 数字
	validate.RegisterValidationCtx("nickname", func(ctx context.Context, fl validator.FieldLevel) bool {
		val := fl.Field().String()
		regx := regexp.MustCompile(`[^\p{Han}a-zA-Z0-9_]+`)
		isMatch := regx.MatchString(val)
		return !isMatch
	})

	//任意 中文 英文 下划线 数字
	validate.RegisterValidationCtx("common_name", func(ctx context.Context, fl validator.FieldLevel) bool {
		val := fl.Field().String()
		regx := regexp.MustCompile(`[^\p{Han}a-zA-Z0-9_]+`)
		isMatch := regx.MatchString(val)
		return !isMatch
	})

	//任意英文
	validate.RegisterValidationCtx("pinyin", func(ctx context.Context, fl validator.FieldLevel) bool {
		val := fl.Field().String()
		regx := regexp.MustCompile(`[^a-zA-Z]+`)
		isMatch := regx.MatchString(val)
		return !isMatch
	})

	//任意 英文 数字 特殊字符
	validate.RegisterValidationCtx("password", func(ctx context.Context, fl validator.FieldLevel) bool {
		val := fl.Field().String()
		regx := regexp.MustCompile(`[^!-~]+`)
		isMatch := regx.MatchString(val)
		return !isMatch
	})

	//强密码 任意 字母 数字 特殊字符 中的两种
	validate.RegisterValidationCtx("stronger_password", func(ctx context.Context, fl validator.FieldLevel) bool {
		val := fl.Field().String()

		// 匹配 字母 数字 特殊字符
		if matched, _ := regexp.MatchString(`[^!-~]+`, val); matched {
			return false
		}

		var typeCount int
		// 是否包含数字
		if matched, _ := regexp.MatchString(`[0-9]+`, val); matched {
			typeCount++
		}

		// 是否包含字母
		if matched, _ := regexp.MatchString(`[a-zA-Z]+`, val); matched {
			typeCount++
		}

		// 是否包含特殊字符
		if matched, _ := regexp.MatchString(`[^0-9a-zA-Z]+`, val); matched {
			typeCount++
		}

		return typeCount >= 2
	})

	//按ut8计算字符串长度utf8Len=最小-最大
	validate.RegisterValidationCtx("utf8Len", func(ctx context.Context, fl validator.FieldLevel) bool {

		val := fl.Field().String()
		runes := []rune(val)
		totalLen := 0

		for i := 0; i < len(runes); i++ {
			if slices.Index(punctuationChars, runes[i]) >= 0 {
				totalLen = totalLen + 1
			} else {
				if runes[i] > 256 {
					totalLen = totalLen + 2
				} else {
					totalLen = totalLen + 1
				}
			}
		}

		params := strings.Split(fl.Param(), "-")
		if len(params) == 2 {
			minLenParam, _ := strconv.ParseInt(params[0], 10, 0)
			maxLenParam, _ := strconv.ParseInt(params[1], 10, 0)
			if totalLen >= int(minLenParam) && totalLen <= int(maxLenParam) {
				return true
			}
			return false
		} else {
			return true
		}
	})

	validate.RegisterValidationCtx("runeLen", func(ctx context.Context, fl validator.FieldLevel) bool {

		val := fl.Field().String()
		runes := []rune(val)
		totalLen := len(runes)

		params := strings.Split(fl.Param(), "-")
		if len(params) == 2 {
			minLenParam, _ := strconv.ParseInt(params[0], 10, 0)
			maxLenParam, _ := strconv.ParseInt(params[1], 10, 0)

			return totalLen >= int(minLenParam) && totalLen <= int(maxLenParam)
		}

		return true
	})
	return validate
}
