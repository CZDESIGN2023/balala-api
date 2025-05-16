package biz_utils

import (
	"cmp"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"slices"
	"strings"

	user_domain "go-cs/internal/domain/user"
)

func ToSimpleUser(u *user_domain.User) *rsp.SimpleUserInfo {
	if u == nil {
		return nil
	}

	return &rsp.SimpleUserInfo{
		Id:           u.Id,
		UserId:       u.Id,
		UserName:     u.UserName,
		UserNickname: u.UserNickname,
		Avatar:       u.Avatar,
	}
}

func GetSortListRankMap[T any](list []T, f func(T) (int64, string)) map[int64]int64 {

	tMap := make(map[int64]string, 0)
	for _, v := range list {
		tKey, tName := f(v)
		tMap[tKey] = tName
	}

	pyMap := make(map[int64]string, 0)
	for k, v := range tMap {
		py := utils.PinyinByFirstChar(v)
		py = strings.ToLower(py)
		pyMap[k] = py
	}

	pyEntries := stream.ToEntries(pyMap)
	slices.SortFunc(pyEntries, func(a, b stream.Entry[int64, string]) int {
		aLen, bLen := len(a.Val), len(b.Val)
		minLen := aLen
		if bLen < aLen {
			minLen = bLen
		}

		av, bv := a.Val, b.Val

		for i := 0; i < minLen; i++ {
			if av[i] != bv[i] {
				if av[i] == '_' {
					return -1
				}
				if bv[i] == '_' {
					return 1
				}
			}

			if av[i] != bv[i] {
				return cmp.Compare(av[i], bv[i])
			}
		}

		switch {
		case aLen > bLen:
			return 1
		case aLen < bLen:
			return -1
		default:
			return cmp.Compare(a.Key, b.Key) //名称都一样，按id升序
		}
	})

	rankMap := stream.ToMap(pyEntries, func(i int, v stream.Entry[int64, string]) (int64, int64) {
		return v.Key, int64(i)
	})

	return rankMap
}
