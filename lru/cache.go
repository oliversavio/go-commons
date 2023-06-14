package lru

import (
	"container/list"
	"errors"
	"sync"
)

// Read-Through Cache
type Cache struct {
	store    map[string]*list.Element
	mx       sync.RWMutex
	l        *list.List
	capacity int
}

type kv struct {
	key   string
	value interface{}
}

var ErrCacheMiss = errors.New("cache miss")

func (c *Cache) Get(key string) (interface{}, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	elem, ok := c.store[key]
	if !ok {
		return nil, ErrCacheMiss
	}
	val := elem.Value.(*kv).value
	return val, nil
}

func (c *Cache) Put(key string, value interface{}) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	var elem *list.Element
	if elem, ok := c.store[key]; ok {
		c.l.Remove(elem)
	}
	elem = c.l.PushFront(&kv{key: key, value: value})
	c.store[key] = elem

	if c.l.Len() > c.capacity {
		back := c.l.Back()
		k := back.Value.(*kv).key
		delete(c.store, k)
		c.l.Remove(back)
	}

	return nil
}

func NewCache(size int) *Cache {
	return &Cache{
		store:    map[string]*list.Element{},
		mx:       sync.RWMutex{},
		l:        list.New(),
		capacity: size,
	}
}
