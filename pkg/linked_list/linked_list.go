package linked_list

import "iter"

type LinkedList[V any] struct {
	head *LinkedListNode[V]
	tail *LinkedListNode[V]
}

func New[V any]() *LinkedList[V] {
	return &LinkedList[V]{}
}

func (ll *LinkedList[V]) PushBack(v V) {
	if ll.head == nil {
		ll.head = NewNode(v)
		ll.tail = ll.head
		return
	}
	ll.tail.Next = NewNode(v)
	ll.tail = ll.tail.Next
}

func (ll *LinkedList[V]) AddChild(child *LinkedList[V]) {
	ll.tail.Child = child.head
}

func (ll *LinkedList[V]) AddChildFromNode(child *LinkedListNode[V]) {
	ll.tail.Child = child
}

func (ll *LinkedList[V]) String() string {
	return ll.head.String()
}

func (ll *LinkedList[V]) Iterate() iter.Seq[*LinkedListNode[V]] {
	return func(yield func(*LinkedListNode[V]) bool) {
		nd := ll.head
		for nd != nil {
			if !yield(nd) {
				return
			}
			nd = nd.Next
		}
	}
}

func (ll *LinkedList[V]) GetHead() *LinkedListNode[V] {
	return ll.head
}
