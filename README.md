# johnnycache

[![Go Report Card](https://goreportcard.com/badge/github.com/xorith/johnnycache)](https://goreportcard.com/report/github.com/xorith/johnnycache)
[![GoDoc](https://godoc.org/github.com/xorith/johnnycache/expiry?status.svg)](https://godoc.org/github.com/xorith/johnnycache/expiry)
[![Build Status](https://travis-ci.org/xorith/johnnycache.svg?branch=master)](https://travis-ci.org/xorith/johnnycache)

Simple cache implementations for Go

* Expiry - A string, interface{} concurrent expiring cache

__Getting Started__

`go get -u github.com/xorith/johnnycache/expiry`

`import "github.com/xorith/johnnycache/expiry"`

__Example__

```go
package main

import (
	"time"
	"fmt"
	"github.com/xorith/johnnycache/expiry"
)

func main() {
	c := expiry.New()
	c.Store("key", "Johnny Cache", time.Minute) // Store the value
	v, ok := c.LoadAsString("key") // LoadAsString performs type-checking for you
	if ok {
		fmt.Println(v)
	}
	// Output: Johnny Cache

	c = expiry.New()
	for i := 0; i < 1000; i++ {
		c.Store(fmt.Sprintf("key %d", i), fmt.Sprintf("value %d", i), time.Nanosecond)
	}
	time.Sleep(time.Nanosecond + 1)
	// Expired elements are only removed when an attempt to load occurs or a Sweep() is called
	fmt.Println(c.Count())

	// Since the life span was set to a nanosecond, these elements should be expired now and Load will return false
	_, ok = c.Load("key 1")
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

```