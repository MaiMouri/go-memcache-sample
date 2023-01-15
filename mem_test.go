package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	ristretto "github.com/dgraph-io/ristretto"
	goCache "github.com/patrickmn/go-cache"
	rainycape "github.com/rainycape/memcache"
)

// gomemcache
func BenchmarkMemLib(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	mc := memcache.New("127.0.0.1:11211")
	for i := 0; i < b.N; i++ {
		mc.Set(&memcache.Item{Key: "test_key", Value: []byte("test value")})
		mc.Get("test_key")
	}
}

// rainycape
func BenchmarkMemLib2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	mc, _ := rainycape.New("127.0.0.1:11211")
	for i := 0; i < b.N; i++ {
		mc.Set(&rainycape.Item{Key: "test_key2", Value: []byte("test value2")})
		mc.Get("test_key2")
	}
}

// go-cache
func BenchmarkMemLib3(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	mc := goCache.New(5*time.Minute, 10*time.Minute)

	for i := 0; i < b.N; i++ {
		mc.Set("test_key3", "test value3", goCache.DefaultExpiration)
		it, _ := mc.Get("test_key3")
		fmt.Println(it)
	}
}

func BenchmarkMemLib4(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	mc, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,
	})

	for i := 0; i < b.N; i++ {
		// 1 → Cost
		mc.Set("test_key4", "val", 1)
		time.Sleep(1 * time.Millisecond)
		// mc.Set("test_key3", "test value3", goCache.DefaultExpiration)
		mc.Get("test_key4")
	}
}
