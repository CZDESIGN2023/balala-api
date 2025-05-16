package pkg

import (
	"cmp"
	"go-cs/pkg/stream"
	"strconv"
	"strings"
)

type Version struct {
	list []int64
}

func NewVersion(version string) Version {
	if version == "" {
		panic("version is empty")
	}

	split := strings.Split(version, ".")

	list := stream.Map(split, func(s string) int64 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}

		return i
	})

	return Version{
		list: list,
	}
}

func (v Version) CompareTo(b Version) int {
	for i := 0; i < len(v.list); i++ {
		if i >= len(b.list) {
			return 1
		}

		av := v.list[i]
		bv := b.list[i]
		if v := cmp.Compare(av, bv); v != 0 {
			return v
		}
	}

	return 0
}

func (v Version) String() string {
	return strings.Join(stream.Map(v.list, func(i int64) string {
		return strconv.FormatInt(i, 10)
	}), ".")
}

func (v Version) IsEmpty() bool {
	return len(v.list) == 0
}
