package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	counter := 0
	for {
		if res := RedisClient.Publish(context.TODO(), "data_baru", counter); res.Err() != nil {
			fmt.Println("error occured, because:", res.Err().Error())
			continue
		}
		fmt.Println("Published to topic:", "data_baru", ", value:", counter)
		counter++
		time.Sleep(3 * time.Second)
	}
}
