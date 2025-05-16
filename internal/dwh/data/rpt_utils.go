package data

import (
	"go-cs/internal/utils/date"
	"slices"
	"time"
)

type TimeRangeItem[T any] interface {
	GetStartDate() time.Time
	GetEndDate() time.Time
	SetStartDate(time.Time)
	SetEndDate(time.Time)
	Clone() T
}

func ExpandDateRangeList[T TimeRangeItem[T]](timeSplitType string, startTime, endTime time.Time, list []T) []T {
	var finalList []T
	for i := 0; i < len(list); i++ {
		finalList = append(finalList, expandDateRange(timeSplitType, startTime, endTime, list[i])...)

		// 如果当前元素和下一个元素的开始时间不相等，填充空隙
		if i+1 < len(list) && list[i].GetEndDate().Before(list[i+1].GetStartDate()) {
			v := list[i].Clone()
			v.SetStartDate(list[i].GetEndDate())
			v.SetEndDate(list[i+1].GetStartDate())

			finalList = append(finalList, expandDateRange(timeSplitType, startTime, endTime, v)...)
		}
	}

	// 使用稳定排序，保证元素顺序
	slices.SortStableFunc(finalList, func(a, b T) int {
		return int(a.GetStartDate().Sub(b.GetStartDate()))
	})

	// slices.CompactFunc 去重保留的是第一个元素，所以倒序，然后再去重
	// 倒序
	slices.Reverse(finalList)

	// 去重
	finalList = slices.CompactFunc(finalList, func(a T, b T) bool {
		return a.GetStartDate().Equal(b.GetStartDate())
	})

	// 倒序
	slices.Reverse(finalList)

	return finalList
}

func expandDateRange[T TimeRangeItem[T]](timeSplitType string, startTime, endTime time.Time, item T) []T {
	var list []T

	if startTime.Before(item.GetStartDate()) {
		startTime = item.GetStartDate()
	}
	if endTime.After(item.GetEndDate()) {
		endTime = item.GetEndDate()
	}

	switch timeSplitType {
	case "Hour":
		begin := date.HourBegin(startTime)
		end := date.HourBegin(endTime)

		for !begin.After(end) {
			newItem := item.Clone()
			newItem.SetStartDate(begin)
			newItem.SetEndDate(begin.Add(time.Hour))

			begin = newItem.GetEndDate()

			list = append(list, newItem)
		}

	case "Day":
		begin := date.DayBegin(startTime)
		end := date.DayBegin(endTime)

		for !begin.After(end) {
			newItem := item.Clone()
			newItem.SetStartDate(begin)
			newItem.SetEndDate(newItem.GetStartDate().AddDate(0, 0, 1))

			list = append(list, newItem)

			begin = newItem.GetEndDate()
		}

	case "Month":
		begin := date.MonthBegin(startTime)
		end := date.MonthBegin(endTime).AddDate(0, 1, 0)

		for !begin.After(end) {
			newItem := item.Clone()
			newItem.SetStartDate(begin)
			newItem.SetEndDate(newItem.GetStartDate().AddDate(0, 1, 0))

			list = append(list, newItem)

			begin = newItem.GetEndDate()
		}
	default:
		list = append(list, item)
	}

	return list
}
