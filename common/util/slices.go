package util

func Map[T any, R any](x []T, mapFunc func(T) R) []R {
	var output []R
	for _, item := range x {
		output = append(output, mapFunc(item))
	}
	return output
}
