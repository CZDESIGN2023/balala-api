package utils

import (
	"fmt"
	"go-cs/internal/utils/date"
	"testing"
	"time"
)

func Test_workingDayCompute(t *testing.T) {
	// 示例配置：2023-01-01起，周一到周五每天两个工作时段
	config1 := WorkConfig{
		StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
		WorkTimes: map[time.Weekday][]TimePart{
			time.Monday:    {{9, 0}, {12, 0}, {13, 0}, {17, 0}},
			time.Tuesday:   {{9, 0}, {12, 0}, {13, 0}, {17, 0}},
			time.Wednesday: {{9, 0}, {12, 0}, {13, 0}, {17, 0}},
			time.Thursday:  {{9, 0}, {12, 0}, {13, 0}, {17, 0}},
			time.Friday:    {{9, 0}, {12, 0}, {13, 0}, {17, 0}},
			time.Saturday:  {{9, 0}, {12, 0}, {13, 0}, {17, 0}},
			time.Sunday:    {{9, 0}, {12, 0}, {13, 0}, {17, 0}},
		},
	}

	// 2023-07-01起，周四增加午休时间调整
	config2 := WorkConfig{
		StartDate: time.Date(2023, 7, 1, 0, 0, 0, 0, time.Local),
		WorkTimes: map[time.Weekday][]TimePart{
			time.Monday:    {{9, 0}, {11, 30}, {13, 30}, {17, 0}}, // 调整后时段
			time.Tuesday:   {{9, 0}, {11, 30}, {13, 30}, {17, 0}}, // 调整后时段
			time.Wednesday: {{9, 0}, {11, 30}, {13, 30}, {17, 0}}, // 调整后时段
			time.Thursday:  {{9, 0}, {11, 30}, {13, 30}, {17, 0}}, // 调整后时段
			time.Friday:    {{9, 0}, {11, 30}, {13, 30}, {17, 0}}, // 调整后时段
			time.Saturday:  {{9, 0}, {11, 30}, {13, 30}, {17, 0}}, // 调整后时段
			time.Sunday:    {{9, 0}, {11, 30}, {13, 30}, {17, 0}}, // 调整后时段
		},
	}

	// 计算2023-05-01 08:00 到 2023-08-01 16:00的工作时长

	start := date.Parse2("2023-06-01 00:00:00")
	end := date.Parse2("2023-08-01 00:00:00")

	total := CalculateWorkingHours(start, end, []WorkConfig{config1, config2})
	fmt.Printf("总工作时长：%v\n", total)
}
