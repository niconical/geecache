package lru

// Elements stored in list
type entry struct {
	k string
	v Value
}

type Value interface {
	Len() int64
}

// LRUCache is the interface for simple LRU cache
type LRUCache interface {
	// Adds a value to the cache, returns true if an eviction occurred
	// and updates the "recently used"-ness of the key
	Add(key string, value Value) bool

	// Returns key's value from the cache and
	// updates the "recently used"-ness of the key. #value, isFound
	Get(key string) (value Value, ok bool)

	// Checks if a key exists in cache without updating the recent-ness.
	Contains(key string) (ok bool)

	// Returns key's value without updating the "recently used"-ness of the key.
	Peek(key string) (value Value, ok bool)

	// Removes a key from the cache
	Remove(key string) bool

	// Removes the oldest entry from cache.
	RemoveOldest() (string, Value, bool)

	// Removes the oldest entry from cache without updating the
	// "recently used"-ness of the key.
	GetOldest() (string, Value, bool)

	// Returns a slice of the keys in the cache, from oldest to newest.
	Keys() []string

	// Returns the number of items in the cache
	Len() int64

	// Clears all cache entries
	Purge()

	// RemaxBytess cache, returning number evicted
	RemaxBytes(int64) (int64, error)
}
