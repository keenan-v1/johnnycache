package expiry

import (
	"sync"
	"time"
)

// member is an internal type used to store data in the cache
type member struct {
	Key        string
	Value      interface{}
	Expiration time.Time
}

// Cache is a concurrent key/value expiring cache
type Cache struct {
	sync.RWMutex
	internal map[string]member
}

// New returns a new key/value concurrent expiring cache
func New() *Cache {
	return &Cache{
		internal: make(map[string]member),
	}
}

// Load fetches an element by key. If the element has expired or does not exist, it returns false
func (c *Cache) Load(key string) (value interface{}, ok bool) {
	c.RLock()
	v, ok := c.internal[key]
	c.RUnlock()
	if !ok {
		return
	}
	if v.Expiration.Before(time.Now()) {
		ok = false
		c.Delete(key)
	}
	value = v.Value
	return
}

// LoadAsString fetches an element by key. If the element has expired, does not exist or is not a string, it returns false
func (c *Cache) LoadAsString(key string) (value string, ok bool) {
	v, ok := c.Load(key)
	if !ok {
		return
	}
	value, ok = v.(string)
	return
}

// Delete removes an element by key
func (c *Cache) Delete(key string) {
	c.Lock()
	delete(c.internal, key)
	c.Unlock()
}

// Store stores a single element with a specified lifespan duration
func (c *Cache) Store(key string, value interface{}, lifespan time.Duration) {
	c.Lock()
	c.internal[key] = member{
		Key:        key,
		Value:      value,
		Expiration: time.Now().Add(lifespan),
	}
	c.Unlock()
}

// Count returns a count of all elements
func (c *Cache) Count() int {
	c.RLock()
	count := len(c.internal)
	c.RUnlock()
	return count
}

// Sweep marks and then sweeps expired elements
func (c *Cache) Sweep() {
	var mark []string
	c.RLock()
	for k, v := range c.internal {
		if v.Expiration.Before(time.Now()) {
			mark = append(mark, k)
		}
	}
	c.RUnlock()
	for _, k := range mark {
		c.Delete(k)
	}
}
