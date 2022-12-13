package util

import "math"

func Abs(v int) int {
	return int(math.Abs(float64(v)))
}

func Min(x int, y int) int {
	return int(math.Min(float64(x), float64(y)))
}
