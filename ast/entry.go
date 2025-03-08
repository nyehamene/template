package ast

func EntrySame[T any](k T, v T) Entry[T, T] {
	return Entry[T, T]{key: k, val: v}
}

func EntryMany[T any](k []T, v T) Entry[[]T, T] {
	return Entry[[]T, T]{key: k, val: v}
}

type Entry[K any, V any] struct {
	key K
	val V
}

func (e Entry[K, _]) Key() K {
	return e.key
}

func (e Entry[_, V]) Val() V {
	return e.val
}
