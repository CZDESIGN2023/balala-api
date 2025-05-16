package tea_im

import (
	"cmp"
	"crypto/md5"
	"encoding/hex"
	"go-cs/pkg/stream"
	"slices"
	"strings"
)

func (c *Client) sign(argsMap map[string]string) string {
	if len(argsMap) == 0 {
		return ""
	}

	entries := stream.ToEntries(argsMap)

	// 排序
	slices.SortFunc(entries, func(a, b stream.Entry[string, string]) int {
		return cmp.Compare(a.Key, b.Key)
	})

	entries = append(entries, stream.Entry[string, string]{
		Key: "pri_key",
		Val: c.privateKey,
	})

	// 转换为key=value格式
	args := stream.Map(entries, func(entry stream.Entry[string, string]) string {
		return entry.Key + "=" + entry.Val
	})

	// 拼接参数
	str := strings.Join(args, "&")

	// 计算MD5
	md5Hash := md5.Sum([]byte(str))
	md5String := hex.EncodeToString(md5Hash[:])

	return md5String
}

func makeArgs(argsMap map[string]string) string {
	if len(argsMap) == 0 {
		return ""
	}

	entries := stream.ToEntries(argsMap)

	// 排序
	slices.SortFunc(entries, func(a, b stream.Entry[string, string]) int {
		return cmp.Compare(a.Key, b.Key)
	})

	// 转换为key=value格式
	args := stream.Map(entries, func(entry stream.Entry[string, string]) string {
		return entry.Key + "=" + entry.Val
	})

	// 拼接参数
	str := strings.Join(args, "&")

	return str
}
