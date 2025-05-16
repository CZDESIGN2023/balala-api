package date

import (
	"time"
)

func HourBegin(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

func WeekBegin(t time.Time) time.Time {
	return dateWeekBegin(t, time.Monday) //周一作为一周的开始
}

func WeekEnd(t time.Time) time.Time {
	return WeekBegin(t).Add(time.Hour*24*7 - time.Second)
}

func WeekBeginEnd(t time.Time) (begin, end time.Time) {
	return WeekBegin(t), WeekEnd(t)
}

func DayBegin(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func DayEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func DayBeginEnd(t time.Time) (begin, end time.Time) {
	return DayBegin(t), DayEnd(t)
}

func MonthBegin(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func MonthEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).Add(-time.Second)
}

func MonthBeginEnd(t time.Time) (begin, end time.Time) {
	return MonthBegin(t), MonthEnd(t)
}

func YearBegin(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func YearEnd(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location()).Add(time.Hour*24*365 - time.Second)
}

func YearBeginEnd(t time.Time) (begin, end time.Time) {
	return YearBegin(t), YearEnd(t)
}

func TodayBegin() time.Time {
	return DayBegin(time.Now())
}

func TodayEnd() time.Time {
	return DayEnd(time.Now())
}

func TodayBeginEnd() (begin, end time.Time) {
	return DayBeginEnd(time.Now())
}

func TomorrowBegin() time.Time {
	return DayBegin(time.Now()).Add(time.Hour * 24)
}

func YesterdayBegin() time.Time {
	return DayBegin(time.Now()).Add(time.Hour * -24)
}

func ThisWeekBegin() time.Time {
	return WeekBegin(time.Now())
}

func ThisWeekBeginEnd() (begin, end time.Time) {
	return WeekBeginEnd(time.Now())
}

func ThisMonthBegin() time.Time {
	return MonthBegin(time.Now())
}

func ThisMonthEnd() time.Time {
	return MonthEnd(time.Now())
}

func ThisMonthBeginEnd() (begin, end time.Time) {
	return MonthBeginEnd(time.Now())
}

func ThisYearBegin() time.Time {
	return YearBegin(time.Now())
}

func ThisYearEnd() time.Time {
	return YearEnd(time.Now())
}

func ThisYearBeginEnd() (begin, end time.Time) {
	return YearBeginEnd(time.Now())
}

func dateWeekBegin(t time.Time, beginWeekDay time.Weekday) time.Time {
	// 默认星期日为一周的开始
	// 0	1 	2	3	4	5	6	索引位
	// 日	一	二	三	四	五	六

	// 如果星期一为一周的开始
	// 0	1 	2	3	4	5	6	索引位
	// 一	二	三	四	五	六	日

	// 获取当前日期是星期几的索引（0表示星期日）
	idx := int(t.Weekday())
	// 获取指定一周开始的索引
	n := int(beginWeekDay)

	// 调整到指定一周开始的日期
	adjustedIdx := (idx - n + 7) % 7

	return time.Date(t.Year(), t.Month(), t.Day()-adjustedIdx, 0, 0, 0, 0, t.Location())
}

// HasInter 区间 [x,y] [a, b] 是否有交集
func HasInter(x, y int64, a, b int64) bool {
	//s1 := a <= x && x <= b
	//s2 := a <= y && y <= b
	//s3 := x <= a && b <= y
	//
	//fmt.Printf("(%v, %v) (%v, %v)", x, y, a, b)
	//fmt.Println(s1, s2, s3)

	return x <= b && a <= y
}

const (
	dateTimeLayout  = "2006/01/02 15:04:05"
	dateTimeLayout2 = "2006-01-02 15:04:05"
)

func Format(t time.Time) string {
	return t.Format(dateTimeLayout)
}

func Parse(s string) time.Time {
	t, _ := time.ParseInLocation(dateTimeLayout, s, time.Local)
	return t
}

func Parse2(s string) time.Time {
	t, _ := time.ParseInLocation(dateTimeLayout2, s, time.Local)
	return t
}

func ParseInLocation(layout string, s string) time.Time {
	t, _ := time.ParseInLocation(layout, s, time.Local)
	return t
}
