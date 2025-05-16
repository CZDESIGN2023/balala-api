package ops_log

import (
	"fmt"
	domain "go-cs/internal/domain/user"
	"slices"
	"strings"
	"time"
)

func Q(v any) string {
	return fmt.Sprintf("<span>「%v」</span>", v)
}

func W(v any) string {
	return fmt.Sprintf("<span>%v</span>", v)
}

func U(v *domain.User) string {
	return fmt.Sprintf("%s（%s）", v.UserNickname, v.UserName)
}

func T(start, end int64) string {

	if start == 0 && end == 0 {
		return ""
	}

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	layout := "2006/01/02"

	startStr := startTime.Format(layout)
	endStr := endTime.Format(layout)

	return fmt.Sprintf("%s ~ %s", startStr, endStr)
}

func parseAttach(val int64, unit string) string {
	if val < 1024 {
		return fmt.Sprintf("%v%s", val, unit)
	}

	units := []string{"B", "KB", "MB", "GB"}
	idx := slices.Index(units, unit)

	if idx < 0 || idx == len(units)-1 {
		return fmt.Sprintf("%v%s", val, unit)
	}

	var fNumber = float64(val)
	for ; idx < len(units)-1 && fNumber >= 1024; idx++ {
		fNumber = float64(val) / 1024
	}

	fNumberStr := fmt.Sprintf(" %.1f", fNumber)
	if strings.HasSuffix(fNumberStr, ".0") {
		fNumberStr = fNumberStr[:len(fNumberStr)-2]
	}

	return fNumberStr + units[idx]
}
