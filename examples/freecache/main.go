// freecache - coocood/freecache
// GCオーバーヘッドゼロのインメモリキャッシュ。
// キーと値は[]byte。外部サーバー不要。
package main

import (
	"fmt"
	"log"

	"github.com/coocood/freecache"
)

func main() {
	// キャッシュサイズ: 100MB
	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)

	// Set (TTL: 300秒)
	err := cache.Set(
		[]byte("user:1"),
		[]byte(`{"name":"Alice","age":30}`),
		300,
	)
	if err != nil {
		log.Fatal("Set error:", err)
	}

	// Get
	val, err := cache.Get([]byte("user:1"))
	if err != nil {
		log.Fatal("Get error:", err)
	}
	fmt.Printf("user:1 = %s\n", string(val))

	// Delete
	cache.Del([]byte("user:1"))

	// Stats
	fmt.Printf("EntryCount: %d, HitRate: %.2f%%\n",
		cache.EntryCount(), cache.HitRate()*100)
}
