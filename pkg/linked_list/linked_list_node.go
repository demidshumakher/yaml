package linked_list

import (
	"fmt"
	"strings"
)

type linkedListNode[V any] struct {
	Value V
	Next  *linkedListNode[V]
	Child *linkedListNode[V]
}

func newNode[V any](value V) *linkedListNode[V] {
	return &linkedListNode[V]{
		Value: value,
	}
}

func (ll *linkedListNode[V]) String() string {
	node := ll
	var b strings.Builder
	b.WriteString(" { ")
	for node != nil {
		b.WriteString(fmt.Sprintf("%s : %s", node.Value, node.Child.String()))
	}
	b.WriteString(" } ")
	return b.String()
}
