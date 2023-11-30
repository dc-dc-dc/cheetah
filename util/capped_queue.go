package util

type CappedQueue[T any] struct {
	cap   int
	queue *Queue[T]
}

func NewCappedQueue[T any](cap int) *CappedQueue[T] {
	return &CappedQueue[T]{
		cap:   cap,
		queue: NewQueue[T](),
	}
}

func (cq *CappedQueue[T]) Push(item T) {
	if cq.queue.Count() == cq.cap {
		cq.queue.PopLeft()
	}
	cq.queue.Push(item)
}

func (cq *CappedQueue[T]) Elements() []T {
	res := make([]T, 0, cq.cap)
	for node := cq.queue.head; node != nil; node = node.next {
		res = append(res, node.value)
	}
	return res
}

func (cq *CappedQueue[T]) Cap() int {
	return cq.cap
}

func (cq *CappedQueue[T]) Full() bool {
	return cq.queue.Count() == cq.cap
}
