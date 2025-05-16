package unit_test

import (
	"fmt"
	// "go-cs/internal/biz/perm"
	"strconv"
	"testing"

	"github.com/casbin/casbin/v2"
)

var supper_manager int64 = 999
var space_manager int64 = 1
var space_editor int64 = 2
var space_watcher int64 = 3

var _node_owner string = "_node_owner"
var _creator string = "_creator"

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func TestCasbinRolePerm(t *testing.T) {

	e, err := casbin.NewEnforcer(
		"/Users/cyt/project/ed-project-manage/server_go/configs/biz_perm/authz_model.conf",
		"/Users/cyt/project/ed-project-manage/server_go/configs/biz_perm/authz_policy.csv",
	)
	if err != nil {
		panic(err)
	}

	check(e, strconv.FormatInt(supper_manager, 10), "PERM_DELETE_SPACE", "*")
	check(e, strconv.FormatInt(space_manager, 10), "PERM_DELETE_SPACE", "*")
	check(e, strconv.FormatInt(space_editor, 10), "PERM_DELETE_SPACE", "*")
	check(e, strconv.FormatInt(space_watcher, 10), "PERM_DELETE_SPACE", "*")
}
