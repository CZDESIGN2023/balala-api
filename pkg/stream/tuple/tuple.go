package tuple

// Pair 二元组
type Pair[F, S any] struct {
	First  F
	Second S
}

// Triple 三元组
type Triple[F, S, T any] struct {
	First  F
	Second S
	Third  T
}

func T2[F, S any](first F, second S) Pair[F, S] {
	return Pair[F, S]{
		First:  first,
		Second: second,
	}
}

func T3[F, S, T any](first F, second S, third T) Triple[F, S, T] {
	return Triple[F, S, T]{
		First:  first,
		Second: second,
		Third:  third,
	}
}
