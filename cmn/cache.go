package cmn

import (
	"container/list"
)

type Cache[K comparable, V any] struct {
	capacity  int
	items     map[K]*list.Element
	evictList *list.List
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		capacity:  capacity,
		items:     make(map[K]*list.Element),
		evictList: list.New(),
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	var zero V
	if elem, found := c.items[key]; found {
		c.evictList.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}
	return zero, false
}

func (c *Cache[K, V]) Set(key K, value V) {
	if elem, found := c.items[key]; found {
		c.evictList.MoveToFront(elem)
		elem.Value.(*entry[K, V]).value = value
		return
	}

	if c.evictList.Len() >= c.capacity {
		c.evict()
	}

	ent := &entry[K, V]{key, value}
	elem := c.evictList.PushFront(ent)
	c.items[key] = elem
}

func (c *Cache[K, V]) evict() {
	elem := c.evictList.Back()
	if elem != nil {
		c.evictList.Remove(elem)
		ent := elem.Value.(*entry[K, V])
		delete(c.items, ent.key)
	}
}
