package date

import (
	"testing"
	"time"
)

func TestDateWeekBegin(t *testing.T) {
	t.Log(WeekBegin(time.Now().AddDate(0, 0, -4)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, -3)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, -2)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, -1)))
	t.Log(WeekBegin(time.Now()))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, 1)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, 2)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, 3)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, 4)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, 5)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, 6)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, 7)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, 8)))
	t.Log(WeekBegin(time.Now().AddDate(0, 0, 9)))
}

func Test_dateWeekBegin(t *testing.T) {
	t.Log(dateWeekBegin(time.Now(), time.Monday))
}

func TestCurWeekBeginEnd(t *testing.T) {
	begin, end := ThisWeekBeginEnd()
	t.Log(begin, end)
}
