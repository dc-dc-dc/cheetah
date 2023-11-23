package util

type LinkedListNode struct {
	value interface{}
	next  *LinkedListNode
}

func NewLinkedListNode(value interface{}) *LinkedListNode {
	return &LinkedListNode{
		value: value,
	}
}

type LinkedList struct {
	head *LinkedListNode
}
