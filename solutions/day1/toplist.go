package day1

import "golang.org/x/exp/constraints"

type TopList[T constraints.Ordered] struct {
	Size int
	data []T
}

func NewTopList[T constraints.Ordered](size int) *TopList[T] {
	if size < 1 {
		return &TopList[T]{Size: 1, data: make([]T, 1)}
	}
	return &TopList[T]{Size: size, data: make([]T, size)}
}

// TryPush will add `val` to the list if it's greater than at least 1 number in the list in order, and returns true.
// Otherwise it does nothing and returns false
func (t *TopList[T]) TryPush(val T) bool {
	if val <= t.data[t.Size-1] {
		return false
	}

	t.data[t.Size-1] = val

	for i := t.Size - 2; i >= 0; i-- {
		if t.data[i] >= t.data[i+1] {
			break
		}
		t.data[i], t.data[i+1] = t.data[i+1], t.data[i]
	}

	return true
}

// Sum will return the sum of all values in the list
func (t *TopList[T]) Sum() (total T) {
	for _, item := range t.data {
		total += item
	}
	return
}
