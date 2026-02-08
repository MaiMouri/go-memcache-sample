// go-cache - patrickmn/go-cache
// インメモリKey-Valueキャッシュ。TTLによる自動期限切れ付き。
// 外部サーバー不要。
package main

import (
	"fmt"
	"time"

	goCache "github.com/patrickmn/go-cache"
)

func main() {
	// デフォルトTTL: 5分、クリーンアップ間隔: 10分
	c := goCache.New(5*time.Minute, 10*time.Minute)

	// Set (デフォルトTTL)
	c.Set("user:1", map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}, goCache.DefaultExpiration)

	// Set (カスタムTTL)
	c.Set("session:abc", "token-xyz", 30*time.Minute)

	// Get
	if val, found := c.Get("user:1"); found {
		fmt.Printf("user:1 = %v\n", val)
	}

	// Delete
	c.Delete("user:1")

	// ItemCount
	fmt.Printf("Cache items: %d\n", c.ItemCount())
}
