package utils

import (
	"sort"
	"time"
)

type TimePart []int

func (t TimePart) Hour() int {
	if len(t) < 1 {
		return 0
	}
	return t[0]
}
func (t TimePart) Minute() int {
	if len(t) < 2 {
		return 0
	}
	return t[1]
}
func (t TimePart) Second() int {
	if len(t) < 3 {
		return 0
	}
	return t[2]
}

// WorkConfig 定义工作日配置
type WorkConfig struct {
	StartDate time.Time // 配置生效日期
	EndDate   time.Time // 配置失效日期，为空表示长期有效
	// 每天的工作时间段，例如：
	// map[time.Monday][][]TimePart{{9,0}, {12,0}}, {{13,0}, {17,0}}
	WorkTimes map[time.Weekday][]TimePart
}

// CalculateWorkingHours 计算动态配置下的工作时长
func CalculateWorkingHours(
	startTime, endTime time.Time,
	configs []WorkConfig) time.Duration {

	if len(configs) == 0 || startTime.After(endTime) {
		return 0
	}

	// 按生效时间排序配置
	sortedConfigs := sortAndMergeConfigs(configs)
	var total time.Duration
	current := startTime

	for _, config := range sortedConfigs {
		// 确定配置有效区间
		configStart := maxTime(config.StartDate, current)
		configEnd := endTime
		if !config.EndDate.IsZero() {
			configEnd = minTime(config.EndDate, endTime)
		}

		if configStart.After(configEnd) {
			continue
		}

		// 计算当前配置区间内的工作时长
		d := computeComplexHours(configStart, configEnd, config.WorkTimes)
		total += d
	}

	return total
}

// computeComplexHours 处理多时段复杂计算
func computeComplexHours(
	start, end time.Time,
	workTimes map[time.Weekday][]TimePart) time.Duration {

	var total time.Duration
	current := start

	for current.Before(end) {
		// 获取当天的工作时段
		day := current.Weekday()
		times, exists := workTimes[day]
		if !exists || len(times) == 0 {
			current = current.Add(24 * time.Hour)
			continue
		}

		// 处理当天的所有工作时段
		for i := 0; i < len(times); i += 2 {
			startHM := times[i]
			endHM := times[i+1]

			// 构建当天的时段
			dayStart := time.Date(current.Year(), current.Month(), current.Day(), startHM.Hour(), startHM.Minute(), startHM.Second(), 0, current.Location())
			dayEnd := time.Date(current.Year(), current.Month(), current.Day(), endHM.Hour(), endHM.Minute(), endHM.Second(), 0, current.Location())

			// 计算有效时间
			effectiveStart := maxTime(current, dayStart)
			effectiveEnd := minTime(end, dayEnd)

			total += calculateIntersectionDuration(effectiveStart, effectiveEnd, dayStart, dayEnd)
		}

		// 移动到下一天
		current = current.Add(24 * time.Hour)
	}

	return total
}

// sortAndMergeConfigs 排序并合并配置区间
func sortAndMergeConfigs(configs []WorkConfig) []WorkConfig {
	// 按生效时间排序
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].StartDate.Before(configs[j].StartDate)
	})

	// 合并配置区间
	merged := make([]WorkConfig, 0)

	for i := 0; i < len(configs); i++ {
		config := configs[i]
		var endDate time.Time
		if i < len(configs)-1 {
			endDate = configs[i+1].StartDate
		}
		merged = append(merged, WorkConfig{
			StartDate: config.StartDate,
			EndDate:   endDate,
			WorkTimes: config.WorkTimes,
		})
	}

	return merged
}

// 时间比较辅助函数
func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}

	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

// 计算两个时间区间的交集时长
func calculateIntersectionDuration(start1, end1, start2, end2 time.Time) time.Duration {
	if start1.After(end1) || start2.After(end2) {
		return 0
	}

	start := maxTime(start1, start2)
	end := minTime(end1, end2)

	if start.After(end) {
		return 0
	}

	return end.Sub(start)
}
