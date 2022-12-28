package util

import "math"

func Abs(v int) int {
	return int(math.Abs(float64(v)))
}

func Min(x int, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

func Max(x int, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func Normalize(x int) int {
	if x == 0 {
		return 0
	}
	return x / int(math.Abs(float64(x)))
}

func GCD(a, b int) int {
	if b == 0 {
		return a
	}
	return GCD(b, a%b)
}

func LCM(a, b int) int {
	return a * b / GCD(a, b)
}
