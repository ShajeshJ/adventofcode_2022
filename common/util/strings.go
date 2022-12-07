package util

import "strconv"

// AtoiNoError calls `strconv.Atoi`, but assumes input `s` is a valid integer
func AtoiNoError(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return val
}
