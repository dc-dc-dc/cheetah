package util

type Queue struct {
	head  *DoublyLinkedListNode
	tail  *DoublyLinkedListNode
	count int
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Push(item interface{}) {
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

func (q *Queue) Elements() []interface{} {
	res := make([]interface{}, 0)
	for node := q.head; node != nil; node = node.next {
		res = append(res, node.value)
	}
	return res
}

func (q *Queue) First() interface{} {
	if q.head == nil {
		return nil
	}
	return q.head.value
}

func (q *Queue) Last() interface{} {
	if q.tail == nil {
		return nil
	}
	return q.tail.value
}

func (q *Queue) PopLeft() {
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

func (q *Queue) Pop() {
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

func (q *Queue) Count() int {
	return q.count
}
