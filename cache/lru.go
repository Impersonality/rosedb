package cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	capacity  int
	cacheMap  map[interface{}]*list.Element
	cacheList *list.List
	mu        sync.Mutex
}

type lruItem struct {
	key   interface{}
	value interface{}
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity:  capacity,
		cacheMap:  make(map[interface{}]*list.Element),
		cacheList: list.New(),
	}
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.get(key)
}

func (c *Cache) Set(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value)
}

func (c *Cache) get(key interface{}) (interface{}, bool) {
	ele, ok := c.cacheMap[key]
	if ok {
		c.cacheList.MoveToFront(ele)
		return ele.Value, true
	}
	return nil, false
}

func (c *Cache) set(key, value interface{}) {
	ele, ok := c.cacheMap[key]
	if ok {
		item := c.cacheMap[key].Value.(*lruItem)
		item.value = value
		// todo c.cacheMap[key].Value = &lruItem{key: key, value: value}
		c.cacheList.MoveToFront(ele)
	} else {
		ele = c.cacheList.PushFront(&lruItem{key: key, value: value})
		c.cacheMap[key] = ele

		if c.cacheList.Len() > c.capacity {
			c.removeOldest()
		}
	}
}

func (c *Cache) removeOldest() {
	ele := c.cacheList.Back()
	c.cacheList.Remove(ele)
	item := ele.Value.(*lruItem)
	delete(c.cacheMap, item.key)
}
