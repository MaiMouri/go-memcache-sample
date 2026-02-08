package main

import (
	"context"
	"fmt"
	"log"
	"io"
	"net"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/coocood/freecache"
	"github.com/dgraph-io/ristretto"
	goCache "github.com/patrickmn/go-cache"
	rainycape "github.com/rainycape/memcache"
)

// bigcache用のサイレント設定を生成するヘルパー
func newBigcacheConfig() bigcache.Config {
	cfg := bigcache.DefaultConfig(10 * time.Minute)
	cfg.Verbose = false
	cfg.Logger = log.New(io.Discard, "", 0)
	return cfg
}

// memcachedが起動しているか確認するヘルパー
func isMemcachedAvailable() bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:11211", 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// =============================================================================
// Memcachedクライアント (外部サーバー必要)
// =============================================================================

// BenchmarkGomemcache - bradfitz/gomemcache
// memcachedサーバーが必要。起動していない場合はスキップされる。
func BenchmarkGomemcache(b *testing.B) {
	if !isMemcachedAvailable() {
		b.Skip("memcached is not running on 127.0.0.1:11211")
	}
	mc := memcache.New("127.0.0.1:11211")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set(&memcache.Item{Key: "bench_key", Value: []byte("bench_value")})
		mc.Get("bench_key")
	}
}

// BenchmarkRainycape - rainycape/memcache
// memcachedサーバーが必要。起動していない場合はスキップされる。
func BenchmarkRainycape(b *testing.B) {
	if !isMemcachedAvailable() {
		b.Skip("memcached is not running on 127.0.0.1:11211")
	}
	mc, err := rainycape.New("127.0.0.1:11211")
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set(&rainycape.Item{Key: "bench_key", Value: []byte("bench_value")})
		mc.Get("bench_key")
	}
}

// =============================================================================
// インメモリキャッシュ (外部サーバー不要)
// =============================================================================

// BenchmarkGoCache - patrickmn/go-cache
func BenchmarkGoCache(b *testing.B) {
	mc := goCache.New(5*time.Minute, 10*time.Minute)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set("bench_key", "bench_value", goCache.DefaultExpiration)
		mc.Get("bench_key")
	}
}

// BenchmarkRistretto - dgraph-io/ristretto
func BenchmarkRistretto(b *testing.B) {
	mc, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set("bench_key", "bench_value", 1)
		mc.Get("bench_key")
	}
}

// BenchmarkRistrettoWithWait - dgraph-io/ristretto (Setの反映を待つ)
func BenchmarkRistrettoWithWait(b *testing.B) {
	mc, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set("bench_key", "bench_value", 1)
		time.Sleep(1 * time.Millisecond)
		mc.Get("bench_key")
	}
}

// BenchmarkFreecache - coocood/freecache
func BenchmarkFreecache(b *testing.B) {
	cacheSize := 100 * 1024 * 1024
	mc := freecache.NewCache(cacheSize)
	expire := 60
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set([]byte("bench_key"), []byte("bench_value"), expire)
		mc.Get([]byte("bench_key"))
	}
}

// BenchmarkBigcache - allegro/bigcache
func BenchmarkBigcache(b *testing.B) {
	mc, err := bigcache.New(context.Background(), newBigcacheConfig())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set("bench_key", []byte("bench_value"))
		mc.Get("bench_key")
	}
}

// =============================================================================
// 並行アクセスベンチマーク (インメモリキャッシュのみ)
// =============================================================================

// BenchmarkGoCacheParallel - go-cacheの並行Set/Get
func BenchmarkGoCacheParallel(b *testing.B) {
	mc := goCache.New(5*time.Minute, 10*time.Minute)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mc.Set("bench_key", "bench_value", goCache.DefaultExpiration)
			mc.Get("bench_key")
		}
	})
}

// BenchmarkRistrettoParallel - ristrettoの並行Set/Get
func BenchmarkRistrettoParallel(b *testing.B) {
	mc, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mc.Set("bench_key", "bench_value", 1)
			mc.Get("bench_key")
		}
	})
}

// BenchmarkBigcacheParallel - bigcacheの並行Set/Get
func BenchmarkBigcacheParallel(b *testing.B) {
	mc, err := bigcache.New(context.Background(), newBigcacheConfig())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mc.Set("bench_key", []byte("bench_value"))
			mc.Get("bench_key")
		}
	})
}

// BenchmarkFreecacheParallel - freecacheの並行Set/Get
func BenchmarkFreecacheParallel(b *testing.B) {
	cacheSize := 100 * 1024 * 1024
	mc := freecache.NewCache(cacheSize)
	expire := 60
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mc.Set([]byte("bench_key"), []byte("bench_value"), expire)
			mc.Get([]byte("bench_key"))
		}
	})
}

// =============================================================================
// 大量キーベンチマーク (インメモリキャッシュのみ)
// =============================================================================

// BenchmarkGoCacheManyKeys - go-cacheでユニークキーのSet/Get
func BenchmarkGoCacheManyKeys(b *testing.B) {
	mc := goCache.New(5*time.Minute, 10*time.Minute)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		mc.Set(key, "bench_value", goCache.DefaultExpiration)
		mc.Get(key)
	}
}

// BenchmarkRistrettoManyKeys - ristrettoでユニークキーのSet/Get
func BenchmarkRistrettoManyKeys(b *testing.B) {
	mc, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		mc.Set(key, "bench_value", 1)
		mc.Get(key)
	}
}

// BenchmarkBigcacheManyKeys - bigcacheでユニークキーのSet/Get
func BenchmarkBigcacheManyKeys(b *testing.B) {
	mc, err := bigcache.New(context.Background(), newBigcacheConfig())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		mc.Set(key, []byte("bench_value"))
		mc.Get(key)
	}
}

// BenchmarkFreecacheManyKeys - freecacheでユニークキーのSet/Get
func BenchmarkFreecacheManyKeys(b *testing.B) {
	cacheSize := 100 * 1024 * 1024
	mc := freecache.NewCache(cacheSize)
	expire := 60
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		mc.Set([]byte(key), []byte("bench_value"), expire)
		mc.Get([]byte(key))
	}
}
