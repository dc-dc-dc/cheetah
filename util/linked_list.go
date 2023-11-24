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

type DoublyLinkedListNode struct {
	value interface{}
	next  *DoublyLinkedListNode
	prev  *DoublyLinkedListNode
}

func NewDoublyLinkedListNode(value interface{}) *DoublyLinkedListNode {
	return &DoublyLinkedListNode{
		value: value,
	}
}

type LinkedList struct {
	head *LinkedListNode
}
