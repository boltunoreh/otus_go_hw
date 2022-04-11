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
	len   int
	front *ListItem
	back  *ListItem
}

func (list list) Len() int {
	return list.len
}

func (list list) Front() *ListItem {
	if list.Len() == 0 {
		return nil
	}

	if list.front != nil {
		return list.front
	}

	return list.back
}

func (list list) Back() *ListItem {
	if list.Len() == 0 {
		return nil
	}

	if list.back != nil {
		return list.back
	}

	return list.front
}

func (list *list) PushFront(v interface{}) *ListItem {
	frontItem := list.Front()

	newItem := ListItem{
		v,
		frontItem,
		nil,
	}

	if frontItem != nil {
		frontItem.Prev = &newItem
	}

	list.front = &newItem
	list.len++

	return list.front
}

func (list *list) PushBack(v interface{}) *ListItem {
	backItem := list.Back()

	newItem := ListItem{
		v,
		nil,
		backItem,
	}

	if backItem != nil {
		backItem.Next = &newItem
	}

	list.back = &newItem
	list.len++

	return list.back
}

func (list *list) Remove(i *ListItem) {
	prevItem := i.Prev
	nextItem := i.Next

	if prevItem != nil {
		prevItem.Next = nextItem
	} else {
		list.front = nextItem
	}
	if nextItem != nil {
		nextItem.Prev = prevItem
	} else {
		list.back = prevItem
	}

	list.len--
}

func (list *list) MoveToFront(i *ListItem) {
	list.Remove(i)
	list.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
