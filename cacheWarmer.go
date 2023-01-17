package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	goCache "github.com/patrickmn/go-cache"
)

func MemberCountForceCache() {
	var members []Member
	fmt.Println("process member_count")
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/yb_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	expire := 60 * 70

	db.Find(&members).Where("gender = ?", 1).Count(&members)

	// cache->set($member_count_key => $count, $expire);
	mc := goCache.New(expire, 10*time.Minute)

	mc.Set("test_key3", "test value3", goCache.DefaultExpiration)
	it, err := mc.Get("test_key3")
	fmt.Println(it, err)
	return count

}

func MemberCountNocache() {

}
