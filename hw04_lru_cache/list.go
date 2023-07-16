package hw04lrucache

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
	count     int
	firstItem *ListItem
	lastItem  *ListItem
}

func (l ListItem) getPrev() *ListItem {
	return l.Prev
}

func (l ListItem) getNext() *ListItem {
	return l.Next
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
	r.count += 1
	r.firstItem.Value = v
	return &item
}

func (r *list) PushBack(v interface{}) *ListItem {
	item := ListItem{}
	if r.count == 0 {
		r.firstItem = &item
		r.lastItem = &item
	} else {
		item.Prev = r.lastItem
		r.lastItem.Next = &item
		r.lastItem = &item
	}
	r.count += 1
	r.lastItem.Value = v
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
	r.count -= 1
}

func (r *list) MoveToFront(i *ListItem) *ListItem {
	if r.count == 0 {
		return nil
	}
	r.Remove(i)
	r.PushFront(i.Value)
	return i
}

func NewList() List {
	l := new(list)
	return l
}
