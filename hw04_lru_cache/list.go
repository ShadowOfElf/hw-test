package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.head,
		Prev:  nil,
	}

	if l.head != nil {
		l.head.Prev = newItem
	} else {
		l.tail = newItem
	}

	l.head = newItem

	l.len++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.tail,
	}

	if l.tail != nil {
		l.tail.Next = newItem
	} else {
		l.head = newItem
	}

	l.tail = newItem

	l.len++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	l.Remove(i)

	i.Next = l.head
	i.Prev = nil
	if l.head != nil {
		l.head.Prev = i
	}
	l.head = i

	l.len++
}

func NewList() List {
	return &list{}
}
