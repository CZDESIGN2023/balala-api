package condition_translater

import (
	v1 "go-cs/api/search/v1"
	"go-cs/internal/utils/date"
	"time"
)

// 将变量转换为具体的值
func (ctx *Ctx) handleVariable(c *v1.Condition) {
	handDateRange(c)
}

func handDateRange(c *v1.Condition) {
	if len(c.Values) != 2 {
		return
	}

	variable := extractVariable(c.Values[0])

	var start time.Time
	var end time.Time

	switch variable {
	case "TODAY":
		start = date.TodayBegin()
		end = start.AddDate(0, 0, 1)
	case "TOMORROW":
		start = date.TomorrowBegin()
		end = start.AddDate(0, 0, 1)
	case "YESTERDAY":
		start = date.YesterdayBegin()
		end = start.AddDate(0, 0, 1)
	case "THIS_WEEK":
		start = date.ThisWeekBegin()
		end = start.AddDate(0, 0, 7)
	case "NEXT_WEEK":
		start = date.ThisWeekBegin().AddDate(0, 0, 7)
		end = start.AddDate(0, 0, 7)
	case "LAST_WEEK":
		start = date.ThisWeekBegin().AddDate(0, 0, -7)
		end = start.AddDate(0, 0, 7)
	case "THIS_MONTH":
		start = date.ThisMonthBegin()
		end = start.AddDate(0, 1, 0)
	case "LAST_MONTH":
		start = date.ThisMonthBegin().AddDate(0, -1, 0)
		end = start.AddDate(0, 1, 0)
	case "NEXT_MONTH":
		start = date.ThisMonthBegin().AddDate(0, 1, 0)
		end = start.AddDate(0, 1, 0)
	case "THIS_YEAR":
		start = date.ThisYearBegin()
		end = start.AddDate(1, 0, 0)
	case "NEXT_YEAR":
		start = date.ThisYearBegin().AddDate(1, 0, 0)
		end = start.AddDate(1, 0, 0)
	case "LAST_YEAR":
		start = date.ThisYearBegin().AddDate(-1, 0, 0)
		end = start.AddDate(1, 0, 0)
	default:
		return
	}

	c.Values[0] = start.Format("2006/01/02 15:04:05")
	c.Values[1] = end.Add(-time.Second).Format("2006/01/02 15:04:05")
}
