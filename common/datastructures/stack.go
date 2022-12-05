package datastructures

type Stack[T any] []T

func (s *Stack[T]) Len() int {
	return len(*s)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Push(val T) {
	*s = append(*s, val)
}

func (s *Stack[T]) Pop() (T, bool) {
	val, ok := s.PopN(1)
	return val[0], ok
}

func (s *Stack[T]) PopN(n int) ([]T, bool) {
	if s.Len() < n {
		return *new([]T), false
	}

	val := (*s)[len(*s)-n:]
	*s = (*s)[:len(*s)-n]
	return val, true
}

func (s *Stack[T]) PushN(val []T) {
	*s = append(*s, val...)
}
