package expiry_test

import (
	"github.com/keenan-v1/johnnycache/expiry"
	"strconv"
	"testing"
	"time"
)

type thing struct {
	Name string
}

func BenchmarkCache_Sweep(b *testing.B) {
	b.ReportAllocs()
	c := expiry.New()
	// Insert 100,000
	for i := 0; i < 100000; i++ {
		c.Store(strconv.Itoa(i), thing{strconv.Itoa(i)}, time.Microsecond*time.Duration(i))
	}
	for i := 0; i < b.N; i++ {
		c.Sweep()
	}
}

func BenchmarkCache_Store(b *testing.B) {
	b.ReportAllocs()
	c := expiry.New()
	// Insert 100,000
	for i := 0; i < 100000; i++ {
		c.Store(strconv.Itoa(i), thing{strconv.Itoa(i)}, time.Minute*1)
	}
	for i := 0; i < b.N; i++ {
		c.Store("abc"+strconv.Itoa(i), thing{strconv.Itoa(i)}, time.Minute*1)
	}
}

func BenchmarkMap_Store(b *testing.B) {
	b.ReportAllocs()
	m := make(map[string]interface{})
	// Insert 100,000
	for i := 0; i < 100000; i++ {
		m[strconv.Itoa(i)] = thing{strconv.Itoa(i)}
	}
	for i := 0; i < b.N; i++ {
		m["abc"+strconv.Itoa(i)] = thing{strconv.Itoa(i)}
	}
}

func BenchmarkCache_Load(b *testing.B) {
	b.ReportAllocs()
	c := expiry.New()
	// Insert 100,000
	for i := 0; i < 100000; i++ {
		c.Store(strconv.Itoa(i), thing{strconv.Itoa(i)}, time.Minute*1)
	}
	for i := 0; i < b.N; i++ {
		_, _ = c.Load("7281")
	}
}

func BenchmarkMap_Load(b *testing.B) {
	b.ReportAllocs()
	m := make(map[string]interface{})
	// Insert 100,000
	for i := 0; i < 100000; i++ {
		m[strconv.Itoa(i)] = thing{strconv.Itoa(i)}
	}
	for i := 0; i < b.N; i++ {
		_ = m["7281"]
	}
}
