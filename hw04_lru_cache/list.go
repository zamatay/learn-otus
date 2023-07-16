package hw04lrucache

import "sync"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem) *ListItem
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	mu        sync.Mutex
	count     int
	firstItem *ListItem
	lastItem  *ListItem
}

func (r *list) Len() int {
	return r.count
}

func (r *list) Front() *ListItem {
	return r.firstItem
}

func (r *list) Back() *ListItem {
	return r.lastItem
}

func (r *list) PushFront(v interface{}) *ListItem {
	item := ListItem{}
	if r.count == 0 {
		r.firstItem = &item
		r.lastItem = &item
	} else {
		item.Next = r.firstItem
		r.firstItem.Prev = &item
		r.firstItem = &item
	}
	r.count++
	r.firstItem.Value = v
	return &item
}

func (r *list) PushBack(v interface{}) *ListItem {
	r.mu.Lock()
	item := ListItem{}
	if r.count == 0 {
		r.firstItem = &item
		r.lastItem = &item
	} else {
		item.Prev = r.lastItem
		r.lastItem.Next = &item
		r.lastItem = &item
	}
	r.count++
	r.lastItem.Value = v
	defer r.mu.Unlock()
	return &item
}

func (r *list) Remove(i *ListItem) {
	if r.count == 0 {
		return
	}
	prevItem := i.Prev
	nextItem := i.Next
	if prevItem != nil {
		prevItem.Next = nextItem
	}
	if nextItem != nil {
		nextItem.Prev = prevItem
	}
	if r.firstItem == i {
		r.firstItem = nextItem
	}
	if r.lastItem == i {
		r.lastItem = prevItem
	}
	r.count--
}

func (r *list) MoveToFront(i *ListItem) *ListItem {
	r.mu.Lock()
	if r.count == 0 {
		return nil
	}
	r.Remove(i)
	r.PushFront(i.Value)
	defer r.mu.Unlock()
	return i
}

func NewList() List {
	return new(list)
}
