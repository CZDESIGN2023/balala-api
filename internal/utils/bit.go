package utils

// GetBit 获取指定位的值 (0 或 1)
// 参数：
// n：待检查的整型值
// pos：待检查的位的索引，从右向左计数，从0开始
func GetBit(n uint32, pos uint32) uint32 {
	val := n & (1 << pos)
	if val > 0 {
		return 1
	}
	return 0
}

// SetBit 将指定位置的值设置为 0 或 1, 然後回傳新值
// 参数：
// n：待检查的整型值
// pos：待检查的位的索引，从右向左计数，从0开始
func SetBit(n uint32, pos uint32, val uint32) uint32 {
	if val == 1 {
		// 将指定位置的值设置为 1
		return n | (1 << pos)
	} else {
		// 将指定位置的值设置为 0
		return n & ^(1 << pos)
	}
}
