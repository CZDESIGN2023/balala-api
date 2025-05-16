package search2

var priority = map[string]int{
	"P0":      0,
	"P1":      1,
	"P2":      2,
	"P3":      3,
	"P4":      4,
	"PENDING": 5,
	"SUSPEND": 6,
}

func GetPriority(p string) int {
	return priority[p]
}
