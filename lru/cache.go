package lru

import (
	"container/list"
	"errors"
	"sync"
)

type Cache struct {
	store map[string]*list.Element
	mu    sync.RWMutex
	items *list.List
	size  uint32
}

type kv struct {
	K string
	V interface{}
}

var ErrCacheMiss = errors.New("cache miss")

func (c *Cache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.store[key]; ok {
		c.items.Remove(elem)
	}
	elem := c.items.PushFront(&kv{
		K: key,
		V: value,
	})
	c.store[key] = elem
	c.limitCacheSize()
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	elem, ok := c.store[key]
	if !ok {
		return nil, ErrCacheMiss
	}

	return elem.Value.(*kv).V, nil
}

func (c *Cache) limitCacheSize() {
	if c.items.Len() > int(c.size) {
		last := c.items.Back()
		lkey := last.Value.(*kv).K

		c.items.Remove(last)
		delete(c.store, lkey)
	}
}

func NewCache(size int) *Cache {
	return &Cache{
		store: map[string]*list.Element{},
		mu:    sync.RWMutex{},
		items: list.New(),
		size:  uint32(size),
	}
}
