package stream

type Stream[T comparable] struct {
	item []T
}

func Of[S []T, T comparable](s S) *Stream[T] {
	return &Stream[T]{
		item: s,
	}
}

func Clone[S []T, T comparable](s S) *Stream[T] {
	return &Stream[T]{
		item: append(S(nil), s...),
	}
}

func (s *Stream[T]) Concat(items ...T) *Stream[T] {
	s.item = append(s.item, items...)
	return s
}

func (s *Stream[T]) Unique() *Stream[T] {
	s.item = Unique(s.item)
	return s
}

func (s *Stream[T]) Remove(item T) *Stream[T] {
	s.item = Remove(s.item, item)
	return s
}

func (s *Stream[T]) Union(items ...T) *Stream[T] {
	s.item = Union(s.item, items)
	return s
}

func (s *Stream[T]) Diff(items ...T) *Stream[T] {
	s.item = Diff(s.item, items)
	return s
}

func (s *Stream[T]) Intersect(items ...T) *Stream[T] {
	s.item = Intersect(s.item, items)
	return s
}

func (s *Stream[T]) List() []T {
	return s.item
}
