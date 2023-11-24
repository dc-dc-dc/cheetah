package util

type LinkedListNode[T any] struct {
	value T
	next  *LinkedListNode[T]
}

func NewLinkedListNode[T any](value T) *LinkedListNode[T] {
	return &LinkedListNode[T]{
		value: value,
	}
}

type DoublyLinkedListNode[T any] struct {
	value T
	next  *DoublyLinkedListNode[T]
	prev  *DoublyLinkedListNode[T]
}

func NewDoublyLinkedListNode[T any](value T) *DoublyLinkedListNode[T] {
	return &DoublyLinkedListNode[T]{
		value: value,
	}
}
