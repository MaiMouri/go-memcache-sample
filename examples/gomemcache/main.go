// gomemcache - bradfitz/gomemcache
// Memcachedサーバーに接続するクライアントライブラリ。
// 事前にmemcachedサーバーの起動が必要。
package main

import (
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	// Memcachedサーバーに接続
	mc := memcache.New("127.0.0.1:11211")

	// Set
	err := mc.Set(&memcache.Item{
		Key:        "user:1",
		Value:      []byte(`{"name":"Alice","age":30}`),
		Expiration: 300, // 秒単位のTTL
	})
	if err != nil {
		log.Fatal("Set error:", err)
	}

	// Get
	item, err := mc.Get("user:1")
	if err != nil {
		log.Fatal("Get error:", err)
	}
	fmt.Printf("Key: %s, Value: %s\n", item.Key, string(item.Value))

	// Delete
	err = mc.Delete("user:1")
	if err != nil {
		log.Fatal("Delete error:", err)
	}
	fmt.Println("Deleted user:1")
}
