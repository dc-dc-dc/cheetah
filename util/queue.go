package util

type Queue struct {
	head  *LinkedListNode
	tail  *LinkedListNode
	count int
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Push(item interface{}) {
	q.count += 1
	if q.head == nil {
		q.head = NewLinkedListNode(item)
		q.tail = q.head
	} else {
		q.tail.next = NewLinkedListNode(item)
		q.tail = q.tail.next
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
	if q.head != nil {
		q.count -= 1
		q.head = q.head.next
	}
}

func (q *Queue) Count() int {
	return q.count
}
