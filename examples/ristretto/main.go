// ristretto - dgraph-io/ristretto
// 高性能なインメモリキャッシュ。TinyLFUベースのアドミッションポリシー付き。
// 外部サーバー不要。
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/ristretto"
)

func main() {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // 追跡するキーの数（MaxCostの10倍を推奨）
		MaxCost:     1 << 30, // キャッシュの最大コスト（1GB）
		BufferItems: 64,      // バッファのアイテム数
	})
	if err != nil {
		log.Fatal(err)
	}

	// Set (cost=1)
	cache.Set("user:1", map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}, 1)

	// SetWithTTL
	cache.SetWithTTL("session:abc", "token-xyz", 1, 30*time.Minute)

	// ristrettoは非同期でSetが反映されるためwaitが必要
	time.Sleep(10 * time.Millisecond)

	// Get
	if val, found := cache.Get("user:1"); found {
		fmt.Printf("user:1 = %v\n", val)
	}

	// Delete
	cache.Del("user:1")

	// Metrics
	fmt.Printf("Hits: %d, Misses: %d\n",
		cache.Metrics.Hits(), cache.Metrics.Misses())
}
