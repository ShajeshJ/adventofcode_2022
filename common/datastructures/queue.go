package datastructures

type Queue[T any] []T

func (s *Queue[T]) Len() int {
	return len(*s)
}

func (s *Queue[T]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Queue[T]) Enqueue(val T) {
	*s = append(*s, val)
}

func (s *Queue[T]) Dequeue() (T, bool) {
	if s.IsEmpty() {
		return *new(T), false
	}

	val := (*s)[0]
	*s = (*s)[1:]
	return val, true
}
