package main

import (
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

func BenchmarkMemLib(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	mc := memcache.New("127.0.0.1:11211")
	for i := 0; i < b.N; i++ {
		mc.Set(&memcache.Item{Key: "test_key", Value: []byte("test value")})
		mc.Get("test_key")
	}
}
