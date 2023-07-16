package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheValue struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		mu:       sync.Mutex{},
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (r *lruCache) Set(key Key, value interface{}) bool {
	cv := cacheValue{key, value}
	r.mu.Lock()
	item, isOk := r.items[key]
	if isOk {
		r.queue.MoveToFront(item).Value = cv
	} else {
		r.items[key] = r.queue.PushFront(cv)
	}
	if r.queue.Len() > r.capacity {
		item = r.queue.Back()
		if item != nil {
			r.remove(item)
		}
	}
	r.mu.Unlock()
	return isOk
}

func (r *lruCache) Get(key Key) (interface{}, bool) {
	r.mu.Lock()
	value, isOk := r.items[key]
	r.mu.Unlock()
	if isOk {
		r.queue.MoveToFront(value)
		i := r.getCacheValue(value)
		return i.value, true
	}
	return nil, false
}

func (r *lruCache) Clear() {
	r.mu.Lock()
	r.items = make(map[Key]*ListItem, 0)
	r.queue = NewList()
	r.mu.Unlock()
}

func (r *lruCache) remove(item *ListItem) {
	v := r.queue.Back()
	if v != nil {
		i := r.getCacheValue(v)
		delete(r.items, i.key)
		r.queue.Remove(item)
	}
}

func (r *lruCache) getCacheValue(v *ListItem) cacheValue {
	return v.Value.(cacheValue)
}
