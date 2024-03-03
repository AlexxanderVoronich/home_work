package hw04lrucache

import (
	"sync"
)

// List is interface for list.
type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
}

// ListItem is node for struct list.
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first  *ListItem
	end    *ListItem
	length int
	mu     sync.RWMutex
}

func (l *list) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.length
}

func (l *list) Front() *ListItem {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.first
}

func (l *list) Back() *ListItem {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.end
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.mu.Lock()
	defer l.mu.Unlock()

	item := &ListItem{Value: v, Next: l.first}
	if l.first != nil {
		l.first.Prev = item
	}
	l.first = item
	if l.end == nil {
		l.end = item
	}
	l.length++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.mu.Lock()
	defer l.mu.Unlock()

	item := &ListItem{Value: v, Prev: l.end}
	if l.end != nil {
		l.end.Next = item
	}
	l.end = item
	if l.first == nil {
		l.first = item
	}
	l.length++
	return item
}

func (l *list) Remove(i *ListItem) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else { // i - first
		l.first = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else { // i - end
		l.end = i.Prev
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if i == l.first {
		return
	}

	// remove old element
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else { // i - first
		l.first = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else { // i - end
		l.end = i.Prev
	}

	// push element front
	i.Prev = nil
	i.Next = l.first
	l.first = i

	if l.first.Next != nil {
		l.first.Next.Prev = l.first
	}
	if l.length == 1 {
		l.end = l.first
	}
}

func (l *list) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.first = nil
	l.end = nil
	l.length = 0
}

// NewList is constructor for list.
func NewList() List {
	return &list{first: nil, end: nil, length: 0}
}
