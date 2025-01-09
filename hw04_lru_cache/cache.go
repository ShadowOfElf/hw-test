package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

type lruItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	item := lruItem{
		key:   key,
		value: value,
	}

	if getItem, ok := l.items[key]; ok {
		getItem.Value = item
		l.queue.MoveToFront(getItem)
		return true
	}

	listItem := l.queue.PushFront(item)
	l.items[key] = listItem

	if l.queue.Len() > l.capacity {
		backItem := l.queue.Back()
		if backItem != nil {
			delete(l.items, backItem.Value.(lruItem).key)
			l.queue.Remove(backItem)
		}
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		return item.Value.(lruItem).value, ok
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.items = map[Key]*ListItem{}
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
