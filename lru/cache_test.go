package lru

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheOps(t *testing.T) {
	key := "hello"
	val := "world"

	c := NewCache(10)

	c.Put(key, val)

	actual, err := c.Get(key)
	assert.Nil(t, err, "Error on Cache Get")
	assert.Equal(t, val, actual)
}

func TestCacheOverWrite(t *testing.T) {
	c := NewCache(10)

	for i := 0; i < 5; i++ {
		c.Put("myKey", i)
	}

	actual, err := c.Get("myKey")
	assert.Nil(t, err)
	assert.Equal(t, 4, actual)

}

func TestCacheMiss(t *testing.T) {
	c := NewCache(10)

	_, err := c.Get("notExist")
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestCacheSize(t *testing.T) {
	c := NewCache(10)

	for i := 0; i < 5; i++ {
		c.Put(fmt.Sprintf("key%d", i), i)
	}

	assert.Equal(t, 5, c.items.Len())
	assert.Equal(t, 5, len(c.store))
}

func TestCacheOverSize(t *testing.T) {
	c := NewCache(5)

	for i := 0; i < 6; i++ {
		c.Put(fmt.Sprintf("key%d", i), i)
	}

	assert.Equal(t, 5, c.items.Len())
	assert.Equal(t, 5, len(c.store))
}

func TestCacheOverWriteValues(t *testing.T) {
	c := NewCache(5)

	for i := 0; i < 100; i++ {
		c.Put(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))
	}

	for i := 0; i < 94; i++ {
		_, err := c.Get(fmt.Sprintf("k%d", i))
		assert.ErrorIs(t, err, ErrCacheMiss)
	}

	for i := 95; i < 100; i++ {
		actual, _ := c.Get(fmt.Sprintf("k%d", i))
		assert.Equal(t, fmt.Sprintf("v%d", i), actual)
	}

	assert.Equal(t, 5, c.items.Len())
	assert.Equal(t, 5, len(c.store))
}

func BenchmarkCache(b *testing.B) {
	c := NewCache(b.N)
	c.Put("key", "myval")
	for i := 0; i < b.N; i++ {
		_, _ = c.Get("key")
	}
}

func BenchmarkCacheMiss(b *testing.B) {
	c := NewCache(b.N)
	for i := 0; i < b.N; i++ {
		_, _ = c.Get("key")
	}
}

func BenchmarkCacheConcurrency(b *testing.B) {
	n := 10000
	cache := NewCache(n)

	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := fmt.Sprintf("hello-%d", i)
		cache.Put(k, nil)
		keys[i] = k
	}

	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for j := 0; j < each; j++ {
				_, err := cache.Get(k)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
}
