package linked_list

type LinkedList[V any] struct {
	head *linkedListNode[V]
	tail *linkedListNode[V]
}

func New[V any]() *LinkedList[V] {
	return &LinkedList[V]{}
}

func (ll *LinkedList[V]) PushBack(v V) {
	if ll.head == nil {
		ll.head = newNode(v)
		ll.tail = ll.head
		return
	}
	ll.tail.Next = newNode(v)
	ll.tail = ll.tail.Next
}

func (ll *LinkedList[V]) AddChild(child *LinkedList[V]) {
	ll.tail.Child = child.head
}

func (ll *LinkedList[V]) String() string {
	return ll.head.String()
}
