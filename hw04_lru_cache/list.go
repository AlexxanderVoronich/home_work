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
	item := new(ListItem)
	item.Value = v
	item.Prev = nil

	l.mu.Lock()
	defer l.mu.Unlock()

	item.Next = l.first
	l.first = item
	l.length++

	if l.first.Next != nil {
		l.first.Next.Prev = l.first
	}
	if l.length == 1 {
		l.end = l.first
	}
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := new(ListItem)
	item.Value = v
	item.Next = nil

	l.mu.Lock()
	defer l.mu.Unlock()

	item.Prev = l.end
	l.end = item
	l.length++

	if l.end.Prev != nil {
		l.end.Prev.Next = l.end
	}
	if l.length == 1 {
		l.first = l.end
	}
	return item
}

func (l *list) Remove(i *ListItem) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.length == 1 && i == l.first {
		l.end = nil
		l.first = nil
		l.length = 0
		return
	}

	prev := i.Prev
	next := i.Next

	if i == l.first {
		l.first = next
	} else if i == l.end {
		l.end = prev
	}

	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
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
	prev := i.Prev
	next := i.Next

	if i == l.first {
		l.first = next
	} else if i == l.end {
		l.end = prev
	}

	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
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
	return &list{
		first:  nil,
		end:    nil,
		length: 0,
	}
}
