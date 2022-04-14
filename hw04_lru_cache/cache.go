package hw04lrucache

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
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	_, itemExists := cache.Get(key)

	if itemExists == false && cache.queue.Len() == cache.capacity {
		backItem := cache.queue.Back()
		cache.queue.Remove(backItem)
	}

	var listItem *ListItem
	if itemExists {
		listItem = cache.items[key]
		listItem.Value = value
		cache.queue.MoveToFront(listItem)
	} else {
		listItem = cache.queue.PushFront(value)
	}

	cache.items[key] = listItem

	return itemExists
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := cache.items[key]

	if ok {
		cache.queue.MoveToFront(item)
		return item.Value, ok
	}

	return nil, ok
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
