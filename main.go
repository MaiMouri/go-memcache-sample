package main

import (
	"fmt"
	"time"

	// "github.com/rainycape/memcache"
	"github.com/dgraph-io/ristretto"
	_ "github.com/dgraph-io/ristretto"
	"github.com/jinzhu/gorm"
	_ "github.com/patrickmn/go-cache"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	dbInit()
	MemberCountForceCache()
	/* ristretto */
	mc, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}

	// 1 → Cost
	mc.Set("test_key4", "val", 1)
	time.Sleep(10 * time.Millisecond)

	if it, found := mc.Get("test_key4"); found {
		fmt.Println(it)
	}

	// mc, _ := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
	// mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

	// it, err := mc.Get("foo")
	// fmt.Println(it, err)

	/* go-cache */
	// mc := goCache.New(5*time.Minute, 10*time.Minute)

	// mc.Set("test_key3", "test value3", goCache.DefaultExpiration)
	// it, err := mc.Get("test_key3")
	// fmt.Println(it, err)
}

type Member struct {
	gorm.Model
	Gender int `gorm:"column:gender"`
	Status int `gorm:"column:status"`
}

func dbInit() {
	fmt.Println("Start dbInit")
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/yb_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	var members []Member
	db.Find(&members)
	fmt.Println(members)

}
