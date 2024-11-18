package linked_list

import (
	"fmt"
	"strings"
)

type LinkedList[V any] struct {
	Value V
	Next  *LinkedList[V]
}

func (ll *LinkedList[V]) String() string {
	var b strings.Builder
	b.Write([]byte{'{', ' '})
	node := ll
	for node != nil {
		b.WriteString(fmt.Sprint(node.Value, ", "))
		node = node.Next
	}
	b.WriteByte('}')
	return b.String()
}

func NewLinkedList[V any](value V) *LinkedList[V] {
	return &LinkedList[V]{
		Value: value,
	}
}

func FromArray[V any](arr []V) *LinkedList[V] {
	if len(arr) == 0 {
		return nil
	}
	root := NewLinkedList(arr[0])
	node := root
	for i := 1; i < len(arr); i++ {
		node.Add(arr[i])
		node = node.Next
	}
	return root
}

func (ll *LinkedList[V]) Add(value V) {
	ll.Next = NewLinkedList(value)
}
