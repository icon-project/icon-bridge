package util

func Max(x, y uint64) uint64 {
	if x < y {
		return y
	}
	return x
}
