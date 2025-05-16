package utils

func ConvertListToMap[T any](s []*T, fn func(*T) string) map[string]*T {
	var maps map[string]*T = make(map[string]*T, 0)
	for i := 0; i < len(s); i++ {
		mapKey := fn(s[i])
		maps[mapKey] = s[i]
	}
	return maps
}
