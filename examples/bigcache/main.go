// bigcache - allegro/bigcache
// GCオーバーヘッドを最小化した高速インメモリキャッシュ。
// キーはstring、値は[]byte。外部サーバー不要。
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
)

func main() {
	// デフォルトTTL: 10分
	cache, err := bigcache.New(context.Background(), bigcache.Config{
		Shards:             1024,              // シャード数 (2の累乗)
		LifeWindow:         10 * time.Minute,  // エントリのTTL
		CleanWindow:        5 * time.Minute,   // 期限切れエントリの削除間隔
		MaxEntriesInWindow: 1000 * 10 * 60,    // LifeWindow内の最大エントリ数 (初期サイズ計算用)
		MaxEntrySize:       500,               // エントリの最大サイズ (バイト、初期サイズ計算用)
		HardMaxCacheSize:   256,               // キャッシュの最大サイズ (MB、0=無制限)
		Verbose:            false,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Set
	err = cache.Set("user:1", []byte(`{"name":"Alice","age":30}`))
	if err != nil {
		log.Fatal("Set error:", err)
	}

	// Get
	val, err := cache.Get("user:1")
	if err != nil {
		log.Fatal("Get error:", err)
	}
	fmt.Printf("user:1 = %s\n", string(val))

	// Delete
	err = cache.Delete("user:1")
	if err != nil {
		log.Fatal("Delete error:", err)
	}

	// Stats
	fmt.Printf("Len: %d\n", cache.Len())
}
