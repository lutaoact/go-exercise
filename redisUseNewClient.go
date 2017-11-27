package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	err = client.Set("hello", "girl", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("hello").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("hello", val)

	val2, err := client.Get("hello2").Result()
	if err == redis.Nil {
		fmt.Println("hello2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}
