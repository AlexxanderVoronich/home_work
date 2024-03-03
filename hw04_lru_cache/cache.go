package hw04lrucache

import "sync"

// Key is a unique identifier for an item.
type Key string

// Cache is interface for lruCache.
type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Len() int
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	v, isInt := value.(int)
	if !isInt {
		return false
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	_, ok := l.items[key]
	if ok {
		l.items[key].Value = v
		l.queue.MoveToFront(l.items[key])
		return true
	}

	item := l.queue.PushFront(v)
	l.items[key] = item
	if l.queue.Len() > l.capacity {
		back := l.queue.Back()
		l.queue.Remove(back)
		for key, val := range l.items {
			if val == back {
				delete(l.items, key)
				break
			}
		}
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	value, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(l.items[key])
		return value.Value, true
	}
	return nil, false
}

func (l *lruCache) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.queue.Len()
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for key := range l.items {
		delete(l.items, key)
	}
	l.queue.Clear()
}

// NewCache is constructor of lruCache.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
