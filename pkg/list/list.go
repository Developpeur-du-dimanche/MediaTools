package list

type List struct {
	items []string
}

func NewList() *List {
	return &List{
		items: make([]string, 0),
	}
}

func (l *List) AddItem(item string) {
	l.items = append(l.items, item)
}

func (l *List) RemoveItem(item string) {
	for i, v := range l.items {
		if v == item {
			l.items = append(l.items[:i], l.items[i+1:]...)
			break
		}
	}
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
	l.items = make([]string, 0)
}

func (l *List) Iterate(yield func(string) bool) {
	for _, item := range l.items {
		if !yield(item) {
			break
		}
	}
}
