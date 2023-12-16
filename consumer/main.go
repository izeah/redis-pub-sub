package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
)

var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	RegisterConsumer(shutdown, RedisClient)

	<-shutdown

	fmt.Println("shutdown")
}

func RegisterConsumer(shutdown chan os.Signal, redisDB *redis.Client) {
	fmt.Println("starting Consumer ...")
	sub := redisDB.Subscribe(context.TODO(), "data_baru")
	defer func() {
		if sub != nil {
			_ = sub.Close()
		}
	}()
	for {
		select {
		case <-shutdown:
			fmt.Println("shutdown signal received, exiting ...")
			return
		default:
			msg, err := sub.ReceiveMessage(context.TODO())
			if err != nil {
				fmt.Println("error occured, because:", err.Error())
				continue
			}
			fmt.Println("consumed data:", msg.Payload)
		}
	}
}
