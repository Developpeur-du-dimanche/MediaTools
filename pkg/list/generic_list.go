package list

import (
	"sync"
)

type List[T comparable] struct {
	items []T
	count int
	mutex *sync.Mutex
}

func NewList[T comparable]() *List[T] {
	return &List[T]{
		items: make([]T, 0),
		mutex: &sync.Mutex{},
		count: 0,
	}
}

func (l *List[T]) AddItem(item T) {
	l.mutex.Lock()
	l.items = append(l.items, item)
	l.count++
	l.mutex.Unlock()
}

func (l *List[T]) RemoveItem(item T) {
	l.mutex.Lock()
	for i, v := range l.items {
		if v == item {
			l.items = append(l.items[:i], l.items[i+1:]...)
			break
		}
	}
	l.mutex.Unlock()
}

func (l *List[T]) RemoveItemAt(index int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.count == 0 {
		return
	}
	l.items = append(l.items[:index], l.items[index+1:]...)
	l.count--

}

func (l *List[T]) GetItems() []T {
	return l.items
}

func (l *List[T]) GetItem(index int) T {
	return l.items[index]
}

func (l *List[T]) GetLength() int {
	return len(l.items)
}

func (l *List[T]) Clear() {
	l.mutex.Lock()
	l.items = make([]T, 0)
	l.count = 0
	l.mutex.Unlock()
}

func (l *List[T]) Iterate(yield func(T) bool) {
	for _, item := range l.items {
		if !yield(item) {
			break
		}
	}
}
