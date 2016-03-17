package main

import (
	"fmt"

	"gopkg.in/redis.v3"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	client.Del("hehe")
	rs := client.LPush("hehe", "fuck", "妇产科")
	if rs.Val() == 0 {
		fmt.Println("插入失败")
		return
	}

	rg := client.LRange("hehe", 0, -1)
	if rg.Err() != nil {
		fmt.Println(rg.Err().Error())
		return
	}
	r, _ := rg.Result()
	fmt.Println(r)
}
