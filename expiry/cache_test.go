package expiry_test

import (
	"fmt"
	assertion "github.com/stretchr/testify/assert"
	"github.com/xorith/johnnycache/expiry"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	assert := assertion.New(t)
	c := expiry.New()
	assert.NotNil(c)
}

func TestCache_Store(t *testing.T) {
	assert := assertion.New(t)
	c := expiry.New()
	assert.NotNil(c)

	key, val := "test_key", "test_val"
	_, ok := c.LoadAsString(key)
	assert.False(ok, "should not be loaded")
	c.Store(key, val, time.Second*10)
	result, ok := c.LoadAsString(key)
	assert.True(ok, "should be loaded")
	assert.Equal(val, result, "should be equal")
}

func TestCache_Load(t *testing.T) {
	assert := assertion.New(t)
	c := expiry.New()
	assert.NotNil(c)

	key, val := "test_key", "test_val"
	_, ok := c.LoadAsString(key)
	assert.False(ok, "should not be loaded")
	c.Store(key, val, time.Millisecond*100)
	result, ok := c.LoadAsString(key)
	assert.True(ok, "should be loaded")
	assert.Equal(val, result, "should be equal")
	time.Sleep(1 * time.Second)
	_, ok = c.LoadAsString(key)
	assert.False(ok, "should have expired")
}

func TestCache_Delete(t *testing.T) {
	assert := assertion.New(t)
	c := expiry.New()
	assert.NotNil(c)

	key, val := "test_key", "test_val"
	_, ok := c.LoadAsString(key)
	assert.False(ok, "should not be loaded")
	c.Store(key, val, time.Second*1)
	_, ok = c.LoadAsString(key)
	assert.True(ok, "should be loaded")
	c.Delete(key)
	_, ok = c.LoadAsString(key)
	assert.False(ok, "should be deleted")
}

func TestCache_Sweep(t *testing.T) {
	assert := assertion.New(t)
	c := expiry.New()
	assert.NotNil(c)

	key, val := "test_key", "test_val"
	_, ok := c.LoadAsString(key)
	assert.False(ok, "should not be loaded")
	c.Store(key, val, time.Millisecond*10)
	count := c.Count()
	assert.NotEqual(0, count, "should have 1 element")
	time.Sleep(1 * time.Second)
	c.Sweep()
	count = c.Count()
	assert.Equal(0, count, "should have 0 elements")
}

func TestConcurrency(t *testing.T) {
	assert := assertion.New(t)
	c := expiry.New()
	assert.NotNil(c)
	for i := 0; i < 100; i++ {
		go func(i int) {
			key, val := fmt.Sprintf("key%d", i), "test"
			time.Sleep(time.Duration(i) * time.Millisecond)
			c.Store(key, val, time.Millisecond*100*time.Duration(i))
			c.Sweep()
			result, ok := c.LoadAsString(key)
			assert.True(ok, "expected result to be in expiry")
			assert.Equal(val, result)
			time.Sleep(time.Millisecond*100*time.Duration(i) + 100)
			_, ok = c.LoadAsString(key)
			assert.False(ok, "should be expired by now")
		}(i)
	}
}
