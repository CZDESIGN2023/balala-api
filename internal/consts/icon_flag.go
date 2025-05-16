package consts

import "slices"

const (
	ICON_FLAG_Pic  = 1 << iota //有图片
	ICON_FLAG_Url  = 1 << iota //有图片
	ICON_FLAG_File = 1 << iota //有附件
)

var allIconFlags = []uint32{
	ICON_FLAG_Pic,
	ICON_FLAG_Url,
	ICON_FLAG_File,
}

func AllIconFlags() []uint32 {
	return allIconFlags
}

func IsValidIconFlags(flags ...uint32) bool {
	for _, flag := range flags {
		if !slices.Contains(allIconFlags, flag) {
			return false
		}
	}

	return true
}

func ParseFlagBit(bits uint32) []uint32 {
	var flags []uint32
	for _, flag := range AllIconFlags() {
		if flag&bits != 0 {
			flags = append(flags, flag)
		}
	}
	return flags
}

func FlagContains(bits uint32, flag uint32) bool {
	return bits&flag > 0
}

func MergeIconFlags(flags []uint32) uint32 {
	var v uint32
	for _, flag := range flags {
		v |= flag
	}

	return v
}
