package util

import "math"

func Abs(v int) int {
	return int(math.Abs(float64(v)))
}
