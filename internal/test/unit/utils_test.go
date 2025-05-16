package unit_test

import (
	"context"
	"fmt"
	"go-cs/internal/utils"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	tplt_config "go-cs/internal/domain/work_flow/flow_tplt_config"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

func TestXxx(t *testing.T) {

	var c chan int = make(chan int)
	var tc <-chan time.Time = time.After(time.Second * 8)

	go func() {
		time.Sleep(time.Second * 10)
		c <- 4
	}()

	select {
	case <-c:
		fmt.Println(1)
	case <-tc:
		fmt.Println(2)
	}
}

func TestCharLength(t *testing.T) {

	// a := "😊🏷️𠮷𠮾我们123abc😈，"
	a := "聊天中的“精选”分享，“来自于...”点击跳转到个人页面的“动态”页签，应该在“精选”页签"
	runes := []rune(a)
	totalLen := 0

	for _, v := range runes {
		if v > 256 {
			totalLen = totalLen + 2
		} else {
			totalLen = totalLen + 1
		}
	}

	fmt.Println(totalLen)

	// len()
}

func TestPwdRegx(t *testing.T) {
	validate := utils.NewValidator()
	if validErr := validate.Var("fan001", "required,utf8Len=2-14,nickname"); validErr != nil {
		fmt.Println(validErr)
	}

}

func TestValid(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	//任意 中文 英文 下划线 数字
	validate.RegisterValidationCtx("utf8Len", func(ctx context.Context, fl validator.FieldLevel) bool {
		val := fl.Field().String()
		utf8Len := utf8.RuneCountInString(val)
		params := strings.Split(fl.Param(), " ")
		if len(params) == 2 {
			minLenParam, _ := strconv.ParseInt(params[0], 10, 0)
			maxLenParam, _ := strconv.ParseInt(params[1], 10, 0)
			if utf8Len >= int(minLenParam) && utf8Len <= int(maxLenParam) {
				return true
			}
			return false
		} else {
			return true
		}
	})

	validate.Var("我们都是龙的传人", "utf8Len=1 10")

}

type MappingCase struct {
	Id   int32
	Name string
}

func TestMapping(t *testing.T) {
	arr := []*MappingCase{
		{1, "a"},
		{2, "b"},
		{3, "c"},
		{4, "d"},
	}

	maps := utils.ConvertListToMap[MappingCase](arr, func(mc *MappingCase) string {
		// return strconv.Itoa(int(mc.Id))
		return mc.Name
	})

	fmt.Println(maps)
}

func TestMustJsonWfConf(t *testing.T) {
	confJson := "{\"terminatedReasonOptions\":[\"取消任务，综合考虑且已同步现在先不做了\",\"任务重复/合并，与其他在做任务一同推进\",\"取消任务，转 其它 任务流程\",\"任务需求已完成/已失效\"],\"enableTerminatedReasonOtherOption\":true,\"rebootReasonOptions\":[\"任务需要继续进行\",\"误操作\"],\"enableRebootReasonOtherOption\":true,\"uuid\":\"3bbb8f3d-a8fa-4fde-8769-e60cc110a6f3\",\"version\":\"0\",\"key\":\"tplt_PesVJ\",\"nodes\":[{\"name\":\"开始\",\"code\":\"started\",\"startMode\":\"pre_node_all_done\",\"belongStatus\":\"started\",\"needDoneOperator\":false,\"doneOperationRole\":[],\"passMode\":\"auto_confirm\",\"onReach\":[],\"onPass\":[],\"enableRollback\":false,\"owner\":{\"forceOwner\":false,\"usageMode\":\"\",\"value\":\"null\",\"ownerRole\":[]},\"doneOperationDisplayName\":\"\",\"enableClose\":false,\"startAt\":\"0\",\"rollbackReasonOptions\":[],\"enableCloseReasonOtherOption\":false,\"restartReasonOptions\":[],\"enableRestartReasonOtherOption\":false},{\"name\":\"未命名\",\"code\":\"state_0\",\"startMode\":\"pre_node_all_done\",\"belongStatus\":\"started\",\"needDoneOperator\":false,\"doneOperationRole\":[\"_node_owner\",\"_space_manager\",\"_creator\",\"_space_editor\"],\"passMode\":\"single_user_confirm\",\"onReach\":[{\"eventType\":\"changeStoryStage\",\"condition\":\"null\",\"targetSubState\":{\"id\":\"1926\",\"key\":\"testing\"}}],\"onPass\":[],\"enableRollback\":false,\"owner\":{\"forceOwner\":true,\"usageMode\":\"none\",\"value\":\"null\",\"ownerRole\":[{\"id\":\"1149\",\"key\":\"_accepter\"}]},\"doneOperationDisplayName\":\"1·2\",\"enableClose\":false,\"startAt\":\"0\",\"rollbackReasonOptions\":[\"开发 基础功能 的完成度，不足以完成期望测试\",\"开发 交互操作逻辑 的完成度，不足以完成期望测试\",\"开发 UI/视觉 的完成度，不足以完成期望测试\"],\"enableCloseReasonOtherOption\":true,\"restartReasonOptions\":[\"任务需要继续进行\",\"误操作\"],\"enableRestartReasonOtherOption\":true},{\"name\":\"完成\",\"code\":\"ended\",\"startMode\":\"pre_node_all_done\",\"belongStatus\":\"started\",\"needDoneOperator\":false,\"doneOperationRole\":[],\"passMode\":\"auto_confirm\",\"onReach\":[{\"eventType\":\"changeStoryStage\",\"condition\":\"null\",\"targetSubState\":{\"id\":\"1924\",\"key\":\"completed\",\"uuid\":\"6e7b0101-15be-4fbc-bb7a-444bfe84a11b\",\"val\":\"2\"}}],\"onPass\":[],\"enableRollback\":false,\"owner\":{\"forceOwner\":false,\"usageMode\":\"\",\"value\":\"null\",\"ownerRole\":[]},\"doneOperationDisplayName\":\"\",\"enableClose\":false,\"startAt\":\"0\",\"rollbackReasonOptions\":[],\"enableCloseReasonOtherOption\":false,\"restartReasonOptions\":[],\"enableRestartReasonOtherOption\":false}],\"connections\":[{\"startNode\":\"started\",\"endNode\":\"state_0\"},{\"startNode\":\"state_0\",\"endNode\":\"ended\"}]}"
	wfConf, _ := tplt_config.MustFormWorkFlowJson(confJson)
	fmt.Println(wfConf)
}

func TestCpyName(t *testing.T) {

	newName := "阿三"
	needReName := false
	nameNoTag := make(map[int32]bool)
	names := []string{"阿三", "阿三(10)", "阿三(1)", "阿三(2)", "阿三(13)", "阿三(3)", "阿三(4)", "阿三(5)", "阿三(6)", "阿三(7)", "阿三(8)", "阿三(9)", "阿三(10)", "阿四5", "ahashdhf", "阿三(1)", "阿三(2)"}
	regx, _ := regexp.Compile(`^` + newName + `\((\d+)\)$`)
	for _, v := range names {
		if newName == v {
			needReName = true
			continue
		}

		if isMatch := regx.Match([]byte(v)); isMatch {
			repeatTag := regx.ReplaceAllString(v, "$1")
			if repeatTag != "" {
				nameNoTag[cast.ToInt32(repeatTag)] = true
			}
		}
	}

	if needReName {
		var repateCount int32 = 1
		for {
			if _, ok := nameNoTag[repateCount]; !ok {
				newName = newName + "(" + cast.ToString(repateCount) + ")"
				break
			}
			repateCount++
		}
	}

	fmt.Println(newName)
}

func TestCast(t *testing.T) {
	v := cast.ToStringSlice([]int64{1, 2, 3, 4})
	fmt.Println(v)
}
