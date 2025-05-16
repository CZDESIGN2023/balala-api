package dim

import (
	"time"
)

type DimDate struct {
	DateId      int
	Date        time.Time
	Year        int
	Quarter     int
	Month       int
	Day         int
	DayOfWeek   int
	WeekOfYear  int
	MonthName   string
	WeekdayName string
}
