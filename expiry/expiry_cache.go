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

// ExpiryCache is a concurrent string key, string value expiring cache
type ExpiryCache struct {
	sync.RWMutex
	internal map[string]member
}

// NewExpiryCache returns a new type-safe string key/value concurrent expiring cache
func NewExpiryCache() *ExpiryCache {
	return &ExpiryCache{
		internal: make(map[string]member),
	}
}

// Load fetches an element by key. If the element has expired or does not exist, it returns false
func (c *ExpiryCache) Load(key string) (value interface{}, ok bool) {
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
func (c *ExpiryCache) LoadAsString(key string) (value string, ok bool) {
	v, ok := c.Load(key)
	if !ok {
		return
	}
	value, ok = v.(string)
	return
}

// Delete removes an element by key
func (c *ExpiryCache) Delete(key string) {
	c.Lock()
	delete(c.internal, key)
	c.Unlock()
}

// Store stores a single element with a specified lifespan duration
func (c *ExpiryCache) Store(key string, value interface{}, lifespan time.Duration) {
	c.Lock()
	c.internal[key] = member{
		Key:        key,
		Value:      value,
		Expiration: time.Now().Add(lifespan),
	}
	c.Unlock()
}

// Count returns a count of all elements
func (c *ExpiryCache) Count() int {
	c.RLock()
	count := len(c.internal)
	c.RUnlock()
	return count
}

// Sweep marks and then sweeps expired elements
func (c *ExpiryCache) Sweep() {
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
