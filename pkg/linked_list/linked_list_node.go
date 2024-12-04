package linked_list

import (
	"fmt"
	"strings"
)

type LinkedListNode[V any] struct {
	Value V
	Next  *LinkedListNode[V]
	Child *LinkedListNode[V]
}

func NewNode[V any](value V) *LinkedListNode[V] {
	return &LinkedListNode[V]{
		Value: value,
	}
}

//func (ll *LinkedListNode[V]) Iter() iter.Seq[V] {
//	return func(yield func(v V) bool) {
//		node := ll
//		for node != nil {
//			if !yield(node.Value) {
//				return
//			}
//			node = node.Next
//		}
//	}
//}

func (ll *LinkedListNode[V]) String() string {
	return ll._string(0)
}

func (ll *LinkedListNode[V]) _string(n int) string {
	node := ll
	if node == nil {
		return ""
	}
	var b strings.Builder
	tabs := strings.Repeat("\t", n)

	b.WriteString(fmt.Sprintf("%s{\n", tabs))
	for node != nil {
		b.WriteString(fmt.Sprintf("\t%s{%s :\n%s}", tabs, node.Value, node.Child._string(n+1)))
		node = node.Next
	}
	b.WriteString(fmt.Sprintf("%s}\n", tabs))
	return b.String()
}
