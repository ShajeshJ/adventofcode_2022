package datastructures

import "golang.org/x/exp/constraints"

type TopList[T constraints.Ordered] struct {
	MaxSize int
	Values  []T
}

func NewTopList[T constraints.Ordered](size int) *TopList[T] {
	if size < 1 {
		size = 1
	}
	return &TopList[T]{MaxSize: size, Values: make([]T, size)}
}

// TryPush will add `val` to the list if it's greater than at least 1 number in the list in order, and returns true.
// Otherwise it does nothing and returns false
func (t *TopList[T]) TryPush(val T) bool {
	if val <= t.Values[t.MaxSize-1] {
		return false
	}

	t.Values[t.MaxSize-1] = val

	for i := t.MaxSize - 2; i >= 0; i-- {
		if t.Values[i] >= t.Values[i+1] {
			break
		}
		t.Values[i], t.Values[i+1] = t.Values[i+1], t.Values[i]
	}

	return true
}
