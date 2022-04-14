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
	mutex    chan struct{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex <- struct{}{}
	defer func() { <-cache.mutex }()

	_, itemExists := cache.items[key]

	if !itemExists && cache.queue.Len() == cache.capacity {
		backItem := cache.queue.Back()
		cache.queue.Remove(backItem)
		delete(cache.items, backItem.Value.(cacheItem).key)
	}

	var listItem *ListItem
	if itemExists {
		listItem = cache.items[key]
		listItem.Value = cacheItem{key, value}
		cache.queue.MoveToFront(listItem)
	} else {
		listItem = cache.queue.PushFront(cacheItem{key, value})
	}

	cache.items[key] = listItem

	return itemExists
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex <- struct{}{}
	defer func() { <-cache.mutex }()

	item, ok := cache.items[key]

	if ok {
		cache.queue.MoveToFront(item)

		return item.Value.(cacheItem).value, ok
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
		mutex:    make(chan struct{}, 1),
	}
}
