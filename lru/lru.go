package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	// memory
	maxBytes int64
	nbytes   int64
	// lru impl
	cache map[string]*list.Element
	ll    *list.List
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	k string
	v Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		d := e.Value.(*entry)
		return d.v, true
	}
	return
}

// removes the oldest item by LRU algorithm
func (c *Cache) RemoveOldest() {
	if e := c.ll.Back(); e != nil {
		d := c.ll.Remove(e).(*entry)
		delete(c.cache, d.k)
		c.nbytes -= int64(len(d.k)) + int64(d.v.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(d.k, d.v)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		d := e.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(d.v.Len())
		d.v = value
	} else {
		e := c.ll.PushFront(&entry{key, value})
		c.cache[key] = e
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}
