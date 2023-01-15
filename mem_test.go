package main

import (
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
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
		mc.Get("test_key3")
	}
}
