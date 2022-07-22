package geecache

import (
	"sync"

	"github.com/niconical/geecache/lru"
)

type concurrentCache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *concurrentCache) add(key string, value ByteView) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru, err = lru.NewCache(c.cacheBytes, nil)
		if err != nil {
			return
		}
	}
	c.lru.Add(key, value)
	return
}

func (c *concurrentCache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
