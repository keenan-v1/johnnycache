package expiry_test

import (
	"fmt"
	"github.com/keenan-v1/johnnycache/expiry"
	"time"
)

func ExampleCache_LoadAsString() {
	c := expiry.New()
	c.Store("key", "Johnny Cache", time.Minute) // Store the value
	v, ok := c.LoadAsString("key")              // LoadAsString performs type-checking for you
	if ok {
		fmt.Println(v)
	}
	// Output: Johnny Cache
}

func ExampleCache_Sweep() {
	c := expiry.New()
	for i := 0; i < 1000; i++ {
		c.Store(fmt.Sprintf("key %d", i), fmt.Sprintf("value %d", i), time.Nanosecond)
	}
	time.Sleep(time.Nanosecond + 1)
	// Expired elements are only removed when an attempt to load occurs or a Sweep() is called
	fmt.Println(c.Count())

	// Since the life span was set to a nanosecond, these elements should be expired now and Load will return false
	_, ok := c.Load("key 1")
	if !ok {
		fmt.Println("expired")
	}
	fmt.Println(c.Count())

	// Sweep can be called directly, or passed to a goroutine for background sweeping
	c.Sweep()
	fmt.Println(c.Count())

	// Output:
	// 1000
	// expired
	// 999
	// 0
}
