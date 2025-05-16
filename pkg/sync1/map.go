package sync1

import "sync"

type Map[K comparable, V any] struct {
	m sync.Map
}

func (t *Map[K, V]) Store(k K, v V) {
	t.m.Store(k, v)
}

func (t *Map[K, V]) Load(k K) (ret V, loaded bool) {
	value, ok := t.m.Load(k)
	if !ok {
		return
	}

	return value.(V), true
}

func (t *Map[K, V]) Delete(k K) {
	t.m.Delete(k)
}

func (t *Map[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	ret, loaded := t.m.LoadOrStore(k, v)
	return ret.(V), loaded
}

func (t *Map[K, V]) LoadAndDelete(k K) (ret V, loaded bool) {
	value, loaded := t.m.LoadAndDelete(k)
	if !loaded {
		return
	}

	return value.(V), true
}

func (t *Map[K, V]) CompareAndDelete(k K, old V) (deleted bool) {
	return t.m.CompareAndDelete(k, old)
}

func (t *Map[K, V]) CompareAndSwap(k K, old, new V) bool {
	return t.m.CompareAndSwap(k, old, new)
}

func (t *Map[K, V]) Swap(k K, v V) (ret V, loaded bool) {
	value, loaded := t.m.Swap(k, v)
	if !loaded {
		return
	}
	return value.(V), true
}

func (t *Map[K, V]) Range(fn func(k K, v V)) {
	t.m.Range(func(key, value any) bool {
		fn(key.(K), value.(V))
		return true
	})
}
