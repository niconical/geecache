package lru

import (
	"container/list"
	"errors"
)

type EvictCallback func(key string, value Value)

var _ LRUCache = (*Cache)(nil)

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	cache    map[string]*list.Element
	ll       *list.List
	maxBytes int64
	nbytes   int64
	// optional and executed when an entry is purged.
	OnEvicted EvictCallback
}

func NewCache(maxBytes int64, onEvicted EvictCallback) (*Cache, error) {
	if maxBytes <= 0 {
		return nil, errors.New("must provide a postive maxBytes")
	}
	return &Cache{
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		maxBytes:  maxBytes,
		OnEvicted: onEvicted,
	}, nil
}

func (c *Cache) Add(key string, value Value) bool {
	// Check for existing item
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		e.Value.(*entry).v = value
		return false
	}

	ent := &entry{
		k: key,
		v: value,
	}
	entry := c.ll.PushFront(ent)
	c.cache[key] = entry
	c.nbytes += int64(len(ent.k)) + int64(ent.v.Len())

	evict := c.nbytes > c.maxBytes

	if evict {
		for c.nbytes != 0 && c.maxBytes < c.nbytes {
			c.RemoveOldest()
		}
	}
	return evict
}

func (c *Cache) Get(key string) (Value, bool) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		return e.Value.(*entry).v, true
	}
	return nil, false
}

func (c *Cache) Contains(key string) (ok bool) {
	_, ok = c.cache[key]
	return
}

func (c *Cache) Peek(key string) (value Value, ok bool) {
	if v, ok := c.cache[key]; ok {
		return v.Value.(*entry).v, true
	}
	return nil, false
}

func (c *Cache) Remove(key string) bool {
	if v, ok := c.cache[key]; ok {
		c.removeElement(v)
		return true
	}
	return false
}

func (c *Cache) RemoveOldest() (string, Value, bool) {
	e := c.ll.Back()
	if e != nil {
		c.removeElement(e)
		ent := e.Value.(*entry)
		return ent.k, ent.v, true
	}
	return "", nil, false
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	ent := e.Value.(*entry)
	delete(c.cache, ent.k)
	c.nbytes -= ent.v.Len() + int64(len(ent.k))
	if c.OnEvicted != nil {
		c.OnEvicted(ent.k, ent.v)
	}
}

func (c *Cache) GetOldest() (string, Value, bool) {
	e := c.ll.Back()
	if e != nil {
		ent := e.Value.(*entry)
		return ent.k, ent.v, true
	}
	return "", nil, false
}

func (c *Cache) Keys() []string {
	keys := make([]string, len(c.cache))
	i := 0
	for iter := c.ll.Back(); iter != nil; iter = iter.Prev() {
		keys[i] = iter.Value.(*entry).k
		i++
	}
	return keys
}

func (c *Cache) Len() int64 {
	return int64(c.nbytes)
}

func (c *Cache) Purge() {
	for k, v := range c.cache {
		if c.OnEvicted != nil {
			c.OnEvicted(k, v.Value.(*entry).v)
		}
		delete(c.cache, k)
	}
	c.ll.Init()
	c.nbytes = 0
}

// RemaxBytes changes the cache maxBytes.
func (c *Cache) RemaxBytes(maxBytes int64) (evicted int64, err error) {
	if maxBytes <= 0 {
		return 0, errors.New("must provide a postive maxBytes")
	}

	diff := int64(c.nbytes) - maxBytes
	if diff < 0 {
		diff = 0
	}
	var i int64 = 0
	for ; c.nbytes > maxBytes; i++ {
		c.RemoveOldest()
	}
	c.maxBytes = maxBytes
	return i, nil
}
