package main

import (
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	memcache2 "github.com/rainycape/memcache"
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

func BenchmarkMemLib2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	mc, _ := memcache2.New("127.0.0.1:11211")
	for i := 0; i < b.N; i++ {
		mc.Set(&memcache2.Item{Key: "test_key2", Value: []byte("test value2")})
		mc.Get("test_key2")
	}
}
