package util

type CappedQueue struct {
	cap   int
	queue *Queue
}

func NewCappedQueue(cap int) *CappedQueue {
	return &CappedQueue{
		cap:   cap,
		queue: NewQueue(),
	}
}

func (cq *CappedQueue) Push(item interface{}) {
	if cq.queue.Count() == cq.cap {
		cq.queue.PopLeft()
	}
	cq.queue.Push(item)
}

func (cq *CappedQueue) Elements() []interface{} {
	return cq.queue.Elements()
}

func (cq *CappedQueue) Cap() int {
	return cq.cap
}

func (cq *CappedQueue) Full() bool {
	return cq.queue.Count() == cq.cap
}
