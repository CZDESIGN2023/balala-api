package group

import (
	"go-cs/pkg/pprint"
	"testing"
)

type User struct {
	Id   int64
	Name string
	Age  int64
	Tags []string
}

func Test(t *testing.T) {
	arr := []*User{
		{1, "a", 100, []string{"1", "2"}},
		{2, "b", 100, []string{"1", "2"}},
		{3, "c", 200, []string{"1", "2"}},
		{4, "d", 200, []string{"1", "2"}},
	}

	root := New(arr)

	root.GroupByToMulti("tag", func(user *User) []string {
		return user.Tags
	})
	//
	//root.GroupBy(func(user *User) string {
	//	return strconv.Itoa(int(user.Id % 2))
	//})
	//
	//root.GroupBy(func(user *User) string {
	//	return user.Name
	//})

	pprint.Println(root)
}
