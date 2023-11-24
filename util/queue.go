package util

type Queue[T any] struct {
	head  *DoublyLinkedListNode[T]
	tail  *DoublyLinkedListNode[T]
	count int
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Push(item T) {
	q.count += 1
	if q.head == nil {
		q.head = NewDoublyLinkedListNode(item)
		q.tail = q.head
	} else {
		last := q.tail
		q.tail.next = NewDoublyLinkedListNode(item)
		q.tail = q.tail.next
		q.tail.prev = last
	}
}

func (q *Queue[T]) Elements() []T {
	res := make([]T, 0)
	for node := q.head; node != nil; node = node.next {
		res = append(res, node.value)
	}
	return res
}

func (q *Queue[T]) First() T {
	if q.head == nil {
		return *new(T)
	}
	return q.head.value
}

func (q *Queue[T]) Last() T {
	if q.tail == nil {
		return *new(T)
	}
	return q.tail.value
}

func (q *Queue[T]) PopLeft() {
	if q.head == nil {
		return
	}
	q.count -= 1
	if q.head == q.tail {
		q.head = nil
		q.tail = nil
		return
	}

	q.head = q.head.next
	q.head.prev = nil
}

func (q *Queue[T]) Pop() {
	if q.head == nil {
		return
	}
	q.count -= 1
	if q.head == q.tail {
		q.head = nil
		q.tail = nil
		return
	}
	q.tail = q.tail.prev
}

func (q *Queue[T]) Count() int {
	return q.count
}
