package main

import (
	"fmt"

	"github.com/rainycape/memcache"
)

func main() {
	mc, _ := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
	mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

	it, err := mc.Get("foo")
	fmt.Println(it, err)
}
