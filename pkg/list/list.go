package list

import "sync"

type List struct {
	items []string
	mutex *sync.Mutex
}

func NewList() *List {
	return &List{
		items: make([]string, 0),
		mutex: &sync.Mutex{},
	}
}

func (l *List) AddItem(item string) {
	l.mutex.Lock()
	l.items = append(l.items, item)
	l.mutex.Unlock()
}

func (l *List) RemoveItem(item string) {
	l.mutex.Lock()
	for i, v := range l.items {
		if v == item {
			l.items = append(l.items[:i], l.items[i+1:]...)
			break
		}
	}
	l.mutex.Unlock()
}

func (l *List) GetItems() []string {
	return l.items
}

func (l *List) GetItem(index int) string {
	return l.items[index]
}

func (l *List) GetLength() int {
	return len(l.items)
}

func (l *List) Clear() {
	l.mutex.Lock()
	l.items = make([]string, 0)
	l.mutex.Unlock()
}

func (l *List) Iterate(yield func(string) bool) {
	for _, item := range l.items {
		if !yield(item) {
			break
		}
	}
}
