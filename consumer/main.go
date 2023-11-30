package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	sub := RedisClient.Subscribe(context.TODO(), "data_baru")
	defer func() {
		if sub != nil {
			_ = sub.Close()
		}
	}()
	fmt.Println("Ready for consuming from topic:", "data_baru")
	for {
		msg, err := sub.ReceiveMessage(context.TODO())
		if err != nil {
			fmt.Println("error occured, because:", err.Error())
			continue
		}
		fmt.Println("Consumed data:", msg.Payload)
	}
}
