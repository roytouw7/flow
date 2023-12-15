package linkedList

type LinkedList[T any] struct {
	Value *T
	next  *LinkedList[T]
}

func (l *LinkedList[T]) HasNext() bool {
	return l.next != nil
}

func (l *LinkedList[T]) Next() *LinkedList[T] {
	next := l.next
	if next == nil {
		panic("end of linked list")
	}
	return next
}

func (l *LinkedList[T]) Push(next T) {
	if l.Value == nil {
		l.Value = &next
		return
	}

	last := getLast(l)
	last.next = &LinkedList[T]{
		Value: &next,
	}
}

func getLast[T any](current *LinkedList[T]) *LinkedList[T] {
	for current.next != nil {
		current = current.next
	}
	return current
}

